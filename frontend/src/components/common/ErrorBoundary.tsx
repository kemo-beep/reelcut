import { Component, type ErrorInfo, type ReactNode } from 'react'
import { Button } from '../ui/button'
import { AlertCircle } from 'lucide-react'

interface Props {
  children: ReactNode
  fallback?: ReactNode
  onError?: (error: Error, errorInfo: ErrorInfo) => void
}

interface State {
  hasError: boolean
  error: Error | null
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false, error: null }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    this.props.onError?.(error, errorInfo)
  }

  render() {
    if (this.state.hasError && this.state.error) {
      if (this.props.fallback) return this.props.fallback
      return (
        <div
          className="flex flex-col items-center justify-center rounded-xl border border-[var(--app-destructive)]/30 bg-[var(--app-destructive-muted)] p-8 text-center"
          role="alert"
        >
          <AlertCircle size={40} className="text-[var(--app-destructive)] mb-4" />
          <h2 className="text-lg font-semibold text-[var(--app-fg)]">Something went wrong</h2>
          <p className="mt-2 text-caption text-[var(--app-fg-muted)] max-w-md">
            {this.state.error.message}
          </p>
          <Button
            variant="outline"
            size="sm"
            className="mt-4"
            onClick={() => this.setState({ hasError: false, error: null })}
          >
            Try again
          </Button>
        </div>
      )
    }
    return this.props.children
  }
}
