import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router'
import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { getUploadUrl, confirmUpload } from '../../../lib/api/videos'
import { listProjects } from '../../../lib/api/projects'
import { useQuery } from '@tanstack/react-query'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import { Label } from '../../../components/ui/label'
import { toast } from 'sonner'
import { ApiError } from '../../../types'

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
  const [file, setFile] = useState<File | null>(null)
  const [dragOver, setDragOver] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [uploadProgress, setUploadProgress] = useState<'idle' | 'getting-url' | 'uploading' | 'confirming' | 'done' | 'error'>('idle')

  const { data: projectsData } = useQuery({
    queryKey: ['projects'],
    queryFn: () => listProjects({ per_page: 100 }),
  })
  const projects = projectsData?.data?.projects ?? []

  const getUrlMutation = useMutation({
    mutationFn: ({ projectId: pid, filename }: { projectId: string; filename: string }) =>
      getUploadUrl({ project_id: pid, filename }),
    onSuccess: async (urlData) => {
      if (!file) return
      setUploadProgress('uploading')
      setError(null)
      try {
        const url = urlData as {
          upload: { upload_url: string; method: string }
          video: { id: string }
        }
        const res = await fetch(url.upload.upload_url, {
          method: url.upload.method,
          body: file,
          headers: { 'Content-Type': file.type },
        })
        if (!res.ok) {
          setUploadProgress('error')
          setError('Upload to storage failed.')
          return
        }
        setUploadProgress('confirming')
        await confirmUpload(url.video.id)
        setUploadProgress('done')
        queryClient.invalidateQueries({ queryKey: ['videos'] })
        queryClient.invalidateQueries({ queryKey: ['projects'] })
        toast.success('Video uploaded')
        setTimeout(() => navigate({ to: '/dashboard/videos' }), 1500)
      } catch (err) {
        setUploadProgress('error')
        const msg = err instanceof ApiError ? err.message : 'Upload failed.'
        setError(msg)
        toast.error(msg)
      }
    },
    onError: (err) => {
      setUploadProgress('error')
      const msg = err instanceof ApiError ? err.message : 'Failed to get upload URL.'
      setError(msg)
      toast.error(msg)
    },
  })

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    if (!projectId.trim()) {
      setError('Select a project.')
      return
    }
    if (!file) {
      setError('Choose a file.')
      return
    }
    setUploadProgress('getting-url')
    getUrlMutation.mutate({ projectId: projectId.trim(), filename: file.name })
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault()
    setDragOver(false)
    const f = e.dataTransfer.files?.[0]
    if (f?.type.startsWith('video/')) setFile(f)
  }

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const f = e.target.files?.[0]
    if (f) setFile(f)
  }

  async function handleUseSampleVideo() {
    setError(null)
    try {
      const res = await fetch('/samples/_samplevideo.mp4')
      if (!res.ok) throw new Error('Sample not found')
      const blob = await res.blob()
      const sampleFile = new File([blob], '_samplevideo.mp4', { type: 'video/mp4' })
      setFile(sampleFile)
      toast.success('Sample video loaded. Select a project and click Upload.')
    } catch {
      setError('Could not load sample. Ensure public/samples/_samplevideo.mp4 exists.')
      toast.error('Sample video not found')
    }
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-h1 text-[var(--app-fg)]">Upload video</h1>
        <p className="text-caption mt-1">Add a video to a project to transcribe and create clips.</p>
      </div>
      <form onSubmit={handleSubmit} className="max-w-lg space-y-5">
        {error && (
          <div className="rounded-lg border border-[var(--app-destructive)]/30 bg-[var(--app-destructive-muted)] px-4 py-3 text-sm text-[var(--app-destructive)]">
            {error}
          </div>
        )}
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
        <div className="space-y-2">
          <Label className="text-label">Video file</Label>
          <div
            onDragOver={(e) => { e.preventDefault(); setDragOver(true) }}
            onDragLeave={() => setDragOver(false)}
            onDrop={handleDrop}
            className={`rounded-xl border-2 border-dashed p-8 text-center transition-[var(--motion-duration-fast)] ${
              dragOver
                ? 'border-[var(--app-accent)] bg-[var(--app-accent-muted)]'
                : 'border-[var(--app-border-strong)] bg-[var(--app-bg-raised)]'
            }`}
          >
            <input
              type="file"
              accept="video/*"
              onChange={handleFileChange}
              className="sr-only"
              id="video-file"
            />
            <label
              htmlFor="video-file"
              className="cursor-pointer text-[var(--app-fg-muted)] hover:text-[var(--app-fg)]"
            >
              {file ? file.name : 'Drop a video or click to browse'}
            </label>
          </div>
          <p className="text-caption text-[var(--app-fg-muted)]">
            Or{' '}
            <button
              type="button"
              onClick={handleUseSampleVideo}
              className="underline hover:text-[var(--app-fg)]"
            >
              use sample video (_samplevideo.mp4)
            </button>
            {' '}to test upload and re-upload.
          </p>
        </div>
        <Button
          type="submit"
          disabled={
            !projectId || !file ||
            ['getting-url', 'uploading', 'confirming'].includes(uploadProgress)
          }
          className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)] focus-visible:ring-[var(--app-accent)]"
        >
          {uploadProgress === 'getting-url' && 'Preparing…'}
          {uploadProgress === 'uploading' && 'Uploading…'}
          {uploadProgress === 'confirming' && 'Finalizing…'}
          {uploadProgress === 'done' && 'Done! Redirecting…'}
          {uploadProgress === 'error' && 'Retry'}
          {uploadProgress === 'idle' && 'Upload'}
        </Button>
      </form>
    </div>
  )
}
