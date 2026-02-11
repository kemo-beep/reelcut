import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Scissors } from 'lucide-react'
import { listClips } from '../../../lib/api/clips'
import { Skeleton } from '../../../components/ui/skeleton'
import { EmptyState } from '../../../components/ui/empty-state'
import { ErrorState } from '../../../components/ui/error-state'
import { Badge } from '../../../components/ui/badge'

export const Route = createFileRoute('/dashboard/clips/')({
  component: ClipsListPage,
})

function statusVariant(
  status: string
): 'default' | 'success' | 'warning' | 'destructive' {
  switch (status) {
    case 'ready':
      return 'success'
    case 'rendering':
      return 'warning'
    case 'failed':
      return 'destructive'
    default:
      return 'default'
  }
}

function ClipsListPage() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['clips'],
    queryFn: () => listClips({ per_page: 50 }),
  })

  if (isLoading) {
    return (
      <div className="space-y-8">
        <div>
          <Skeleton className="mb-2 h-8 w-40" />
          <Skeleton className="h-4 w-56" />
        </div>
        <ul className="space-y-3">
          {Array.from({ length: 5 }).map((_, i) => (
            <li key={i} className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4 shadow-card">
              <Skeleton className="mb-2 h-5 w-1/2" />
              <Skeleton className="h-4 w-1/3" />
            </li>
          ))}
        </ul>
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-8">
        <div>
          <h1 className="text-h1 text-[var(--app-fg)]">Clips</h1>
          <p className="text-caption mt-1">Your created clips.</p>
        </div>
        <ErrorState message="Failed to load clips." onRetry={() => refetch()} />
      </div>
    )
  }

  const clips = data?.data?.clips ?? []

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-h1 text-[var(--app-fg)]">Clips</h1>
        <p className="text-caption mt-1">Your created clips.</p>
      </div>

      {clips.length === 0 ? (
        <EmptyState
          icon={<Scissors size={28} />}
          title="No clips yet"
          description="Create clips from your videos to export short-form content."
        />
      ) : (
        <ul className="space-y-3">
          {clips.map((c) => (
            <li key={c.id}>
              <Link
                to="/dashboard/clips/$clipId"
                params={{ clipId: c.id }}
                className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4 shadow-card transition-[var(--motion-duration-fast)] hover:border-[var(--app-border-strong)] hover:shadow-lg focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]"
              >
                <span className="font-medium text-[var(--app-fg)]">{c.name}</span>
                <div className="flex items-center gap-2">
                  <span className="text-caption">
                    {c.duration_seconds != null
                      ? `${Math.round(c.duration_seconds)}s`
                      : '—'}{' '}
                    · {c.aspect_ratio}
                  </span>
                  <Badge variant={statusVariant(c.status)}>{c.status}</Badge>
                </div>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
