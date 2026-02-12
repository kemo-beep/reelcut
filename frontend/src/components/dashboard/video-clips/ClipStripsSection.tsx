import { useEffect, useRef } from 'react'
import { ClipStripCard } from '../../clip/ClipStripCard'
import type { Clip } from '../../../types'

export interface ClipStripsSectionProps {
  clips: Clip[]
  selectedClipId: string | null
  onSelectClip: (clipId: string) => void
}

export function ClipStripsSection({ clips, selectedClipId, onSelectClip }: ClipStripsSectionProps) {
  const selectedClipRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (selectedClipId) {
      selectedClipRef.current?.scrollIntoView({ block: 'nearest', behavior: 'smooth' })
    }
  }, [selectedClipId])

  const clipCount = clips.length

  return (
    <section aria-label="Clips" className="w-full">
      <h2 className="mb-1 text-lg font-semibold text-[var(--app-fg)]">
        Clips {clipCount > 0 && `(${clipCount})`}
      </h2>
      <p className="mb-3 text-sm text-[var(--app-fg-muted)]">
        Click to select and play
      </p>
      <div className="relative -mx-1 px-1 pb-2">
        <div className="flex gap-4 overflow-x-auto scroll-smooth snap-x snap-mandatory pb-1 [scrollbar-gutter:stable]">
          {clips.map((clip) => {
            const isSelected = clip.id === selectedClipId
            return (
              <div
                key={clip.id}
                ref={isSelected ? selectedClipRef : undefined}
                className="snap-start shrink-0"
              >
                <ClipStripCard
                  clip={clip}
                  isSelected={isSelected}
                  onSelect={() => onSelectClip(clip.id)}
                />
              </div>
            )
          })}
        </div>
        {/* Scroll hint: gradient fade on the right when content overflows */}
        <div
          aria-hidden
          className="pointer-events-none absolute right-0 top-0 bottom-3 w-8 shrink-0 bg-gradient-to-l from-[var(--app-bg)] to-transparent lg:bottom-4"
        />
      </div>
    </section>
  )
}
