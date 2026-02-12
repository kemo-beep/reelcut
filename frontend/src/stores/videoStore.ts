import { create } from 'zustand'
import type { Video } from '../types'

interface VideoState {
  selectedVideoId: string | null
  setSelectedVideo: (video: Video | null) => void
}

export const useVideoStore = create<VideoState>((set) => ({
  selectedVideoId: null,
  setSelectedVideo: (video) =>
    set({ selectedVideoId: video?.id ?? null }),
}))
