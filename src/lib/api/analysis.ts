import { get, post } from './client'

export async function analyzeVideo(videoId: string): Promise<unknown> {
  return post(`/api/v1/analysis/videos/${videoId}`)
}

export async function getAnalysisByVideoId(videoId: string): Promise<unknown> {
  return get(`/api/v1/analysis/videos/${videoId}`)
}

export async function suggestClips(videoId: string, body?: Record<string, unknown>): Promise<unknown> {
  return post(`/api/v1/analysis/videos/${videoId}/suggest-clips`, body)
}
