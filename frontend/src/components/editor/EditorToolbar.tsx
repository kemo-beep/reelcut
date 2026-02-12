import { Play, Pause, SkipBack, SkipForward } from 'lucide-react'
import { Button } from '../ui/button'
import { cn } from '../../lib/utils'

export interface EditorToolbarProps {
  isPlaying: boolean
  onPlayPause: () => void
  onStepBack?: () => void
  onStepForward?: () => void
  className?: string
}

export function EditorToolbar({
  isPlaying,
  onPlayPause,
  onStepBack,
  onStepForward,
  className,
}: EditorToolbarProps) {
  return (
    <div
      className={cn(
        'flex items-center gap-2 rounded-lg border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-2',
        className
      )}
    >
      {onStepBack && (
        <Button variant="outline" size="icon" onClick={onStepBack} aria-label="Step back">
          <SkipBack size={18} />
        </Button>
      )}
      <Button
        variant="outline"
        size="icon"
        onClick={onPlayPause}
        aria-label={isPlaying ? 'Pause' : 'Play'}
      >
        {isPlaying ? <Pause size={18} /> : <Play size={18} />}
      </Button>
      {onStepForward && (
        <Button variant="outline" size="icon" onClick={onStepForward} aria-label="Step forward">
          <SkipForward size={18} />
        </Button>
      )}
    </div>
  )
}
