import { Link } from '@tanstack/react-router'
import { Plus, Sparkles } from 'lucide-react'
import { Button } from '../../ui/button'

export interface VideoClipsPageHeaderProps {
  videoId: string
  videoFilename: string
  onSuggestClips: () => void
  isSuggesting: boolean
}

export function VideoClipsPageHeader({
  videoId,
  videoFilename,
  onSuggestClips,
  isSuggesting,
}: VideoClipsPageHeaderProps) {
  return (
    <div className="flex items-center justify-between">
      <div>
        <Link to="/dashboard/videos" className="text-sm text-[var(--app-accent)] hover:underline">
          Videos
        </Link>
        <span className="mx-2 text-[var(--app-fg-muted)]">/</span>
        <Link
          to="/dashboard/videos/$videoId"
          params={{ videoId }}
          className="text-sm text-[var(--app-accent)] hover:underline"
        >
          {videoFilename}
        </Link>
      </div>
      <div className="flex gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={onSuggestClips}
          disabled={isSuggesting}
          className="border-[var(--app-border)]"
        >
          <Sparkles size={16} className="mr-1" />
          {isSuggesting ? 'Analyzing…' : 'AI suggest clips'}
        </Button>
        <Link to="/dashboard/clips" search={{ videoId }}>
          <Button size="sm" className="bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]">
            <Plus size={16} className="mr-1" />
            Create clip
          </Button>
        </Link>
      </div>
    </div>
  )
}
