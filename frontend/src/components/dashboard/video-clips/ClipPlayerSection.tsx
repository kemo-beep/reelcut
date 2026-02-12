import { RefObject } from 'react'
import { VideoPlayer } from '../../video/VideoPlayer'
import { ClipTimeline, type ClipTimelineSegment } from '../../clip/ClipTimeline'
import type { Clip } from '../../../types'

export interface ClipPlayerSectionProps {
  playerRef: RefObject<{ seek: (time: number) => void; getCurrentTime: () => number } | null>
  playbackUrl: string | null
  playbackLoading: boolean
  selectedClip: Clip | null
  selectedClipId: string | null
  clipPlaybackCurrentTime: number
  onTimeUpdate: (time: number) => void
  durationSec: number
  timelineSegments: ClipTimelineSegment[]
  timelineCurrentTime: number
  onTimelineSeek: (time: number) => void
  onSegmentClick: (segmentId: string) => void
  onSegmentChange: (id: string, payload: { start_time: number; end_time: number }) => void
}

export function ClipPlayerSection({
  playerRef,
  playbackUrl,
  playbackLoading,
  selectedClip,
  selectedClipId,
  clipPlaybackCurrentTime,
  onTimeUpdate,
  durationSec,
  timelineSegments,
  timelineCurrentTime,
  onTimelineSeek,
  onSegmentClick,
  onSegmentChange,
}: ClipPlayerSectionProps) {
  return (
    <div className="space-y-4">
      <div className="rounded-xl overflow-hidden border border-[var(--app-border)] bg-[var(--app-bg)]">
        {playbackLoading && (
          <div className="aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
            Loading…
          </div>
        )}
        {!playbackLoading && playbackUrl && (
          <VideoPlayer
            key={selectedClipId ?? ''}
            ref={playerRef}
            src={playbackUrl}
            initialTime={clipPlaybackCurrentTime}
            onTimeUpdate={onTimeUpdate}
          />
        )}
        {!playbackLoading && selectedClip && !playbackUrl && (
          <div className="aspect-video flex flex-col items-center justify-center text-center p-6 text-[var(--app-fg-muted)]">
            <p className="font-medium text-[var(--app-fg)]">Clip has no cut file yet</p>
            <p className="text-sm mt-1">
              This clip was created without a pre-cut video. Use the editor to render it, or run Auto cut on the video
              to generate cut files for all clips.
            </p>
          </div>
        )}
        {!playbackLoading && !selectedClip && (
          <div className="aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
            Select a clip from the timeline
          </div>
        )}
      </div>

      {durationSec > 0 && (
        <ClipTimeline
          duration={durationSec}
          segments={timelineSegments}
          currentTime={timelineCurrentTime}
          onSeek={onTimelineSeek}
          onSegmentChange={onSegmentChange}
          onSegmentClick={onSegmentClick}
        />
      )}
    </div>
  )
}
