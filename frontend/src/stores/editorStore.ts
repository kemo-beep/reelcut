import { create } from 'zustand'
import type { Clip, ClipStyle } from '../types'

interface EditorState {
  currentClip: Clip | null
  currentTime: number
  isPlaying: boolean
  zoom: number
  selectedLayer: string | null
  setClip: (clip: Clip | null) => void
  setCurrentTime: (time: number) => void
  play: () => void
  pause: () => void
  setZoom: (zoom: number) => void
  updateStyle: (style: Partial<ClipStyle>) => void
}

export const useEditorStore = create<EditorState>((set, get) => ({
  currentClip: null,
  currentTime: 0,
  isPlaying: false,
  zoom: 1,
  selectedLayer: null,
  setClip: (clip) => set({ currentClip: clip }),
  setCurrentTime: (time) => set({ currentTime: time }),
  play: () => set({ isPlaying: true }),
  pause: () => set({ isPlaying: false }),
  setZoom: (zoom) => set({ zoom }),
  updateStyle: (style) => {
    const { currentClip } = get()
    if (!currentClip) return
    set({
      currentClip: {
        ...currentClip,
        style: { ...(currentClip.style ?? {}), ...style } as ClipStyle,
      },
    })
  },
}))
