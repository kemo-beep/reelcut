import { useCallback } from 'react'
import { useAuthStore } from '../stores/authStore'
import * as authApi from '../lib/api/auth'

export function useAuth() {
  const user = useAuthStore((s) => s.user)
  const accessToken = useAuthStore((s) => s.accessToken)
  const setAuth = useAuthStore((s) => s.setAuth)
  const clearAuth = useAuthStore((s) => s.clearAuth)

  const login = useCallback(
    async (email: string, password: string) => {
      const res = await authApi.login({ email, password })
      setAuth(res.user, res.token.access_token, res.token.refresh_token)
      return res
    },
    [setAuth]
  )

  const register = useCallback(
    async (email: string, password: string, fullName?: string) => {
      const res = await authApi.register({
        email,
        password,
        ...(fullName != null && fullName !== '' ? { full_name: fullName } : {}),
      })
      setAuth(res.user, res.token.access_token, res.token.refresh_token)
      return res
    },
    [setAuth]
  )

  const logout = useCallback(() => {
    clearAuth()
  }, [clearAuth])

  const refreshToken = useCallback(async () => {
    const refresh = useAuthStore.getState().refreshToken
    if (!refresh) throw new Error('No refresh token')
    const res = await authApi.refreshToken(refresh)
    setAuth(res.user, res.token.access_token, res.token.refresh_token)
    return res
  }, [setAuth])

  return {
    user,
    login,
    register,
    logout,
    refreshToken,
    isAuthenticated: !!accessToken,
  }
}
