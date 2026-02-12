import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { CreditCard } from 'lucide-react'
import { getMySubscription, cancelSubscription, createSubscription } from '../../../lib/api/subscriptions'
import { Skeleton } from '../../../components/ui/skeleton'
import { Button } from '../../../components/ui/button'
import { ErrorState } from '../../../components/ui/error-state'
import { toast } from 'sonner'

export const Route = createFileRoute('/dashboard/settings/billing')({
  component: BillingPage,
})

function BillingPage() {
  const queryClient = useQueryClient()
  const { data: sub, isLoading, error, refetch } = useQuery({
    queryKey: ['subscription'],
    queryFn: async () => {
      try {
        return await getMySubscription()
      } catch (e: unknown) {
        if (e && typeof e === 'object' && 'status' in e && (e as { status: number }).status === 404) {
          return null
        }
        throw e
      }
    },
  })
  const cancelMut = useMutation({
    mutationFn: cancelSubscription,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['subscription'] })
      toast.success('Subscription canceled')
    },
    onError: () => toast.error('Failed to cancel subscription'),
  })
  const createMut = useMutation({
    mutationFn: (tier: string) => createSubscription(tier),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['subscription'] })
      toast.success('Subscription updated')
    },
    onError: () => toast.error('Failed to update subscription'),
  })

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-32" />
        <Skeleton className="h-24 w-full max-w-md" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-6">
        <h1 className="text-2xl font-bold text-[var(--app-fg)]">Billing</h1>
        <ErrorState message="Failed to load subscription." onRetry={() => refetch()} />
      </div>
    )
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-[var(--app-fg)]">Billing</h1>
        <p className="text-caption mt-1 text-[var(--app-fg-muted)]">
          Manage your subscription and payment.
        </p>
      </div>
      <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-6 max-w-lg">
        {!sub ? (
          <div className="space-y-4">
            <p className="text-[var(--app-fg-muted)]">You are on the free plan.</p>
            <Button
              onClick={() => createMut.mutate('pro')}
              disabled={createMut.isPending}
              className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]"
            >
              {createMut.isPending ? 'Processing…' : 'Upgrade to Pro'}
            </Button>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <CreditCard size={20} className="text-[var(--app-fg-muted)]" />
              <span className="font-medium text-[var(--app-fg)] capitalize">{sub.tier}</span>
              <span className="text-caption text-[var(--app-fg-muted)]">({sub.status})</span>
            </div>
            {sub.current_period_end && (
              <p className="text-caption text-[var(--app-fg-muted)]">
                Current period ends: {new Date(sub.current_period_end).toLocaleDateString()}
              </p>
            )}
            {sub.status === 'active' && (
              <Button
                variant="outline"
                size="sm"
                className="border-[var(--app-destructive)]/50 text-[var(--app-destructive)] hover:bg-[var(--app-destructive-muted)]"
                onClick={() => cancelMut.mutate()}
                disabled={cancelMut.isPending}
              >
                {cancelMut.isPending ? 'Canceling…' : 'Cancel subscription'}
              </Button>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
