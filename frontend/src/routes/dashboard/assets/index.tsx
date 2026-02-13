import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Video, Image } from 'lucide-react'
import { listVideos } from '../../../lib/api/videos'
import { listProjects } from '../../../lib/api/projects'
import { Skeleton } from '../../../components/ui/skeleton'
import { EmptyState } from '../../../components/ui/empty-state'
import { ErrorState } from '../../../components/ui/error-state'
import { Badge } from '../../../components/ui/badge'
import type { Project } from '../../../types'

export const Route = createFileRoute('/dashboard/assets/')({
    component: AssetsPage,
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

function AssetsPage() {
    const { data: projectsData } = useQuery({
        queryKey: ['projects'],
        queryFn: () => listProjects({ per_page: 100 }),
    })
    const projects = projectsData?.data?.projects ?? []
    const projectMap = projects.reduce<Record<string, Project>>((acc, p) => {
        acc[p.id] = p
        return acc
    }, {})

    const { data, isLoading, error, refetch } = useQuery({
        queryKey: ['videos', 'all-assets'],
        queryFn: () => listVideos({ per_page: 100 }),
        refetchInterval: (query) => {
            const videos = (query.state.data as { data?: { videos?: { status: string }[] } })?.data?.videos ?? []
            const hasProcessing = videos.some((v) => v.status === 'processing' || v.status === 'uploading')
            return hasProcessing ? 3000 : false
        },
    })

    if (isLoading) {
        return (
            <div className="space-y-8">
                <div>
                    <Skeleton className="mb-2 h-8 w-40" />
                    <Skeleton className="h-4 w-64" />
                </div>
                <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                    {Array.from({ length: 8 }).map((_, i) => (
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
                    <h1 className="text-h1 text-[var(--app-fg)]">Assets</h1>
                    <p className="text-caption mt-1">All your uploaded media across projects.</p>
                </div>
                <ErrorState message="Failed to load assets." onRetry={() => refetch()} />
            </div>
        )
    }

    const videos = data?.data?.videos ?? []

    return (
        <div className="space-y-8">
            <div>
                <h1 className="text-h1 text-[var(--app-fg)]">Assets</h1>
                <p className="text-caption mt-1">
                    All your uploaded media across every project. Click to view details and clips.
                </p>
            </div>

            {videos.length === 0 ? (
                <EmptyState
                    icon={<Image size={28} />}
                    title="No assets yet"
                    description="Upload videos to any project to see them here."
                />
            ) : (
                <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
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
                                        <Video size={36} aria-hidden />
                                    </div>
                                )}
                                <div className="p-3.5">
                                    <p className="text-sm font-medium text-[var(--app-fg)] truncate">
                                        {v.original_filename}
                                    </p>
                                    <div className="mt-2 flex flex-wrap items-center gap-1.5">
                                        <span className="text-caption text-xs">
                                            {v.duration_seconds != null
                                                ? `${Math.round(v.duration_seconds)}s`
                                                : '—'}
                                        </span>
                                        <Badge variant={statusVariant(v.status)}>{v.status}</Badge>
                                        {projectMap[v.project_id] && (
                                            <Badge variant="default">{projectMap[v.project_id].name}</Badge>
                                        )}
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
