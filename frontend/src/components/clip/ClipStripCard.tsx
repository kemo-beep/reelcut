import { useQuery } from '@tanstack/react-query'
import { getClipPlaybackUrl } from '../../lib/api/clips'
import { cn } from '../../lib/utils'
import type { Clip } from '../../types'

function formatDuration(sec: number): string {
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

export interface ClipStripCardProps {
  clip: Clip
  isSelected: boolean
  onSelect: () => void
  className?: string
}

export function ClipStripCard({ clip, isSelected, onSelect, className }: ClipStripCardProps) {
  const hasVideo = clip.storage_path != null && clip.storage_path !== ''
  const { data: playback } = useQuery({
    queryKey: ['clip-playback', clip.id],
    queryFn: () => getClipPlaybackUrl(clip.id),
    enabled: hasVideo,
  })
  const thumbnailUrl = playback?.url ?? null
  const duration = clip.duration_seconds != null ? formatDuration(clip.duration_seconds) : '—'

  return (
    <button
      type="button"
      onClick={onSelect}
      className={cn(
        'flex-shrink-0 w-[180px] rounded-xl border-2 overflow-hidden text-left transition-all focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]',
        isSelected
          ? 'border-[var(--app-accent)] shadow-lg shadow-[var(--app-accent)]/20'
          : 'border-[var(--app-border)] hover:border-[var(--app-border-strong)] hover:shadow-md',
        className
      )}
      aria-pressed={isSelected}
      aria-label={`Clip: ${clip.name}. ${duration}. Click to play.`}
    >
      <div className="aspect-[9/16] w-full bg-[var(--app-bg)] relative">
        {thumbnailUrl ? (
          <video
            src={thumbnailUrl}
            className="absolute inset-0 h-full w-full object-cover"
            muted
            playsInline
            preload="metadata"
            aria-hidden
          />
        ) : (
          <div className="absolute inset-0 flex items-center justify-center text-[var(--app-fg-muted)] text-sm">
            No preview
          </div>
        )}
        <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent px-2 py-2">
          <span className="text-xs font-medium text-white/90">{duration}</span>
        </div>
      </div>
      <div className="p-2.5 bg-[var(--app-bg-raised)] border-t border-[var(--app-border)]">
        <p className="text-sm font-medium text-[var(--app-fg)] truncate" title={clip.name}>
          {clip.name || `Clip`}
        </p>
      </div>
    </button>
  )
}
