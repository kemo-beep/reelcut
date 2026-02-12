import { get, post } from './client'

export interface Subscription {
  id: string
  user_id: string
  tier: string
  status: string
  stripe_subscription_id?: string | null
  current_period_start?: string | null
  current_period_end?: string | null
  created_at: string
  updated_at: string
}

export async function getMySubscription(): Promise<Subscription> {
  return get<Subscription>('/api/v1/users/me/subscription')
}

export async function createSubscription(tier: string): Promise<Subscription> {
  return post<Subscription>('/api/v1/subscriptions/create', { tier })
}

export async function cancelSubscription(): Promise<{ message: string }> {
  return post<{ message: string }>('/api/v1/subscriptions/cancel', {})
}

export async function updateSubscription(tier: string): Promise<Subscription> {
  return post<Subscription>('/api/v1/subscriptions/update', { tier })
}
