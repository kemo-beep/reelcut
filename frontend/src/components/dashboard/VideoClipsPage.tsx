import { Link, useParams, useRouterState } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState, useEffect, useRef, useCallback } from 'react'
import { Scissors, Plus, Sparkles, Check, CheckCheck } from 'lucide-react'
import { getVideo, getPlaybackUrl } from '../../lib/api/videos'
import { listClips, createClip, getClipPlaybackUrl, updateClip } from '../../lib/api/clips'
import { getTranscriptionByVideoId } from '../../lib/api/transcriptions'
import { suggestClips, type ClipSuggestion } from '../../lib/api/analysis'
import { Button } from '../ui/button'
import { Skeleton } from '../ui/skeleton'
import { ErrorState } from '../ui/error-state'
import { VideoPlayer } from '../video/VideoPlayer'
import { ClipTimeline, type ClipTimelineSegment } from '../clip/ClipTimeline'
import { ClipStripCard } from '../clip/ClipStripCard'
import { TranscriptViewer } from '../transcription/TranscriptViewer'
import { toast } from 'sonner'
import { cn } from '../../lib/utils'
import type { Clip, Video } from '../../types'
import type { TranscriptSegment } from '../../types'

function MainVideoView({
  videoId,
  video,
  segments,
}: {
  videoId: string
  video: Video
  segments: TranscriptSegment[]
}) {
  const playerRef = useRef<{ seek: (time: number) => void; getCurrentTime: () => number }>(null)
  const [currentTime, setCurrentTime] = useState(0)

  const { data: playbackData, isFetching: playbackLoading } = useQuery({
    queryKey: ['video-playback', videoId],
    queryFn: () => getPlaybackUrl(videoId),
    enabled: !!videoId,
  })
  const playbackUrl = playbackData?.url ?? null

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[1fr_minmax(320px,400px)] gap-6 items-start">
      <div className="rounded-xl overflow-hidden border border-[var(--app-border)] bg-[var(--app-bg)]">
        {playbackLoading && (
          <div className="aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
            Loading…
          </div>
        )}
        {!playbackLoading && playbackUrl && (
          <VideoPlayer
            ref={playerRef}
            src={playbackUrl}
            initialTime={currentTime}
            onTimeUpdate={setCurrentTime}
          />
        )}
        {!playbackLoading && !playbackUrl && (
          <div className="aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
            Playback not available for this video.
          </div>
        )}
      </div>
      <div
        className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4 lg:sticky lg:top-6"
        id="transcript-panel"
      >
        <h2 className="font-semibold text-[var(--app-fg)] mb-3">Transcript: {video.original_filename}</h2>
        {segments.length > 0 ? (
          <TranscriptViewer
            segments={segments}
            currentTime={currentTime}
            onSeek={(t) => playerRef.current?.seek(t)}
            className="max-h-[50vh]"
          />
        ) : (
          <p className="text-caption text-[var(--app-fg-muted)]">No transcript available.</p>
        )}
      </div>
    </div>
  )
}

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

function segmentsForClip(
  segments: TranscriptSegment[],
  clipStart: number,
  clipEnd: number
): TranscriptSegment[] {
  return segments
    .filter((s) => s.end_time > clipStart && s.start_time < clipEnd)
    .map((s) => ({
      ...s,
      start_time: Math.max(0, s.start_time - clipStart),
      end_time: Math.min(clipEnd - clipStart, s.end_time - clipStart),
    }))
}

