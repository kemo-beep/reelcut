import { cn } from '../../lib/utils'

export interface ProgressBarProps {
  value: number
  max?: number
  className?: string
  showLabel?: boolean
}

export function ProgressBar({
  value,
  max = 100,
  className,
  showLabel = false,
}: ProgressBarProps) {
  const pct = Math.min(100, Math.max(0, max > 0 ? (value / max) * 100 : 0))
  return (
    <div className={cn('w-full', className)} role="progressbar" aria-valuenow={value} aria-valuemin={0} aria-valuemax={max}>
      <div className="h-2 w-full overflow-hidden rounded-full bg-[var(--app-bg)] border border-[var(--app-border)]">
        <div
          className="h-full rounded-full bg-[var(--app-accent)] transition-[width] duration-200"
          style={{ width: `${pct}%` }}
        />
      </div>
      {showLabel && (
        <p className="mt-1 text-caption text-[var(--app-fg-muted)]">{Math.round(pct)}%</p>
      )}
    </div>
  )
}
