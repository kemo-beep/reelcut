import { createFileRoute, Outlet, redirect, useNavigate } from '@tanstack/react-router'
import { useEffect } from 'react'
import { shouldRedirectToLogin, useAuthHasHydrated, hasStoredAuth } from '../stores/authStore'

export const Route = createFileRoute('/editor')({
  beforeLoad: () => {
    if (shouldRedirectToLogin()) {
      throw redirect({ to: '/auth/login', search: { redirectTo: '/dashboard' } })
    }
  },
  component: EditorLayout,
})

function EditorLayout() {
  const navigate = useNavigate()
  const hasHydrated = useAuthHasHydrated()

  useEffect(() => {
    if (typeof window === 'undefined') return
    if (!hasHydrated) return
    if (!hasStoredAuth()) {
      navigate({ to: '/auth/login', search: { redirectTo: '/dashboard' } })
    }
  }, [hasHydrated, navigate])

  return (
    <div className="min-h-screen bg-[var(--app-bg)]">
      <Outlet />
    </div>
  )
}
