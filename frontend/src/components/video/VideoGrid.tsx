import { VideoCard } from './VideoCard'
import type { Video } from '../../types'
import { cn } from '../../lib/utils'

export interface VideoGridProps {
  videos: Video[]
  className?: string
}

export function VideoGrid({ videos, className }: VideoGridProps) {
  return (
    <ul
      className={cn(
        'grid gap-4 sm:grid-cols-2 lg:grid-cols-3',
        className
      )}
      role="list"
    >
      {videos.map((v) => (
        <li key={v.id}>
          <VideoCard video={v} />
        </li>
      ))}
    </ul>
  )
}
