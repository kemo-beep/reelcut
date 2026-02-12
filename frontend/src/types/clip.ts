export interface ClipStyle {
  caption_enabled: boolean
  caption_font: string
  caption_size: number
  caption_color: string
  caption_bg_color?: string | null
  caption_position: 'top' | 'center' | 'bottom'
  caption_animation?: string | null
  brand_logo_url?: string | null
  brand_logo_position?: string | null
  overlay_template?: string | null
  background_music_url?: string | null
  background_music_volume: number
}

export interface Clip {
  id: string
  video_id: string
  user_id: string
  name: string
  start_time: number
  end_time: number
  duration_seconds?: number | null
  aspect_ratio: '9:16' | '1:1' | '16:9'
  virality_score?: number | null
  status: 'draft' | 'rendering' | 'ready' | 'failed'
  storage_path?: string | null
  thumbnail_url?: string | null
  is_ai_suggested: boolean
  style?: ClipStyle | null
  created_at: string
  updated_at: string
}
