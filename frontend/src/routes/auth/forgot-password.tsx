import { createFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { AlertCircle, CheckCircle } from 'lucide-react'
import { forgotPassword } from '../../lib/api/auth'
import { ApiError } from '../../types'
import { Button } from '../../components/ui/button'
import { Input } from '../../components/ui/input'
import { Label } from '../../components/ui/label'

export const Route = createFileRoute('/auth/forgot-password')({
  component: ForgotPasswordPage,
})

function ForgotPasswordPage() {
  const [email, setEmail] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [sent, setSent] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await forgotPassword(email)
      setSent(true)
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.details?.length ? err.details.map((d) => d.message).join(', ') : err.message)
      } else {
        setError('Request failed. Please try again.')
      }
    } finally {
      setLoading(false)
    }
  }

  if (sent) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[var(--app-bg)] px-4 py-12">
        <div className="w-full max-w-[400px] space-y-8 rounded-2xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 shadow-modal">
          <div className="space-y-2 text-center">
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-[var(--app-success-muted)]">
              <CheckCircle size={24} className="text-[var(--app-success)]" />
            </div>
            <h1 className="text-h1 text-[var(--app-fg)]">Check your email</h1>
            <p className="text-caption text-[var(--app-fg-muted)]">
              If an account exists for {email}, we&apos;ve sent a link to reset your password.
            </p>
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
      <div className="w-full max-w-[400px] space-y-8 rounded-2xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 shadow-modal">
        <div className="space-y-2 text-center">
          <h1 className="text-h1 text-[var(--app-fg)]">Forgot password</h1>
          <p className="text-caption text-[var(--app-fg-muted)]">
            Enter your email and we&apos;ll send a link to reset your password.
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
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoComplete="email"
              className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)]"
            />
          </div>
          <Button
            type="submit"
            className="h-11 w-full bg-[var(--app-accent)] font-semibold text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]"
            disabled={loading}
          >
            {loading ? 'Sending…' : 'Send reset link'}
          </Button>
        </form>
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
