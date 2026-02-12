import { get, post } from './client'

export interface ClipSuggestion {
  start_time: number
  end_time: number
  duration?: number
  virality_score?: number
  reason?: string
  highlights?: string[]
  transcript?: string
}

export interface SuggestClipsResponse {
  suggestions: ClipSuggestion[]
}

export async function analyzeVideo(videoId: string): Promise<unknown> {
  return post(`/api/v1/analysis/videos/${videoId}`)
}

export async function getAnalysisByVideoId(videoId: string): Promise<unknown> {
  return get(`/api/v1/analysis/videos/${videoId}`)
}

export async function suggestClips(
  videoId: string,
  body?: Record<string, unknown>
): Promise<SuggestClipsResponse> {
  return post<SuggestClipsResponse>(`/api/v1/analysis/videos/${videoId}/suggest-clips`, body)
}
