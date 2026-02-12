import { get, put, del } from './client'
import type { User } from '../../types'
import { useAuthStore } from '../../stores/authStore'

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

export async function changePassword(currentPassword: string, newPassword: string): Promise<void> {
  await put('/api/v1/users/me/password', {
    current_password: currentPassword,
    new_password: newPassword,
  })
}

export async function uploadAvatar(file: File): Promise<{ user: User }> {
  const formData = new FormData()
  formData.append('file', file)
  const token = useAuthStore.getState().getAccessToken()
  const base = (typeof import.meta !== 'undefined' && import.meta.env?.VITE_API_URL) ?? 'http://localhost:8080'
  const url = `${base.replace(/\/$/, '')}/api/v1/users/me/avatar`
  const res = await fetch(url, {
    method: 'POST',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    body: formData,
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data?.error?.message ?? 'Upload failed')
  return data
}

export async function deleteAccount(): Promise<void> {
  await del('/api/v1/users/me')
}
