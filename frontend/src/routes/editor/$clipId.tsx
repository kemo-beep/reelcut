import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useRef, useState, useEffect, useCallback } from 'react'
import { getClip, getClipStyle, updateClipStyle, renderClip, getClipDownloadUrl, updateClip } from '../../lib/api/clips'
import { getPlaybackUrl } from '../../lib/api/videos'
import { Button } from '../../components/ui/button'
import { Skeleton } from '../../components/ui/skeleton'
import { ErrorState } from '../../components/ui/error-state'
import { VideoPlayer, type VideoPlayerHandle } from '../../components/video/VideoPlayer'
import { TimelineEditor } from '../../components/editor/TimelineEditor'
import { StylePanel } from '../../components/editor/StylePanel'
import { ExportPanel } from '../../components/editor/ExportPanel'
import type { ClipStyle } from '../../types'

export const Route = createFileRoute('/editor/$clipId')({
  component: EditorPage,
})

function EditorPage() {
  const { clipId } = Route.useParams()
  const queryClient = useQueryClient()
  const videoRef = useRef<VideoPlayerHandle>(null)
  const [clipTime, setClipTime] = useState(0)

  const { data: clipData, isLoading: clipLoading, error: clipError } = useQuery({
    queryKey: ['clip', clipId],
    queryFn: () => getClip(clipId),
    refetchInterval: (query) => {
      const status = (query.state.data as { clip?: { status: string } })?.clip?.status
      return status === 'rendering' ? 2000 : false
    },
  })
  const clip = clipData?.clip
  const videoId = clip?.video_id

  const { data: playbackData } = useQuery({
    queryKey: ['video-playback', videoId],
    queryFn: () => getPlaybackUrl(videoId!),
    enabled: !!videoId,
  })
  const playbackUrl = playbackData?.url ?? null

  const { data: styleData } = useQuery({
    queryKey: ['clip-style', clipId],
    queryFn: () => getClipStyle(clipId),
    enabled: !!clipId,
  })
  const fetchedStyle = styleData?.style ?? null

  const [localStyle, setLocalStyle] = useState<Partial<ClipStyle> | null>(null)
  const style = localStyle != null ? { ...defaultStyle(), ...fetchedStyle, ...localStyle } : (fetchedStyle ?? defaultStyle())

  useEffect(() => {
    if (fetchedStyle) setLocalStyle(null)
  }, [fetchedStyle])

  const updateStyleMut = useMutation({
    mutationFn: (updates: Partial<ClipStyle>) => updateClipStyle(clipId, updates),
    onSuccess: (_, updates) => {
      queryClient.invalidateQueries({ queryKey: ['clip-style', clipId] })
      setLocalStyle((prev) => (prev ? { ...prev, ...updates } : updates))
    },
  })

  const renderMut = useMutation({
    mutationFn: () => renderClip(clipId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['clip', clipId] }),
  })

  const updateClipMut = useMutation({
    mutationFn: (body: Parameters<typeof updateClip>[1]) => updateClip(clipId, body),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['clip', clipId] }),
  })

  const handleSeek = useCallback(
    (time: number) => {
      if (!clip) return
      const t = Math.max(0, Math.min(clip.end_time - clip.start_time, time))
      setClipTime(t)
      videoRef.current?.seek(clip.start_time + t)
    },
    [clip]
  )

  const handleTimeUpdate = useCallback(
    (videoCurrentTime: number) => {
      if (!clip) return
      const start = clip.start_time
      const end = clip.end_time
      const duration = end - start
      if (videoCurrentTime < start || videoCurrentTime > end) return
      setClipTime(videoCurrentTime - start)
    },
    [clip]
  )

  useEffect(() => {
    if (!clip || !playbackUrl) return
    videoRef.current?.seek(clip.start_time)
    setClipTime(0)
  }, [clip?.id, playbackUrl])

  if (clipLoading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen p-8">
        <Skeleton className="h-64 w-full max-w-4xl rounded-xl" />
        <Skeleton className="h-12 w-48 mt-6" />
      </div>
    )
  }

  if (clipError || !clip) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen p-8">
        <ErrorState message="Clip not found." onRetry={() => queryClient.invalidateQueries({ queryKey: ['clip', clipId] })} />
        <Link to="/dashboard/clips" className="mt-4">
          <Button variant="outline">Back to Clips</Button>
        </Link>
      </div>
    )
  }

  const duration = clip.end_time - clip.start_time
  const renderStatus =
    clip.status === 'rendering'
      ? 'rendering'
      : clip.status === 'ready'
        ? 'ready'
        : clip.status === 'failed'
          ? 'failed'
          : 'idle'
  const downloadUrl = clip.status === 'ready' && clip.storage_path ? getClipDownloadUrl(clipId) : null

  const handleRender = (aspectRatio: string) => {
    if (aspectRatio !== clip.aspect_ratio) {
      updateClipMut.mutate({ aspect_ratio: aspectRatio as '9:16' | '1:1' | '16:9' }, {
        onSuccess: () => renderMut.mutate(),
      })
    } else {
      renderMut.mutate()
    }
  }

  return (
    <div className="flex flex-col min-h-screen bg-[var(--app-bg)]">
      <header className="flex items-center justify-between border-b border-[var(--app-border)] bg-[var(--app-bg-raised)] px-4 py-3 shrink-0">
        <Link
          to="/dashboard/clips/$clipId"
          params={{ clipId }}
          className="text-sm text-[var(--app-fg-muted)] hover:text-[var(--app-accent)]"
        >
          ← Back to clip
        </Link>
        <h1 className="font-semibold text-[var(--app-fg)] truncate max-w-md">{clip.name}</h1>
        <div className="w-24" />
      </header>

      <main className="flex-1 flex flex-col lg:flex-row gap-6 p-4 md:p-6 min-h-0">
        <div className="flex-1 flex flex-col min-w-0">
          <div className="rounded-xl overflow-hidden border border-[var(--app-border)] bg-[var(--app-bg-raised)]">
            {playbackUrl ? (
              <VideoPlayer
                ref={videoRef}
                src={playbackUrl}
                onTimeUpdate={handleTimeUpdate}
              />
            ) : (
              <div className="aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
                Loading video…
              </div>
            )}
          </div>
          <div className="mt-4">
            <TimelineEditor
              duration={duration}
              currentTime={clipTime}
              onSeek={handleSeek}
            />
          </div>
        </div>

        <aside className="w-full lg:w-80 shrink-0 space-y-6 overflow-y-auto">
          <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4">
            <StylePanel
              style={style}
              onChange={(updates) => {
                setLocalStyle((prev) => ({ ...prev, ...updates }))
                updateStyleMut.mutate(updates)
              }}
            />
          </div>
          <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4">
            <ExportPanel
              aspectRatios={['9:16', '1:1', '16:9']}
              onRender={handleRender}
              renderStatus={renderStatus}
              renderProgress={0}
              downloadUrl={downloadUrl}
            />
          </div>
        </aside>
      </main>
    </div>
  )
}

function defaultStyle(): Partial<ClipStyle> {
  return {
    caption_enabled: true,
    caption_font: 'Inter',
    caption_size: 48,
    caption_color: '#FFFFFF',
    caption_position: 'bottom',
    background_music_volume: 0.3,
  }
}
