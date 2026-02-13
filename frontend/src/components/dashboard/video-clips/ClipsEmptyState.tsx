import { Link } from '@tanstack/react-router'
import { Scissors } from 'lucide-react'
import { Button } from '../../ui/button'

export interface ClipsEmptyStateProps {
  onSuggestClips: () => void
  isSuggesting: boolean
  /** When set, "Create clip" links to the Clips page with this video pre-selected (matches header action). */
  videoId?: string
}

export function ClipsEmptyState({
  onSuggestClips,
  isSuggesting,
  videoId,
}: ClipsEmptyStateProps) {
  return (
    <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] px-8 py-10 text-center">
      <Scissors size={40} className="mx-auto text-[var(--app-fg-muted)]" aria-hidden />
      <p className="mt-2 font-medium text-[var(--app-fg)]">No clips yet</p>
      <p className="mt-1 text-sm text-[var(--app-fg-muted)]">
        Use AI suggest clips above, or create a clip from the Clips page.
      </p>
      <div className="mt-6 flex flex-wrap justify-center gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={onSuggestClips}
          disabled={isSuggesting}
          className="border-[var(--app-border)]"
        >
          AI suggest clips
        </Button>
        {videoId ? (
          <Link to="/dashboard/videos/$videoId/clips" params={{ videoId }}>
            <Button
              size="sm"
              className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]"
            >
              Create clip
            </Button>
          </Link>
        ) : (
          <Link to="/dashboard/videos">
            <Button
              size="sm"
              className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]"
            >
              Create clip
            </Button>
          </Link>
        )}
      </div>
    </div>
  )
}
