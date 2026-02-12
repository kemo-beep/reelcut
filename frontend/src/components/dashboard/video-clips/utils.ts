import type { ClipSuggestion } from '../../../lib/api/analysis'
import type { TranscriptSegment } from '../../../types'

export function formatTimeForLabel(sec: number): string {
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

export function clipNameFromSuggestion(suggestion: ClipSuggestion, index: number): string {
  const start = formatTimeForLabel(suggestion.start_time)
  if (suggestion.transcript && suggestion.transcript.length > 0) {
    const truncated = suggestion.transcript.slice(0, 40)
    return truncated.length < suggestion.transcript.length ? `${truncated}…` : truncated
  }
  return `Clip ${index + 1} (${start})`
}

export function segmentsForClip(
  segments: TranscriptSegment[],
  clipStart: number,
  clipEnd: number
): TranscriptSegment[] {
  return segments
    .filter((s) => s.end_time > clipStart && s.start_time < clipEnd)
    .map((s) => ({
      ...s,
      start_time: Math.max(0, s.start_time - clipStart),
      end_time: Math.min(clipEnd - clipStart, s.end_time - clipStart),
    }))
}
