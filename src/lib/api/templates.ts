import { get, post, put, del } from './client'
import type { Template } from '../../types'
import type { PaginatedResponse } from '../../types'

export interface CreateTemplateInput {
  name: string
  category?: string | null
  is_public?: boolean
  style_config: Record<string, unknown>
}

export interface UpdateTemplateInput {
  name?: string
  category?: string | null
  is_public?: boolean
  style_config?: Record<string, unknown>
}

export interface ListTemplatesResponse extends PaginatedResponse<{ templates: Template[] }> {}

export async function listTemplates(params?: { page?: number; per_page?: number }): Promise<ListTemplatesResponse> {
  const search = new URLSearchParams()
  if (params?.page != null) search.set('page', String(params.page))
  if (params?.per_page != null) search.set('per_page', String(params.per_page))
  const q = search.toString()
  return get(`/api/v1/templates${q ? `?${q}` : ''}`)
}

export async function getPublicTemplates(): Promise<{ templates: Template[] }> {
  return get('/api/v1/templates/public')
}

export async function createTemplate(input: CreateTemplateInput): Promise<{ template: Template }> {
  return post('/api/v1/templates', input)
}

export async function getTemplate(id: string): Promise<{ template: Template }> {
  return get(`/api/v1/templates/${id}`)
}

export async function updateTemplate(id: string, input: UpdateTemplateInput): Promise<{ template: Template }> {
  return put(`/api/v1/templates/${id}`, input)
}

export async function deleteTemplate(id: string): Promise<void> {
  await del(`/api/v1/templates/${id}`)
}
