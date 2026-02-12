import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getClip, renderClip, getClipDownloadUrl } from '../../../lib/api/clips'
import { Button } from '../../../components/ui/button'
import { Download, Film, Pencil } from 'lucide-react'
import { toast } from 'sonner'

export const Route = createFileRoute('/dashboard/clips/$clipId')({
  component: ClipDetailPage,
})

function ClipDetailPage() {
  const { clipId } = Route.useParams()
  const queryClient = useQueryClient()
  const { data, isLoading, error } = useQuery({
    queryKey: ['clip', clipId],
    queryFn: () => getClip(clipId),
  })
  const renderMutation = useMutation({
    mutationFn: () => renderClip(clipId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clip', clipId] })
      toast.success('Render started')
    },
    onError: () => toast.error('Failed to start render'),
  })

  if (isLoading) return <p className="text-slate-400">Loading...</p>
  if (error || !data?.clip) return <p className="text-red-400">Clip not found.</p>

  const clip = data.clip
  const canDownload = clip.status === 'ready' && clip.storage_path

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <Link
          to="/dashboard/clips"
          className="text-cyan-400 hover:underline text-sm"
        >
          Back to clips
        </Link>
        <div className="flex gap-2">
          <Link to="/editor/$clipId" params={{ clipId }}>
            <Button variant="outline" size="sm">
              <Pencil size={16} className="mr-1" />
              Edit
            </Button>
          </Link>
          {clip.status === 'draft' && (
            <Button
              size="sm"
              className="bg-[var(--app-accent)] text-[#0a0a0b]"
              onClick={() => renderMutation.mutate()}
              disabled={renderMutation.isPending}
            >
              <Film size={16} className="mr-1" />
              {renderMutation.isPending ? 'Starting…' : 'Render'}
            </Button>
          )}
          {clip.status === 'rendering' && (
            <span className="text-caption text-[var(--app-fg-muted)]">Rendering…</span>
          )}
          {canDownload && (
            <a
              href={getClipDownloadUrl(clipId)}
              download
              target="_blank"
              rel="noreferrer"
            >
              <Button variant="outline" size="sm">
                <Download size={16} className="mr-1" />
                Download
              </Button>
            </a>
          )}
        </div>
      </div>
      <h1 className="text-2xl font-bold text-[var(--app-fg)]">{clip.name}</h1>
      <p className="text-caption text-[var(--app-fg-muted)]">
        {clip.duration_seconds != null ? `${Math.round(clip.duration_seconds)}s` : '—'} · {clip.aspect_ratio} · {clip.status}
      </p>
    </div>
  )
}
