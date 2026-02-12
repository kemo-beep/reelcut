import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/dashboard/videos/$videoId/transcription')({
  beforeLoad: ({ params }) => {
    throw redirect({
      to: '/dashboard/videos/$videoId',
      params: { videoId: params.videoId },
    })
  },
  component: () => null,
})
