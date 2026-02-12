import { useRef, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getPlaybackUrl } from '../../../lib/api/videos'
import { VideoPlayer } from '../../video/VideoPlayer'
import { TranscriptViewer } from '../../transcription/TranscriptViewer'
import type { Video } from '../../../types'
import type { TranscriptSegment } from '../../../types'

export interface MainVideoViewProps {
  videoId: string
  video: Video
  segments: TranscriptSegment[]
}

export function MainVideoView({ videoId, video, segments }: MainVideoViewProps) {
  const playerRef = useRef<{ seek: (time: number) => void; getCurrentTime: () => number }>(null)
  const [currentTime, setCurrentTime] = useState(0)

  const { data: playbackData, isFetching: playbackLoading } = useQuery({
    queryKey: ['video-playback', videoId],
    queryFn: () => getPlaybackUrl(videoId),
    enabled: !!videoId,
  })
  const playbackUrl = playbackData?.url ?? null

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[1fr_minmax(320px,400px)] gap-6 items-start">
      <div className="rounded-xl overflow-hidden border border-[var(--app-border)] bg-[var(--app-bg)]">
        {playbackLoading && (
          <div className="aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
            Loading…
          </div>
        )}
        {!playbackLoading && playbackUrl && (
          <VideoPlayer
            ref={playerRef}
            src={playbackUrl}
            initialTime={currentTime}
            onTimeUpdate={setCurrentTime}
          />
        )}
        {!playbackLoading && !playbackUrl && (
          <div className="aspect-video flex items-center justify-center text-[var(--app-fg-muted)]">
            Playback not available for this video.
          </div>
        )}
      </div>
      <div
        className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-4 lg:sticky lg:top-6"
        id="transcript-panel"
      >
        <h2 className="font-semibold text-[var(--app-fg)] mb-3">Transcript: {video.original_filename}</h2>
        {segments.length > 0 ? (
          <TranscriptViewer
            segments={segments}
            currentTime={currentTime}
            onSeek={(t) => playerRef.current?.seek(t)}
            className="max-h-[50vh]"
          />
        ) : (
          <p className="text-caption text-[var(--app-fg-muted)]">No transcript available.</p>
        )}
      </div>
    </div>
  )
}
