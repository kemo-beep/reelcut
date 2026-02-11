import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getVideo, getPlaybackUrl, deleteVideo } from '../../../../lib/api/videos'
import { Button } from '../../../../components/ui/button'
import { Trash2 } from 'lucide-react'
import { toast } from 'sonner'

export const Route = createFileRoute('/dashboard/videos/$videoId/')({
  component: VideoDetailPage,
})

function VideoDetailPage() {
  const { videoId } = Route.useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { data, isLoading, error } = useQuery({
    queryKey: ['video', videoId],
    queryFn: () => getVideo(videoId),
  })
  const { data: playbackData, error: playbackError, isFetching: playbackLoading } = useQuery({
    queryKey: ['video-playback', videoId],
    queryFn: () => getPlaybackUrl(videoId),
    enabled: !!data?.video,
    retry: false,
  })
  const deleteMutation = useMutation({
    mutationFn: () => deleteVideo(videoId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['videos'] })
      toast.success('Video deleted')
      navigate({ to: '/dashboard/videos' })
    },
    onError: () => toast.error('Failed to delete video'),
  })

  if (isLoading) return <p className="text-slate-400">Loading...</p>
  if (error || !data?.video) return <p className="text-red-400">Video not found.</p>

  const video = data.video
  const playbackUrl = playbackData?.url

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <Link
          to="/dashboard/videos"
          className="text-cyan-400 hover:underline text-sm"
        >
          Back to videos
        </Link>
        <Button
          variant="outline"
          size="sm"
          className="text-red-400 border-red-400/50 hover:bg-red-500/10"
          onClick={() => deleteMutation.mutate()}
          disabled={deleteMutation.isPending}
        >
          <Trash2 size={16} className="mr-1" />
          Delete
        </Button>
      </div>
      <h1 className="text-2xl font-bold text-[var(--app-fg)]">{video.original_filename}</h1>
      <p className="text-caption text-[var(--app-fg-muted)]">
        Status: {video.status} · Duration: {video.duration_seconds != null ? `${Math.round(video.duration_seconds)}s` : '—'}
      </p>

      {playbackLoading && (
        <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg)] aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
          Loading playback…
        </div>
      )}
      {!playbackLoading && playbackUrl && (
        <div className="rounded-xl overflow-hidden border border-[var(--app-border)] bg-[var(--app-bg)]">
          <video
            src={playbackUrl}
            controls
            className="w-full max-h-[70vh]"
            preload="metadata"
          >
            Your browser does not support the video tag.
          </video>
        </div>
      )}
      {!playbackLoading && !playbackUrl && !playbackError && (
        <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg)] aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
          Video is still processing. Check back in a moment.
        </div>
      )}
      {!playbackLoading && playbackError && (
        <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg)] aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
          Unable to load playback.
        </div>
      )}

      <div className="flex gap-3">
        <Link
          to="/dashboard/clips"
          search={{ videoId: video.id }}
          className="inline-flex items-center justify-center rounded-lg bg-[var(--app-accent)] px-4 py-2 text-sm font-semibold text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]"
        >
          View clips
        </Link>
      </div>
    </div>
  )
}
