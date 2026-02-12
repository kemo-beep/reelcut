import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { AlertCircle } from 'lucide-react'
import { useAuth } from '../../hooks/useAuth'
import { ApiError } from '../../types'
import { Button } from '../../components/ui/button'
import { Input } from '../../components/ui/input'
import { Label } from '../../components/ui/label'

export const Route = createFileRoute('/auth/register')({
  component: RegisterPage,
})

function RegisterPage() {
  const { register } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [fullName, setFullName] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await register(email, password, fullName || undefined)
      navigate({ to: '/dashboard' })
    } catch (err) {
      if (err instanceof ApiError) {
        const msg = err.details?.length
          ? err.details.map((d) => d.message).join(', ')
          : err.message
        setError(msg)
      } else {
        setError('Registration failed. Please try again.')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-[var(--app-bg)] px-4 py-12">
      <div className="w-full max-w-[400px] space-y-8 rounded-2xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-8 shadow-modal">
        <div className="space-y-2 text-center">
          <h1 className="text-h1 text-[var(--app-fg)]">Create account</h1>
          <p className="text-caption">
            Create your Reelcut account
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
            <Label htmlFor="fullName" className="text-label">
              Full name (optional)
            </Label>
            <Input
              id="fullName"
              type="text"
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
              autoComplete="name"
              className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)] placeholder:text-[var(--app-fg-subtle)] focus-visible:ring-[var(--app-accent)]"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password" className="text-label">
              Password
            </Label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength={8}
              autoComplete="new-password"
              className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)] placeholder:text-[var(--app-fg-subtle)] focus-visible:ring-[var(--app-accent)]"
            />
            <p className="text-caption text-[var(--app-fg-subtle)]">
              At least 8 characters, with uppercase, lowercase and a digit
            </p>
          </div>
          <Button
            type="submit"
            className="h-11 w-full bg-[var(--app-accent)] font-semibold text-[#0a0a0b] hover:bg-[var(--app-accent-hover)] focus-visible:ring-[var(--app-accent)]"
            disabled={loading}
          >
            {loading ? 'Creating account…' : 'Create account'}
          </Button>
        </form>
        <p className="text-center text-caption">
          Already have an account?{' '}
          <Link
            to="/auth/login"
            className="font-medium text-[var(--app-accent)] hover:text-[var(--app-accent-hover)] focus-visible:underline focus-visible:outline-none"
          >
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}
