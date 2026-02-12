import { AlertCircle } from 'lucide-react'
import { cn } from '@/lib/utils'

interface ErrorStateProps {
  message: string
  onRetry?: () => void
  className?: string
}

export function ErrorState({ message, onRetry, className }: ErrorStateProps) {
  return (
    <div
      className={cn(
        'flex flex-col items-center justify-center gap-4 rounded-xl border border-[var(--app-destructive)]/30 bg-[var(--app-destructive-muted)] px-6 py-8 text-center',
        className
      )}
      role="alert"
    >
      <div className="flex h-12 w-12 items-center justify-center rounded-full bg-[var(--app-destructive)]/20 text-[var(--app-destructive)]">
        <AlertCircle className="h-6 w-6" aria-hidden />
      </div>
      <p className="text-body text-[var(--app-fg)]">{message}</p>
      {onRetry && (
        <button
          type="button"
          onClick={onRetry}
          className="rounded-lg border border-[var(--app-border-strong)] bg-[var(--app-bg-raised)] px-4 py-2 text-sm font-medium text-[var(--app-fg)] transition-[var(--motion-duration-fast)] hover:bg-[var(--app-bg-overlay)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]"
        >
          Retry
        </button>
      )}
    </div>
  )
}
