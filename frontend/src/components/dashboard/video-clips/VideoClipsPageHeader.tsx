import { Link } from '@tanstack/react-router'
import { ChevronRight, Plus, Sparkles } from 'lucide-react'
import { Button } from '../../ui/button'

export interface VideoClipsPageHeaderProps {
  videoId: string
  videoFilename: string
  onSuggestClips: () => void
  isSuggesting: boolean
}

export function VideoClipsPageHeader({
  videoId,
  videoFilename,
  onSuggestClips,
  isSuggesting,
}: VideoClipsPageHeaderProps) {
  return (
    <header className="flex flex-wrap items-start justify-between gap-4">
      <div className="min-w-0">
        <nav aria-label="Breadcrumb" className="flex items-center gap-1.5 text-sm">
          <Link
            to="/dashboard/videos"
            className="text-[var(--app-accent)] hover:underline focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)] rounded"
          >
            Videos
          </Link>
          <ChevronRight
            size={16}
            className="shrink-0 text-[var(--app-fg-muted)]"
            aria-hidden
          />
          <Link
            to="/dashboard/videos/$videoId"
            params={{ videoId }}
            className="truncate text-[var(--app-accent)] hover:underline focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)] rounded"
          >
            {videoFilename}
          </Link>
        </nav>
        <h1 className="mt-1.5 text-2xl font-bold text-[var(--app-fg)]">Clips</h1>
        <p className="mt-0.5 text-sm text-[var(--app-fg-muted)]">
          Review and trim clips from this video
        </p>
      </div>
      <div className="flex shrink-0 flex-wrap items-center gap-2 sm:gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={onSuggestClips}
          disabled={isSuggesting}
          className="border-[var(--app-border)]"
        >
          <Sparkles size={16} className="mr-1" />
          {isSuggesting ? 'Analyzing…' : 'AI suggest clips'}
        </Button>
        <Link to="/dashboard/videos/$videoId/clips" params={{ videoId }}>
          <Button size="sm" className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]">
            <Plus size={16} className="mr-1" />
            Create clip
          </Button>
        </Link>
      </div>
    </header>
  )
}
