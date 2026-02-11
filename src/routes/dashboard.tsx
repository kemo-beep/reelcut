import { createFileRoute, redirect } from '@tanstack/react-router'
import DashboardLayout from '../components/layout/DashboardLayout'
import { useAuthStore } from '../stores/authStore'

export const Route = createFileRoute('/dashboard')({
  beforeLoad: () => {
    const token = useAuthStore.getState().getAccessToken()
    if (!token) {
      throw redirect({ to: '/auth/login', search: { redirectTo: '/dashboard' } })
    }
  },
  component: DashboardLayout,
})