export function VideoClipsPage() {
  const { videoId } = useParams({ from: '/dashboard/videos/$videoId/clips' })
  const queryClient = useQueryClient()
  const locationState = useRouterState({ select: (s) => s.location.state })
  const playerRef = useRef<{ seek: (time: number) => void; getCurrentTime: () => number }>(null)

  const [suggestions, setSuggestions] = useState<ClipSuggestion[]>([])
  const [selectedClipId, setSelectedClipId] = useState<string | null>(null)
  const [clipPlaybackCurrentTime, setClipPlaybackCurrentTime] = useState(0)
  type ViewMode = 'clips' | 'main-video'
  const [viewMode, setViewMode] = useState<ViewMode>('clips')

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
  const { data: transData } = useQuery({
    queryKey: ['transcription', videoId],
    queryFn: () => getTranscriptionByVideoId(videoId),
    enabled: !!videoData?.video,
  })

  const clips = (clipsData?.data?.clips ?? []) as Clip[]
  const sortedClips = [...clips].sort((a, b) => a.start_time - b.start_time)
  const selectedClip = selectedClipId ? clips.find((c) => c.id === selectedClipId) ?? null : null

  const { data: playbackData, isFetching: playbackLoading } = useQuery({
    queryKey: ['clip-playback', selectedClipId],
    queryFn: () => getClipPlaybackUrl(selectedClipId!),
    enabled: !!selectedClipId && !!(selectedClip?.storage_path != null && selectedClip.storage_path !== ''),
  })
  const playbackUrl = playbackData?.url ?? null

  useEffect(() => {
    if (clips.length > 0 && !selectedClipId) {
      const firstWithVideo = clips.find((c) => c.storage_path != null && c.storage_path !== '')
      setSelectedClipId((firstWithVideo ?? clips[0]).id)
    }
  }, [clips.length, selectedClipId, clips])

  const suggestMut = useMutation({
    mutationFn: () =>
      suggestClips(videoId, {
        min_duration: 7,
        max_duration: 60,
        max_suggestions: 20,
      }),
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
      setViewMode('clips')
      toast.success('All suggestions added as clips')
    },
    onError: () => toast.error('Failed to create some clips'),
  })

  const updateClipMut = useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: { start_time: number; end_time: number } }) =>
      updateClip(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clips', { video_id: videoId }] })
      toast.success('Clip updated')
    },
    onError: () => toast.error('Failed to update clip'),
  })

  const handleTimelineSeek = useCallback(
    (time: number) => {
      const clip = sortedClips.find((c) => time >= c.start_time && time <= c.end_time)
      if (clip) {
        setSelectedClipId(clip.id)
        const offset = time - clip.start_time
        setClipPlaybackCurrentTime(offset)
        playerRef.current?.seek(offset)
      }
    },
    [sortedClips]
  )

  const handleSegmentClick = useCallback((segmentId: string) => {
    setSelectedClipId(segmentId)
    setClipPlaybackCurrentTime(0)
    playerRef.current?.seek(0)
  }, [])

  const handleSegmentChange = useCallback(
    (id: string, payload: { start_time: number; end_time: number }) => {
      updateClipMut.mutate({ id, payload })
    },
    [updateClipMut]
  )

  const transcription = transData?.transcription ?? null
  const segments = transcription?.segments ?? []
  const clipTranscriptSegments = selectedClip
    ? segmentsForClip(segments, selectedClip.start_time, selectedClip.end_time)
    : []

  const video = videoData?.video
  const durationSec = video?.duration_seconds ?? 0
  const timelineSegments: ClipTimelineSegment[] = sortedClips.map((c) => ({
    id: c.id,
    start_time: c.start_time,
    end_time: c.end_time,
    virality_score: c.virality_score ?? undefined,
    isSuggestion: false,
  }))
  const timelineCurrentTime =
    selectedClip != null
      ? selectedClip.start_time + clipPlaybackCurrentTime
      : 0

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
        <ErrorState
          message="Video not found."
          onRetry={() => queryClient.invalidateQueries({ queryKey: ['video', videoId] })}
        />
      </div>
    )
  }

  return (
    <div className="space-y-6">
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
            {video!.original_filename}
          </Link>
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

      <div className="flex items-center gap-2">
        <span className="text-sm font-medium text-[var(--app-fg-muted)]">View:</span>
        <div
          role="tablist"
          aria-label="Switch between Clips and Main video"
          className="inline-flex rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] p-0.5"
        >
          <button
            type="button"
            role="tab"
            aria-selected={viewMode === 'clips'}
            onClick={() => setViewMode('clips')}
            className={cn(
              'rounded-md px-4 py-2 text-sm font-medium transition-colors',
              viewMode === 'clips'
                ? 'bg-[var(--app-accent)] text-[#0a0a0b]'
                : 'text-[var(--app-fg-muted)] hover:text-[var(--app-fg)] hover:bg-[var(--app-bg-raised)]'
            )}
          >
            Clips
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={viewMode === 'main-video'}
            onClick={() => setViewMode('main-video')}
            className={cn(
              'rounded-md px-4 py-2 text-sm font-medium transition-colors',
              viewMode === 'main-video'
                ? 'bg-[var(--app-accent)] text-[#0a0a0b]'
                : 'text-[var(--app-fg-muted)] hover:text-[var(--app-fg)] hover:bg-[var(--app-bg-raised)]'
            )}
          >
            Main video
          </button>
        </div>
      </div>

      {viewMode === 'main-video' && (
        <MainVideoView videoId={videoId} video={video!} segments={segments} />
      )}

      {viewMode === 'clips' && (clipsLoading ? (
        <Skeleton className="h-64 w-full rounded-xl" />
      ) : clips.length === 0 ? (
        <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 text-center">
          <Scissors size={40} className="mx-auto text-[var(--app-fg-muted)]" />
          <p className="mt-2 font-medium text-[var(--app-fg)]">No clips yet</p>
          <p className="text-caption text-[var(--app-fg-muted)] mt-1">
            Use AI suggest clips above, or create a clip from the Clips page.
          </p>
          <div className="mt-4 flex justify-center gap-2">
            <Button variant="outline" size="sm" onClick={() => suggestMut.mutate()} disabled={suggestMut.isPending}>
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
        <div className="space-y-6">
          <section aria-label="Clips" className="w-full">
            <h2 className="text-sm font-medium text-[var(--app-fg-muted)] mb-3">Clips — click to select and play</h2>
            <div className="flex gap-4 overflow-x-auto pb-2 -mx-1 px-1 scroll-smooth snap-x snap-mandatory">
              {sortedClips.map((clip) => (
                <ClipStripCard
                  key={clip.id}
                  clip={clip}
                  isSelected={clip.id === selectedClipId}
                  onSelect={() => {
                    setSelectedClipId(clip.id)
                    setClipPlaybackCurrentTime(0)
                    playerRef.current?.seek(0)
                  }}
                  className="snap-start"
                />
              ))}
            </div>
          </section>

          <div className="grid grid-cols-1 lg:grid-cols-[1fr_minmax(320px,400px)] gap-6 items-start">
            <div className="space-y-4">
              <div className="rounded-xl overflow-hidden border border-[var(--app-border)] bg-[var(--app-bg)]">
              {playbackLoading && (
                <div className="aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
                  Loading…
                </div>
              )}
              {!playbackLoading && playbackUrl && (
                <VideoPlayer
                  key={selectedClipId ?? ''}
                  ref={playerRef}
                  src={playbackUrl}
                  initialTime={clipPlaybackCurrentTime}
                  onTimeUpdate={setClipPlaybackCurrentTime}
                />
              )}
              {!playbackLoading && selectedClip && !playbackUrl && (
                <div className="aspect-video flex flex-col items-center justify-center text-center p-6 text-[var(--app-fg-muted)]">
                  <p className="font-medium text-[var(--app-fg)]">Clip has no cut file yet</p>
                  <p className="text-sm mt-1">
                    This clip was created without a pre-cut video. Use the editor to render it, or run Auto cut on the
                    video to generate cut files for all clips.
                  </p>
                </div>
              )}
              {!playbackLoading && !selectedClip && (
                <div className="aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
                  Select a clip from the timeline
                </div>
              )}
            </div>

            {durationSec > 0 && (
              <ClipTimeline
                duration={durationSec}
                segments={timelineSegments}
                currentTime={timelineCurrentTime}
                onSeek={handleTimelineSeek}
                onSegmentChange={handleSegmentChange}
                onSegmentClick={handleSegmentClick}
              />
            )}
          </div>

          <div
            className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4 lg:sticky lg:top-6"
            id="transcript-panel"
          >
            <h2 className="font-semibold text-[var(--app-fg)] mb-3">
              {selectedClip ? `Transcript: ${selectedClip.name}` : 'Transcript'}
            </h2>
            {selectedClip && clipTranscriptSegments.length > 0 ? (
              <TranscriptViewer
                segments={clipTranscriptSegments}
                currentTime={clipPlaybackCurrentTime}
                onSeek={(t) => playerRef.current?.seek(t)}
                className="max-h-[50vh]"
              />
            ) : selectedClip ? (
              <p className="text-caption text-[var(--app-fg-muted)]">No transcript for this clip range.</p>
            ) : (
              <p className="text-caption text-[var(--app-fg-muted)]">Select a clip to see its transcript.</p>
            )}
          </div>
          </div>
        </div>
      ))}
      {viewMode === 'clips' && !clipsLoading && suggestions.length > 0 && (
        <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4">
          <div className="flex items-center justify-between mb-3">
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
          <ul className="flex flex-wrap gap-2">
            {suggestions.map((suggestion, index) => (
              <li
                key={`${suggestion.start_time}-${suggestion.end_time}-${index}`}
                className="flex items-center gap-2 rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-2"
              >
                <span className="text-caption text-[var(--app-fg-muted)]">
                  {formatTimeForLabel(suggestion.start_time)}–{formatTimeForLabel(suggestion.end_time)}
                </span>
                {suggestion.transcript && (
                  <span className="max-w-[200px] truncate text-sm text-[var(--app-fg-muted)]" title={suggestion.transcript}>
                    {suggestion.transcript}
                  </span>
                )}
                <Button
                  variant="outline"
                  size="sm"
                  className="shrink-0"
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
    </div>
  )
}
