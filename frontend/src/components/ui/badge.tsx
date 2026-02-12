import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const badgeVariants = cva(
  'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium transition-[var(--motion-duration-fast)]',
  {
    variants: {
      variant: {
        default:
          'border border-[var(--app-border)] bg-[var(--app-bg-overlay)] text-[var(--app-fg-muted)]',
        success:
          'border border-[var(--app-success)]/30 bg-[var(--app-success)]/15 text-[var(--app-success)]',
        warning:
          'border border-[var(--app-warning)]/30 bg-[var(--app-warning)]/15 text-[var(--app-warning)]',
        destructive:
          'border border-[var(--app-destructive)]/30 bg-[var(--app-destructive-muted)] text-[var(--app-destructive)]',
        accent:
          'border border-[var(--app-accent)]/30 bg-[var(--app-accent-muted)] text-[var(--app-accent)]',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  }
)

export interface BadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <span
      className={cn(badgeVariants({ variant }), className)}
      {...props}
    />
  )
}

export { Badge, badgeVariants }
