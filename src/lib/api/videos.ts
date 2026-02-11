import { get, post, del } from './client'
import type { Video } from '../../types'
import type { PaginatedResponse } from '../../types'

export interface GetUploadUrlInput {
  project_id: string
  filename: string
}

export interface GetUploadUrlResponse {
  video: Pick<Video, 'id' | 'project_id' | 'original_filename' | 'status'>
  upload: { upload_url: string; method: string }
}

export interface ListVideosParams {
  page?: number
  per_page?: number
  project_id?: string
  status?: string
  sort_by?: string
  sort_order?: string
}

export interface ListVideosResponse extends PaginatedResponse<{ videos: Video[] }> {}

export async function getUploadUrl(input: GetUploadUrlInput): Promise<GetUploadUrlResponse> {
  return post('/api/v1/videos/upload', input)
}

export async function confirmUpload(videoId: string): Promise<{ message: string; video_id: string }> {
  return post(`/api/v1/videos/${videoId}/confirm`)
}

export async function uploadFromUrl(_url: string): Promise<unknown> {
  return post('/api/v1/videos/upload/url', { url: _url })
}

export async function listVideos(params?: ListVideosParams): Promise<ListVideosResponse> {
  const search = new URLSearchParams()
  if (params?.page != null) search.set('page', String(params.page))
  if (params?.per_page != null) search.set('per_page', String(params.per_page))
  if (params?.project_id) search.set('project_id', params.project_id)
  if (params?.status) search.set('status', params.status)
  if (params?.sort_by) search.set('sort_by', params.sort_by)
  if (params?.sort_order) search.set('sort_order', params.sort_order)
  const q = search.toString()
  return get(`/api/v1/videos${q ? `?${q}` : ''}`)
}

export async function getVideo(id: string): Promise<{ video: Video }> {
  return get(`/api/v1/videos/${id}`)
}

export async function deleteVideo(id: string): Promise<void> {
  await del(`/api/v1/videos/${id}`)
}

export async function getVideoMetadata(id: string): Promise<unknown> {
  return get(`/api/v1/videos/${id}/metadata`)
}

export function getThumbnailUrl(id: string, baseUrl?: string): string {
  const base = baseUrl ?? (typeof import.meta !== 'undefined' && import.meta.env?.VITE_API_URL) ?? 'http://localhost:8080'
  return `${base.replace(/\/$/, '')}/api/v1/videos/${id}/thumbnail`
}

export async function getPlaybackUrl(id: string): Promise<{ url: string | null; status?: string }> {
  return get(`/api/v1/videos/${id}/playback-url`)
}
