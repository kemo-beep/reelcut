import { get, put, del } from './client'
import type { User } from '../../types'

export interface UserProfileUpdate {
  full_name?: string | null
  avatar_url?: string | null
}

export async function getProfile(): Promise<{ user: User }> {
  return get('/api/v1/users/me')
}

export async function updateProfile(body: UserProfileUpdate): Promise<{ user: User }> {
  return put('/api/v1/users/me', body)
}

export async function getUsageStats(params?: { page?: number; per_page?: number }): Promise<unknown> {
  const search = new URLSearchParams()
  if (params?.page != null) search.set('page', String(params.page))
  if (params?.per_page != null) search.set('per_page', String(params.per_page))
  const q = search.toString()
  return get(`/api/v1/users/me/usage${q ? `?${q}` : ''}`)
}

export async function deleteAccount(): Promise<void> {
  await del('/api/v1/users/me')
}
