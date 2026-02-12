import { Check, CheckCheck } from 'lucide-react'
import { Button } from '../../ui/button'
import { ViralityScore } from '../../clip/ViralityScore'
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
  const suggestionCount = suggestions.length

  return (
    <section
      aria-label="AI suggestions"
      className="border-t border-[var(--app-border)] pt-6"
    >
      <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4">
        <div className="mb-3 flex items-center justify-between">
          <div>
            <h2 className="font-semibold text-[var(--app-fg)]">
              AI suggestions{suggestionCount > 0 && ` (${suggestionCount})`}
            </h2>
            <p className="mt-0.5 text-sm text-[var(--app-fg-muted)]">
              Add these as clips with one click
            </p>
          </div>
          <Button
            size="sm"
            onClick={onAcceptAll}
            disabled={isAcceptingAll}
            className="shrink-0 bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]"
          >
            <CheckCheck size={14} className="mr-1" />
            {isAcceptingAll ? 'Creating…' : 'Accept all'}
          </Button>
        </div>
        <ul className="flex flex-wrap gap-2" aria-live="polite">
          {suggestions.map((suggestion, index) => (
            <li
              key={`${suggestion.start_time}-${suggestion.end_time}-${index}`}
              className="flex flex-wrap items-center gap-2 rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-2"
            >
              <span className="shrink-0 text-sm font-medium text-[var(--app-fg-muted)]">
                {formatTimeForLabel(suggestion.start_time)}–{formatTimeForLabel(suggestion.end_time)}
              </span>
              {suggestion.virality_score != null && (
                <span className="shrink-0">
                  <ViralityScore score={suggestion.virality_score} />
                </span>
              )}
              {suggestion.transcript && (
                <span
                  className="max-w-[240px] line-clamp-2 text-sm text-[var(--app-fg-muted)] sm:max-w-[320px]"
                  title={suggestion.transcript}
                >
                  {suggestion.transcript}
                </span>
              )}
              <Button
                variant="outline"
                size="sm"
                className="ml-auto shrink-0"
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
    </section>
  )
}
