import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useRef, useState, useEffect, useCallback } from 'react'
import { getClip, getClipStyle, updateClipStyle, renderClip, getClipDownloadUrl, updateClip, listBrollSegments, addBrollSegment, deleteBrollSegment, type ClipBrollSegment } from '../../lib/api/clips'
import { listBrollAssets, uploadBrollAsset } from '../../lib/api/broll'
import { getPlaybackUrl } from '../../lib/api/videos'
import { Button } from '../../components/ui/button'
import { Skeleton } from '../../components/ui/skeleton'
import { Label } from '../../components/ui/label'
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
  const { data: brollData } = useQuery({
    queryKey: ['clip-broll', clipId],
    queryFn: () => listBrollSegments(clipId),
    enabled: !!clipId,
  })
  const { data: brollAssetsData } = useQuery({
    queryKey: ['broll-assets'],
    queryFn: () => listBrollAssets({ limit: 100 }),
  })
  const fetchedStyle = styleData?.style ?? null
  const brollSegments = brollData?.segments ?? []
  const brollAssets = brollAssetsData?.assets ?? []

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
    mutationFn: (opts?: { preset?: string }) => renderClip(clipId, opts),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['clip', clipId] }),
  })

  const updateClipMut = useMutation({
    mutationFn: (body: Parameters<typeof updateClip>[1]) => updateClip(clipId, body),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['clip', clipId] }),
  })
  const addBrollMut = useMutation({
    mutationFn: (body: { broll_asset_id: string; start_time: number; end_time: number; position?: string; scale?: number; opacity?: number }) =>
      addBrollSegment(clipId, body),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['clip-broll', clipId] }),
  })
  const deleteBrollMut = useMutation({
    mutationFn: (segmentId: string) => deleteBrollSegment(clipId, segmentId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['clip-broll', clipId] }),
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
        <Link to="/dashboard/videos" className="mt-4">
          <Button variant="outline">Back to Videos</Button>
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

  const handleRender = (options: { preset?: string; aspectRatio?: string }) => {
    const targetAspect = options.aspectRatio ?? (options.preset ? undefined : clip.aspect_ratio)
    if (targetAspect && targetAspect !== clip.aspect_ratio) {
      updateClipMut.mutate(
        { aspect_ratio: targetAspect as '9:16' | '1:1' | '16:9' },
        { onSuccess: () => renderMut.mutate({ preset: options.preset }) }
      )
    } else {
      renderMut.mutate({ preset: options.preset })
    }
  }

  return (
    <div className="flex flex-col min-h-screen bg-[var(--app-bg)]">
      <header className="flex items-center justify-between border-b border-[var(--app-border)] bg-[var(--app-bg-raised)] px-4 py-3 shrink-0">
        <Link
          to="/dashboard/videos/$videoId"
          params={{ videoId: videoId ?? '' }}
          className="text-sm text-[var(--app-fg-muted)] hover:text-[var(--app-accent)]"
        >
          ← Back to video
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
              videoId={videoId}
              onChange={(updates) => {
                setLocalStyle((prev) => ({ ...prev, ...updates }))
                updateStyleMut.mutate(updates)
              }}
            />
          </div>
          <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4">
            <BrollPanel
              clipDuration={duration}
              segments={brollSegments}
              assets={brollAssets}
              onAdd={(body) => addBrollMut.mutate(body)}
              onDelete={(segmentId) => deleteBrollMut.mutate(segmentId)}
              addPending={addBrollMut.isPending}
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

function BrollPanel({
  clipDuration,
  segments,
  assets,
  onAdd,
  onDelete,
  addPending,
}: {
  clipDuration: number
  segments: ClipBrollSegment[]
  assets: { id: string; original_filename: string }[]
  onAdd: (body: { broll_asset_id: string; start_time: number; end_time: number; position?: string; scale?: number; opacity?: number }) => void
  onDelete: (segmentId: string) => void
  addPending: boolean
}) {
  const [assetId, setAssetId] = useState('')
  const [start, setStart] = useState(0)
  const [end, setEnd] = useState(Math.min(5, clipDuration))
  const [position, setPosition] = useState<'overlay' | 'cut_in'>('overlay')
  const handleAdd = () => {
    if (!assetId || end <= start) return
    onAdd({ broll_asset_id: assetId, start_time: start, end_time: end, position, scale: 0.5, opacity: 1 })
    setStart(end)
    setEnd(Math.min(end + 5, clipDuration))
  }
  return (
    <div className="space-y-4">
      <h3 className="font-semibold text-[var(--app-fg)]">B-roll</h3>
      <ul className="space-y-2">
        {segments.map((seg) => (
          <li key={seg.id} className="flex items-center justify-between rounded border border-[var(--app-border)] bg-[var(--app-bg)] px-2 py-1.5 text-sm">
            <span className="truncate">
              {seg.asset?.original_filename ?? seg.broll_asset_id} · {seg.start_time.toFixed(1)}–{seg.end_time.toFixed(1)}s · {seg.position}
            </span>
            <button
              type="button"
              className="text-[var(--app-fg-muted)] hover:text-[var(--app-destructive)]"
              onClick={() => onDelete(seg.id)}
              aria-label="Remove B-roll segment"
            >
              Remove
            </button>
          </li>
        ))}
      </ul>
      <div className="space-y-2">
        <Label>Add B-roll</Label>
        <select
          className="w-full rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-2 text-sm text-[var(--app-fg)]"
          value={assetId}
          onChange={(e) => setAssetId(e.target.value)}
        >
          <option value="">Select asset</option>
          {assets.map((a) => (
            <option key={a.id} value={a.id}>{a.original_filename}</option>
          ))}
        </select>
        <div className="grid grid-cols-2 gap-2">
          <div>
            <Label className="text-xs">Start (s)</Label>
            <input
              type="number"
              min={0}
              max={clipDuration}
              step={0.5}
              className="w-full rounded border border-[var(--app-border)] bg-[var(--app-bg)] px-2 py-1 text-sm"
              value={start}
              onChange={(e) => setStart(Number(e.target.value))}
            />
          </div>
          <div>
            <Label className="text-xs">End (s)</Label>
            <input
              type="number"
              min={0}
              max={clipDuration}
              step={0.5}
              className="w-full rounded border border-[var(--app-border)] bg-[var(--app-bg)] px-2 py-1 text-sm"
              value={end}
              onChange={(e) => setEnd(Number(e.target.value))}
            />
          </div>
        </div>
        <select
          className="w-full rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-2 text-sm text-[var(--app-fg)]"
          value={position}
          onChange={(e) => setPosition(e.target.value as 'overlay' | 'cut_in')}
        >
          <option value="overlay">Overlay (PIP)</option>
          <option value="cut_in">Cut in</option>
        </select>
        <Button size="sm" className="w-full" onClick={handleAdd} disabled={!assetId || end <= start || addPending}>
          {addPending ? 'Adding…' : 'Add B-roll'}
        </Button>
      </div>
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
