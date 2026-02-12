import { useParams, useRouterState } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState, useEffect, useRef, useCallback } from 'react'
import { getVideo, getPlaybackUrl } from '../../../lib/api/videos'
import { listClips, createClip, getClipPlaybackUrl, updateClip } from '../../../lib/api/clips'
import { getTranscriptionByVideoId } from '../../../lib/api/transcriptions'
import { suggestClips, type ClipSuggestion } from '../../../lib/api/analysis'
import { toast } from 'sonner'
import { clipNameFromSuggestion, segmentsForClip } from './utils'
import type { Clip, Video } from '../../../types'
import type { ClipTimelineSegment } from '../../clip/ClipTimeline'
import type { ViewMode } from './ViewModeTabs'

export function useVideoClipsPage() {
  const { videoId } = useParams({ from: '/dashboard/videos/$videoId/clips' })
  const queryClient = useQueryClient()
  const locationState = useRouterState({ select: (s) => s.location.state })
  const playerRef = useRef<{ seek: (time: number) => void; getCurrentTime: () => number }>(null)

  const [suggestions, setSuggestions] = useState<ClipSuggestion[]>([])
  const [selectedClipId, setSelectedClipId] = useState<string | null>(null)
  const [clipPlaybackCurrentTime, setClipPlaybackCurrentTime] = useState(0)
  const [viewMode, setViewMode] = useState<ViewMode>('clips')

  useEffect(() => {
    const fromState = (locationState as { suggestions?: ClipSuggestion[] })?.suggestions
    if (Array.isArray(fromState) && fromState.length > 0) setSuggestions(fromState)
  }, [locationState])

  const { data: videoData, isLoading: videoLoading, error: videoError } = useQuery({
    queryKey: ['video', videoId],
    queryFn: () => getVideo(videoId),
  })
  const { data: clipsData, isLoading: clipsLoading } = useQuery({
    queryKey: ['clips', { video_id: videoId }],
    queryFn: () => listClips({ video_id: videoId, per_page: 100 }),
    enabled: !!videoData?.video,
  })
  const { data: transData } = useQuery({
    queryKey: ['transcription', videoId],
    queryFn: () => getTranscriptionByVideoId(videoId),
    enabled: !!videoData?.video,
  })

  const clips = (clipsData?.data?.clips ?? []) as Clip[]
  const sortedClips = [...clips].sort((a, b) => a.start_time - b.start_time)
  const selectedClip = selectedClipId ? clips.find((c) => c.id === selectedClipId) ?? null : null

  const { data: playbackData, isFetching: playbackLoading } = useQuery({
    queryKey: ['clip-playback', selectedClipId],
    queryFn: () => getClipPlaybackUrl(selectedClipId!),
    enabled:
      !!selectedClipId &&
      !!(selectedClip?.storage_path != null && selectedClip.storage_path !== ''),
  })
  const playbackUrl = playbackData?.url ?? null

  useEffect(() => {
    if (clips.length > 0 && !selectedClipId) {
      const firstWithVideo = clips.find((c) => c.storage_path != null && c.storage_path !== '')
      setSelectedClipId((firstWithVideo ?? clips[0]).id)
    }
  }, [clips.length, selectedClipId, clips])

  const suggestMut = useMutation({
    mutationFn: () =>
      suggestClips(videoId, {
        min_duration: 7,
        max_duration: 60,
        max_suggestions: 20,
      }),
    onSuccess: (data) => {
      const list = data.suggestions ?? []
      setSuggestions(list)
      if (list.length > 0) {
        toast.success(`${list.length} clip suggestions ready`)
      } else {
        toast.info('No suggestions found for this video')
      }
    },
    onError: () => toast.error('Failed to get suggestions'),
  })
  const acceptClipMut = useMutation({
    mutationFn: ({ suggestion, index }: { suggestion: ClipSuggestion; index: number }) =>
      createClip({
        video_id: videoId,
        name: clipNameFromSuggestion(suggestion, index),
        start_time: suggestion.start_time,
        end_time: suggestion.end_time,
        aspect_ratio: '9:16',
        virality_score: suggestion.virality_score ?? undefined,
        from_suggestion: true,
      }),
    onSuccess: (_, { suggestion }) => {
      queryClient.invalidateQueries({ queryKey: ['clips', { video_id: videoId }] })
      setSuggestions((prev) => prev.filter((s) => s !== suggestion))
      toast.success('Clip created')
    },
    onError: () => toast.error('Failed to create clip'),
  })
  const acceptAllMut = useMutation({
    mutationFn: async () => {
      const list = suggestions
      for (let i = 0; i < list.length; i++) {
        const s = list[i]
        await createClip({
          video_id: videoId,
          name: clipNameFromSuggestion(s, i),
          start_time: s.start_time,
          end_time: s.end_time,
          aspect_ratio: '9:16',
          virality_score: s.virality_score ?? undefined,
          from_suggestion: true,
        })
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clips', { video_id: videoId }] })
      setSuggestions([])
      setViewMode('clips')
      toast.success('All suggestions added as clips')
    },
    onError: () => toast.error('Failed to create some clips'),
  })

  const updateClipMut = useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string
      payload: { start_time: number; end_time: number }
    }) => updateClip(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clips', { video_id: videoId }] })
      toast.success('Clip updated')
    },
    onError: () => toast.error('Failed to update clip'),
  })

  const handleTimelineSeek = useCallback(
    (time: number) => {
      const clip = sortedClips.find((c) => time >= c.start_time && time <= c.end_time)
      if (clip) {
        setSelectedClipId(clip.id)
        const offset = time - clip.start_time
        setClipPlaybackCurrentTime(offset)
        playerRef.current?.seek(offset)
      }
    },
    [sortedClips]
  )
  const handleSegmentClick = useCallback((segmentId: string) => {
    setSelectedClipId(segmentId)
    setClipPlaybackCurrentTime(0)
    playerRef.current?.seek(0)
  }, [])
  const handleSegmentChange = useCallback(
    (id: string, payload: { start_time: number; end_time: number }) => {
      updateClipMut.mutate({ id, payload })
    },
    [updateClipMut]
  )

  const transcription = transData?.transcription ?? null
  const segments = transcription?.segments ?? []
  const clipTranscriptSegments = selectedClip
    ? segmentsForClip(segments, selectedClip.start_time, selectedClip.end_time)
    : []

  const video = videoData?.video
  const durationSec = video?.duration_seconds ?? 0
  const timelineSegments: ClipTimelineSegment[] = sortedClips.map((c) => ({
    id: c.id,
    start_time: c.start_time,
    end_time: c.end_time,
    virality_score: c.virality_score ?? undefined,
    isSuggestion: false,
  }))
  const timelineCurrentTime =
    selectedClip != null ? selectedClip.start_time + clipPlaybackCurrentTime : 0

  const handleSelectClip = useCallback((clipId: string) => {
    setSelectedClipId(clipId)
    setClipPlaybackCurrentTime(0)
    playerRef.current?.seek(0)
  }, [])

  return {
    videoId,
    video,
    videoLoading,
    videoError,
    clipsLoading,
    clips,
    sortedClips,
    selectedClipId,
    selectedClip,
    setSelectedClipId,
    clipPlaybackCurrentTime,
    setClipPlaybackCurrentTime,
    viewMode,
    setViewMode,
    segments,
    playerRef,
    playbackUrl,
    playbackLoading,
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
  }
}
