import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { FolderPlus } from 'lucide-react'
import { listProjects, createProject } from '../../lib/api/projects'
import { Button } from '../ui/button'
import { Input } from '../ui/input'
import { Label } from '../ui/label'
import { Modal } from '../ui/modal'
import { Skeleton } from '../ui/skeleton'
import { EmptyState } from '../ui/empty-state'
import { ErrorState } from '../ui/error-state'
import { toast } from 'sonner'
import { ApiError } from '../../types'

export function ProjectsList() {
  const queryClient = useQueryClient()
  const [modalOpen, setModalOpen] = useState(false)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [createError, setCreateError] = useState<string | null>(null)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['projects'],
    queryFn: async () => {
      const res = await listProjects({ per_page: 50 })
      return res
    },
  })

  const createMutation = useMutation({
    mutationFn: (input: { name: string; description?: string }) =>
      createProject(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      setModalOpen(false)
      setName('')
      setDescription('')
      setCreateError(null)
      toast.success('Project created')
    },
    onError: (err) => {
      const msg = err instanceof ApiError ? err.message : 'Failed to create project'
      setCreateError(msg)
      toast.error(msg)
    },
  })

  function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    setCreateError(null)
    createMutation.mutate({ name, description: description || undefined })
  }

  if (isLoading) {
    return (
      <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <li key={i} className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-5 shadow-card">
            <Skeleton className="mb-3 h-5 w-3/4" />
            <Skeleton className="h-4 w-full" />
          </li>
        ))}
      </ul>
    )
  }

  if (error) {
    return (
      <ErrorState
        message="Failed to load projects."
        onRetry={() => refetch()}
      />
    )
  }

  const projects = data?.data?.projects ?? []

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <Button
          onClick={() => setModalOpen(true)}
          className="bg-[var(--app-accent)] font-semibold text-[#0a0a0b] hover:bg-[var(--app-accent-hover)] focus-visible:ring-[var(--app-accent)]"
        >
          New project
        </Button>
      </div>

      <Modal
        open={modalOpen}
        onClose={() => {
          setModalOpen(false)
          setCreateError(null)
        }}
        title="Create project"
      >
        <form onSubmit={handleCreate} className="space-y-4">
          {createError && (
            <div className="rounded-lg border border-[var(--app-destructive)]/30 bg-[var(--app-destructive-muted)] px-3 py-2 text-sm text-[var(--app-destructive)]">
              {createError}
            </div>
          )}
          <div className="space-y-2">
            <Label className="text-label">Name</Label>
            <Input
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)] focus-visible:ring-[var(--app-accent)]"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-label">Description (optional)</Label>
            <Input
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)] focus-visible:ring-[var(--app-accent)]"
            />
          </div>
          <div className="flex gap-3 pt-2">
            <Button
              type="submit"
              disabled={createMutation.isPending}
              className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)] focus-visible:ring-[var(--app-accent)]"
            >
              {createMutation.isPending ? 'Creating…' : 'Create'}
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => setModalOpen(false)}
              className="border-[var(--app-border-strong)] focus-visible:ring-[var(--app-accent)]"
            >
              Cancel
            </Button>
          </div>
        </form>
      </Modal>

      {projects.length === 0 ? (
        <EmptyState
          icon={<FolderPlus size={28} />}
          title="No projects yet"
          description="Create a project to organize your videos and clips."
          action={
            <Button
              onClick={() => setModalOpen(true)}
              className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)] focus-visible:ring-[var(--app-accent)]"
            >
              New project
            </Button>
          }
        />
      ) : (
        <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {projects.map((p) => (
            <li key={p.id}>
              <Link
                to="/dashboard/projects/$projectId"
                params={{ projectId: p.id }}
                className="block rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-5 shadow-card transition-[var(--motion-duration-fast)] hover:border-[var(--app-border-strong)] hover:shadow-lg focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)]"
              >
                <span className="font-medium text-[var(--app-fg)]">{p.name}</span>
                {p.description && (
                  <p className="text-caption mt-1">{p.description}</p>
                )}
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
