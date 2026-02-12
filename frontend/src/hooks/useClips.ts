import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import * as clipsApi from '../lib/api/clips'
import type { ListClipsParams } from '../lib/api/clips'
import type { ClipStyle } from '../types'

export function useClips(params?: ListClipsParams) {
  const queryClient = useQueryClient()
  const query = useQuery({
    queryKey: ['clips', params],
    queryFn: () => clipsApi.listClips(params),
  })
  const createMutation = useMutation({
    mutationFn: clipsApi.createClip,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clips'] })
    },
  })
  const deleteMutation = useMutation({
    mutationFn: (id: string) => clipsApi.deleteClip(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clips'] })
    },
  })
  return {
    ...query,
    clips: query.data?.data?.clips ?? [],
    total: query.data?.data?.total ?? 0,
    createClip: createMutation.mutateAsync,
    deleteClip: deleteMutation.mutateAsync,
    isCreating: createMutation.isPending,
    isDeleting: deleteMutation.isPending,
  }
}

export function useClip(clipId: string | undefined, enabled = true) {
  const queryClient = useQueryClient()
  const query = useQuery({
    queryKey: ['clip', clipId],
    queryFn: () => clipsApi.getClip(clipId!),
    enabled: !!clipId && enabled,
  })
  const updateMutation = useMutation({
    mutationFn: ({ id, body }: { id: string; body: Parameters<typeof clipsApi.updateClip>[1] }) =>
      clipsApi.updateClip(id, body),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['clip', id] })
      queryClient.invalidateQueries({ queryKey: ['clips'] })
    },
  })
  const updateStyleMutation = useMutation({
    mutationFn: ({ clipId, style }: { clipId: string; style: Partial<ClipStyle> }) =>
      clipsApi.updateClipStyle(clipId, style),
    onSuccess: (_, { clipId }) => {
      queryClient.invalidateQueries({ queryKey: ['clip', clipId] })
    },
  })
  const renderMutation = useMutation({
    mutationFn: (id: string) => clipsApi.renderClip(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['clip', id] })
      queryClient.invalidateQueries({ queryKey: ['clips'] })
    },
  })
  return {
    ...query,
    clip: query.data?.clip,
    updateClip: updateMutation.mutateAsync,
    updateStyle: updateStyleMutation.mutateAsync,
    renderClip: renderMutation.mutateAsync,
    isUpdating: updateMutation.isPending,
    isRendering: renderMutation.isPending,
  }
}
