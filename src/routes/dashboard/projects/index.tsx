import { createFileRoute } from '@tanstack/react-router'
import { ProjectsList } from '../../../components/dashboard/ProjectsList'

export const Route = createFileRoute('/dashboard/projects/')({
  component: ProjectsPage,
})

function ProjectsPage() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-h1 text-[var(--app-fg)]">Projects</h1>
        <p className="text-caption mt-1">Organize your videos and clips in projects.</p>
      </div>
      <ProjectsList />
    </div>
  )
}
