import { Skeleton } from '../../ui/skeleton'

/**
 * Skeleton for the initial page load (video loading).
 * Mirrors the layout: header, view tabs, and content area.
 */
export function VideoClipsPageSkeleton() {
  return (
    <div className="space-y-6" aria-busy="true" aria-label="Loading clips page">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <Skeleton className="h-4 w-14" />
          <Skeleton className="h-4 w-4 rounded-full" />
          <Skeleton className="h-4 w-32" />
        </div>
        <div className="flex gap-2">
          <Skeleton className="h-8 w-28" />
          <Skeleton className="h-8 w-24" />
        </div>
      </div>
      <div className="flex items-center gap-2">
        <Skeleton className="h-9 w-20" />
        <Skeleton className="h-9 w-24" />
      </div>
      <div className="space-y-4">
        <Skeleton className="aspect-video w-full rounded-xl" />
        <Skeleton className="h-12 w-full rounded-lg" />
      </div>
    </div>
  )
}

/**
 * Skeleton for the clips view while clips are loading.
 * Mirrors ClipsViewContent: player area, timeline, and clip strips.
 */
export function ClipsViewSkeleton() {
  return (
    <div className="space-y-6" aria-busy="true" aria-label="Loading clips">
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_minmax(320px,400px)] items-start">
        <div className="space-y-4">
          <Skeleton className="aspect-video w-full rounded-xl" />
          <Skeleton className="h-12 w-full rounded-lg" />
        </div>
        <Skeleton className="min-h-[200px] rounded-xl lg:min-h-[300px]" />
      </div>
      <section className="w-full">
        <Skeleton className="mb-3 h-4 w-24" />
        <div className="flex gap-4 overflow-hidden">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-[280px] w-[180px] shrink-0 rounded-xl" />
          ))}
        </div>
      </section>
    </div>
  )
}
