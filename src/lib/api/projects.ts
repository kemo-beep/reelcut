import { get, post, put, del } from './client'
import type { Project } from '../../types'
import type { PaginatedResponse } from '../../types'

export interface CreateProjectInput {
  name: string
  description?: string | null
}

export interface UpdateProjectInput {
  name?: string
  description?: string | null
}

export interface ListProjectsResponse extends PaginatedResponse<{ projects: Project[] }> {}

export async function listProjects(params?: {
  page?: number
  per_page?: number
}): Promise<ListProjectsResponse> {
  const search = new URLSearchParams()
  if (params?.page != null) search.set('page', String(params.page))
  if (params?.per_page != null) search.set('per_page', String(params.per_page))
  const q = search.toString()
  return get(`/api/v1/projects${q ? `?${q}` : ''}`)
}

export async function createProject(input: CreateProjectInput): Promise<{ project: Project }> {
  return post('/api/v1/projects', input)
}

export async function getProject(id: string): Promise<{ project: Project }> {
  return get(`/api/v1/projects/${id}`)
}

export async function updateProject(id: string, input: UpdateProjectInput): Promise<{ project: Project }> {
  return put(`/api/v1/projects/${id}`, input)
}

export async function deleteProject(id: string): Promise<void> {
  await del(`/api/v1/projects/${id}`)
}
