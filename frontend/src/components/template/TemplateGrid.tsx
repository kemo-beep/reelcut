import { TemplateCard } from './TemplateCard'
import type { Template } from '../../types'
import { cn } from '../../lib/utils'

export interface TemplateGridProps {
  templates: Template[]
  className?: string
}

export function TemplateGrid({ templates, className }: TemplateGridProps) {
  return (
    <ul
      className={cn(
        'grid gap-3 sm:grid-cols-2 lg:grid-cols-3',
        className
      )}
      role="list"
    >
      {templates.map((t) => (
        <li key={t.id}>
          <TemplateCard template={t} />
        </li>
      ))}
    </ul>
  )
}
