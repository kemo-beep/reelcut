import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { AlertCircle } from 'lucide-react'
import { useAuth } from '../../hooks/useAuth'
import { ApiError } from '../../types'
import { Button } from '../../components/ui/button'
import { Input } from '../../components/ui/input'
import { Label } from '../../components/ui/label'

export const Route = createFileRoute('/auth/login')({
  component: LoginPage,
})

function LoginPage() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await login(email, password)
      navigate({ to: '/dashboard' })
    } catch (err) {
      if (err instanceof ApiError) {
        const msg = err.details?.length
          ? err.details.map((d) => d.message).join(', ')
          : err.message
        setError(msg)
      } else {
        setError('Login failed. Please try again.')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-[var(--app-bg)] px-4 py-12">
      <div className="w-full max-w-[400px] space-y-8 rounded-2xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 shadow-modal">
        <div className="space-y-2 text-center">
          <h1 className="text-h1 text-[var(--app-fg)]">Sign in</h1>
          <p className="text-caption">
            Sign in to your Reelcut account
          </p>
        </div>
        <form onSubmit={handleSubmit} className="space-y-5">
          {error && (
            <div
              className="flex items-center gap-2 rounded-lg border border-[var(--app-destructive)]/30 bg-[var(--app-destructive-muted)] px-4 py-3 text-sm text-[var(--app-destructive)]"
              role="alert"
            >
              <AlertCircle size={18} className="flex-shrink-0" aria-hidden />
              <span>{error}</span>
            </div>
          )}
          <div className="space-y-2">
            <Label htmlFor="email" className="text-label">
              Email
            </Label>
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoComplete="email"
              className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)] placeholder:text-[var(--app-fg-subtle)] focus-visible:ring-[var(--app-accent)]"
            />
          </div>
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label htmlFor="password" className="text-label">
                Password
              </Label>
              <Link
                to="/auth/forgot-password"
                className="text-xs font-medium text-[var(--app-fg-muted)] hover:text-[var(--app-accent)]"
              >
                Forgot password?
              </Link>
            </div>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              autoComplete="current-password"
              className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)] placeholder:text-[var(--app-fg-subtle)] focus-visible:ring-[var(--app-accent)]"
            />
          </div>
          <Button
            type="submit"
            className="h-11 w-full bg-[var(--app-accent)] font-semibold text-[#0a0a0b] hover:bg-[var(--app-accent-hover)] focus-visible:ring-[var(--app-accent)]"
            disabled={loading}
          >
            {loading ? 'Signing in…' : 'Sign in'}
          </Button>
        </form>
        <p className="text-center text-caption">
          Don&apos;t have an account?{' '}
          <Link
            to="/auth/register"
            className="font-medium text-[var(--app-accent)] hover:text-[var(--app-accent-hover)] focus-visible:underline focus-visible:outline-none"
          >
            Register
          </Link>
        </p>
      </div>
    </div>
  )
}
