export interface Template {
  id: string
  user_id?: string | null
  name: string
  category?: string | null
  is_public: boolean
  preview_url?: string | null
  style_config: Record<string, unknown>
  usage_count: number
  created_at: string
  updated_at: string
}
