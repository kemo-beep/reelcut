import { get, post, put, del } from './client'
import type { Clip, ClipStyle } from '../../types'
import type { PaginatedResponse } from '../../types'

export interface CreateClipInput {
  video_id: string
  name: string
  start_time: number
  end_time: number
  aspect_ratio?: string
  virality_score?: number | null
  from_suggestion?: boolean
}

export interface ListClipsParams {
  page?: number
  per_page?: number
  video_id?: string
  status?: string
  sort_by?: string
  sort_order?: string
}

export interface ListClipsResponse extends PaginatedResponse<{ clips: Clip[] }> {}

export async function listClips(params?: ListClipsParams): Promise<ListClipsResponse> {
  const search = new URLSearchParams()
  if (params?.page != null) search.set('page', String(params.page))
  if (params?.per_page != null) search.set('per_page', String(params.per_page))
  if (params?.video_id) search.set('video_id', params.video_id)
  if (params?.status) search.set('status', params.status)
  if (params?.sort_by) search.set('sort_by', params.sort_by)
  if (params?.sort_order) search.set('sort_order', params.sort_order)
  const q = search.toString()
  return get(`/api/v1/clips${q ? `?${q}` : ''}`)
}

export async function createClip(input: CreateClipInput): Promise<{ clip: Clip }> {
  return post('/api/v1/clips', {
    ...input,
    from_suggestion: input.from_suggestion ?? false,
  })
}

export async function getClip(id: string): Promise<{ clip: Clip }> {
  return get(`/api/v1/clips/${id}`)
}

export async function updateClip(
  id: string,
  body: Partial<Pick<Clip, 'name' | 'start_time' | 'end_time' | 'aspect_ratio'>>
): Promise<{ clip: Clip }> {
  return put(`/api/v1/clips/${id}`, body)
}

export async function deleteClip(id: string): Promise<void> {
  await del(`/api/v1/clips/${id}`)
}

export async function renderClip(id: string): Promise<unknown> {
  return post(`/api/v1/clips/${id}/render`)
}

export async function getRenderStatus(id: string): Promise<unknown> {
  return get(`/api/v1/clips/${id}/status`)
}

export function getClipDownloadUrl(id: string, baseUrl?: string): string {
  const base = baseUrl ?? (typeof import.meta !== 'undefined' && import.meta.env?.VITE_API_URL) ?? 'http://localhost:8080'
  return `${base.replace(/\/$/, '')}/api/v1/clips/${id}/download`
}

export async function duplicateClip(id: string): Promise<{ clip: Clip }> {
  return post(`/api/v1/clips/${id}/duplicate`)
}

export async function getClipStyle(clipId: string): Promise<{ style: ClipStyle }> {
  return get(`/api/v1/clips/${clipId}/style`)
}

export async function updateClipStyle(clipId: string, style: Partial<ClipStyle>): Promise<{ style: ClipStyle }> {
  return put(`/api/v1/clips/${clipId}/style`, style)
}

export async function applyTemplateToClip(clipId: string, templateId: string): Promise<unknown> {
  return post(`/api/v1/clips/${clipId}/style/apply-template/${templateId}`)
}
