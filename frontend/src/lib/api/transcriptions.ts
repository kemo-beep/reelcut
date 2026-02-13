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

/** Get transcription by video ID. Optional language for multi-language (ISO 639-1). Treats 404 as "no transcription" and returns null. */
export async function getTranscriptionByVideoId(
  videoId: string,
  opts?: { language?: string }
): Promise<{ transcription: Transcription | null }> {
  try {
    const url = opts?.language
      ? `/api/v1/transcriptions/videos/${videoId}?language=${encodeURIComponent(opts.language)}`
      : `/api/v1/transcriptions/videos/${videoId}`
    return await get<{ transcription: Transcription | null }>(url)
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) {
      return { transcription: null }
    }
    throw e
  }
}

/** List all completed transcriptions for a video (for caption language selection). */
export async function listTranscriptionsByVideo(videoId: string): Promise<{ transcriptions: Transcription[] }> {
  return get(`/api/v1/transcriptions/videos/${videoId}/list`)
}

/** Translate a completed transcription to another language (same timestamps). */
export async function translateTranscription(
  transcriptionId: string,
  body: { target_language: string }
): Promise<{ transcription: Transcription }> {
  return post(`/api/v1/transcriptions/${transcriptionId}/translate`, body)
}

export async function updateSegment(
  transcriptionId: string,
  segmentId: string,
  body: Partial<Pick<TranscriptSegment, 'text'> & { start_time?: number; end_time?: number }>
): Promise<unknown> {
  return put(`/api/v1/transcriptions/${transcriptionId}/segments/${segmentId}`, body)
}
