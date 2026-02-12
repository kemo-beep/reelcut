export interface TranscriptWord {
  id: string
  word: string
  start_time: number
  end_time: number
  confidence: number
}

export interface TranscriptSegment {
  id: string
  start_time: number
  end_time: number
  text: string
  confidence: number
  speaker_id?: number | null
  words?: TranscriptWord[]
}

export interface Transcription {
  id: string
  video_id: string
  language: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  error_message?: string | null
  segments?: TranscriptSegment[]
  created_at: string
}
