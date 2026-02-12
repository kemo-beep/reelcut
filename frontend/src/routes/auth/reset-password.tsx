import { createFileRoute, Link, useNavigate, useSearch } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { AlertCircle } from 'lucide-react'
import { resetPassword } from '../../lib/api/auth'
import { ApiError } from '../../types'
import { Button } from '../../components/ui/button'
import { Input } from '../../components/ui/input'
import { Label } from '../../components/ui/label'

export const Route = createFileRoute('/auth/reset-password')({
  validateSearch: (s): { token?: string } => ({
    token: typeof s?.token === 'string' ? s.token : undefined,
  }),
  component: ResetPasswordPage,
})

function ResetPasswordPage() {
  const navigate = useNavigate()
  const search = useSearch({ from: '/auth/reset-password' })
  const [token, setToken] = useState(search.token ?? '')
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)

  useEffect(() => {
    if (search.token) setToken(search.token)
  }, [search.token])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    if (password !== confirm) {
      setError('Passwords do not match.')
      return
    }
    if (!token.trim()) {
      setError('Reset link is invalid or missing. Request a new one from the forgot password page.')
      return
    }
    setLoading(true)
    try {
      await resetPassword({ token: token.trim(), new_password: password })
      setSuccess(true)
      setTimeout(() => navigate({ to: '/auth/login' }), 2000)
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.details?.length ? err.details.map((d) => d.message).join(', ') : err.message)
      } else {
        setError('Failed to reset password. The link may have expired.')
      }
    } finally {
      setLoading(false)
    }
  }

  if (success) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[var(--app-bg)] px-4 py-12">
        <div className="w-full max-w-[400px] space-y-8 rounded-2xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 shadow-modal text-center">
          <h1 className="text-h1 text-[var(--app-fg)]">Password reset</h1>
          <p className="text-caption text-[var(--app-fg-muted)]">
            Your password has been reset. Redirecting to sign in…
          </p>
          <Link to="/auth/login" className="font-medium text-[var(--app-accent)]">
            Sign in
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-[var(--app-bg)] px-4 py-12">
      <div className="w-full max-w-[400px] space-y-8 rounded-2xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 shadow-modal">
        <div className="space-y-2 text-center">
          <h1 className="text-h1 text-[var(--app-fg)]">Set new password</h1>
          <p className="text-caption text-[var(--app-fg-muted)]">
            Enter your new password below.
          </p>
        </div>
        <form onSubmit={handleSubmit} className="space-y-5">
          {error && (
            <div
              className="flex items-center gap-2 rounded-lg border border-[var(--app-destructive)]/30 bg-[var(--app-destructive-muted)] px-4 py-3 text-sm text-[var(--app-destructive)]"
              role="alert"
            >
              <AlertCircle size={18} className="flex-shrink-0" />
              <span>{error}</span>
            </div>
          )}
          {!search.token && (
            <div className="space-y-2">
              <Label htmlFor="token">Reset token</Label>
              <Input
                id="token"
                type="text"
                value={token}
                onChange={(e) => setToken(e.target.value)}
                placeholder="Paste the token from your email"
                className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)] font-mono text-sm"
              />
            </div>
          )}
          <div className="space-y-2">
            <Label htmlFor="password">New password</Label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength={8}
              autoComplete="new-password"
              className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)]"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="confirm">Confirm password</Label>
            <Input
              id="confirm"
              type="password"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              required
              minLength={8}
              autoComplete="new-password"
              className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)]"
            />
          </div>
          <Button
            type="submit"
            className="h-11 w-full bg-[var(--app-accent)] font-semibold text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]"
            disabled={loading}
          >
            {loading ? 'Resetting…' : 'Reset password'}
          </Button>
        </form>
        <p className="text-center text-caption">
          <Link
            to="/auth/forgot-password"
            className="font-medium text-[var(--app-accent)] hover:text-[var(--app-accent-hover)]"
          >
            Request a new link
          </Link>
        </p>
      </div>
    </div>
  )
}
