import { ClipStripCard } from '../../clip/ClipStripCard'
import type { Clip } from '../../../types'

export interface ClipStripsSectionProps {
  clips: Clip[]
  selectedClipId: string | null
  onSelectClip: (clipId: string) => void
}

export function ClipStripsSection({ clips, selectedClipId, onSelectClip }: ClipStripsSectionProps) {
  return (
    <section aria-label="Clips" className="w-full">
      <h2 className="text-sm font-medium text-[var(--app-fg-muted)] mb-3">Clips — click to select and play</h2>
      <div className="flex gap-4 overflow-x-auto pb-2 -mx-1 px-1 scroll-smooth snap-x snap-mandatory">
        {clips.map((clip) => (
          <ClipStripCard
            key={clip.id}
            clip={clip}
            isSelected={clip.id === selectedClipId}
            onSelect={() => onSelectClip(clip.id)}
            className="snap-start"
          />
        ))}
      </div>
    </section>
  )
}
