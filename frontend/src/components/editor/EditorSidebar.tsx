import { cn } from '../../lib/utils'

export interface EditorSidebarProps {
  children: React.ReactNode
  className?: string
}

export function EditorSidebar({ children, className }: EditorSidebarProps) {
  return (
    <aside
      className={cn(
        'w-72 shrink-0 border-l border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4 overflow-y-auto',
        className
      )}
    >
      {children}
    </aside>
  )
}
