import type { TranscriptSegment } from '../../types'
import { cn } from '../../lib/utils'

function formatTime(sec: number): string {
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

export interface TranscriptViewerProps {
  segments: TranscriptSegment[]
  currentTime?: number
  onSeek?: (time: number) => void
  className?: string
}

export function TranscriptViewer({
  segments,
  currentTime = 0,
  onSeek,
  className,
}: TranscriptViewerProps) {
  return (
    <div
      className={cn('space-y-2 max-h-[60vh] overflow-y-auto', className)}
      role="region"
      aria-label="Transcript"
    >
      {segments.map((seg) => {
        const isActive =
          currentTime >= seg.start_time && currentTime <= seg.end_time
        return (
          <button
            key={seg.id}
            type="button"
            onClick={() => onSeek?.(seg.start_time)}
            className={cn(
              'w-full rounded-lg border px-4 py-2 text-left text-sm transition-colors',
              isActive
                ? 'border-[var(--app-accent)] bg-[var(--app-accent-muted)]'
                : 'border-[var(--app-border)] bg-[var(--app-bg)] hover:border-[var(--app-border-strong)]'
            )}
          >
            <span className="text-caption text-[var(--app-fg-muted)]">
              {formatTime(seg.start_time)} – {formatTime(seg.end_time)}
            </span>
            <p className="text-[var(--app-fg)] mt-0.5">{seg.text}</p>
          </button>
        )
      })}
    </div>
  )
}
