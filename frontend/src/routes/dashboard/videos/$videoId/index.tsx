import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useEffect, useRef, useState } from 'react'
import { getVideo, getPlaybackUrl, deleteVideo, triggerAutoCut } from '../../../../lib/api/videos'
import {
  createTranscription,
  getTranscriptionByVideoId,
  listTranscriptionsByVideo,
  translateTranscription,
} from '../../../../lib/api/transcriptions'
import { analyzeVideo, suggestClips } from '../../../../lib/api/analysis'
import { listClips, updateClip } from '../../../../lib/api/clips'
import { Button } from '../../../../components/ui/button'
import { VideoPlayer } from '../../../../components/video/VideoPlayer'
import { TranscriptViewer } from '../../../../components/transcription/TranscriptViewer'
import { ClipTimeline, type ClipTimelineSegment } from '../../../../components/clip/ClipTimeline'
import { Trash2, Sparkles, Loader2, FileText } from 'lucide-react'
import { toast } from 'sonner'
import { ApiError } from '../../../../types'

export const Route = createFileRoute('/dashboard/videos/$videoId/')({
  component: VideoDetailPage,
})

function VideoDetailPage() {
  const { videoId } = Route.useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const playerRef = useRef<{ seek: (time: number) => void; getCurrentTime: () => number }>(null)
  const [currentTime, setCurrentTime] = useState(0)
  const [createLanguage, setCreateLanguage] = useState('en')
  const [viewLanguage, setViewLanguage] = useState<string | null>(null)
  const [translateTarget, setTranslateTarget] = useState('es')

  const { data, isLoading, error } = useQuery({
    queryKey: ['video', videoId],
    queryFn: () => getVideo(videoId),
    refetchInterval: (query) =>
      query.state.data?.video?.status === 'processing' ? 3000 : false,
  })
  const { data: playbackData, error: playbackError, isFetching: playbackLoading } = useQuery({
    queryKey: ['video-playback', videoId],
    queryFn: () => getPlaybackUrl(videoId),
    enabled: !!data?.video,
    retry: false,
  })
  const { data: transData, isLoading: transLoading } = useQuery({
    queryKey: ['transcription', videoId, viewLanguage ?? 'default'],
    queryFn: () =>
      getTranscriptionByVideoId(videoId, viewLanguage ? { language: viewLanguage } : undefined),
    enabled: !!data?.video,
    refetchOnMount: 'always',
    refetchOnWindowFocus: true,
    refetchInterval: (query) =>
      query.state.data?.transcription?.status === 'processing' ? 2500 : false,
  })
  const { data: listData } = useQuery({
    queryKey: ['transcriptions-list', videoId],
    queryFn: () => listTranscriptionsByVideo(videoId),
    enabled: !!data?.video && transData?.transcription?.status === 'completed',
    staleTime: 30 * 1000,
  })
  const transcriptionsList = listData?.transcriptions ?? []
  const { data: clipsData } = useQuery({
    queryKey: ['clips', { video_id: videoId }],
    queryFn: () => listClips({ video_id: videoId, per_page: 100 }),
    enabled: !!data?.video,
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
  const autoCutMut = useMutation({
    mutationFn: () => triggerAutoCut(videoId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clips', { video_id: videoId }] })
      toast.success('Auto-cut started; clips will appear shortly.')
    },
    onError: () => toast.error('Failed to start auto-cut'),
  })
  const autoStartedRef = useRef(false)
  const transcriptionLoaded = transData !== undefined
  useEffect(() => {
    const video = data?.video
    const transcription = transData?.transcription
    if (!video || video.status !== 'ready' || autoStartedRef.current) return
    if (!transcriptionLoaded) return
    if (transcription != null) return
    autoStartedRef.current = true
    createTranscription(videoId, { language: createLanguage })
      .then(() => {
        queryClient.invalidateQueries({ queryKey: ['transcription', videoId] })
        queryClient.invalidateQueries({ queryKey: ['transcriptions-list', videoId] })
        toast.info('Transcription started automatically')
      })
      .catch(() => { autoStartedRef.current = false })
  }, [data?.video, transData?.transcription, transcriptionLoaded, videoId, queryClient, createLanguage])

  const deleteMutation = useMutation({
    mutationFn: () => deleteVideo(videoId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['videos'] })
      toast.success('Video deleted')
      navigate({ to: '/dashboard/videos' })
    },
    onError: () => toast.error('Failed to delete video'),
  })

  const createTranscriptionMut = useMutation({
    mutationFn: () => createTranscription(videoId, { language: createLanguage }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transcription', videoId] })
      queryClient.invalidateQueries({ queryKey: ['transcriptions-list', videoId] })
      toast.success('Transcription started')
    },
    onError: (err) => {
      const message = err instanceof ApiError ? err.message : 'Failed to start transcription'
      toast.error(message)
    },
  })

  const translateMut = useMutation({
    mutationFn: ({ transcriptionId, targetLanguage }: { transcriptionId: string; targetLanguage: string }) =>
      translateTranscription(transcriptionId, { target_language: targetLanguage }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transcription', videoId] })
      queryClient.invalidateQueries({ queryKey: ['transcriptions-list', videoId] })
      toast.success('Translation started; new language will appear shortly.')
    },
    onError: () => toast.error('Translation failed'),
  })

  const transcribeAndSuggestMut = useMutation({
    mutationFn: async () => {
      const { transcription } = await createTranscription(videoId, { language: createLanguage })
      const pollMs = 3000
      const maxWait = 10 * 60 * 1000
      const start = Date.now()
      while (Date.now() - start < maxWait) {
        await new Promise((r) => setTimeout(r, pollMs))
        const { transcription: t } = await getTranscriptionByVideoId(videoId)
        if (!t) continue
        if (t.status === 'failed') throw new Error('Transcription failed')
        if (t.status === 'completed') break
      }
      const { transcription: t } = await getTranscriptionByVideoId(videoId)
      if (!t || t.status !== 'completed') throw new Error('Transcription did not complete in time')
      await analyzeVideo(videoId).catch(() => {})
      const { suggestions } = await suggestClips(videoId)
      return { suggestions }
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['transcription', videoId] })
      queryClient.invalidateQueries({ queryKey: ['clips', { video_id: videoId }] })
      toast.success(data.suggestions?.length ? `${data.suggestions.length} clip suggestions ready` : 'Transcription complete')
      navigate({
        to: '/dashboard/videos/$videoId/clips',
        params: { videoId },
        state: { suggestions: data.suggestions ?? [] },
      })
    },
    onError: (e) => toast.error(e instanceof Error ? e.message : 'Failed to transcribe and suggest clips'),
  })

  if (isLoading) return <p className="text-slate-400">Loading...</p>
  if (error || !data?.video) return <p className="text-red-400">Video not found.</p>

  const video = data.video
  const playbackUrl = playbackData?.url
  const transcription = transData?.transcription ?? null
  const transLoadingState = transLoading && transData === undefined
  const noTranscriptionYet = !transLoadingState && transData && transcription === null
  const videoReady = video.status === 'ready'
  const clips = clipsData?.data?.clips ?? []
  const durationSec = video.duration_seconds ?? 0
  const timelineSegments: ClipTimelineSegment[] = clips.map((c) => ({
    id: c.id,
    start_time: c.start_time,
    end_time: c.end_time,
    virality_score: c.virality_score ?? undefined,
    isSuggestion: false,
  }))
  const handleSegmentChange = (id: string, payload: { start_time: number; end_time: number }) => {
    if (id.startsWith('suggestion-')) return
    updateClipMut.mutate({ id, payload })
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <Link
          to="/dashboard/videos"
          className="text-cyan-400 hover:underline text-sm"
        >
          Back to videos
        </Link>
        <Button
          variant="outline"
          size="sm"
          className="text-red-400 border-red-400/50 hover:bg-red-500/10"
          onClick={() => deleteMutation.mutate()}
          disabled={deleteMutation.isPending}
        >
          <Trash2 size={16} className="mr-1" />
          Delete
        </Button>
      </div>
      <h1 className="text-2xl font-bold text-[var(--app-fg)]">{video.original_filename}</h1>
      <p className="text-caption text-[var(--app-fg-muted)]">
        Status: {video.status} · Duration: {video.duration_seconds != null ? `${Math.round(video.duration_seconds)}s` : '—'}
      </p>

      <div className="grid grid-cols-1 lg:grid-cols-[1fr_minmax(320px,400px)] gap-6 items-start">
        <div className="space-y-4">
          {playbackLoading && (
            <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg)] aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
              Loading playback…
            </div>
          )}
          {!playbackLoading && playbackUrl && (
            <div className="rounded-xl overflow-hidden border border-[var(--app-border)] bg-[var(--app-bg)]">
              <VideoPlayer
                ref={playerRef}
                src={playbackUrl}
                onTimeUpdate={setCurrentTime}
              />
            </div>
          )}
          {!playbackLoading && !playbackUrl && !playbackError && (
            <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg)] aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
              Video is still processing. Check back in a moment.
            </div>
          )}
          {!playbackLoading && playbackError && (
            <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg)] aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
              Unable to load playback.
            </div>
          )}

          {durationSec > 0 && (
            <ClipTimeline
              duration={durationSec}
              segments={timelineSegments}
              currentTime={currentTime}
              onSeek={(t) => playerRef.current?.seek(t)}
              onSegmentChange={handleSegmentChange}
            />
          )}

          <div className="flex flex-wrap gap-3">
            <Button
              size="sm"
              className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]"
              onClick={() => transcribeAndSuggestMut.mutate()}
              disabled={transcribeAndSuggestMut.isPending || video.status !== 'ready'}
            >
              {transcribeAndSuggestMut.isPending ? (
                <>
                  <Loader2 size={16} className="mr-1 animate-spin" />
                  Transcribing…
                </>
              ) : (
                <>
                  <Sparkles size={16} className="mr-1" />
                  Transcribe & suggest clips
                </>
              )}
            </Button>
            {transcription?.status === 'completed' && (
              <Button
                size="sm"
                variant="outline"
                className="border-[var(--app-border)]"
                onClick={() => autoCutMut.mutate()}
                disabled={autoCutMut.isPending || video.status !== 'ready'}
              >
                {autoCutMut.isPending ? (
                  <>
                    <Loader2 size={16} className="mr-1 animate-spin" />
                    Starting…
                  </>
                ) : (
                  'Auto cut'
                )}
              </Button>
            )}
            <Link
              to="/dashboard/videos/$videoId/clips"
              params={{ videoId: video.id }}
              className="inline-flex items-center justify-center rounded-lg border border-[var(--app-border)] bg-[var(--app-bg-raised)] px-4 py-2 text-sm font-medium text-[var(--app-fg)] hover:bg-[var(--app-bg-overlay)]"
            >
              Clips from this video
            </Link>
            <Link
              to="/dashboard/clips"
              search={{ videoId: video.id }}
              className="inline-flex items-center justify-center rounded-lg border border-[var(--app-border)] bg-[var(--app-bg-raised)] px-4 py-2 text-sm font-medium text-[var(--app-fg)] hover:bg-[var(--app-bg-overlay)]"
            >
              View all clips
            </Link>
          </div>
        </div>

        <div
          className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4 lg:sticky lg:top-6"
          id="transcript-panel"
        >
          <h2 className="font-semibold text-[var(--app-fg)] mb-3">Transcript</h2>

          {transLoadingState && (
            <div className="flex flex-col items-center justify-center py-8 text-center">
              <Loader2 size={24} className="animate-spin text-[var(--app-fg-muted)]" />
              <p className="mt-2 text-caption text-[var(--app-fg-muted)]">Loading transcription…</p>
            </div>
          )}

          {noTranscriptionYet && (
            <div className="space-y-3">
              <div className="flex items-center gap-2 text-[var(--app-fg-muted)]">
                <FileText size={20} />
                <span className="text-sm">No transcription yet</span>
              </div>
              <p className="text-caption text-[var(--app-fg-muted)]">
                {videoReady
                  ? 'Start transcription to generate captions and enable AI clip suggestions.'
                  : 'Video is still processing. Transcription can be started when the video is ready.'}
              </p>
              <div className="flex flex-wrap items-center gap-2">
                <label className="text-sm text-[var(--app-fg-muted)]">Source language:</label>
                <select
                  className="rounded border border-[var(--app-border)] bg-[var(--app-bg)] px-2 py-1 text-sm text-[var(--app-fg)]"
                  value={createLanguage}
                  onChange={(e) => setCreateLanguage(e.target.value)}
                >
                  <option value="en">English</option>
                  <option value="es">Spanish</option>
                  <option value="fr">French</option>
                  <option value="de">German</option>
                  <option value="it">Italian</option>
                  <option value="pt">Portuguese</option>
                  <option value="ja">Japanese</option>
                  <option value="zh">Chinese</option>
                </select>
                <Button
                  onClick={() => createTranscriptionMut.mutate()}
                  disabled={createTranscriptionMut.isPending || !videoReady}
                  size="sm"
                  className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)] disabled:opacity-50"
                >
                  {createTranscriptionMut.isPending ? 'Starting…' : 'Start transcription'}
                </Button>
              </div>
            </div>
          )}

          {transcription && transcription.status === 'processing' && (
            <div className="space-y-3">
              <div className="flex items-center gap-2 text-caption text-[var(--app-fg-muted)]">
                <Loader2 size={16} className="animate-spin" />
                Transcription in progress… Segments will appear as they’re ready.
              </div>
              {transcription.segments && transcription.segments.length > 0 && (
                <TranscriptViewer
                  segments={transcription.segments}
                  currentTime={currentTime}
                  onSeek={(t) => playerRef.current?.seek(t)}
                  className="max-h-[50vh]"
                />
              )}
            </div>
          )}

          {transcription &&
            transcription.status === 'completed' &&
            transcription.segments &&
            transcription.segments.length > 0 && (
              <>
                {transcriptionsList.length > 0 && (
                  <div className="flex flex-wrap items-center gap-2 mb-3">
                    <span className="text-sm text-[var(--app-fg-muted)]">View:</span>
                    <button
                      type="button"
                      className={viewLanguage === null ? 'font-medium text-[var(--app-accent)]' : 'text-sm text-[var(--app-fg-muted)] hover:text-[var(--app-fg)]'}
                      onClick={() => setViewLanguage(null)}
                    >
                      Default
                    </button>
                    {transcriptionsList.map((t) => (
                      <button
                        key={t.id}
                        type="button"
                        className={viewLanguage === t.language ? 'font-medium text-[var(--app-accent)]' : 'text-sm text-[var(--app-fg-muted)] hover:text-[var(--app-fg)]'}
                        onClick={() => setViewLanguage(t.language)}
                      >
                        {t.language === 'en' ? 'English' : t.language === 'es' ? 'Spanish' : t.language === 'fr' ? 'French' : t.language === 'de' ? 'German' : t.language}
                      </button>
                    ))}
                    <span className="text-sm text-[var(--app-fg-muted)] ml-2">|</span>
                    <span className="text-sm text-[var(--app-fg-muted)]">Translate to:</span>
                    <select
                      className="rounded border border-[var(--app-border)] bg-[var(--app-bg)] px-2 py-1 text-sm text-[var(--app-fg)]"
                      value={translateTarget}
                      onChange={(e) => setTranslateTarget(e.target.value)}
                    >
                      <option value="es">Spanish</option>
                      <option value="fr">French</option>
                      <option value="de">German</option>
                      <option value="it">Italian</option>
                      <option value="pt">Portuguese</option>
                      <option value="ja">Japanese</option>
                      <option value="zh">Chinese</option>
                    </select>
                    <Button
                      size="sm"
                      variant="outline"
                      className="border-[var(--app-border)]"
                      disabled={translateMut.isPending || transcriptionsList.some((t) => t.language === translateTarget)}
                      onClick={() => transcription && translateMut.mutate({ transcriptionId: transcription.id, targetLanguage: translateTarget })}
                    >
                      {translateMut.isPending ? 'Translating…' : 'Translate'}
                    </Button>
                  </div>
                )}
                <TranscriptViewer
                  segments={transcription.segments}
                  currentTime={currentTime}
                  onSeek={(t) => playerRef.current?.seek(t)}
                  className="max-h-[50vh]"
                />
              </>
            )}

          {transcription &&
            transcription.status === 'completed' &&
            (!transcription.segments || transcription.segments.length === 0) && (
              <div className="space-y-3 py-4">
                <p className="text-caption text-[var(--app-fg-muted)]">
                  No segments were saved. This can happen if the transcription job failed to persist results.
                </p>
                <Button
                  onClick={() => createTranscriptionMut.mutate()}
                  disabled={createTranscriptionMut.isPending || !videoReady}
                  size="sm"
                  className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)] disabled:opacity-50"
                >
                  {createTranscriptionMut.isPending ? 'Starting…' : 'Transcribe again'}
                </Button>
              </div>
            )}

          {transcription && transcription.status === 'failed' && (
            <>
              <div className="text-[var(--app-destructive)] text-sm py-4 space-y-1">
                <p>Transcription failed. Try starting again from the actions above.</p>
                {transcription.error_message && (
                  <p className="text-[var(--app-muted-foreground)] font-mono text-xs mt-2 break-all">
                    {transcription.error_message}
                  </p>
                )}
              </div>
              {transcription.segments && transcription.segments.length > 0 && (
                <TranscriptViewer
                  segments={transcription.segments}
                  currentTime={currentTime}
                  onSeek={(t) => playerRef.current?.seek(t)}
                  className="max-h-[50vh]"
                />
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
