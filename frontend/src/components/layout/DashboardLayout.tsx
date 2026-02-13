import { Link, Outlet, useNavigate } from '@tanstack/react-router'
import { useEffect } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { useWebSocket } from '../../hooks/useWebSocket'
import { useAuthHasHydrated, hasStoredAuth } from '../../stores/authStore'
import { ProjectSwitcher } from '../dashboard/ProjectSwitcher'
import {
  LayoutDashboard,
  Video,
  Image,
  LayoutTemplate,
  Settings,
  CreditCard,
  BarChart3,
} from 'lucide-react'

const workspaceNav = [
  { to: '/dashboard', label: 'Dashboard', icon: LayoutDashboard, exact: true },
  { to: '/dashboard/videos', label: 'Videos', icon: Video, exact: false },
  { to: '/dashboard/assets', label: 'Assets', icon: Image, exact: false },
  { to: '/dashboard/templates', label: 'Templates', icon: LayoutTemplate, exact: false },
]

const settingsNav = [
  { to: '/dashboard/settings/profile', label: 'Profile', icon: Settings, exact: false },
  { to: '/dashboard/settings/billing', label: 'Billing', icon: CreditCard, exact: false },
  { to: '/dashboard/settings/usage', label: 'Usage', icon: BarChart3, exact: false },
]

export default function DashboardLayout() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const hasHydrated = useAuthHasHydrated()

  useEffect(() => {
    if (typeof window === 'undefined') return
    if (!hasHydrated) return
    if (!hasStoredAuth()) {
      navigate({ to: '/auth/login', search: { redirectTo: '/dashboard' } })
    }
  }, [hasHydrated, navigate])

  useWebSocket({
    enabled: true,
    onMessage: (data: unknown) => {
      const d = data as { type?: string; job?: { entity_type?: string; entity_id?: string; status?: string } }
      if (d?.type === 'job_updated' && d.job) {
        queryClient.invalidateQueries({ queryKey: ['jobs'] })
        if (d.job.entity_type === 'clip' && d.job.entity_id) {
          queryClient.invalidateQueries({ queryKey: ['clip', d.job.entity_id] })
          queryClient.invalidateQueries({ queryKey: ['clips'] })
        }
        if (d.job.entity_type === 'video' && d.job.entity_id) {
          queryClient.invalidateQueries({ queryKey: ['video', d.job.entity_id] })
          queryClient.invalidateQueries({ queryKey: ['videos'] })
        }
      }
    },
  })

  return (
    <div className="flex min-h-screen bg-[var(--app-bg)]">
      <aside className="flex w-64 flex-shrink-0 flex-col border-r border-[var(--app-border)] bg-[var(--app-bg-raised)] shadow-card">
        {/* Project switcher */}
        <div className="border-b border-[var(--app-border)] p-3">
          <ProjectSwitcher />
        </div>

        {/* Workspace navigation */}
        <nav className="flex-1 space-y-4 p-3">
          <div className="space-y-0.5">
            <p className="mb-1.5 px-3 text-[11px] font-semibold uppercase tracking-wider text-[var(--app-fg-subtle)]">
              Workspace
            </p>
            {workspaceNav.map(({ to, label, icon: Icon, exact }) => (
              <Link
                key={to}
                to={to}
                activeOptions={{ exact }}
                className="flex items-center gap-3 rounded-lg px-3 py-2.5 text-[var(--app-fg-muted)] transition-[var(--motion-duration-fast)] hover:bg-[var(--app-bg-overlay)] hover:text-[var(--app-fg)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg-raised)]"
                activeProps={{
                  className:
                    'border-l-2 border-[var(--app-accent)] bg-[var(--app-accent-muted)] text-[var(--app-accent)] hover:bg-[var(--app-accent-muted)] hover:text-[var(--app-accent)] -ml-px pl-[13px]',
                }}
              >
                <Icon size={20} aria-hidden />
                <span className="font-medium">{label}</span>
              </Link>
            ))}
          </div>

          <div className="space-y-0.5">
            <p className="mb-1.5 px-3 text-[11px] font-semibold uppercase tracking-wider text-[var(--app-fg-subtle)]">
              Settings
            </p>
            {settingsNav.map(({ to, label, icon: Icon }) => (
              <Link
                key={to}
                to={to}
                className="flex items-center gap-3 rounded-lg px-3 py-2.5 text-[var(--app-fg-muted)] transition-[var(--motion-duration-fast)] hover:bg-[var(--app-bg-overlay)] hover:text-[var(--app-fg)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg-raised)]"
                activeProps={{
                  className:
                    'border-l-2 border-[var(--app-accent)] bg-[var(--app-accent-muted)] text-[var(--app-accent)] hover:bg-[var(--app-accent-muted)] hover:text-[var(--app-accent)] -ml-px pl-[13px]',
                }}
              >
                <Icon size={20} aria-hidden />
                <span className="font-medium">{label}</span>
              </Link>
            ))}
          </div>
        </nav>
      </aside>

      <main className="min-w-0 flex-1 overflow-auto">
        <div className="mx-auto max-w-6xl px-6 py-8">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
