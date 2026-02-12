import { JobStatus } from './JobStatus'
import type { ProcessingJob } from '../../types'

export interface ProcessingQueueProps {
  jobs: ProcessingJob[]
  isLoading?: boolean
  className?: string
}

export function ProcessingQueue({
  jobs,
  isLoading,
  className = '',
}: ProcessingQueueProps) {
  if (isLoading) {
    return (
      <div className={`space-y-2 ${className}`}>
        <div className="h-14 rounded-lg border border-[var(--app-border)] bg-[var(--app-bg-raised)] animate-pulse" />
        <div className="h-14 rounded-lg border border-[var(--app-border)] bg-[var(--app-bg-raised)] animate-pulse" />
      </div>
    )
  }
  if (jobs.length === 0) {
    return (
      <div className={`rounded-lg border border-[var(--app-border)] bg-[var(--app-bg-raised)] px-4 py-6 text-center text-caption text-[var(--app-fg-muted)] ${className}`}>
        No active or recent jobs
      </div>
    )
  }
  return (
    <ul className={`space-y-2 ${className}`} role="list">
      {jobs.map((job) => (
        <li key={job.id}>
          <JobStatus
            status={job.status}
            progress={job.progress ?? 0}
            error={job.error_message ?? undefined}
            label={`${job.job_type} · ${job.entity_type}`}
          />
        </li>
      ))}
    </ul>
  )
}
