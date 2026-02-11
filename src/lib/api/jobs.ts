import { get, post } from './client'
import type { ProcessingJob } from '../../types'
import type { PaginatedResponse } from '../../types'

export interface ListJobsResponse extends PaginatedResponse<{ jobs: ProcessingJob[] }> {}

export async function listJobs(params?: {
  page?: number
  per_page?: number
  status?: string
}): Promise<ListJobsResponse> {
  const search = new URLSearchParams()
  if (params?.page != null) search.set('page', String(params.page))
  if (params?.per_page != null) search.set('per_page', String(params.per_page))
  if (params?.status) search.set('status', params.status)
  const q = search.toString()
  return get(`/api/v1/jobs${q ? `?${q}` : ''}`)
}

export async function getJob(id: string): Promise<{ job: ProcessingJob }> {
  return get(`/api/v1/jobs/${id}`)
}

export async function cancelJob(id: string): Promise<unknown> {
  return post(`/api/v1/jobs/${id}/cancel`)
}
