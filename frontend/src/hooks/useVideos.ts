import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import * as videosApi from '../lib/api/videos'
import type { ListVideosParams } from '../lib/api/videos'

export function useVideos(params?: ListVideosParams) {
  const queryClient = useQueryClient()
  const query = useQuery({
    queryKey: ['videos', params],
    queryFn: () => videosApi.listVideos(params),
  })
  const deleteMutation = useMutation({
    mutationFn: (id: string) => videosApi.deleteVideo(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['videos'] })
    },
  })
  return {
    ...query,
    videos: query.data?.data?.videos ?? [],
    total: query.data?.data?.total ?? 0,
    deleteVideo: deleteMutation.mutateAsync,
    isDeleting: deleteMutation.isPending,
  }
}

export function useVideo(videoId: string | undefined, enabled = true) {
  const query = useQuery({
    queryKey: ['video', videoId],
    queryFn: () => videosApi.getVideo(videoId!),
    enabled: !!videoId && enabled,
  })
  return {
    ...query,
    video: query.data?.video,
  }
}

export function usePlaybackUrl(videoId: string | undefined, enabled = true) {
  const query = useQuery({
    queryKey: ['video-playback', videoId],
    queryFn: () => videosApi.getPlaybackUrl(videoId!),
    enabled: !!videoId && enabled,
  })
  return {
    ...query,
    playbackUrl: query.data?.url,
  }
}
