import { useState, useCallback, useEffect } from 'react'
import { Search } from 'lucide-react'
import { Input } from '../ui/input'
import { cn } from '../../lib/utils'

export interface SearchBarProps {
  value?: string
  onChange: (value: string) => void
  placeholder?: string
  debounceMs?: number
  className?: string
}

export function SearchBar({
  value: controlledValue,
  onChange,
  placeholder = 'Search…',
  debounceMs = 300,
  className,
}: SearchBarProps) {
  const [local, setLocal] = useState(controlledValue ?? '')
  const isControlled = controlledValue !== undefined

  useEffect(() => {
    if (isControlled) setLocal(controlledValue)
  }, [isControlled, controlledValue])

  useEffect(() => {
    if (debounceMs <= 0) {
      onChange(local)
      return
    }
    const t = setTimeout(() => onChange(local), debounceMs)
    return () => clearTimeout(t)
  }, [local, debounceMs, onChange])

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const v = e.target.value
      setLocal(v)
      if (debounceMs <= 0) onChange(v)
    },
    [debounceMs, onChange]
  )

  return (
    <div className={cn('relative', className)}>
      <Search
        size={18}
        className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--app-fg-muted)] pointer-events-none"
        aria-hidden
      />
      <Input
        type="search"
        value={local}
        onChange={handleChange}
        placeholder={placeholder}
        className="pl-9 h-10 border-[var(--app-border)] bg-[var(--app-bg)] text-[var(--app-fg)] placeholder:text-[var(--app-fg-subtle)]"
        aria-label="Search"
      />
    </div>
  )
}
