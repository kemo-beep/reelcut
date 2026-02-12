import { create } from 'zustand'
import { persist } from 'zustand/middleware'

const STORAGE_KEY = 'reelcut-theme'

export type Theme = 'light' | 'dark' | 'system'
export type EffectiveTheme = 'light' | 'dark'

function getSystemPrefersDark(): boolean {
  if (typeof window === 'undefined') return false
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

export function getEffectiveTheme(theme: Theme): EffectiveTheme {
  if (theme === 'system') return getSystemPrefersDark() ? 'dark' : 'light'
  return theme
}

interface ThemeState {
  theme: Theme
  setTheme: (value: Theme) => void
  getEffectiveTheme: () => EffectiveTheme
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set, get) => ({
      theme: 'system',
      setTheme: (value) => set({ theme: value }),
      getEffectiveTheme: () => getEffectiveTheme(get().theme),
    }),
    { name: STORAGE_KEY }
  )
)
