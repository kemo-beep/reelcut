import { Link } from '@tanstack/react-router'
import { Video } from 'lucide-react'
import { Badge } from '../ui/badge'
import type { Video as VideoType } from '../../types'
import { getThumbnailUrl } from '../../lib/api/videos'
import { cn } from '../../lib/utils'

export interface VideoCardProps {
  video: VideoType
  className?: string
}

function statusVariant(
  status: string
): 'default' | 'success' | 'warning' | 'destructive' {
  switch (status) {
    case 'ready':
      return 'success'
    case 'processing':
    case 'uploading':
      return 'warning'
    case 'failed':
      return 'destructive'
    default:
      return 'default'
  }
}

export function VideoCard({ video, className }: VideoCardProps) {
  const thumbUrl =
    video.thumbnail_display_url?.startsWith('http') ?
      video.thumbnail_display_url
      : getThumbnailUrl(video.id)

  return (
    <Link
      to="/dashboard/videos/$videoId"
      params={{ videoId: video.id }}
      className={cn(
        'block overflow-hidden rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] shadow-card transition-[var(--motion-duration-fast)] hover:border-[var(--app-border-strong)] hover:shadow-lg focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]',
        className
      )}
    >
      <div className="aspect-video w-full bg-[var(--app-bg)]">
        {video.status === 'ready' ? (
          <img
            src={thumbUrl}
            alt=""
            className="h-full w-full object-cover"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center text-[var(--app-fg-subtle)]">
            <Video size={40} aria-hidden />
          </div>
        )}
      </div>
      <div className="p-4">
        <p className="font-medium text-[var(--app-fg)] truncate">
          {video.original_filename}
        </p>
        <div className="mt-2 flex items-center gap-2">
          <span className="text-caption text-[var(--app-fg-muted)]">
            {video.duration_seconds != null
              ? `${Math.round(video.duration_seconds)}s`
              : '—'}
          </span>
          <Badge variant={statusVariant(video.status)}>{video.status}</Badge>
        </div>
      </div>
    </Link>
  )
}
