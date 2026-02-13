import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface ActiveProjectState {
  activeProjectId: string | null
  setActiveProject: (id: string | null) => void
  clearActiveProject: () => void
}

export const useActiveProject = create<ActiveProjectState>()(
  persist(
    (set) => ({
      activeProjectId: null,
      setActiveProject: (id) => set({ activeProjectId: id }),
      clearActiveProject: () => set({ activeProjectId: null }),
    }),
    {
      name: 'reelcut-active-project',
    }
  )
)
