import { cn } from '@/lib/utils'

function Skeleton({
  className,
  ...props
}: React.ComponentProps<'div'>) {
  return (
    <div
      className={cn(
        'animate-pulse rounded-lg bg-[var(--app-border-strong)]',
        className
      )}
      {...props}
    />
  )
}

export { Skeleton }
