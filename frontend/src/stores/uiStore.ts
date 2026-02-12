import { create } from 'zustand'

interface ModalState {
  id: string
  data?: unknown
}

interface UIState {
  modals: ModalState[]
  openModal: (id: string, data?: unknown) => void
  closeModal: (id: string) => void
  closeAllModals: () => void
}

export const useUIStore = create<UIState>((set, get) => ({
  modals: [],
  openModal: (id, data) =>
    set((s) => ({
      modals: [...s.modals.filter((m) => m.id !== id), { id, data }],
    })),
  closeModal: (id) =>
    set((s) => ({ modals: s.modals.filter((m) => m.id !== id) })),
  closeAllModals: () => set({ modals: [] }),
}))
