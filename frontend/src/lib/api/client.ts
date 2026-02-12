import { useAuthStore } from '../../stores/authStore'
import { ApiError } from '../../types'

function getBaseUrl(): string {
  const url = import.meta.env.VITE_API_URL
  if (typeof url === 'string' && url) return url.replace(/\/$/, '')
  return 'http://localhost:8080'
}

export interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  body?: unknown
  headers?: Record<string, string>
  skipAuth?: boolean
}

export async function request<T>(
  path: string,
  options: RequestOptions = {}
): Promise<T> {
  const { method = 'GET', body, headers = {}, skipAuth = false } = options
  const url = `${getBaseUrl()}${path.startsWith('/') ? path : `/${path}`}`

  const h: Record<string, string> = {
    'Content-Type': 'application/json',
    ...headers,
  }

  if (!skipAuth) {
    const token = useAuthStore.getState().getAccessToken()
    if (token) h['Authorization'] = `Bearer ${token}`
  }

  const init: RequestInit = {
    method,
    headers: h,
  }
  if (body !== undefined && method !== 'GET') {
    init.body = JSON.stringify(body)
  }

  const res = await fetch(url, init)
  const text = await res.text()
  let data: unknown = null
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      // leave data null
    }
  }

  if (!res.ok) {
    const err = data as {
      error?: { code?: string; message?: string; details?: Array<{ field?: string; message: string }> }
    } | null
    const code = err?.error?.code ?? 'UNKNOWN'
    const message = err?.error?.message ?? (res.statusText || 'Request failed')
    const details = err?.error?.details ?? []
    throw new ApiError(code, message, res.status, details)
  }

  return (data ?? null) as T
}

export function get<T>(path: string, options?: Omit<RequestOptions, 'method' | 'body'>): Promise<T> {
  return request<T>(path, { ...options, method: 'GET' })
}

export function post<T>(path: string, body?: unknown, options?: Omit<RequestOptions, 'method' | 'body'>): Promise<T> {
  return request<T>(path, { ...options, method: 'POST', body })
}

export function put<T>(path: string, body?: unknown, options?: Omit<RequestOptions, 'method' | 'body'>): Promise<T> {
  return request<T>(path, { ...options, method: 'PUT', body })
}

export function del<T>(path: string, options?: Omit<RequestOptions, 'method' | 'body'>): Promise<T> {
  return request<T>(path, { ...options, method: 'DELETE' })
}
