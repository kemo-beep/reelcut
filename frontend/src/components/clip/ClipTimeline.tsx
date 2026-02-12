import { useRef, useState, useCallback, useEffect } from 'react'
import { cn } from '../../lib/utils'

export interface ClipTimelineSegment {
  id: string
  start_time: number
  end_time: number
  virality_score?: number | null
  reason?: string | null
  /** If true, segment is a suggestion (not yet saved); style can differ */
  isSuggestion?: boolean
}

export interface ClipTimelineProps {
  duration: number
  segments: ClipTimelineSegment[]
  currentTime: number
  onSeek: (time: number) => void
  onSegmentChange?: (id: string, payload: { start_time: number; end_time: number }) => void
  /** When user clicks a segment block (e.g. to select that clip for the viewer) */
  onSegmentClick?: (segmentId: string) => void
  className?: string
}

const MIN_SEGMENT_SEC = 5

export function ClipTimeline({
  duration,
  segments,
  currentTime,
  onSeek,
  onSegmentChange,
  onSegmentClick,
  className,
}: ClipTimelineProps) {
  const trackRef = useRef<HTMLDivElement>(null)
  const [dragState, setDragState] = useState<{
    segmentId: string
    edge: 'left' | 'right'
    startX: number
    startTime: number
    endTime: number
  } | null>(null)

  const pxToTime = useCallback(
    (clientX: number): number => {
      const el = trackRef.current
      if (!el || duration <= 0) return 0
      const rect = el.getBoundingClientRect()
      const x = clientX - rect.left
      const fraction = Math.max(0, Math.min(1, x / rect.width))
      return fraction * duration
    },
    [duration]
  )

  const handleTrackClick = useCallback(
    (e: React.MouseEvent<HTMLDivElement>) => {
      if (trackRef.current && duration > 0) {
        const time = pxToTime(e.clientX)
        onSeek(time)
      }
    },
    [duration, onSeek, pxToTime]
  )

  const handleResizeStart = useCallback(
    (e: React.MouseEvent, segmentId: string, edge: 'left' | 'right', startTime: number, endTime: number) => {
      e.stopPropagation()
      if (!onSegmentChange) return
      setDragState({ segmentId, edge, startX: e.clientX, startTime, endTime })
    },
    [onSegmentChange]
  )

  const handleMouseMove = useCallback(
    (e: MouseEvent) => {
      if (!dragState || !trackRef.current) return
      const el = trackRef.current
      const rect = el.getBoundingClientRect()
      const deltaX = e.clientX - dragState.startX
      const deltaSec = (deltaX / rect.width) * duration

      let newStart = dragState.startTime
      let newEnd = dragState.endTime
      if (dragState.edge === 'left') {
        newStart = Math.max(0, Math.min(dragState.startTime + deltaSec, dragState.endTime - MIN_SEGMENT_SEC))
      } else {
        newEnd = Math.min(duration, Math.max(dragState.endTime + deltaSec, dragState.startTime + MIN_SEGMENT_SEC))
      }
      setDragState((prev) => prev && { ...prev, startTime: newStart, endTime: newEnd, startX: e.clientX })
    },
    [dragState, duration]
  )

  const handleMouseUp = useCallback(() => {
    setDragState((prev) => {
      if (prev && onSegmentChange) {
        onSegmentChange(prev.segmentId, { start_time: prev.startTime, end_time: prev.endTime })
      }
      return null
    })
  }, [onSegmentChange])

  useEffect(() => {
    if (!dragState) return
    window.addEventListener('mousemove', handleMouseMove)
    window.addEventListener('mouseup', handleMouseUp)
    return () => {
      window.removeEventListener('mousemove', handleMouseMove)
      window.removeEventListener('mouseup', handleMouseUp)
    }
  }, [dragState, handleMouseMove, handleMouseUp])

  if (duration <= 0) {
    return (
      <div className={cn('rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] p-3', className)}>
        <p className="text-caption text-[var(--app-fg-muted)]">Timeline (no duration)</p>
      </div>
    )
  }

  const playheadPct = (currentTime / duration) * 100

  return (
    <div className={cn('space-y-2', className)}>
      <div
        ref={trackRef}
        role="slider"
        aria-valuenow={currentTime}
        aria-valuemin={0}
        aria-valuemax={duration}
        tabIndex={0}
        onClick={handleTrackClick}
        className="relative h-12 w-full cursor-pointer overflow-hidden rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)]"
      >
        {/* Segment blocks */}
        {segments.map((seg) => {
          const isDragging = dragState?.segmentId === seg.id
          const start = isDragging ? dragState.startTime : seg.start_time
          const end = isDragging ? dragState.endTime : seg.end_time
          const left = (start / duration) * 100
          const width = ((end - start) / duration) * 100
          return (
            <div
              key={seg.id}
              role="button"
              tabIndex={0}
              className={cn(
                'absolute top-1 bottom-1 rounded border border-[var(--app-border)] cursor-pointer',
                seg.isSuggestion
                  ? 'bg-[var(--app-accent)]/40 border-[var(--app-accent)]'
                  : 'bg-[var(--app-accent)] border-[var(--app-accent)]'
              )}
              style={{ left: `${left}%`, width: `${width}%` }}
              onClick={(e) => {
                e.stopPropagation()
                onSegmentClick?.(seg.id)
              }}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault()
                  e.stopPropagation()
                  onSegmentClick?.(seg.id)
                }
              }}
            >
              {/* Left resize handle */}
              {onSegmentChange && (
                <>
                  <div
                    className="absolute left-0 top-0 bottom-0 w-2 cursor-ew-resize shrink-0"
                    onMouseDown={(e) => handleResizeStart(e, seg.id, 'left', seg.start_time, seg.end_time)}
                    aria-label="Resize start"
                  />
                  <div
                    className="absolute right-0 top-0 bottom-0 w-2 cursor-ew-resize shrink-0"
                    onMouseDown={(e) => handleResizeStart(e, seg.id, 'right', seg.start_time, seg.end_time)}
                    aria-label="Resize end"
                  />
                </>
              )}
            </div>
          )
        })}
        {/* Playhead */}
        <div
          className="pointer-events-none absolute top-0 bottom-0 w-0.5 bg-red-500 z-10"
          style={{ left: `${playheadPct}%`, transform: 'translateX(-50%)' }}
        />
      </div>
      <p className="text-caption text-[var(--app-fg-muted)]">
        {formatTime(currentTime)} / {formatTime(duration)}
      </p>
    </div>
  )
}

function formatTime(sec: number): string {
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}
