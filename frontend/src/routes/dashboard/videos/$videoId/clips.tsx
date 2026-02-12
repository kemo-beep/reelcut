import { createFileRoute, Link, useRouterState } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState, useEffect } from 'react'
import { Scissors, Plus, Sparkles, Check, CheckCheck } from 'lucide-react'
import { getVideo } from '../../../../lib/api/videos'
import { listClips, createClip } from '../../../../lib/api/clips'
import { suggestClips, type ClipSuggestion } from '../../../../lib/api/analysis'
import { Button } from '../../../../components/ui/button'
import { Skeleton } from '../../../../components/ui/skeleton'
import { ErrorState } from '../../../../components/ui/error-state'
import { Badge } from '../../../../components/ui/badge'
import { toast } from 'sonner'
import type { Clip } from '../../../../types'

export const Route = createFileRoute('/dashboard/videos/$videoId/clips')({
  component: VideoClipsPage,
})

function formatTimeForLabel(sec: number): string {
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

function clipNameFromSuggestion(suggestion: ClipSuggestion, index: number): string {
  const start = formatTimeForLabel(suggestion.start_time)
  if (suggestion.transcript && suggestion.transcript.length > 0) {
    const truncated = suggestion.transcript.slice(0, 40)
    return truncated.length < suggestion.transcript.length ? `${truncated}…` : truncated
  }
  return `Clip ${index + 1} (${start})`
}

function VideoClipsPage() {
  const { videoId } = Route.useParams()
  const queryClient = useQueryClient()
  const locationState = useRouterState({ select: (s) => s.location.state })
  const [suggestions, setSuggestions] = useState<ClipSuggestion[]>([])
  useEffect(() => {
    const fromState = (locationState as { suggestions?: ClipSuggestion[] })?.suggestions
    if (Array.isArray(fromState) && fromState.length > 0) setSuggestions(fromState)
  }, [locationState])
  const { data: videoData, isLoading: videoLoading, error: videoError } = useQuery({
    queryKey: ['video', videoId],
    queryFn: () => getVideo(videoId),
  })
  const { data: clipsData, isLoading: clipsLoading } = useQuery({
    queryKey: ['clips', { video_id: videoId }],
    queryFn: () => listClips({ video_id: videoId, per_page: 100 }),
    enabled: !!videoData?.video,
  })
  const suggestMut = useMutation({
    mutationFn: () => suggestClips(videoId),
    onSuccess: (data) => {
      const list = data.suggestions ?? []
      setSuggestions(list)
      if (list.length > 0) {
        toast.success(`${list.length} clip suggestions ready`)
      } else {
        toast.info('No suggestions found for this video')
      }
    },
    onError: () => toast.error('Failed to get suggestions'),
  })
  const acceptClipMut = useMutation({
    mutationFn: ({ suggestion, index }: { suggestion: ClipSuggestion; index: number }) =>
      createClip({
        video_id: videoId,
        name: clipNameFromSuggestion(suggestion, index),
        start_time: suggestion.start_time,
        end_time: suggestion.end_time,
        aspect_ratio: '9:16',
        virality_score: suggestion.virality_score ?? undefined,
        from_suggestion: true,
      }),
    onSuccess: (_, { suggestion }) => {
      queryClient.invalidateQueries({ queryKey: ['clips', { video_id: videoId }] })
      setSuggestions((prev) => prev.filter((s) => s !== suggestion))
      toast.success('Clip created')
    },
    onError: () => toast.error('Failed to create clip'),
  })
  const acceptAllMut = useMutation({
    mutationFn: async () => {
      const list = suggestions
      for (let i = 0; i < list.length; i++) {
        const s = list[i]
        await createClip({
          video_id: videoId,
          name: clipNameFromSuggestion(s, i),
          start_time: s.start_time,
          end_time: s.end_time,
          aspect_ratio: '9:16',
          virality_score: s.virality_score ?? undefined,
          from_suggestion: true,
        })
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clips', { video_id: videoId }] })
      setSuggestions([])
      toast.success('All suggestions added as clips')
    },
    onError: () => toast.error('Failed to create some clips'),
  })

  if (videoLoading || !videoData?.video) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-32 w-full" />
      </div>
    )
  }
  if (videoError) {
    return (
      <div className="space-y-6">
        <h1 className="text-2xl font-bold text-[var(--app-fg)]">Clips</h1>
        <ErrorState message="Video not found." onRetry={() => queryClient.invalidateQueries({ queryKey: ['video', videoId] })} />
      </div>
    )
  }

  const video = videoData.video
  const clips = clipsData?.data?.clips ?? []

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <Link to="/dashboard/videos" className="text-sm text-[var(--app-accent)] hover:underline">
            Videos
          </Link>
          <span className="mx-2 text-[var(--app-fg-muted)]">/</span>
          <Link
            to="/dashboard/videos/$videoId"
            params={{ videoId }}
            className="text-sm text-[var(--app-accent)] hover:underline"
          >
            {video.original_filename}
          </Link>
          <h1 className="text-2xl font-bold text-[var(--app-fg)] mt-1">Clips from this video</h1>
        </div>
        <div className="flex gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => suggestMut.mutate()}
            disabled={suggestMut.isPending}
            className="border-[var(--app-border)]"
          >
            <Sparkles size={16} className="mr-1" />
            {suggestMut.isPending ? 'Analyzing…' : 'AI suggest clips'}
          </Button>
          <Link to="/dashboard/clips" search={{ videoId }}>
            <Button size="sm" className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]">
              <Plus size={16} className="mr-1" />
              Create clip
            </Button>
          </Link>
        </div>
      </div>

      {suggestions.length > 0 && (
        <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="font-semibold text-[var(--app-fg)]">AI suggestions</h2>
            <Button
              size="sm"
              onClick={() => acceptAllMut.mutate()}
              disabled={acceptAllMut.isPending}
              className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]"
            >
              <CheckCheck size={14} className="mr-1" />
              {acceptAllMut.isPending ? 'Creating…' : 'Accept all'}
            </Button>
          </div>
          <ul className="space-y-3">
            {suggestions.map((suggestion, index) => (
              <li
                key={`${suggestion.start_time}-${suggestion.end_time}-${index}`}
                className="flex items-center justify-between rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] p-3"
              >
                <div className="min-w-0 flex-1">
                  <p className="text-caption text-[var(--app-fg-muted)]">
                    {formatTimeForLabel(suggestion.start_time)} – {formatTimeForLabel(suggestion.end_time)}
                    {suggestion.duration != null && ` · ${Math.round(suggestion.duration)}s`}
                    {suggestion.virality_score != null && ` · Score ${Math.round(suggestion.virality_score)}`}
                  </p>
                  {suggestion.transcript && (
                    <p className="mt-0.5 text-sm text-[var(--app-fg)] truncate" title={suggestion.transcript}>
                      {suggestion.transcript}
                    </p>
                  )}
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  className="ml-3 shrink-0"
                  onClick={() => acceptClipMut.mutate({ suggestion, index })}
                  disabled={acceptClipMut.isPending}
                >
                  <Check size={14} className="mr-1" />
                  Accept
                </Button>
              </li>
            ))}
          </ul>
        </div>
      )}

      {clipsLoading ? (
        <ul className="space-y-3">
          {[1, 2, 3].map((i) => (
            <li key={i}>
              <Skeleton className="h-16 w-full rounded-xl" />
            </li>
          ))}
        </ul>
      ) : clips.length === 0 ? (
        <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 text-center">
          <Scissors size={40} className="mx-auto text-[var(--app-fg-muted)]" />
          <p className="mt-2 font-medium text-[var(--app-fg)]">No clips yet</p>
          <p className="text-caption text-[var(--app-fg-muted)] mt-1">
            Use AI suggest or create a clip manually from the Clips page.
          </p>
          <div className="mt-4 flex justify-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => suggestMut.mutate()}
              disabled={suggestMut.isPending}
            >
              AI suggest clips
            </Button>
            <Link to="/dashboard/clips">
              <Button size="sm" className="bg-[var(--app-accent)] text-[#0a0a0b]">
                Go to Clips
              </Button>
            </Link>
          </div>
        </div>
      ) : (
        <ul className="space-y-3">
          {clips.map((clip: Clip) => (
            <li
              key={clip.id}
              className="flex items-center justify-between rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4"
            >
              <div>
                <Link
                  to="/dashboard/clips/$clipId"
                  params={{ clipId: clip.id }}
                  className="font-medium text-[var(--app-fg)] hover:text-[var(--app-accent)]"
                >
                  {clip.name}
                </Link>
                <p className="text-caption text-[var(--app-fg-muted)] mt-0.5">
                  {clip.duration_seconds != null ? `${Math.round(clip.duration_seconds)}s` : '—'} · {clip.aspect_ratio} · {clip.status}
                </p>
              </div>
              <div className="flex items-center gap-2">
                {clip.is_ai_suggested && (
                  <Badge variant="secondary" className="text-xs">AI</Badge>
                )}
                <Link to="/editor/$clipId" params={{ clipId: clip.id }}>
                  <Button variant="outline" size="sm">Edit</Button>
                </Link>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
