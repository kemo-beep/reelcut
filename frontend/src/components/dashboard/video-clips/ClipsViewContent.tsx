import { RefObject } from 'react'
import { ClipStripsSection } from './ClipStripsSection'
import { ClipPlayerSection } from './ClipPlayerSection'
import { ClipTranscriptPanel } from './ClipTranscriptPanel'
import type { Clip } from '../../../types'
import type { ClipTimelineSegment } from '../../clip/ClipTimeline'
import type { TranscriptSegment } from '../../../types'

export interface ClipsViewContentProps {
  sortedClips: Clip[]
  selectedClipId: string | null
  selectedClip: Clip | null
  onSelectClip: (clipId: string) => void
  playerRef: RefObject<{ seek: (time: number) => void; getCurrentTime: () => number } | null>
  playbackUrl: string | null
  playbackLoading: boolean
  clipPlaybackCurrentTime: number
  onClipTimeUpdate: (time: number) => void
  durationSec: number
  timelineSegments: ClipTimelineSegment[]
  timelineCurrentTime: number
  onTimelineSeek: (time: number) => void
  onSegmentClick: (segmentId: string) => void
  onSegmentChange: (id: string, payload: { start_time: number; end_time: number }) => void
  clipTranscriptSegments: TranscriptSegment[]
}

export function ClipsViewContent({
  sortedClips,
  selectedClipId,
  selectedClip,
  onSelectClip,
  playerRef,
  playbackUrl,
  playbackLoading,
  clipPlaybackCurrentTime,
  onClipTimeUpdate,
  durationSec,
  timelineSegments,
  timelineCurrentTime,
  onTimelineSeek,
  onSegmentClick,
  onSegmentChange,
  clipTranscriptSegments,
}: ClipsViewContentProps) {
  const handleSeek = (time: number) => playerRef.current?.seek(time)

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 lg:grid-cols-[1fr_minmax(320px,400px)] gap-6 items-start">
        <ClipPlayerSection
          playerRef={playerRef}
          playbackUrl={playbackUrl}
          playbackLoading={playbackLoading}
          selectedClip={selectedClip}
          selectedClipId={selectedClipId}
          clipPlaybackCurrentTime={clipPlaybackCurrentTime}
          onTimeUpdate={onClipTimeUpdate}
          durationSec={durationSec}
          timelineSegments={timelineSegments}
          timelineCurrentTime={timelineCurrentTime}
          onTimelineSeek={onTimelineSeek}
          onSegmentClick={onSegmentClick}
          onSegmentChange={onSegmentChange}
        />
        <ClipTranscriptPanel
          clip={selectedClip}
          segments={clipTranscriptSegments}
          currentTime={clipPlaybackCurrentTime}
          onSeek={handleSeek}
        />
      </div>
      <ClipStripsSection
        clips={sortedClips}
        selectedClipId={selectedClipId}
        onSelectClip={onSelectClip}
      />
      
    </div>
  )
}
