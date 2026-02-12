import { createFileRoute, redirect } from '@tanstack/react-router'
import DashboardLayout from '../components/layout/DashboardLayout'
import { shouldRedirectToLogin } from '../stores/authStore'

export const Route = createFileRoute('/dashboard')({
  // Disable SSR for dashboard to avoid 500 when backend is unreachable during server render
  // (e.g. fetch to API throws HTTPError in some runtimes). Dashboard loads fine on client.
  ssr: false,
  beforeLoad: () => {
    if (shouldRedirectToLogin()) {
      throw redirect({ to: '/auth/login', search: { redirectTo: '/dashboard' } })
    }
  },
  component: DashboardLayout,
})
