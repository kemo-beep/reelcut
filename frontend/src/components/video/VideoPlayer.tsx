import { useRef, useEffect, forwardRef, useImperativeHandle } from 'react'
import { cn } from '../../lib/utils'

export interface VideoPlayerProps {
  src: string
  className?: string
  poster?: string
  /** Seek to this time once the video has loaded (e.g. when switching clips) */
  initialTime?: number
  onTimeUpdate?: (currentTime: number) => void
  onPlay?: () => void
  onPause?: () => void
}

export interface VideoPlayerHandle {
  seek: (time: number) => void
  getCurrentTime: () => number
}

export const VideoPlayer = forwardRef<VideoPlayerHandle, VideoPlayerProps>(function VideoPlayer(
  { src, className, poster, initialTime, onTimeUpdate, onPlay, onPause },
  ref
) {
  const videoRef = useRef<HTMLVideoElement>(null)
  useImperativeHandle(ref, () => ({
    seek: (time: number) => {
      const el = videoRef.current
      if (el) el.currentTime = time
    },
    getCurrentTime: () => videoRef.current?.currentTime ?? 0,
  }))

  useEffect(() => {
    const el = videoRef.current
    if (!el) return
    const handleTimeUpdate = () => onTimeUpdate?.(el.currentTime)
    el.addEventListener('timeupdate', handleTimeUpdate)
    return () => el.removeEventListener('timeupdate', handleTimeUpdate)
  }, [onTimeUpdate])

  useEffect(() => {
    const el = videoRef.current
    if (!el || src === '' || initialTime == null || initialTime === 0) return
    const seekWhenReady = () => {
      el.currentTime = initialTime
      onTimeUpdate?.(initialTime)
    }
    if (el.readyState >= 2) seekWhenReady()
    else el.addEventListener('loadedmetadata', seekWhenReady, { once: true })
    return () => el.removeEventListener('loadedmetadata', seekWhenReady)
  }, [src, initialTime, onTimeUpdate])

  return (
    <div className={cn('overflow-hidden rounded-xl bg-[var(--app-bg)]', className)}>
      <video
        ref={videoRef}
        src={src}
        poster={poster}
        controls
        className="w-full max-h-[70vh]"
        preload="metadata"
        onPlay={onPlay}
        onPause={onPause}
      >
        Your browser does not support the video tag.
      </video>
    </div>
  )
})
