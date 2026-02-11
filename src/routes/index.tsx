import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useEffect } from 'react'
import { useAuthStore } from '../stores/authStore'

export const Route = createFileRoute('/')({ component: HomePage })

function HomePage() {
  const accessToken = useAuthStore((s) => s.accessToken)
  const navigate = useNavigate()

  useEffect(() => {
    if (accessToken) {
      navigate({ to: '/dashboard' })
    }
  }, [accessToken, navigate])

  if (accessToken) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[var(--app-bg)]">
        <div className="flex flex-col items-center gap-3">
          <div className="h-8 w-8 animate-pulse rounded-full bg-[var(--app-accent-muted)]" />
          <p className="text-caption">Redirecting to dashboard...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-[var(--app-bg)]">
      <section className="relative overflow-hidden px-6 py-24 text-center md:py-32">
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_80%_50%_at_50%_-20%,var(--app-accent-muted),transparent)]" />
        <div className="relative mx-auto max-w-4xl">
          <h1 className="text-display mb-6 text-[var(--app-fg)] md:mb-8">
            <span className="bg-gradient-to-r from-[var(--app-accent)] to-cyan-300 bg-clip-text text-transparent">
              Reelcut
            </span>
          </h1>
          <p className="text-body mx-auto mb-10 max-w-2xl text-[var(--app-fg-muted)] md:mb-12 md:text-lg">
            AI-powered video clip generation. Turn long-form videos into engaging short clips for social media.
          </p>
          <div className="flex flex-wrap items-center justify-center gap-4">
            <Link
              to="/auth/register"
              className="inline-flex items-center justify-center rounded-lg bg-[var(--app-accent)] px-8 py-3.5 text-base font-semibold text-[#0a0a0b] shadow-card transition-[var(--motion-duration-fast)] hover:bg-[var(--app-accent-hover)] hover:shadow-lg focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]"
            >
              Create account
            </Link>
            <Link
              to="/auth/login"
              className="inline-flex items-center justify-center rounded-lg border border-[var(--app-border-strong)] bg-[var(--app-bg-raised)] px-8 py-3.5 text-base font-medium text-[var(--app-fg)] transition-[var(--motion-duration-fast)] hover:bg-[var(--app-bg-overlay)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]"
            >
              Sign in
            </Link>
          </div>
          <p className="text-caption mt-12 text-[var(--app-fg-subtle)]">
            Trusted by creators and teams to scale short-form content.
          </p>
        </div>
      </section>
    </div>
  )
}
