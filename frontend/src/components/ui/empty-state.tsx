import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'

interface EmptyStateProps {
  icon?: ReactNode
  title: string
  description?: string
  action?: ReactNode
  className?: string
}

export function EmptyState({
  icon,
  title,
  description,
  action,
  className,
}: EmptyStateProps) {
  return (
    <div
      className={cn(
        'flex flex-col items-center justify-center rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] px-8 py-16 text-center shadow-card',
        className
      )}
      role="status"
      aria-label={`Empty: ${title}`}
    >
      {icon && (
        <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-[var(--app-accent-muted)] text-[var(--app-accent)]">
          {icon}
        </div>
      )}
      <h3 className="text-h3 mb-2 font-semibold text-[var(--app-fg)]">
        {title}
      </h3>
      {description && (
        <p className="text-caption mb-6 max-w-sm text-[var(--app-fg-muted)]">
          {description}
        </p>
      )}
      {action && <div>{action}</div>}
    </div>
  )
}
