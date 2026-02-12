import { createFileRoute } from '@tanstack/react-router'
import { VideoClipsPage } from '../../../../components/dashboard/VideoClipsPage'

export const Route = createFileRoute('/dashboard/videos/$videoId/clips')({
  component: VideoClipsPage,
})
