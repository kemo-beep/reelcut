import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { getClip } from '../../../lib/api/clips'

export const Route = createFileRoute('/dashboard/clips/$clipId')({
  component: ClipDetailPage,
})

function ClipDetailPage() {
  const { clipId } = Route.useParams()
  const { data, isLoading, error } = useQuery({
    queryKey: ['clip', clipId],
    queryFn: () => getClip(clipId),
  })

  if (isLoading) return <p className="text-slate-400">Loading...</p>
  if (error || !data?.clip) return <p className="text-red-400">Clip not found.</p>

  const clip = data.clip
  return (
    <div className="space-y-6">
      <Link
        to="/dashboard/clips"
        className="text-cyan-400 hover:underline text-sm"
      >
        Back to clips
      </Link>
      <h1 className="text-2xl font-bold text-white">{clip.name}</h1>
      <p className="text-slate-400">
        {clip.duration_seconds != null ? `${Math.round(clip.duration_seconds)}s` : '—'} · {clip.aspect_ratio} · {clip.status}
      </p>
    </div>
  )
}
