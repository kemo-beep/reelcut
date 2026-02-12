import { VideoPlayer } from '../video/VideoPlayer'
import { cn } from '../../lib/utils'

export interface EditorCanvasProps {
  playbackUrl: string
  poster?: string
  onTimeUpdate?: (time: number) => void
  onPlay?: () => void
  onPause?: () => void
  className?: string
}

export function EditorCanvas({
  playbackUrl,
  poster,
  onTimeUpdate,
  onPlay,
  onPause,
  className,
}: EditorCanvasProps) {
  return (
    <div
      className={cn(
        'rounded-xl border border-[var(--app-border)] bg-[var(--app-bg)] overflow-hidden',
        className
      )}
    >
      <VideoPlayer
        src={playbackUrl}
        poster={poster}
        onTimeUpdate={onTimeUpdate}
        onPlay={onPlay}
        onPause={onPause}
      />
    </div>
  )
}
