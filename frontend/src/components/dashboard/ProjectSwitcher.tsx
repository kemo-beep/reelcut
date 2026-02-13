import { useState, useRef, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
    listProjects,
    createProject,
    updateProject,
    deleteProject,
} from '../../lib/api/projects'
import { useActiveProject } from '../../stores/useActiveProject'
import { Button } from '../ui/button'
import { Input } from '../ui/input'
import { Label } from '../ui/label'
import { Modal } from '../ui/modal'
import { toast } from 'sonner'
import { ApiError } from '../../types'
import type { Project } from '../../types'
import {
    ChevronDown,
    FolderOpen,
    Plus,
    Pencil,
    Trash2,
    Check,
    MoreHorizontal,
    Layers,
} from 'lucide-react'

export function ProjectSwitcher() {
    const queryClient = useQueryClient()
    const { activeProjectId, setActiveProject } = useActiveProject()
    const [dropdownOpen, setDropdownOpen] = useState(false)
    const [createOpen, setCreateOpen] = useState(false)
    const [editProject, setEditProject] = useState<Project | null>(null)
    const [deleteTarget, setDeleteTarget] = useState<Project | null>(null)
    const [contextMenu, setContextMenu] = useState<string | null>(null)
    const dropdownRef = useRef<HTMLDivElement>(null)

    const { data, isLoading } = useQuery({
        queryKey: ['projects'],
        queryFn: () => listProjects({ per_page: 100 }),
    })
    const projects = data?.data?.projects ?? []

    // Auto-select first project if none selected
    useEffect(() => {
        if (!activeProjectId && projects.length > 0) {
            setActiveProject(projects[0].id)
        }
    }, [activeProjectId, projects, setActiveProject])

    const activeProject = projects.find((p) => p.id === activeProjectId)

    // Close dropdown on click outside
    useEffect(() => {
        if (!dropdownOpen) return
        const handler = (e: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
                setDropdownOpen(false)
                setContextMenu(null)
            }
        }
        document.addEventListener('mousedown', handler)
        return () => document.removeEventListener('mousedown', handler)
    }, [dropdownOpen])

    function handleSelect(id: string | null) {
        setActiveProject(id)
        setDropdownOpen(false)
        setContextMenu(null)
    }

    return (
        <div className="relative" ref={dropdownRef}>
            {/* Switcher trigger */}
            <button
                type="button"
                onClick={() => {
                    setDropdownOpen((o) => !o)
                    setContextMenu(null)
                }}
                className="flex w-full items-center gap-2.5 rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-2.5 text-left text-sm transition-all hover:border-[var(--app-border-strong)] hover:bg-[var(--app-bg-overlay)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)]"
                aria-expanded={dropdownOpen}
                aria-haspopup="listbox"
            >
                <div className="flex h-7 w-7 items-center justify-center rounded-md bg-[var(--app-accent-muted)] text-[var(--app-accent)]">
                    {activeProject ? (
                        <FolderOpen size={14} />
                    ) : (
                        <Layers size={14} />
                    )}
                </div>
                <div className="min-w-0 flex-1">
                    <p className="truncate font-medium text-[var(--app-fg)] text-sm leading-tight">
                        {isLoading
                            ? 'Loading…'
                            : activeProject
                                ? activeProject.name
                                : 'All Projects'}
                    </p>
                    <p className="text-[11px] leading-tight text-[var(--app-fg-subtle)]">
                        Workspace
                    </p>
                </div>
                <ChevronDown
                    size={16}
                    className={`flex-shrink-0 text-[var(--app-fg-muted)] transition-transform ${dropdownOpen ? 'rotate-180' : ''}`}
                />
            </button>

            {/* Dropdown */}
            {dropdownOpen && (
                <div className="absolute left-0 right-0 top-full z-50 mt-1.5 overflow-hidden rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] shadow-modal">
                    {/* All projects option */}
                    <div className="border-b border-[var(--app-border)] p-1.5">
                        <button
                            type="button"
                            onClick={() => handleSelect(null)}
                            className={`flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm transition-colors ${!activeProjectId
                                    ? 'bg-[var(--app-accent-muted)] text-[var(--app-accent)]'
                                    : 'text-[var(--app-fg-muted)] hover:bg-[var(--app-bg-overlay)] hover:text-[var(--app-fg)]'
                                }`}
                        >
                            <Layers size={15} />
                            <span className="flex-1 text-left font-medium">All Projects</span>
                            {!activeProjectId && <Check size={15} />}
                        </button>
                    </div>

                    {/* Project list */}
                    <div className="max-h-52 overflow-y-auto p-1.5">
                        {projects.length === 0 && !isLoading && (
                            <p className="px-2.5 py-3 text-center text-xs text-[var(--app-fg-subtle)]">
                                No projects yet
                            </p>
                        )}
                        {projects.map((p) => (
                            <div key={p.id} className="group relative">
                                <button
                                    type="button"
                                    onClick={() => handleSelect(p.id)}
                                    className={`flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm transition-colors ${activeProjectId === p.id
                                            ? 'bg-[var(--app-accent-muted)] text-[var(--app-accent)]'
                                            : 'text-[var(--app-fg)] hover:bg-[var(--app-bg-overlay)]'
                                        }`}
                                >
                                    <FolderOpen size={15} className="flex-shrink-0" />
                                    <span className="flex-1 truncate text-left font-medium">{p.name}</span>
                                    {activeProjectId === p.id && <Check size={15} className="flex-shrink-0" />}
                                </button>
                                {/* Context menu trigger */}
                                <button
                                    type="button"
                                    onClick={(e) => {
                                        e.stopPropagation()
                                        setContextMenu(contextMenu === p.id ? null : p.id)
                                    }}
                                    className="absolute right-1.5 top-1/2 -translate-y-1/2 rounded-md p-1 text-[var(--app-fg-subtle)] opacity-0 transition-opacity hover:bg-[var(--app-bg-overlay)] hover:text-[var(--app-fg)] group-hover:opacity-100"
                                    aria-label={`Options for ${p.name}`}
                                >
                                    <MoreHorizontal size={14} />
                                </button>
                                {/* Context menu */}
                                {contextMenu === p.id && (
                                    <div className="absolute right-0 top-full z-10 mt-0.5 w-36 overflow-hidden rounded-lg border border-[var(--app-border)] bg-[var(--app-bg-overlay)] py-1 shadow-modal">
                                        <button
                                            type="button"
                                            onClick={(e) => {
                                                e.stopPropagation()
                                                setEditProject(p)
                                                setContextMenu(null)
                                                setDropdownOpen(false)
                                            }}
                                            className="flex w-full items-center gap-2 px-3 py-1.5 text-xs text-[var(--app-fg)] hover:bg-[var(--app-bg-raised)]"
                                        >
                                            <Pencil size={13} /> Edit
                                        </button>
                                        <button
                                            type="button"
                                            onClick={(e) => {
                                                e.stopPropagation()
                                                setDeleteTarget(p)
                                                setContextMenu(null)
                                                setDropdownOpen(false)
                                            }}
                                            className="flex w-full items-center gap-2 px-3 py-1.5 text-xs text-[var(--app-destructive)] hover:bg-[var(--app-bg-raised)]"
                                        >
                                            <Trash2 size={13} /> Delete
                                        </button>
                                    </div>
                                )}
                            </div>
                        ))}
                    </div>

                    {/* Create new project button */}
                    <div className="border-t border-[var(--app-border)] p-1.5">
                        <button
                            type="button"
                            onClick={() => {
                                setCreateOpen(true)
                                setDropdownOpen(false)
                                setContextMenu(null)
                            }}
                            className="flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm font-medium text-[var(--app-fg-muted)] transition-colors hover:bg-[var(--app-bg-overlay)] hover:text-[var(--app-fg)]"
                        >
                            <Plus size={15} />
                            New project
                        </button>
                    </div>
                </div>
            )}

            {/* Create project modal */}
            <CreateProjectModal
                open={createOpen}
                onClose={() => setCreateOpen(false)}
                onCreated={(project) => {
                    setActiveProject(project.id)
                }}
            />

            {/* Edit project modal */}
            {editProject && (
                <EditProjectModal
                    project={editProject}
                    open={!!editProject}
                    onClose={() => setEditProject(null)}
                />
            )}

            {/* Delete confirmation modal */}
            {deleteTarget && (
                <DeleteProjectModal
                    project={deleteTarget}
                    open={!!deleteTarget}
                    onClose={() => setDeleteTarget(null)}
                    onDeleted={() => {
                        if (activeProjectId === deleteTarget.id) {
                            const remaining = projects.filter((p) => p.id !== deleteTarget.id)
                            setActiveProject(remaining.length > 0 ? remaining[0].id : null)
                        }
                    }}
                />
            )}
        </div>
    )
}

