import { createFileRoute, Link, useSearch } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { AlertCircle, CheckCircle } from 'lucide-react'
import { verifyEmail } from '../../lib/api/auth'
import { ApiError } from '../../types'
import { Button } from '../../components/ui/button'

export const Route = createFileRoute('/auth/verify-email')({
  validateSearch: (s): { token?: string } => ({
    token: typeof s?.token === 'string' ? s.token : undefined,
  }),
  component: VerifyEmailPage,
})

function VerifyEmailPage() {
  const search = useSearch({ from: '/auth/verify-email' })
  const [status, setStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle')
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const token = search.token
    if (!token) {
      setStatus('error')
      setError('Verification link is missing. Check your email for the full link.')
      return
    }
    let cancelled = false
    setStatus('loading')
    verifyEmail(token)
      .then(() => {
        if (!cancelled) setStatus('success')
      })
      .catch((err) => {
        if (!cancelled) {
          setStatus('error')
          setError(err instanceof ApiError ? err.message : 'Verification failed. The link may have expired.')
        }
      })
    return () => { cancelled = true }
  }, [search.token])

  if (status === 'success') {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[var(--app-bg)] px-4 py-12">
        <div className="w-full max-w-[400px] space-y-8 rounded-2xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 shadow-modal text-center">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-[var(--app-success-muted)]">
            <CheckCircle size={24} className="text-[var(--app-success)]" />
          </div>
          <h1 className="text-h1 text-[var(--app-fg)]">Email verified</h1>
          <p className="text-caption text-[var(--app-fg-muted)]">
            Your email has been verified. You can now sign in.
          </p>
          <Link to="/auth/login">
            <Button className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]">
              Sign in
            </Button>
          </Link>
        </div>
      </div>
    )
  }

  if (status === 'error') {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[var(--app-bg)] px-4 py-12">
        <div className="w-full max-w-[400px] space-y-8 rounded-2xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 shadow-modal">
          <div className="flex items-center gap-2 rounded-lg border border-[var(--app-destructive)]/30 bg-[var(--app-destructive-muted)] px-4 py-3 text-sm text-[var(--app-destructive)]">
            <AlertCircle size={18} className="flex-shrink-0" />
            <span>{error}</span>
          </div>
          <p className="text-center text-caption">
            <Link
              to="/auth/login"
              className="font-medium text-[var(--app-accent)] hover:text-[var(--app-accent-hover)]"
            >
              Back to sign in
            </Link>
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-[var(--app-bg)] px-4 py-12">
      <div className="w-full max-w-[400px] space-y-8 rounded-2xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 shadow-modal text-center">
        <h1 className="text-h1 text-[var(--app-fg)]">Verifying your email…</h1>
        <p className="text-caption text-[var(--app-fg-muted)]">
          {status === 'loading' ? 'Please wait.' : 'Redirecting…'}
        </p>
      </div>
    </div>
  )
}
