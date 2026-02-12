import { useEffect } from 'react'
import {
  HeadContent,
  Scripts,
  createRootRouteWithContext,
  useNavigate,
} from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'
import { toast } from 'sonner'

import Header from '../components/Header'
import { ThemeProvider } from '../components/ThemeProvider'
import { useAuthStore } from '../stores/authStore'

import StoreDevtools from '../lib/demo-store-devtools'

import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'

import appCss from '../styles.css?url'
import { Toaster } from 'sonner'

import type { QueryClient } from '@tanstack/react-query'

const THEME_SCRIPT = `(function(){
  try {
    var raw = localStorage.getItem('reelcut-theme');
    var t = 'system';
    if (raw) {
      var data = JSON.parse(raw);
      if (data && data.state && data.state.theme) t = data.state.theme;
    }
    var dark = t === 'dark' || (t === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches);
    document.documentElement.classList.toggle('dark', dark);
  } catch (e) {}
})();`

interface MyRouterContext {
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  // Disable SSR app-wide to avoid 500 "HTTPError" during server render (e.g. Cloudflare/workerd or backend unreachable).
  ssr: false,
  head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        title: 'Reelcut',
      },
    ],
    links: [
      {
        rel: 'preconnect',
        href: 'https://fonts.googleapis.com',
      },
      {
        rel: 'preconnect',
        href: 'https://fonts.gstatic.com',
        crossOrigin: 'anonymous',
      },
      {
        rel: 'stylesheet',
        href: 'https://fonts.googleapis.com/css2?family=DM+Sans:ital,opsz,wght@0,9..40,400;0,9..40,500;0,9..40,600;0,9..40,700;1,9..40,400&family=Inter:wght@400;500;600;700&display=swap',
      },
      {
        rel: 'stylesheet',
        href: appCss,
      },
    ],
  }),

  shellComponent: RootDocument,
})

function RootDocument({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate()
  useEffect(() => {
    useAuthStore.getState().setOnSessionExpired(() => {
      toast.error('Session expired. Please log in again.')
      navigate({ to: '/login' })
    })
    return () => useAuthStore.getState().setOnSessionExpired(null)
  }, [navigate])

  return (
    <html lang="en">
      <head>
        <script dangerouslySetInnerHTML={{ __html: THEME_SCRIPT }} />
        <HeadContent />
      </head>
      <body className="min-h-screen bg-[var(--app-bg)] text-[var(--app-fg)] font-sans antialiased">
        <ThemeProvider>
          <Header />
          {children}
          <Toaster
          position="bottom-right"
          toastOptions={{
            style: {
              background: 'var(--app-bg-raised)',
              border: '1px solid var(--app-border)',
              color: 'var(--app-fg)',
            },
          }}
        />
        <TanStackDevtools
          config={{
            position: 'bottom-right',
          }}
          plugins={[
            {
              name: 'Tanstack Router',
              render: <TanStackRouterDevtoolsPanel />,
            },
            StoreDevtools,
            TanStackQueryDevtools,
          ]}
        />
        </ThemeProvider>
        <Scripts />
      </body>
    </html>
  )
}