/* ─── Create Project Modal ────────────────────────────────────── */

function CreateProjectModal({
    open,
    onClose,
    onCreated,
}: {
    open: boolean
    onClose: () => void
    onCreated: (project: Project) => void
}) {
    const queryClient = useQueryClient()
    const [name, setName] = useState('')
    const [description, setDescription] = useState('')
    const [error, setError] = useState<string | null>(null)

    const mutation = useMutation({
        mutationFn: (input: { name: string; description?: string }) =>
            createProject(input),
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ['projects'] })
            toast.success('Project created')
            onCreated(data.project)
            handleClose()
        },
        onError: (err) => {
            const msg = err instanceof ApiError ? err.message : 'Failed to create project'
            setError(msg)
            toast.error(msg)
        },
    })

    function handleClose() {
        onClose()
        setName('')
        setDescription('')
        setError(null)
    }

    function handleSubmit(e: React.FormEvent) {
        e.preventDefault()
        setError(null)
        mutation.mutate({ name, description: description || undefined })
    }

    return (
        <Modal open={open} onClose={handleClose} title="Create project">
            <form onSubmit={handleSubmit} className="space-y-4">
                {error && (
                    <div className="rounded-lg border border-[var(--app-destructive)]/30 bg-[var(--app-destructive-muted)] px-3 py-2 text-sm text-[var(--app-destructive)]">
                        {error}
                    </div>
                )}
                <div className="space-y-2">
                    <Label className="text-label">Name</Label>
                    <Input
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                        autoFocus
                        placeholder="My project"
                        className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)] focus-visible:ring-[var(--app-accent)]"
                    />
                </div>
                <div className="space-y-2">
                    <Label className="text-label">Description (optional)</Label>
                    <Input
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        placeholder="A brief description"
                        className="h-11 border-[var(--app-border-strong)] bg-[var(--app-bg)] text-[var(--app-fg)] focus-visible:ring-[var(--app-accent)]"
                    />
                </div>
                <div className="flex gap-3 pt-2">
                    <Button
                        type="submit"
                        disabled={mutation.isPending}
                        className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)] focus-visible:ring-[var(--app-accent)]"
                    >
                        {mutation.isPending ? 'Creating…' : 'Create'}
                    </Button>
                    <Button
                        type="button"
                        variant="outline"
                        onClick={handleClose}
                        className="border-[var(--app-border-strong)] focus-visible:ring-[var(--app-accent)]"
                    >
                        Cancel
                    </Button>
                </div>
            </form>
        </Modal>
    )
}

