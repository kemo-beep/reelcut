import { useEffect } from 'react'
import { useThemeStore } from '../stores/themeStore'

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const theme = useThemeStore((s) => s.theme)

  useEffect(() => {
    const apply = () => {
      const effective = useThemeStore.getState().getEffectiveTheme()
      document.documentElement.classList.toggle('dark', effective === 'dark')
    }
    apply()
    if (theme === 'system') {
      const m = window.matchMedia('(prefers-color-scheme: dark)')
      m.addEventListener('change', apply)
      return () => m.removeEventListener('change', apply)
    }
  }, [theme])

  return <>{children}</>
}
