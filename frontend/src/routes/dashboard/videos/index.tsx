import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Video, Upload } from 'lucide-react'
import { listVideos } from '../../../lib/api/videos'
import { Skeleton } from '../../../components/ui/skeleton'
import { EmptyState } from '../../../components/ui/empty-state'
import { ErrorState } from '../../../components/ui/error-state'
import { Badge } from '../../../components/ui/badge'

export const Route = createFileRoute('/dashboard/videos/')({
  component: VideosListPage,
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

function VideosListPage() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['videos'],
    queryFn: () => listVideos({ per_page: 50 }),
    refetchInterval: (query) => {
      const videos = (query.state.data as { data?: { videos?: { status: string }[] } })?.data?.videos ?? []
      const hasProcessing = videos.some((v) => v.status === 'processing' || v.status === 'uploading')
      return hasProcessing ? 3000 : false
    },
  })

  if (isLoading) {
    return (
      <div className="space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <Skeleton className="mb-2 h-8 w-48" />
            <Skeleton className="h-4 w-64" />
          </div>
        </div>
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
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-8">
        <div>
          <h1 className="text-h1 text-[var(--app-fg)]">Videos</h1>
          <p className="text-caption mt-1">Your uploaded videos.</p>
        </div>
        <ErrorState message="Failed to load videos." onRetry={() => refetch()} />
      </div>
    )
  }

  const videos = data?.data?.videos ?? []

  return (
    <div className="space-y-8">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-h1 text-[var(--app-fg)]">Videos</h1>
          <p className="text-caption mt-1">Your uploaded videos.</p>
        </div>
        <Link
          to="/dashboard/videos/upload"
          className="inline-flex items-center justify-center gap-2 rounded-lg bg-[var(--app-accent)] px-4 py-2.5 text-sm font-semibold text-[#0a0a0b] shadow-card transition-[var(--motion-duration-fast)] hover:bg-[var(--app-accent-hover)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]"
        >
          <Upload size={18} />
          Upload video
        </Link>
      </div>

      {videos.length === 0 ? (
        <EmptyState
          icon={<Video size={28} />}
          title="No videos yet"
          description="Upload a video to get started with transcription and clips."
          action={
            <Link
              to="/dashboard/videos/upload"
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
  )
}
