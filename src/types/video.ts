export interface Video {
  id: string
  project_id: string
  user_id: string
  original_filename: string
  storage_path: string
  thumbnail_url?: string | null
  /** Presigned URL for list display (from API list endpoint) */
  thumbnail_display_url?: string | null
  duration_seconds?: number | null
  width?: number | null
  height?: number | null
  file_size_bytes?: number | null
  status: 'uploading' | 'processing' | 'ready' | 'failed'
  created_at: string
  updated_at: string
}
