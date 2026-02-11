import { Sun, Moon, Monitor } from 'lucide-react'
import { useThemeStore, type Theme } from '../stores/themeStore'

const cycle: Theme[] = ['light', 'dark', 'system']
const labels: Record<Theme, string> = {
  light: 'Light',
  dark: 'Dark',
  system: 'System',
}

function ThemeIcon({ theme }: { theme: Theme }) {
  switch (theme) {
    case 'light':
      return <Sun size={18} />
    case 'dark':
      return <Moon size={18} />
    case 'system':
      return <Monitor size={18} />
  }
}

export function ThemeToggle() {
  const theme = useThemeStore((s) => s.theme)
  const setTheme = useThemeStore((s) => s.setTheme)

  const cycleTheme = () => {
    const i = cycle.indexOf(theme)
    setTheme(cycle[(i + 1) % cycle.length])
  }

  return (
    <button
      type="button"
      onClick={cycleTheme}
      className="p-2 rounded-lg text-[var(--app-fg-muted)] transition-[var(--motion-duration-fast)] hover:bg-[var(--app-bg-overlay)] hover:text-[var(--app-fg)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg-raised)]"
      aria-label={`Theme: ${labels[theme]}. Click to switch.`}
      title={labels[theme]}
    >
      <ThemeIcon theme={theme} />
    </button>
  )
}
