import { Link } from '@tanstack/react-router'
import { Scissors } from 'lucide-react'
import { Button } from '../../ui/button'

export interface ClipsEmptyStateProps {
  onSuggestClips: () => void
  isSuggesting: boolean
}

export function ClipsEmptyState({ onSuggestClips, isSuggesting }: ClipsEmptyStateProps) {
  return (
    <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 text-center">
      <Scissors size={40} className="mx-auto text-[var(--app-fg-muted)]" />
      <p className="mt-2 font-medium text-[var(--app-fg)]">No clips yet</p>
      <p className="text-caption text-[var(--app-fg-muted)] mt-1">
        Use AI suggest clips above, or create a clip from the Clips page.
      </p>
      <div className="mt-4 flex justify-center gap-2">
        <Button variant="outline" size="sm" onClick={onSuggestClips} disabled={isSuggesting}>
          AI suggest clips
        </Button>
        <Link to="/dashboard/clips">
          <Button size="sm" className="bg-[var(--app-accent)] text-[#0a0a0b]">
            Go to Clips
          </Button>
        </Link>
      </div>
    </div>
  )
}
