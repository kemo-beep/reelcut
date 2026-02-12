import { Loader2, CheckCircle, XCircle, Clock } from 'lucide-react'
import { ProgressBar } from './ProgressBar'
import { cn } from '../../lib/utils'

export interface JobStatusProps {
  status: string
  progress?: number
  error?: string | null
  label?: string
  className?: string
}

export function JobStatus({
  status,
  progress = 0,
  error,
  label,
  className,
}: JobStatusProps) {
  const isPending = status === 'pending' || status === 'processing'
  const isComplete = status === 'completed'
  const isFailed = status === 'failed'

  return (
    <div className={cn('flex items-center gap-3 rounded-lg border border-[var(--app-border)] bg-[var(--app-bg-raised)] px-4 py-3', className)}>
      {isPending && <Loader2 size={20} className="animate-spin text-[var(--app-accent)] shrink-0" />}
      {isComplete && <CheckCircle size={20} className="text-[var(--app-success)] shrink-0" />}
      {isFailed && <XCircle size={20} className="text-[var(--app-destructive)] shrink-0" />}
      {!isPending && !isComplete && !isFailed && <Clock size={20} className="text-[var(--app-fg-muted)] shrink-0" />}
      <div className="min-w-0 flex-1">
        {label && <p className="text-sm font-medium text-[var(--app-fg)]">{label}</p>}
        <p className="text-caption text-[var(--app-fg-muted)] capitalize">{status}</p>
        {isPending && (progress > 0 || progress === 0) && (
          <ProgressBar value={progress} className="mt-2" showLabel />
        )}
        {isFailed && error && (
          <p className="mt-1 text-sm text-[var(--app-destructive)]">{error}</p>
        )}
      </div>
    </div>
  )
}
