import { useState, useCallback } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { useAuthStore } from '../stores/authStore'
import { getUploadUrl, confirmUpload } from '../lib/api/videos'
import type { ApiError } from '../types'

export interface UseUploadOptions {
  projectId: string
  onSuccess?: (videoId: string) => void
  onError?: (error: Error | ApiError) => void
}

export interface UploadState {
  status: 'idle' | 'getting-url' | 'uploading' | 'confirming' | 'done' | 'error'
  progress: number
  error: string | null
  videoId: string | null
}

/**
 * Handles video upload: get presigned URL, upload file (with progress), confirm.
 * Progress is approximate (0 before upload, 90 during upload, 100 after confirm).
 */
export function useUpload(options: UseUploadOptions) {
  const { projectId, onSuccess, onError } = options
  const queryClient = useQueryClient()
  const [state, setState] = useState<UploadState>({
    status: 'idle',
    progress: 0,
    error: null,
    videoId: null,
  })
  const getAccessToken = useAuthStore((s) => s.getAccessToken)

  const upload = useCallback(
    async (file: File) => {
      setState({ status: 'getting-url', progress: 0, error: null, videoId: null })
      try {
        const urlData = await getUploadUrl({
          project_id: projectId,
          filename: file.name,
        })
        const data = urlData as {
          upload: { upload_url: string; method: string }
          video: { id: string }
        }
        setState((s) => ({ ...s, status: 'uploading', progress: 5 }))

        const xhr = new XMLHttpRequest()
        await new Promise<void>((resolve, reject) => {
          xhr.upload.addEventListener('progress', (e) => {
            if (e.lengthComputable) {
              const pct = 5 + Math.round((e.loaded / e.total) * 85)
              setState((s) => ({ ...s, progress: pct }))
            }
          })
          xhr.addEventListener('load', () => {
            if (xhr.status >= 200 && xhr.status < 300) resolve()
            else reject(new Error(`Upload failed: ${xhr.status}`))
          })
          xhr.addEventListener('error', () => reject(new Error('Upload failed')))
          xhr.open(data.upload.method, data.upload.upload_url)
          xhr.setRequestHeader('Content-Type', file.type)
          xhr.send(file)
        })

        setState((s) => ({ ...s, status: 'confirming', progress: 92 }))
        await confirmUpload(data.video.id)
        setState({
          status: 'done',
          progress: 100,
          error: null,
          videoId: data.video.id,
        })
        queryClient.invalidateQueries({ queryKey: ['video', data.video.id] })
        queryClient.invalidateQueries({ queryKey: ['transcription', data.video.id] })
        onSuccess?.(data.video.id)
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Upload failed'
        setState({
          status: 'error',
          progress: 0,
          error: message,
          videoId: null,
        })
        onError?.(err instanceof Error ? err : new Error(message))
      }
    },
    [projectId, onSuccess, onError, queryClient]
  )

  const reset = useCallback(() => {
    setState({ status: 'idle', progress: 0, error: null, videoId: null })
  }, [])

  return { ...state, upload, reset }
}
