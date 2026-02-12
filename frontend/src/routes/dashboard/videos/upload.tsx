import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router'
import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { listProjects } from '../../../lib/api/projects'
import { useQuery } from '@tanstack/react-query'
import { Label } from '../../../components/ui/label'
import { VideoUploader } from '../../../components/video/VideoUploader'
import { toast } from 'sonner'

export const Route = createFileRoute('/dashboard/videos/upload')({
  validateSearch: (s): { projectId?: string } => ({
    projectId: typeof s?.projectId === 'string' ? s.projectId : undefined,
  }),
  component: VideoUploadPage,
})

function VideoUploadPage() {
  const navigate = useNavigate()
  const search = useSearch({ from: '/dashboard/videos/upload' })
  const queryClient = useQueryClient()
  const [projectId, setProjectId] = useState(search.projectId ?? '')

  const { data: projectsData } = useQuery({
    queryKey: ['projects'],
    queryFn: () => listProjects({ per_page: 100 }),
  })
  const projects = projectsData?.data?.projects ?? []

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-h1 text-[var(--app-fg)]">Upload video</h1>
        <p className="text-caption mt-1">Add a video to a project to transcribe and create clips.</p>
      </div>
      <div className="max-w-lg space-y-5">
        <div className="space-y-2">
          <Label className="text-label">Project</Label>
          <select
            value={projectId}
            onChange={(e) => setProjectId(e.target.value)}
            required
            className="h-11 w-full rounded-lg border border-[var(--app-border-strong)] bg-[var(--app-bg)] px-3 py-2 text-[var(--app-fg)] focus:outline-none focus:ring-2 focus:ring-[var(--app-accent)]"
          >
            <option value="">Select project</option>
            {projects.map((p) => (
              <option key={p.id} value={p.id}>
                {p.name}
              </option>
            ))}
          </select>
        </div>
        {projectId ? (
          <VideoUploader
            projectId={projectId}
            onSuccess={() => {
              queryClient.invalidateQueries({ queryKey: ['videos'] })
              queryClient.invalidateQueries({ queryKey: ['projects'] })
              toast.success('Video uploaded')
              setTimeout(() => navigate({ to: '/dashboard/videos' }), 1500)
            }}
            onError={() => {}}
          />
        ) : (
          <p className="text-caption text-[var(--app-fg-muted)]">Select a project to upload a video.</p>
        )}
      </div>
    </div>
  )
}
