import { get } from './client'

export async function getCaptionFonts(): Promise<{ fonts: string[] }> {
  return get('/api/v1/config/caption-fonts')
}

export interface ExportPreset {
  id: string
  name: string
  description: string
  aspect_ratio: string
  width: number
  height: number
  video_bitrate_kbps: number
  audio_bitrate_kbps: number
  fps: number
}

export async function getExportPresets(): Promise<{ presets: ExportPreset[] }> {
  return get('/api/v1/config/export-presets')
}
