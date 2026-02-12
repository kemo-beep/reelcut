import { useEffect, useRef, useState, useCallback } from 'react'
import { useAuthStore } from '../stores/authStore'

const MIN_RECONNECT_MS = 2000
const MAX_RECONNECT_MS = 30000
const INITIAL_RECONNECT_MS = 3000

function getWsUrl(): string {
  const base = (typeof import.meta !== 'undefined' && import.meta.env?.VITE_API_URL) ?? 'http://localhost:8080'
  const url = base.replace(/\/$/, '').replace(/^http/, 'ws')
  return `${url}/ws`
}

export interface UseWebSocketOptions {
  onMessage?: (data: unknown) => void
  onOpen?: () => void
  onClose?: () => void
  onError?: (event: Event) => void
  enabled?: boolean
}

/**
 * Connects to the backend WebSocket for real-time job updates (e.g. upload, transcription, render progress).
 * Sends token via query param for auth. Uses refs for callbacks so we don't reconnect on every parent render.
 */
export function useWebSocket(options: UseWebSocketOptions = {}) {
  const { enabled = true } = options
  const [connected, setConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const reconnectDelayRef = useRef(INITIAL_RECONNECT_MS)
  const getAccessToken = useAuthStore((s) => s.getAccessToken)

  const onMessageRef = useRef(options.onMessage)
  const onOpenRef = useRef(options.onOpen)
  const onCloseRef = useRef(options.onClose)
  const onErrorRef = useRef(options.onError)
  onMessageRef.current = options.onMessage
  onOpenRef.current = options.onOpen
  onCloseRef.current = options.onClose
  onErrorRef.current = options.onError

  const connect = useCallback(() => {
    const token = getAccessToken()
    if (!token || !enabled) return
    if (wsRef.current?.readyState === WebSocket.OPEN) return
    const url = `${getWsUrl()}?token=${encodeURIComponent(token)}`
    const ws = new WebSocket(url)
    wsRef.current = ws

    ws.onopen = () => {
      reconnectDelayRef.current = INITIAL_RECONNECT_MS
      setConnected(true)
      onOpenRef.current?.()
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data as string) as unknown
        onMessageRef.current?.(data)
      } catch {
        onMessageRef.current?.(event.data)
      }
    }

    ws.onclose = () => {
      setConnected(false)
      wsRef.current = null
      onCloseRef.current?.()
      if (enabled && token) {
        const delay = reconnectDelayRef.current
        reconnectDelayRef.current = Math.min(MAX_RECONNECT_MS, delay * 2)
        reconnectTimeoutRef.current = setTimeout(() => connect(), Math.max(MIN_RECONNECT_MS, delay))
      }
    }

    ws.onerror = (event) => {
      onErrorRef.current?.(event)
    }
  }, [enabled, getAccessToken])

  useEffect(() => {
    connect()
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
        reconnectTimeoutRef.current = null
      }
      if (wsRef.current) {
        wsRef.current.close()
        wsRef.current = null
      }
      setConnected(false)
    }
  }, [connect])

  return { connected }
}
