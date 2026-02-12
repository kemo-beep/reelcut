import { createFileRoute, redirect } from '@tanstack/react-router'
import DashboardLayout from '../components/layout/DashboardLayout'
import { shouldRedirectToLogin } from '../stores/authStore'

export const Route = createFileRoute('/dashboard')({
  beforeLoad: () => {
    if (shouldRedirectToLogin()) {
      throw redirect({ to: '/auth/login', search: { redirectTo: '/dashboard' } })
    }
  },
  component: DashboardLayout,
})
