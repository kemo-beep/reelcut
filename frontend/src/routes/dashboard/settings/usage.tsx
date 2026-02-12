import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { BarChart3 } from 'lucide-react'
import { getUsageStats } from '../../../lib/api/users'
import { Skeleton } from '../../../components/ui/skeleton'
import { ErrorState } from '../../../components/ui/error-state'

export const Route = createFileRoute('/dashboard/settings/usage')({
  component: UsagePage,
})

function UsagePage() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['usage'],
    queryFn: () => getUsageStats({ per_page: 50 }),
  })

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-28" />
        <Skeleton className="h-40 w-full max-w-2xl" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-6">
        <h1 className="text-2xl font-bold text-[var(--app-fg)]">Usage</h1>
        <ErrorState message="Failed to load usage." onRetry={() => refetch()} />
      </div>
    )
  }

  const usage = data as {
    data?: { logs?: Array<{ action: string; credits_used?: number; created_at?: string }> }
    total?: number
  } | undefined
  const logs = usage?.data?.logs ?? usage?.logs ?? []

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-[var(--app-fg)]">Usage</h1>
        <p className="text-caption mt-1 text-[var(--app-fg-muted)]">
          View your usage and credit history.
        </p>
      </div>
      <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-6">
        <div className="flex items-center gap-2 mb-4">
          <BarChart3 size={20} className="text-[var(--app-fg-muted)]" />
          <span className="font-medium text-[var(--app-fg)]">Recent activity</span>
        </div>
        {logs.length === 0 ? (
          <p className="text-caption text-[var(--app-fg-muted)]">No usage records yet.</p>
        ) : (
          <ul className="space-y-2">
            {logs.slice(0, 20).map((log: { action?: string; credits_used?: number; created_at?: string }, i: number) => (
              <li
                key={i}
                className="flex items-center justify-between rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] px-4 py-2 text-sm"
              >
                <span className="text-[var(--app-fg)]">{log.action ?? '—'}</span>
                <span className="text-caption text-[var(--app-fg-muted)]">
                  {log.credits_used != null ? `${log.credits_used} credits` : ''}
                  {log.created_at ? ` · ${new Date(log.created_at).toLocaleString()}` : ''}
                </span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}
