import { TrendingUp } from 'lucide-react'
import { cn } from '../../lib/utils'

export interface ViralityScoreProps {
  score: number
  className?: string
}

export function ViralityScore({ score, className }: ViralityScoreProps) {
  const pct = Math.round(Math.min(100, Math.max(0, score)))
  const variant =
    pct >= 70 ? 'success'
    : pct >= 40 ? 'warning'
    : 'muted'
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 text-caption',
        variant === 'success' && 'text-[var(--app-success)]',
        variant === 'warning' && 'text-[var(--app-fg-muted)]',
        variant === 'muted' && 'text-[var(--app-fg-subtle)]',
        className
      )}
      title="Virality score"
    >
      <TrendingUp size={14} />
      {pct}
    </span>
  )
}
