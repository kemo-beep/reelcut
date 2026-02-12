import { Link, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useAuthStore } from '../stores/authStore'
import { ThemeToggle } from './ThemeToggle'
import { Menu, X, LayoutDashboard, Settings, LogOut, ChevronDown } from 'lucide-react'

function HeaderAuth() {
  const accessToken = useAuthStore((s) => s.accessToken)
  const user = useAuthStore((s) => s.user)
  const clearAuth = useAuthStore((s) => s.clearAuth)
  const navigate = useNavigate()
  const [dropdownOpen, setDropdownOpen] = useState(false)

  if (accessToken && user) {
    return (
      <div className="relative">
        <button
          type="button"
          onClick={() => setDropdownOpen((o) => !o)}
          className="flex items-center gap-2 rounded-lg px-3 py-2 text-sm text-[var(--app-fg-muted)] transition-[var(--motion-duration-fast)] hover:bg-[var(--app-bg-overlay)] hover:text-[var(--app-fg)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg-raised)]"
          aria-expanded={dropdownOpen}
          aria-haspopup="true"
        >
          <span className="max-w-[140px] truncate">
            {user.full_name || user.email}
          </span>
          <ChevronDown
            size={16}
            className={`transition-transform ${dropdownOpen ? 'rotate-180' : ''}`}
          />
        </button>
        {dropdownOpen && (
          <>
            <div
              className="fixed inset-0 z-40"
              aria-hidden
              onClick={() => setDropdownOpen(false)}
            />
            <div
              className="absolute right-0 top-full z-50 mt-1 w-48 rounded-lg border border-[var(--app-border)] bg-[var(--app-bg-overlay)] py-1 shadow-modal"
              role="menu"
            >
              <Link
                to="/dashboard"
                onClick={() => setDropdownOpen(false)}
                className="flex items-center gap-2 px-3 py-2 text-sm text-[var(--app-fg)] hover:bg-[var(--app-bg-raised)]"
                role="menuitem"
              >
                <LayoutDashboard size={16} />
                Dashboard
              </Link>
              <Link
                to="/dashboard/settings/profile"
                onClick={() => setDropdownOpen(false)}
                className="flex items-center gap-2 px-3 py-2 text-sm text-[var(--app-fg)] hover:bg-[var(--app-bg-raised)]"
                role="menuitem"
              >
                <Settings size={16} />
                Settings
              </Link>
              <button
                type="button"
                onClick={() => {
                  setDropdownOpen(false)
                  clearAuth()
                  navigate({ to: '/auth/login' })
                }}
                className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm text-[var(--app-fg)] hover:bg-[var(--app-bg-raised)]"
                role="menuitem"
              >
                <LogOut size={16} />
                Logout
              </button>
            </div>
          </>
        )}
      </div>
    )
  }
  return (
    <div className="flex items-center gap-2">
      <Link
        to="/auth/login"
        className="rounded-lg px-3 py-2 text-sm font-medium text-[var(--app-fg-muted)] transition-[var(--motion-duration-fast)] hover:bg-[var(--app-bg-overlay)] hover:text-[var(--app-fg)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg-raised)]"
      >
        Sign in
      </Link>
      <Link
        to="/auth/register"
        className="rounded-lg bg-[var(--app-accent)] px-4 py-2 text-sm font-semibold text-[#0a0a0b] transition-[var(--motion-duration-fast)] hover:bg-[var(--app-accent-hover)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg-raised)]"
      >
        Register
      </Link>
    </div>
  )
}

export default function Header() {
  const [isOpen, setIsOpen] = useState(false)

  return (
    <>
      <header className="sticky top-0 z-30 flex items-center gap-4 border-b border-[var(--app-border)] bg-[var(--app-bg-raised)] px-4 py-3 shadow-card">
        <button
          onClick={() => setIsOpen(true)}
          className="p-2 rounded-lg text-[var(--app-fg-muted)] transition-[var(--motion-duration-fast)] hover:bg-[var(--app-bg-overlay)] hover:text-[var(--app-fg)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)]"
          aria-label="Open menu"
        >
          <Menu size={22} />
        </button>
        <Link
          to="/"
          className="text-h3 font-semibold text-[var(--app-fg)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 rounded"
        >
          Reelcut
        </Link>
        <div className="ml-auto flex items-center gap-2">
          <ThemeToggle />
          <HeaderAuth />
        </div>
      </header>

      {/* Mobile / overlay nav */}
      <div
        className={`fixed inset-0 z-50 bg-black/50 transition-opacity duration-200 ${isOpen ? 'opacity-100' : 'pointer-events-none opacity-0'}`}
        aria-hidden
        onClick={() => setIsOpen(false)}
      />
      <aside
        className={`fixed top-0 left-0 z-50 flex h-full w-72 flex-col border-r border-[var(--app-border)] bg-[var(--app-bg-raised)] shadow-modal transition-transform duration-200 ease-out ${isOpen ? 'translate-x-0' : '-translate-x-full'}`}
      >
        <div className="flex items-center justify-between border-b border-[var(--app-border)] p-4">
          <span className="text-h3 font-semibold text-[var(--app-fg)]">
            Menu
          </span>
          <button
            onClick={() => setIsOpen(false)}
            className="p-2 rounded-lg text-[var(--app-fg-muted)] hover:bg-[var(--app-bg-overlay)] hover:text-[var(--app-fg)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)]"
            aria-label="Close menu"
          >
            <X size={22} />
          </button>
        </div>
        <nav className="flex-1 overflow-y-auto p-3">
          <Link
            to="/"
            onClick={() => setIsOpen(false)}
            className="flex items-center gap-3 rounded-lg px-3 py-2.5 text-[var(--app-fg-muted)] transition-[var(--motion-duration-fast)] hover:bg-[var(--app-bg-overlay)] hover:text-[var(--app-fg)]"
          >
            <span className="font-medium">Home</span>
          </Link>
          <div className="flex items-center gap-2 border-t border-[var(--app-border)] pt-3 mt-3">
            <ThemeToggle />
            <HeaderAuth />
          </div>
        </nav>
      </aside>
    </>
  )
}
