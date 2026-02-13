import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Video, Upload, Sparkles, FolderOpen } from 'lucide-react'
import { listVideos } from '../../lib/api/videos'
import { useActiveProject } from '../../stores/useActiveProject'
import { Skeleton } from '../../components/ui/skeleton'
import { EmptyState } from '../../components/ui/empty-state'
import { ErrorState } from '../../components/ui/error-state'
import { Badge } from '../../components/ui/badge'
import { listProjects } from '../../lib/api/projects'

export const Route = createFileRoute('/dashboard/')({
  component: DashboardHome,
})

function statusVariant(
  status: string
): 'default' | 'success' | 'warning' | 'destructive' {
  switch (status) {
    case 'ready':
      return 'success'
    case 'processing':
    case 'uploading':
      return 'warning'
    case 'failed':
      return 'destructive'
    default:
      return 'default'
  }
}

function DashboardHome() {
  const { activeProjectId } = useActiveProject()

  const { data: projectsData } = useQuery({
    queryKey: ['projects'],
    queryFn: () => listProjects({ per_page: 100 }),
  })
  const projects = projectsData?.data?.projects ?? []
  const activeProject = projects.find((p) => p.id === activeProjectId)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['videos', { project_id: activeProjectId }],
    queryFn: () =>
      listVideos({
        per_page: 12,
        ...(activeProjectId ? { project_id: activeProjectId } : {}),
      }),
  })
  const videos = data?.data?.videos ?? []

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-h1 text-[var(--app-fg)]">
            {activeProject ? activeProject.name : 'All Projects'}
          </h1>
          <p className="text-caption mt-1">
            {activeProject
              ? activeProject.description || 'Manage your videos and clips.'
              : 'Showing content across all your projects.'}
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Link
            to="/dashboard/videos/upload"
            search={activeProjectId ? { projectId: activeProjectId } : {}}
            className="inline-flex items-center justify-center gap-2 rounded-lg bg-[var(--app-accent)] px-4 py-2.5 text-sm font-semibold text-[#0a0a0b] shadow-card transition-[var(--motion-duration-fast)] hover:bg-[var(--app-accent-hover)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]"
          >
            <Upload size={18} />
            Upload video
          </Link>
        </div>
      </div>

      {/* Quick stats */}
      <div className="grid gap-4 sm:grid-cols-3">
        <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-5 shadow-card">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-[var(--app-accent-muted)] text-[var(--app-accent)]">
              <Video size={20} />
            </div>
            <div>
              <p className="text-2xl font-bold text-[var(--app-fg)]">{isLoading ? '—' : videos.length}</p>
              <p className="text-caption">Videos</p>
            </div>
          </div>
        </div>
        <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-5 shadow-card">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-[var(--app-accent-muted)] text-[var(--app-accent)]">
              <Sparkles size={20} />
            </div>
            <div>
              <p className="text-2xl font-bold text-[var(--app-fg)]">
                {isLoading ? '—' : videos.filter((v) => v.status === 'ready').length}
              </p>
              <p className="text-caption">Ready</p>
            </div>
          </div>
        </div>
        <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-5 shadow-card">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-[var(--app-accent-muted)] text-[var(--app-accent)]">
              <FolderOpen size={20} />
            </div>
            <div>
              <p className="text-2xl font-bold text-[var(--app-fg)]">{projects.length}</p>
              <p className="text-caption">Projects</p>
            </div>
          </div>
        </div>
      </div>

      {/* Recent videos */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-h3 text-[var(--app-fg)]">Recent Videos</h2>
          <Link
            to="/dashboard/videos"
            className="text-sm font-medium text-[var(--app-accent)] hover:underline"
          >
            View all →
          </Link>
        </div>

        {isLoading ? (
          <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <li key={i} className="overflow-hidden rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] shadow-card">
                <Skeleton className="aspect-video w-full" />
                <div className="p-4">
                  <Skeleton className="mb-2 h-5 w-full" />
                  <Skeleton className="h-4 w-2/3" />
                </div>
              </li>
            ))}
          </ul>
        ) : error ? (
          <ErrorState message="Failed to load videos." onRetry={() => refetch()} />
        ) : videos.length === 0 ? (
          <EmptyState
            icon={<Video size={28} />}
            title="No videos yet"
            description={
              activeProject
                ? `Upload a video to "${activeProject.name}" to get started.`
                : 'Upload a video to get started with transcription and clips.'
            }
            action={
              <Link
                to="/dashboard/videos/upload"
                search={activeProjectId ? { projectId: activeProjectId } : {}}
                className="inline-flex items-center justify-center gap-2 rounded-lg bg-[var(--app-accent)] px-4 py-2.5 text-sm font-semibold text-[#0a0a0b] hover:bg-[var(--app-accent-hover)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)]"
              >
                <Upload size={18} />
                Upload video
              </Link>
            }
          />
        ) : (
          <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {videos.map((v) => (
              <li key={v.id}>
                <Link
                  to="/dashboard/videos/$videoId"
                  params={{ videoId: v.id }}
                  className="block overflow-hidden rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] shadow-card transition-[var(--motion-duration-fast)] hover:border-[var(--app-border-strong)] hover:shadow-lg focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]"
                >
                  {v.thumbnail_display_url && v.thumbnail_display_url.startsWith('http') ? (
                    <img
                      src={v.thumbnail_display_url}
                      alt=""
                      className="aspect-video w-full object-cover bg-[var(--app-bg)]"
                    />
                  ) : (
                    <div className="flex aspect-video w-full items-center justify-center bg-[var(--app-bg)] text-[var(--app-fg-subtle)]">
                      <Video size={40} aria-hidden />
                    </div>
                  )}
                  <div className="p-4">
                    <p className="font-medium text-[var(--app-fg)] truncate">
                      {v.original_filename}
                    </p>
                    <div className="mt-2 flex items-center gap-2">
                      <span className="text-caption">
                        {v.duration_seconds != null
                          ? `${Math.round(v.duration_seconds)}s`
                          : '—'}
                      </span>
                      <Badge variant={statusVariant(v.status)}>{v.status}</Badge>
                    </div>
                  </div>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}
