import type { TranscriptWord } from '../../types'
import { cn } from '../../lib/utils'

export interface WordTimelineProps {
  words: TranscriptWord[]
  duration: number
  currentTime: number
  onSeek?: (time: number) => void
  className?: string
}

export function WordTimeline({
  words,
  duration,
  currentTime,
  onSeek,
  className,
}: WordTimelineProps) {
  if (duration <= 0) return null
  return (
    <div
      className={cn('flex flex-wrap gap-1 text-sm', className)}
      role="region"
      aria-label="Word timeline"
    >
      {words.map((w) => {
        const isActive = currentTime >= w.start_time && currentTime <= w.end_time
        return (
          <button
            key={w.id}
            type="button"
            onClick={() => onSeek?.(w.start_time)}
            className={cn(
              'rounded px-1 py-0.5 transition-colors',
              isActive
                ? 'bg-[var(--app-accent)] text-[#0a0a0b]'
                : 'text-[var(--app-fg)] hover:bg-[var(--app-bg-overlay)]'
            )}
          >
            {w.word}
          </button>
        )
      })}
    </div>
  )
}
