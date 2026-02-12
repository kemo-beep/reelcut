import { create } from 'zustand'
import type { Clip } from '../types'

interface ClipState {
  selectedClipId: string | null
  setSelectedClip: (clip: Clip | null) => void
}

export const useClipStore = create<ClipState>((set) => ({
  selectedClipId: null,
  setSelectedClip: (clip) =>
    set({ selectedClipId: clip?.id ?? null }),
}))
