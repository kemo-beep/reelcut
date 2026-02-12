export interface ProcessingJob {
  id: string
  user_id: string
  job_type: string
  entity_type: string
  entity_id: string
  priority: number
  status: string
  progress: number
  error_message?: string | null
  retry_count: number
  max_retries: number
  metadata?: Record<string, unknown> | null
  started_at?: string | null
  completed_at?: string | null
  created_at: string
  updated_at: string
}