/* ─── Edit Project Modal ──────────────────────────────────────── */

function EditProjectModal({
    project,
    open,
    onClose,
}: {
    project: Project
    open: boolean
    onClose: () => void
}) {
    const queryClient = useQueryClient()
    const [name, setName] = useState(project.name)
    const [description, setDescription] = useState(project.description ?? '')
    const [error, setError] = useState<string | null>(null)

    // Sync when project changes
    useEffect(() => {
        setName(project.name)
        setDescription(project.description ?? '')
    }, [project])

    const mutation = useMutation({
        mutationFn: (input: { name?: string; description?: string | null }) =>
            updateProject(project.id, input),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['projects'] })
            toast.success('Project updated')
            onClose()
        },
        onError: (err) => {
            const msg = err instanceof ApiError ? err.message : 'Failed to update project'
            setError(msg)
            toast.error(msg)
        },
    })

    function handleSubmit(e: React.FormEvent) {
        e.preventDefault()
        setError(null)
        mutation.mutate({ name, description: description || null })
    }

    return (
        <Modal open={open} onClose={onClose} title="Edit project">
            <form onSubmit={handleSubmit} className="space-y-4">
                {error && (
                    <div className="rounded-lg border border-[var(--app-destructive)]/30 bg-[var(--app-destructive-muted)] px-3 py-2 text-sm text-[var(--app-destructive)]">
                        {error}
                    </div>
                )}
                <div className="space-y-2">
                    <Label className="text-label">Name</Label>
                    <Input
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                        autoFocus
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
                        disabled={mutation.isPending}
                        className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)] focus-visible:ring-[var(--app-accent)]"
                    >
                        {mutation.isPending ? 'Saving…' : 'Save'}
                    </Button>
                    <Button
                        type="button"
                        variant="outline"
                        onClick={onClose}
                        className="border-[var(--app-border-strong)] focus-visible:ring-[var(--app-accent)]"
                    >
                        Cancel
                    </Button>
                </div>
            </form>
        </Modal>
    )
}

/* ─── Delete Project Modal ────────────────────────────────────── */

function DeleteProjectModal({
    project,
    open,
    onClose,
    onDeleted,
}: {
    project: Project
    open: boolean
    onClose: () => void
    onDeleted: () => void
}) {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: () => deleteProject(project.id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['projects'] })
            toast.success('Project deleted')
            onDeleted()
            onClose()
        },
        onError: (err) => {
            const msg = err instanceof ApiError ? err.message : 'Failed to delete project'
            toast.error(msg)
        },
    })

    return (
        <Modal open={open} onClose={onClose} title="Delete project">
            <div className="space-y-4">
                <p className="text-sm text-[var(--app-fg-muted)]">
                    Are you sure you want to delete <strong className="text-[var(--app-fg)]">{project.name}</strong>?
                    This action cannot be undone. All videos and clips in this project will be permanently deleted.
                </p>
                <div className="flex gap-3 pt-2">
                    <Button
                        type="button"
                        onClick={() => mutation.mutate()}
                        disabled={mutation.isPending}
                        className="bg-[var(--app-destructive)] text-white hover:bg-[var(--app-destructive)]/90 focus-visible:ring-[var(--app-destructive)]"
                    >
                        {mutation.isPending ? 'Deleting…' : 'Delete'}
                    </Button>
                    <Button
                        type="button"
                        variant="outline"
                        onClick={onClose}
                        className="border-[var(--app-border-strong)] focus-visible:ring-[var(--app-accent)]"
                    >
                        Cancel
                    </Button>
                </div>
            </div>
        </Modal>
    )
}
