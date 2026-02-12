import { LayoutTemplate } from 'lucide-react'
import { cn } from '../../lib/utils'

export interface TemplatePreviewProps {
  previewUrl?: string | null
  name?: string
  className?: string
}

export function TemplatePreview({
  previewUrl,
  name,
  className,
}: TemplatePreviewProps) {
  if (previewUrl?.startsWith('http')) {
    return (
      <img
        src={previewUrl}
        alt={name ?? 'Template preview'}
        className={cn('rounded-lg object-cover aspect-video w-full', className)}
      />
    )
  }
  return (
    <div
      className={cn(
        'flex aspect-video w-full items-center justify-center rounded-lg bg-[var(--app-bg)] text-[var(--app-fg-muted)]',
        className
      )}
    >
      <LayoutTemplate size={40} />
    </div>
  )
}
