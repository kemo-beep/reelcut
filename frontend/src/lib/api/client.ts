import { useAuthStore } from '../../stores/authStore'
import { ApiError } from '../../types'
import type { AuthResponse } from '../../types'

function getBaseUrl(): string {
  const url = import.meta.env.VITE_API_URL
  if (typeof url === 'string' && url) return url.replace(/\/$/, '')
  return 'http://localhost:8080'
}

/** Single in-flight refresh so concurrent 401s share one refresh call. */
let refreshPromise: Promise<boolean> | null = null

async function runRefresh(): Promise<boolean> {
  const refresh = useAuthStore.getState().refreshToken
  if (!refresh) {
    useAuthStore.getState().clearAuth()
    return false
  }
  try {
    const r = await fetch(`${getBaseUrl()}/api/v1/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refresh }),
    })
    const text = await r.text()
    let data: unknown = null
    if (text) {
      try {
        data = JSON.parse(text)
      } catch {
        // leave data null
      }
    }
    if (!r.ok) {
      useAuthStore.getState().clearAuth()
      return false
    }
    const res = data as AuthResponse
    if (res?.user && res?.token?.access_token) {
      useAuthStore.getState().setAuth(res.user, res.token.access_token, res.token.refresh_token)
      return true
    }
    useAuthStore.getState().clearAuth()
    return false
  } catch {
    useAuthStore.getState().clearAuth()
    return false
  }
}

async function getOrRunRefresh(): Promise<boolean> {
  if (refreshPromise) return refreshPromise
  refreshPromise = runRefresh().finally(() => {
    refreshPromise = null
  })
  return refreshPromise
}

export interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  body?: unknown
  headers?: Record<string, string>
  skipAuth?: boolean
}

export async function request<T>(
  path: string,
  options: RequestOptions = {},
  retried = false
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
    const isAuthPath = path.includes('auth/refresh')
    if (res.status === 401 && !skipAuth && !isAuthPath && !retried) {
      const refreshed = await getOrRunRefresh()
      if (refreshed) return request<T>(path, options, true)
      useAuthStore.getState().onSessionExpired?.()
      const err = data as {
        error?: { code?: string; message?: string; details?: Array<{ field?: string; message: string }> }
      } | null
      const code = err?.error?.code ?? 'UNAUTHORIZED'
      const message = err?.error?.message ?? 'Session expired. Please log in again.'
      const details = err?.error?.details ?? []
      throw new ApiError(code, message, res.status, details)
    }
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
