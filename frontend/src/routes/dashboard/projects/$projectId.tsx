import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { getProject } from '../../../lib/api/projects'

export const Route = createFileRoute('/dashboard/projects/$projectId')({
  component: ProjectDetailPage,
})

function ProjectDetailPage() {
  const { projectId } = Route.useParams()
  const { data, isLoading, error } = useQuery({
    queryKey: ['project', projectId],
    queryFn: () => getProject(projectId),
  })

  if (isLoading) return <p className="text-slate-400">Loading...</p>
  if (error || !data?.project) return <p className="text-red-400">Project not found.</p>

  const project = data.project
  return (
    <div className="space-y-6">
      <Link
        to="/dashboard/projects"
        className="text-cyan-400 hover:underline text-sm"
      >
        Back to projects
      </Link>
      <h1 className="text-2xl font-bold text-white">{project.name}</h1>
      {project.description && (
        <p className="text-slate-400">{project.description}</p>
      )}
      <Link
        to="/dashboard/videos/upload"
        search={{ projectId: project.id }}
        className="inline-block rounded-lg bg-cyan-600 px-4 py-2 text-white hover:bg-cyan-500"
      >
        Upload video to this project
      </Link>
    </div>
  )
}
