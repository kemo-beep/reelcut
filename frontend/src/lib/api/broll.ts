import { get } from './client'
import type { BrollAsset } from './clips'

export async function listBrollAssets(params?: { project_id?: string; limit?: number; offset?: number }): Promise<{ assets: BrollAsset[]; total: number }> {
  const search = new URLSearchParams()
  if (params?.project_id) search.set('project_id', params.project_id)
  if (params?.limit != null) search.set('limit', String(params.limit))
  if (params?.offset != null) search.set('offset', String(params.offset))
  const q = search.toString()
  return get(`/api/v1/broll/assets${q ? `?${q}` : ''}`)
}

export async function uploadBrollAsset(file: File, projectId?: string): Promise<{ asset: BrollAsset }> {
  const form = new FormData()
  form.append('file', file)
  if (projectId) form.append('project_id', projectId)
  const base = (typeof import.meta !== 'undefined' && import.meta.env?.VITE_API_URL) ?? 'http://localhost:8080'
  const token = localStorage.getItem('access_token')
  const res = await fetch(`${base.replace(/\/$/, '')}/api/v1/broll/assets`, {
    method: 'POST',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    body: form,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err?.error?.message ?? `Upload failed: ${res.status}`)
  }
  return res.json()
}
