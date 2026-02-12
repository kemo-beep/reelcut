import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import * as transcriptionsApi from '../lib/api/transcriptions'
import type { TranscriptSegment } from '../types'

export function useTranscriptionByVideoId(videoId: string | undefined, enabled = true) {
  const queryClient = useQueryClient()
  const query = useQuery({
    queryKey: ['transcription', videoId],
    queryFn: () => transcriptionsApi.getTranscriptionByVideoId(videoId!),
    enabled: !!videoId && enabled,
    retry: (_, err: { status?: number }) => (err?.status === 404 ? false : true),
  })
  const createMutation = useMutation({
    mutationFn: () => transcriptionsApi.createTranscription(videoId!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transcription', videoId] })
    },
  })
  const updateSegmentMutation = useMutation({
    mutationFn: ({
      transcriptionId,
      segmentId,
      text,
    }: {
      transcriptionId: string
      segmentId: string
      text: string
    }) =>
      transcriptionsApi.updateSegment(transcriptionId, segmentId, { text }),
    onSuccess: (_, { transcriptionId }) => {
      queryClient.invalidateQueries({ queryKey: ['transcription'] })
    },
  })
  return {
    ...query,
    transcription: query.data?.transcription,
    createTranscription: createMutation.mutateAsync,
    updateSegment: updateSegmentMutation.mutateAsync,
    isCreating: createMutation.isPending,
    isUpdatingSegment: updateSegmentMutation.isPending,
  }
}

export function useTranscription(transcriptionId: string | undefined, enabled = true) {
  const query = useQuery({
    queryKey: ['transcription', transcriptionId],
    queryFn: () => transcriptionsApi.getTranscription(transcriptionId!),
    enabled: !!transcriptionId && enabled,
  })
  return {
    ...query,
    transcription: query.data?.transcription,
  }
}
