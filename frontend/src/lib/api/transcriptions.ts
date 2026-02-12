import { get, post, put } from './client'
import type { Transcription, TranscriptSegment } from '../../types'
import { ApiError } from '../../types'

export interface CreateTranscriptionInput {
  language?: string
  enable_diarization?: boolean
}

export async function createTranscription(
  videoId: string,
  input?: CreateTranscriptionInput
): Promise<{ transcription: Transcription }> {
  return post(`/api/v1/transcriptions/videos/${videoId}`, input ?? {})
}

export async function getTranscription(id: string): Promise<{ transcription: Transcription }> {
  return get(`/api/v1/transcriptions/${id}`)
}

/** Get transcription by video ID. Treats 404 as "no transcription" and returns null so the UI can show empty state. */
export async function getTranscriptionByVideoId(videoId: string): Promise<{ transcription: Transcription | null }> {
  try {
    return await get<{ transcription: Transcription | null }>(`/api/v1/transcriptions/videos/${videoId}`)
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) {
      return { transcription: null }
    }
    throw e
  }
}

export async function updateSegment(
  transcriptionId: string,
  segmentId: string,
  body: Partial<Pick<TranscriptSegment, 'text'>>
): Promise<unknown> {
  return put(`/api/v1/transcriptions/${transcriptionId}/segments/${segmentId}`, body)
}
