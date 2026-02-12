import { cn } from '../../../lib/utils'

export type ViewMode = 'clips' | 'main-video'

export interface ViewModeTabsProps {
  value: ViewMode
  onChange: (mode: ViewMode) => void
}

export function ViewModeTabs({ value, onChange }: ViewModeTabsProps) {
  return (
    <div
      className="flex items-center gap-2 border-t border-[var(--app-border)] pt-6"
      role="group"
      aria-labelledby="view-mode-label"
    >
      <span id="view-mode-label" className="text-sm font-medium text-[var(--app-fg-muted)]">
        View:
      </span>
      <div
        role="tablist"
        aria-label="Switch between Clips and Main video"
        className="inline-flex rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] p-0.5"
      >
        <button
          type="button"
          role="tab"
          aria-selected={value === 'clips'}
          onClick={() => onChange('clips')}
          className={cn(
            'rounded-md px-4 py-2 text-sm font-medium transition-colors',
            value === 'clips'
              ? 'bg-[var(--app-accent)] text-[#0a0a0b]'
              : 'text-[var(--app-fg-muted)] hover:text-[var(--app-fg)] hover:bg-[var(--app-bg-raised)]'
          )}
        >
          Clips
        </button>
        <button
          type="button"
          role="tab"
          aria-selected={value === 'main-video'}
          onClick={() => onChange('main-video')}
          className={cn(
            'rounded-md px-4 py-2 text-sm font-medium transition-colors',
            value === 'main-video'
              ? 'bg-[var(--app-accent)] text-[#0a0a0b]'
              : 'text-[var(--app-fg-muted)] hover:text-[var(--app-fg)] hover:bg-[var(--app-bg-raised)]'
          )}
        >
          Main video
        </button>
      </div>
    </div>
  )
}
