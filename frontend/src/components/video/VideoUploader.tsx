import { useState, useCallback } from 'react'
import { Upload } from 'lucide-react'
import { Button } from '../ui/button'
import { ProgressBar } from '../processing/ProgressBar'
import { useUpload } from '../../hooks/useUpload'
import { cn } from '../../lib/utils'

export interface VideoUploaderProps {
  projectId: string
  onSuccess?: (videoId: string) => void
  onError?: (error: Error) => void
  className?: string
  accept?: string
  maxSizeBytes?: number
}

const DEFAULT_ACCEPT = 'video/mp4,video/webm,video/quicktime'
const DEFAULT_MAX_BYTES = 2 * 1024 * 1024 * 1024 // 2GB

export function VideoUploader({
  projectId,
  onSuccess,
  onError,
  className,
  accept = DEFAULT_ACCEPT,
  maxSizeBytes = DEFAULT_MAX_BYTES,
}: VideoUploaderProps) {
  const [dragOver, setDragOver] = useState(false)
  const [fileError, setFileError] = useState<string | null>(null)
  const { status, progress, error, upload, reset } = useUpload({
    projectId,
    onSuccess,
    onError,
  })

  const validate = useCallback(
    (file: File): string | null => {
      const types = accept.split(',').map((t) => t.trim())
      if (!types.some((t) => file.type === t || (t.endsWith('/*') && file.type.startsWith(t.slice(0, -1))))) {
        return 'Invalid file type'
      }
      if (file.size > maxSizeBytes) {
        return `File too large (max ${Math.round(maxSizeBytes / 1024 / 1024)}MB)`
      }
      return null
    },
    [accept, maxSizeBytes]
  )

  const handleFile = useCallback(
    (file: File) => {
      setFileError(null)
      const err = validate(file)
      if (err) {
        setFileError(err)
        return
      }
      upload(file)
    },
    [validate, upload]
  )

  const onDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault()
      setDragOver(false)
      const file = e.dataTransfer.files[0]
      if (file) handleFile(file)
    },
    [handleFile]
  )

  const onDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setDragOver(true)
  }, [])

  const onDragLeave = useCallback(() => {
    setDragOver(false)
  }, [])

  const onInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0]
      if (file) handleFile(file)
      e.target.value = ''
    },
    [handleFile]
  )

  const displayError = fileError ?? error

  return (
    <div className={cn('space-y-4', className)}>
      {(status === 'idle' || status === 'done' || status === 'error') && (
        <label
          htmlFor="video-upload-input"
          onDrop={onDrop}
          onDragOver={onDragOver}
          onDragLeave={onDragLeave}
          className={cn(
            'flex flex-col items-center justify-center rounded-xl border-2 border-dashed p-8 transition-colors cursor-pointer',
            dragOver
              ? 'border-[var(--app-accent)] bg-[var(--app-accent-muted)]'
              : 'border-[var(--app-border)] bg-[var(--app-bg-raised)]'
          )}
        >
          <input
            type="file"
            accept={accept}
            onChange={onInputChange}
            className="sr-only"
            id="video-upload-input"
          />
          <Upload size={48} className="text-[var(--app-fg-muted)] mb-4" />
          <p className="text-[var(--app-fg)] font-medium">Drop a video here or click to browse</p>
          <p className="text-caption text-[var(--app-fg-muted)] mt-1">
            MP4, WebM, MOV · max {Math.round(maxSizeBytes / 1024 / 1024)}MB
          </p>
          <Button type="button" variant="outline" className="mt-4" tabIndex={-1}>
            Select file
          </Button>
        </label>
      )}

      {(status === 'getting-url' || status === 'uploading' || status === 'confirming') && (
        <div className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-6">
          <p className="text-sm font-medium text-[var(--app-fg)]">
            {status === 'getting-url' && 'Preparing…'}
            {status === 'uploading' && 'Uploading…'}
            {status === 'confirming' && 'Finalizing…'}
          </p>
          <ProgressBar value={progress} className="mt-2" showLabel />
        </div>
      )}

      {status === 'done' && (
        <div className="rounded-xl border border-[var(--app-success)]/30 bg-[var(--app-success-muted)] p-4">
          <p className="text-sm font-medium text-[var(--app-fg)]">Upload complete</p>
          <Button variant="outline" size="sm" onClick={reset} className="mt-2">
            Upload another
          </Button>
        </div>
      )}

      {displayError && (
        <div className="rounded-lg border border-[var(--app-destructive)]/30 bg-[var(--app-destructive-muted)] px-4 py-3 text-sm text-[var(--app-destructive)]">
          {displayError}
        </div>
      )}
    </div>
  )
}
