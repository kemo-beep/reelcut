import { Check, CheckCheck } from 'lucide-react'
import { Button } from '../../ui/button'
import { formatTimeForLabel } from './utils'
import type { ClipSuggestion } from '../../../lib/api/analysis'

export interface ClipSuggestionsPanelProps {
  suggestions: ClipSuggestion[]
  onAccept: (suggestion: ClipSuggestion, index: number) => void
  onAcceptAll: () => void
  isAccepting: boolean
  isAcceptingAll: boolean
}

export function ClipSuggestionsPanel({
  suggestions,
  onAccept,
  onAcceptAll,
  isAccepting,
  isAcceptingAll,
}: ClipSuggestionsPanelProps) {
  return (
    <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4">
      <div className="flex items-center justify-between mb-3">
        <h2 className="font-semibold text-[var(--app-fg)]">AI suggestions</h2>
        <Button
          size="sm"
          onClick={onAcceptAll}
          disabled={isAcceptingAll}
          className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]"
        >
          <CheckCheck size={14} className="mr-1" />
          {isAcceptingAll ? 'Creating…' : 'Accept all'}
        </Button>
      </div>
      <ul className="flex flex-wrap gap-2">
        {suggestions.map((suggestion, index) => (
          <li
            key={`${suggestion.start_time}-${suggestion.end_time}-${index}`}
            className="flex items-center gap-2 rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-2"
          >
            <span className="text-caption text-[var(--app-fg-muted)]">
              {formatTimeForLabel(suggestion.start_time)}–{formatTimeForLabel(suggestion.end_time)}
            </span>
            {suggestion.transcript && (
              <span
                className="max-w-[200px] truncate text-sm text-[var(--app-fg-muted)]"
                title={suggestion.transcript}
              >
                {suggestion.transcript}
              </span>
            )}
            <Button
              variant="outline"
              size="sm"
              className="shrink-0"
              onClick={() => onAccept(suggestion, index)}
              disabled={isAccepting}
            >
              <Check size={14} className="mr-1" />
              Accept
            </Button>
          </li>
        ))}
      </ul>
    </div>
  )
}
