import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useEffect } from 'react'
import { useAuthStore } from '../stores/authStore'
import { Hero } from '../components/landing/Hero'
import { Features } from '../components/landing/Features'
import { HowItWorks } from '../components/landing/HowItWorks'
import { Testimonials } from '../components/landing/Testimonials'
import { Footer } from '../components/landing/Footer'

export const Route = createFileRoute('/')({ component: HomePage })

function HomePage() {
  const accessToken = useAuthStore((s) => s.accessToken)
  const navigate = useNavigate()

  useEffect(() => {
    if (accessToken) {
      navigate({ to: '/dashboard' })
    }
  }, [accessToken, navigate])

  if (accessToken) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[var(--app-bg)]">
        <div className="flex flex-col items-center gap-3">
          <div className="h-8 w-8 animate-pulse rounded-full bg-[var(--app-accent-muted)]" />
          <p className="text-caption">Redirecting to dashboard...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-[var(--app-bg)] selection:bg-[var(--app-accent)] selection:text-white">
      {/* Main Content Wrapper for Footer Reveal Effect */}
      <div className="relative z-10 bg-[var(--app-bg)] mb-[400px] shadow-2xl">
        <Hero />
        <Features />
        <HowItWorks />
        <Testimonials />

        {/* CTA Section */}
        <section className="py-32 text-center bg-[#0a0a0b] text-white">
          <div className="mx-auto max-w-4xl px-6">
            <h2 className="text-display mb-8 text-5xl font-bold tracking-tight md:text-7xl">
              Ready to go viral?
            </h2>
            <p className="text-xl text-gray-400 mx-auto mb-12 max-w-2xl">
              Join the automated video revolution and start creating content that converts.
            </p>
            <a
              href="/auth/register"
              className="inline-flex items-center justify-center rounded-full bg-white px-10 py-5 text-xl font-bold text-black shadow-lg transition-transform hover:scale-105 hover:bg-gray-100"
            >
              Get Started for Free
            </a>
          </div>
        </section>
      </div>

      <Footer />
    </div>
  )
}
