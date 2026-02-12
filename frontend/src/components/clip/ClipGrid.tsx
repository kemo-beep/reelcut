import { ClipCard } from './ClipCard'
import type { Clip } from '../../types'
import { cn } from '../../lib/utils'

export interface ClipGridProps {
  clips: Clip[]
  className?: string
}

export function ClipGrid({ clips, className }: ClipGridProps) {
  return (
    <ul
      className={cn('space-y-3', className)}
      role="list"
    >
      {clips.map((c) => (
        <li key={c.id}>
          <ClipCard clip={c} />
        </li>
      ))}
    </ul>
  )
}
