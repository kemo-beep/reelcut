import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/dashboard/')({
  component: DashboardHome,
})

function DashboardHome() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-h1 text-[var(--app-fg)]">Dashboard</h1>
        <p className="text-caption mt-1">
          Welcome to Reelcut. Create projects, upload videos, and generate short clips.
        </p>
      </div>
      <div className="flex flex-wrap gap-4">
        <Link
          to="/dashboard/projects"
          className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] px-5 py-3 font-medium text-[var(--app-fg)] shadow-card transition-[var(--motion-duration-fast)] hover:border-[var(--app-border-strong)] hover:shadow-lg focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]"
        >
          View projects
        </Link>
        <Link
          to="/dashboard/videos"
          className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] px-5 py-3 font-medium text-[var(--app-fg)] shadow-card transition-[var(--motion-duration-fast)] hover:border-[var(--app-border-strong)] hover:shadow-lg focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]"
        >
          View videos
        </Link>
        <Link
          to="/dashboard/videos/upload"
          className="rounded-xl bg-[var(--app-accent)] px-5 py-3 font-semibold text-[#0a0a0b] shadow-card transition-[var(--motion-duration-fast)] hover:bg-[var(--app-accent-hover)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]"
        >
          Upload video
        </Link>
      </div>
    </div>
  )
}
