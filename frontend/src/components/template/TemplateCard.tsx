import { Link } from '@tanstack/react-router'
import { LayoutTemplate } from 'lucide-react'
import type { Template } from '../../types'
import { cn } from '../../lib/utils'

export interface TemplateCardProps {
  template: Template
  className?: string
}

export function TemplateCard({ template, className }: TemplateCardProps) {
  return (
    <Link
      to="/dashboard/templates"
      search={{}}
      className={cn(
        'block rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-5 shadow-card transition-[var(--motion-duration-fast)] hover:border-[var(--app-border-strong)] focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)]',
        className
      )}
    >
      <div className="flex items-start gap-3">
        <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-[var(--app-bg)]">
          <LayoutTemplate size={24} className="text-[var(--app-fg-muted)]" />
        </div>
        <div className="min-w-0 flex-1">
          <p className="font-medium text-[var(--app-fg)]">{template.name}</p>
          {template.category && (
            <p className="text-caption text-[var(--app-fg-muted)] mt-0.5">{template.category}</p>
          )}
        </div>
      </div>
    </Link>
  )
}
