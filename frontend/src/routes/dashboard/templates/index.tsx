import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { LayoutTemplate } from 'lucide-react'
import { listTemplates } from '../../../lib/api/templates'
import { Skeleton } from '../../../components/ui/skeleton'
import { EmptyState } from '../../../components/ui/empty-state'
import { ErrorState } from '../../../components/ui/error-state'

export const Route = createFileRoute('/dashboard/templates/')({
  component: TemplatesListPage,
})

function TemplatesListPage() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['templates'],
    queryFn: () => listTemplates({ per_page: 50 }),
  })

  if (isLoading) {
    return (
      <div className="space-y-8">
        <div>
          <Skeleton className="mb-2 h-8 w-44" />
          <Skeleton className="h-4 w-52" />
        </div>
        <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <li key={i} className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-5 shadow-card">
              <Skeleton className="mb-2 h-5 w-2/3" />
              <Skeleton className="h-4 w-1/2" />
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
          <h1 className="text-h1 text-[var(--app-fg)]">Templates</h1>
          <p className="text-caption mt-1">Style templates for your clips.</p>
        </div>
        <ErrorState message="Failed to load templates." onRetry={() => refetch()} />
      </div>
    )
  }

  const templates = data?.data?.templates ?? []

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-h1 text-[var(--app-fg)]">Templates</h1>
        <p className="text-caption mt-1">Style templates for your clips.</p>
      </div>

      {templates.length === 0 ? (
        <EmptyState
          icon={<LayoutTemplate size={28} />}
          title="No templates yet"
          description="Create or use templates to apply consistent styles to your clips."
        />
      ) : (
        <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {templates.map((t) => (
            <li
              key={t.id}
              className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-5 shadow-card"
            >
              <span className="font-medium text-[var(--app-fg)]">{t.name}</span>
              {t.category && (
                <p className="text-caption mt-1">{t.category}</p>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
