import type { TranscriptSegment } from '../../types'
import { cn } from '../../lib/utils'

export interface SpeakerLabelsProps {
  segments: TranscriptSegment[]
  currentTime?: number
  className?: string
}

export function SpeakerLabels({
  segments,
  currentTime = 0,
  className,
}: SpeakerLabelsProps) {
  const activeSegment = segments.find(
    (s) => currentTime >= s.start_time && currentTime <= s.end_time
  )
  const speakerId = activeSegment?.speaker_id
  if (speakerId == null) return null
  return (
    <span
      className={cn(
        'inline-block rounded bg-[var(--app-bg-overlay)] px-2 py-0.5 text-caption text-[var(--app-fg-muted)]',
        className
      )}
    >
      Speaker {speakerId + 1}
    </span>
  )
}
