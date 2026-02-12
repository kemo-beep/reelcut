import { Layers } from 'lucide-react'
import { cn } from '../../lib/utils'

export interface LayerItem {
  id: string
  label: string
  type: 'video' | 'caption' | 'overlay' | 'audio'
}

export interface LayersPanelProps {
  layers: LayerItem[]
  selectedId?: string | null
  onSelect?: (id: string) => void
  className?: string
}

export function LayersPanel({
  layers,
  selectedId,
  onSelect,
  className,
}: LayersPanelProps) {
  return (
    <div className={cn('space-y-2', className)}>
      <div className="flex items-center gap-2 font-medium text-[var(--app-fg)]">
        <Layers size={18} />
        Layers
      </div>
      <ul className="space-y-1" role="list">
        {layers.map((layer) => (
          <li key={layer.id}>
            <button
              type="button"
              onClick={() => onSelect?.(layer.id)}
              className={cn(
                'w-full rounded-lg border px-3 py-2 text-left text-sm transition-colors',
                selectedId === layer.id
                  ? 'border-[var(--app-accent)] bg-[var(--app-accent-muted)] text-[var(--app-fg)]'
                  : 'border-[var(--app-border)] bg-[var(--app-bg)] text-[var(--app-fg-muted)] hover:border-[var(--app-border-strong)]'
              )}
            >
              {layer.label}
            </button>
          </li>
        ))}
      </ul>
    </div>
  )
}
