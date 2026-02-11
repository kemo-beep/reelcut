export interface User {
  id: string
  email: string
  full_name?: string | null
  avatar_url?: string | null
  subscription_tier: string
  credits_remaining: number
  email_verified: boolean
  created_at: string
  updated_at: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_in: number
  expires_at: string
}

export interface AuthResponse {
  user: User
  token: TokenPair
}
