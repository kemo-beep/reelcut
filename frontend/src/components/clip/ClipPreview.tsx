import { VideoPlayer } from '../video/VideoPlayer'
import { cn } from '../../lib/utils'

export interface ClipPreviewProps {
  playbackUrl: string
  poster?: string
  className?: string
}

export function ClipPreview({
  playbackUrl,
  poster,
  className,
}: ClipPreviewProps) {
  return (
    <div className={cn('overflow-hidden rounded-xl bg-[var(--app-bg)]', className)}>
      <VideoPlayer src={playbackUrl} poster={poster} />
    </div>
  )
}
