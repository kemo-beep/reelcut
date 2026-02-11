import { post } from './client'
import type { AuthResponse, User } from '../../types'

export interface RegisterInput {
  email: string
  password: string
  full_name?: string
}

export interface LoginInput {
  email: string
  password: string
}

export interface ResetPasswordInput {
  token: string
  new_password: string
}

export async function register(input: RegisterInput): Promise<AuthResponse> {
  return post<AuthResponse>('/api/v1/auth/register', input, { skipAuth: true })
}

export async function login(input: LoginInput): Promise<AuthResponse> {
  return post<AuthResponse>('/api/v1/auth/login', input, { skipAuth: true })
}

export async function refreshToken(refreshToken: string): Promise<AuthResponse> {
  return post<AuthResponse>(
    '/api/v1/auth/refresh',
    { refresh_token: refreshToken },
    { skipAuth: true }
  )
}

export async function forgotPassword(email: string): Promise<{ message: string }> {
  return post<{ message: string }>('/api/v1/auth/forgot-password', { email }, { skipAuth: true })
}

export async function resetPassword(input: ResetPasswordInput): Promise<{ message: string }> {
  return post<{ message: string }>('/api/v1/auth/reset-password', input, { skipAuth: true })
}

export async function verifyEmail(token: string): Promise<{ message: string }> {
  return post<{ message: string }>('/api/v1/auth/verify-email', { token }, { skipAuth: true })
}
