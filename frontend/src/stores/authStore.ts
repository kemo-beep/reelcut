import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../types'

const STORAGE_KEY = 'reelcut-auth'

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  _hasHydrated: boolean
  onSessionExpired: (() => void) | null
  setAuth: (user: User, accessToken: string, refreshToken: string) => void
  clearAuth: () => void
  setOnSessionExpired: (fn: (() => void) | null) => void
  getAccessToken: () => string | null
}

/**
 * Reads the access token from persisted storage (localStorage).
 * Use in route guards so we don't redirect to login before zustand persist has rehydrated.
 */
export function getStoredAccessToken(): string | null {
  if (typeof window === 'undefined') return null
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    const data = JSON.parse(raw) as { state?: { accessToken?: string } }
    const token = data?.state?.accessToken
    return typeof token === 'string' ? token : null
  } catch {
    return null
  }
}

/**
 * Returns true when we're in an environment where we cannot know auth (e.g. server).
 * In that case route guards should not redirect (let the route load; client will re-check).
 */
function isClient(): boolean {
  return typeof window !== 'undefined'
}

/**
 * Returns true if the user has a token (in memory or in persisted storage).
 * Use in beforeLoad so refresh doesn't redirect to login before rehydration.
 * On the server this returns false (no localStorage), so guards must only redirect on the client.
 */
export function hasStoredAuth(): boolean {
  return !!(useAuthStore.getState().getAccessToken() ?? getStoredAccessToken())
}

/**
 * Use in route beforeLoad: only redirect to login when we're on the client and there's no auth.
 * On the server we skip the redirect so SSR doesn't send everyone to login (client will re-check).
 */
export function shouldRedirectToLogin(): boolean {
  return isClient() && !hasStoredAuth()
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      _hasHydrated: false,
      onSessionExpired: null,
      setAuth: (user, accessToken, refreshToken) =>
        set({ user, accessToken, refreshToken }),
      clearAuth: () =>
        set({ user: null, accessToken: null, refreshToken: null }),
      setOnSessionExpired: (fn) => set({ onSessionExpired: fn }),
      getAccessToken: () => get().accessToken,
    }),
    {
      name: STORAGE_KEY,
      partialize: (s) => ({
        user: s.user,
        accessToken: s.accessToken,
        refreshToken: s.refreshToken,
      }),
      onRehydrateStorage: () => () => {
        useAuthStore.setState({ _hasHydrated: true })
      },
    }
  )
)

/** For use in components: true once persisted auth has been rehydrated from storage. */
export function useAuthHasHydrated(): boolean {
  return useAuthStore((s) => s._hasHydrated)
}
