import { useRef } from 'react'
import { cn } from '../../lib/utils'

export interface TimelineEditorProps {
  duration: number
  currentTime: number
  onSeek: (time: number) => void
  className?: string
}

export function TimelineEditor({
  duration,
  currentTime,
  onSeek,
  className,
}: TimelineEditorProps) {
  const barRef = useRef<HTMLDivElement>(null)
  const pct = duration > 0 ? (currentTime / duration) * 100 : 0

  const handleClick = (e: React.MouseEvent<HTMLDivElement>) => {
    const el = barRef.current
    if (!el) return
    const rect = el.getBoundingClientRect()
    const x = e.clientX - rect.left
    const fraction = x / rect.width
    const time = Math.max(0, Math.min(duration, fraction * duration))
    onSeek(time)
  }

  return (
    <div className={cn('space-y-1', className)}>
      <div
        ref={barRef}
        role="slider"
        aria-valuenow={currentTime}
        aria-valuemin={0}
        aria-valuemax={duration}
        tabIndex={0}
        onClick={handleClick}
        className="h-3 w-full cursor-pointer overflow-hidden rounded-full bg-[var(--app-bg)] border border-[var(--app-border)]"
      >
        <div
          className="h-full rounded-full bg-[var(--app-accent)] transition-[width] duration-100"
          style={{ width: `${pct}%` }}
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
