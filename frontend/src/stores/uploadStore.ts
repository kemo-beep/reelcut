import { create } from 'zustand'

export interface UploadProgress {
  fileId: string
  fileName: string
  status: 'pending' | 'uploading' | 'confirming' | 'done' | 'error'
  progress: number
  videoId: string | null
  error: string | null
}

interface UploadState {
  uploads: UploadProgress[]
  addUpload: (fileId: string, fileName: string) => void
  updateUpload: (
    fileId: string,
    update: Partial<Omit<UploadProgress, 'fileId' | 'fileName'>>
  ) => void
  removeUpload: (fileId: string) => void
  clearCompleted: () => void
}

export const useUploadStore = create<UploadState>((set, get) => ({
  uploads: [],
  addUpload: (fileId, fileName) =>
    set((s) => ({
      uploads: [
        ...s.uploads.filter((u) => u.fileId !== fileId),
        {
          fileId,
          fileName,
          status: 'pending',
          progress: 0,
          videoId: null,
          error: null,
        },
      ],
    })),
  updateUpload: (fileId, update) =>
    set((s) => ({
      uploads: s.uploads.map((u) =>
        u.fileId === fileId ? { ...u, ...update } : u
      ),
    })),
  removeUpload: (fileId) =>
    set((s) => ({ uploads: s.uploads.filter((u) => u.fileId !== fileId) })),
  clearCompleted: () =>
    set((s) => ({
      uploads: s.uploads.filter(
        (u) => u.status !== 'done' && u.status !== 'error'
      ),
    })),
}))
