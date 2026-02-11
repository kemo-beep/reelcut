import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../types'

const STORAGE_KEY = 'reelcut-auth'

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  setAuth: (user: User, accessToken: string, refreshToken: string) => void
  clearAuth: () => void
  getAccessToken: () => string | null
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      setAuth: (user, accessToken, refreshToken) =>
        set({ user, accessToken, refreshToken }),
      clearAuth: () =>
        set({ user: null, accessToken: null, refreshToken: null }),
      getAccessToken: () => get().accessToken,
    }),
    { name: STORAGE_KEY }
  )
)
