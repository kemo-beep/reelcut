import { useState } from 'react'
import type { TranscriptSegment } from '../../types'
import { cn } from '../../lib/utils'

function formatTime(sec: number): string {
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

export interface TranscriptEditorProps {
  segments: TranscriptSegment[]
  onSegmentChange?: (segmentId: string, text: string) => void
  className?: string
}

export function TranscriptEditor({
  segments,
  onSegmentChange,
  className,
}: TranscriptEditorProps) {
  const [editingId, setEditingId] = useState<string | null>(null)
  const [draft, setDraft] = useState('')

  return (
    <div className={cn('space-y-2', className)}>
      {segments.map((seg) => {
        const isEditing = editingId === seg.id
        const text = isEditing ? draft : seg.text
        return (
          <div
            key={seg.id}
            className="rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] px-4 py-2"
          >
            <span className="text-caption text-[var(--app-fg-muted)]">
              {formatTime(seg.start_time)} – {formatTime(seg.end_time)}
            </span>
            {isEditing ? (
              <textarea
                value={draft}
                onChange={(e) => setDraft(e.target.value)}
                onBlur={() => {
                  onSegmentChange?.(seg.id, draft)
                  setEditingId(null)
                }}
                className="mt-1 w-full resize-none rounded border border-[var(--app-border)] bg-[var(--app-bg-raised)] px-2 py-1 text-sm text-[var(--app-fg)]"
                rows={2}
                autoFocus
              />
            ) : (
              <p
                className="mt-0.5 cursor-pointer text-sm text-[var(--app-fg)] hover:underline"
                onClick={() => {
                  setEditingId(seg.id)
                  setDraft(seg.text)
                }}
              >
                {text}
              </p>
            )}
          </div>
        )
      })}
    </div>
  )
}
