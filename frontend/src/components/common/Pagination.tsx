import { Button } from '../ui/button'
import { ChevronLeft, ChevronRight } from 'lucide-react'

export interface PaginationProps {
  page: number
  perPage: number
  total: number
  onPageChange: (page: number) => void
  className?: string
}

export function Pagination({
  page,
  perPage,
  total,
  onPageChange,
  className = '',
}: PaginationProps) {
  const totalPages = Math.max(1, Math.ceil(total / perPage))
  const from = total === 0 ? 0 : (page - 1) * perPage + 1
  const to = Math.min(page * perPage, total)

  return (
    <div
      className={`flex items-center justify-between gap-4 ${className}`}
      aria-label="Pagination"
    >
      <p className="text-caption text-[var(--app-fg-muted)]">
        Showing {from}–{to} of {total}
      </p>
      <div className="flex items-center gap-1">
        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange(page - 1)}
          disabled={page <= 1}
          className="h-9 w-9 p-0"
          aria-label="Previous page"
        >
          <ChevronLeft size={18} />
        </Button>
        <span className="px-2 text-sm text-[var(--app-fg-muted)]" aria-current="page">
          Page {page} of {totalPages}
        </span>
        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange(page + 1)}
          disabled={page >= totalPages}
          className="h-9 w-9 p-0"
          aria-label="Next page"
        >
          <ChevronRight size={18} />
        </Button>
      </div>
    </div>
  )
}
