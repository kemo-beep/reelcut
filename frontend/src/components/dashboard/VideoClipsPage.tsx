import { useQueryClient } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { ErrorState } from '../ui/error-state'
import { MainVideoView } from './video-clips/MainVideoView'
import { VideoClipsPageHeader } from './video-clips/VideoClipsPageHeader'
import { ViewModeTabs } from './video-clips/ViewModeTabs'
import { ClipsEmptyState } from './video-clips/ClipsEmptyState'
import { ClipsViewContent } from './video-clips/ClipsViewContent'
import { ClipSuggestionsPanel } from './video-clips/ClipSuggestionsPanel'
import { useVideoClipsPage } from './video-clips/useVideoClipsPage'
import { ClipsViewSkeleton, VideoClipsPageSkeleton } from './video-clips/VideoClipsPageSkeletons'

export function VideoClipsPage() {
  const queryClient = useQueryClient()
  const {
    videoId,
    video,
    videoLoading,
    videoError,
    clipsLoading,
    sortedClips,
    selectedClipId,
    selectedClip,
    viewMode,
    setViewMode,
    segments,
    playerRef,
    playbackUrl,
    playbackLoading,
    clipPlaybackCurrentTime,
    setClipPlaybackCurrentTime,
    suggestions,
    suggestMut,
    acceptClipMut,
    acceptAllMut,
    handleTimelineSeek,
    handleSegmentClick,
    handleSegmentChange,
    handleSelectClip,
    clipTranscriptSegments,
    durationSec,
    timelineSegments,
    timelineCurrentTime,
  } = useVideoClipsPage()

  if (videoLoading || !video) {
    return <VideoClipsPageSkeleton />
  }

  if (videoError) {
    return (
      <div className="space-y-6">
        <Link
          to="/dashboard/videos"
          className="text-sm text-[var(--app-accent)] hover:underline focus-visible:outline focus-visible:ring-2 focus-visible:ring-[var(--app-accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--app-bg)] rounded"
        >
          Back to Videos
        </Link>
        <ErrorState
          message="Video not found."
          onRetry={() => queryClient.invalidateQueries({ queryKey: ['video', videoId] })}
        />
      </div>
    )
  }

  return (
    <main className="space-y-6" aria-label="Clips">
      <VideoClipsPageHeader
        videoId={videoId}
        videoFilename={video.original_filename}
        onSuggestClips={() => suggestMut.mutate()}
        isSuggesting={suggestMut.isPending}
      />

      <ViewModeTabs value={viewMode} onChange={setViewMode} />

      {viewMode === 'main-video' && (
        <MainVideoView videoId={videoId} video={video} segments={segments} />
      )}

      {viewMode === 'clips' &&
        (clipsLoading ? (
          <ClipsViewSkeleton />
        ) : sortedClips.length === 0 ? (
          <ClipsEmptyState
            onSuggestClips={() => suggestMut.mutate()}
            isSuggesting={suggestMut.isPending}
            videoId={videoId}
          />
        ) : (
          <ClipsViewContent
            sortedClips={sortedClips}
            selectedClipId={selectedClipId}
            selectedClip={selectedClip}
            onSelectClip={handleSelectClip}
            playerRef={playerRef}
            playbackUrl={playbackUrl}
            playbackLoading={playbackLoading}
            clipPlaybackCurrentTime={clipPlaybackCurrentTime}
            onClipTimeUpdate={setClipPlaybackCurrentTime}
            durationSec={durationSec}
            timelineSegments={timelineSegments}
            timelineCurrentTime={timelineCurrentTime}
            onTimelineSeek={handleTimelineSeek}
            onSegmentClick={handleSegmentClick}
            onSegmentChange={handleSegmentChange}
            clipTranscriptSegments={clipTranscriptSegments}
          />
        ))}

      {viewMode === 'clips' && !clipsLoading && suggestions.length > 0 && (
        <ClipSuggestionsPanel
          suggestions={suggestions}
          onAccept={(suggestion, index) => acceptClipMut.mutate({ suggestion, index })}
          onAcceptAll={() => acceptAllMut.mutate()}
          isAccepting={acceptClipMut.isPending}
          isAcceptingAll={acceptAllMut.isPending}
        />
      )}
    </main>
  )
}
