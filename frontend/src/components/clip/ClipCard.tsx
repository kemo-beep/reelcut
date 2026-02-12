import { Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Badge } from '../ui/badge'
import { ViralityScore } from './ViralityScore'
import { getClipPlaybackUrl } from '../../lib/api/clips'
import type { Clip } from '../../types'
import { cn } from '../../lib/utils'

export interface ClipCardProps {
  clip: Clip
  className?: string
}

function statusVariant(
  status: string
): 'default' | 'success' | 'warning' | 'destructive' {
  switch (status) {
    case 'ready':
      return 'success'
    case 'rendering':
      return 'warning'
    case 'failed':
      return 'destructive'
    default:
      return 'default'
  }
}

export function ClipCard({ clip, className }: ClipCardProps) {
  const hasVideo = clip.storage_path != null && clip.storage_path !== ''
  const { data: playback } = useQuery({
    queryKey: ['clip-playback', clip.id],
    queryFn: () => getClipPlaybackUrl(clip.id),
    enabled: hasVideo,
  })
  const playbackUrl = playback?.url ?? null

  return (
    <Link
      to="/dashboard/clips/$clipId"
      params={{ clipId: clip.id }}
      className={cn(
        'block rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] overflow-hidden shadow-card transition-[var(--motion-duration-fast)] hover:border-[var(--app-border-strong)] hover:shadow-lg focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)]',
        className
      )}
    >
      {hasVideo && playbackUrl && (
        <div className="aspect-video w-full bg-[var(--app-bg)]">
          <video
            src={playbackUrl}
            className="h-full w-full object-contain"
            muted
            playsInline
            preload="metadata"
            aria-label={`Preview: ${clip.name}`}
          />
        </div>
      )}
      <div className="p-4 flex items-start justify-between gap-2">
        <div className="min-w-0 flex-1">
          <p className="font-medium text-[var(--app-fg)] truncate">{clip.name}</p>
          <div className="mt-1 flex flex-wrap items-center gap-2 text-caption text-[var(--app-fg-muted)]">
            <span>
              {clip.duration_seconds != null
                ? `${Math.round(clip.duration_seconds)}s`
                : '—'}
            </span>
            <span>·</span>
            <span>{clip.aspect_ratio}</span>
            {clip.virality_score != null && (
              <>
                <span>·</span>
                <ViralityScore score={clip.virality_score} />
              </>
            )}
          </div>
        </div>
        <div className="flex items-center gap-1 shrink-0">
          {clip.is_ai_suggested && (
            <Badge variant="secondary" className="text-xs">AI</Badge>
          )}
          <Badge variant={statusVariant(clip.status)}>{clip.status}</Badge>
        </div>
      </div>
    </Link>
  )
}
