import { TranscriptViewer } from '../../transcription/TranscriptViewer'
import type { Clip } from '../../../types'
import type { TranscriptSegment } from '../../../types'

export interface ClipTranscriptPanelProps {
  clip: Clip | null
  segments: TranscriptSegment[]
  currentTime: number
  onSeek: (time: number) => void
}

export function ClipTranscriptPanel({ clip, segments, currentTime, onSeek }: ClipTranscriptPanelProps) {
  return (
    <div
      className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4 lg:sticky lg:top-6"
      id="transcript-panel"
    >
      <h2 className="font-semibold text-[var(--app-fg)] mb-3">
        {clip ? `Transcript: ${clip.name}` : 'Transcript'}
      </h2>
      {clip && segments.length > 0 ? (
        <TranscriptViewer
          segments={segments}
          currentTime={currentTime}
          onSeek={onSeek}
          className="max-h-[50vh]"
        />
      ) : clip ? (
        <p className="text-caption text-[var(--app-fg-muted)]">No transcript for this clip range.</p>
      ) : (
        <p className="text-caption text-[var(--app-fg-muted)]">Select a clip to see its transcript.</p>
      )}
    </div>
  )
}
