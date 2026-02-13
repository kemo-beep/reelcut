import { useRef, MouseEvent } from 'react'
import { Link } from '@tanstack/react-router'
import { gsap } from 'gsap'
import { useGSAP } from '@gsap/react'
import { ArrowRight, Play, Sparkles } from 'lucide-react'

export function Hero() {
    const containerRef = useRef<HTMLDivElement>(null)
    const visualRef = useRef<HTMLDivElement>(null)
    const titleRef = useRef<HTMLHeadingElement>(null)

    // Mouse move effect for the visual (3D tilt) - Desktop only
    const { contextSafe } = useGSAP({ scope: containerRef })

    const handleMouseMove = contextSafe((e: MouseEvent) => {
        if (!visualRef.current || window.matchMedia('(hover: none)').matches) return
        const { clientX, clientY, currentTarget } = e
        const { width, height, left, top } = currentTarget.getBoundingClientRect()

        // Calculate mouse position relative to center of the container (-1 to 1)
        const x = (clientX - left - width / 2) / (width / 2)
        const y = (clientY - top - height / 2) / (height / 2)

        gsap.to(visualRef.current, {
            x: x * 20, // Move slightly
            y: y * 20,
            rotateY: x * 10, // Tilt
            rotateX: -y * 10,
            duration: 0.5,
            ease: 'power2.out',
        })
    })

    const handleMouseLeave = contextSafe(() => {
        if (!visualRef.current) return
        gsap.to(visualRef.current, {
            x: 0,
            y: 0,
            rotateY: 0,
            rotateX: 0,
            duration: 1,
            ease: 'elastic.out(1, 0.5)',
        })
    })

    useGSAP(
        () => {
            const tl = gsap.timeline()

            tl.from(titleRef.current, {
                y: 50,
                opacity: 0,
                duration: 1.2,
                ease: 'power4.out',
                stagger: 0.1 // This might not stagger if it's a single element ref, but keeps logic consistent
            })
                .from('.hero-element', {
                    y: 30,
                    opacity: 0,
                    duration: 1,
                    ease: 'power3.out',
                    stagger: 0.1
                }, '-=0.8')
                .from(visualRef.current, {
                    y: 50,
                    opacity: 0,
                    scale: 0.9,
                    duration: 1.5,
                    ease: 'power3.out'
                }, '-=1')
        },
        { scope: containerRef }
    )

    return (
        <section
            ref={containerRef}
            className="relative min-h-[90vh] flex flex-col items-center justify-center overflow-hidden bg-[var(--app-bg)] pt-24 pb-16 md:pt-32 md:pb-20"
            onMouseMove={handleMouseMove}
            onMouseLeave={handleMouseLeave}
        >
            {/* Dynamic Background Noise/Gradient */}
            <div className="absolute inset-0 pointer-events-none">
                <div className="absolute top-[-20%] left-[-10%] w-[100vw] h-[100vw] md:w-[70vw] md:h-[70vw] rounded-full bg-purple-500/20 blur-[80px] md:blur-[120px] mix-blend-screen animate-pulse opacity-50 dark:opacity-30" />
                <div className="absolute bottom-[-20%] right-[-10%] w-[100vw] h-[100vw] md:w-[70vw] md:h-[70vw] rounded-full bg-blue-500/20 blur-[80px] md:blur-[120px] mix-blend-screen animate-pulse delay-1000 opacity-50 dark:opacity-30" />
            </div>

            <div className="relative z-10 mx-auto max-w-7xl px-4 sm:px-6 text-center w-full">
                {/* Badge */}
                <div className="hero-element mb-6 inline-flex items-center gap-2 rounded-full border border-[var(--app-border)] bg-[var(--app-bg-raised)]/50 backdrop-blur-md px-3 py-1 md:px-4 md:py-1.5 shadow-sm transition-transform hover:scale-105">
                    <Sparkles className="h-3.5 w-3.5 md:h-4 md:w-4 text-[var(--app-accent)]" />
                    <span className="text-xs md:text-sm font-medium text-[var(--app-fg-muted)]">
                        v2.0 is now live
                    </span>
                </div>

                {/* Heading */}
                <h1
                    ref={titleRef}
                    className="text-display max-w-5xl mx-auto text-4xl sm:text-5xl md:text-7xl lg:text-8xl font-bold tracking-tight text-[var(--app-fg)] mb-6 md:mb-8 leading-[1.1]"
                >
                    Create Viral <br className="hidden sm:block" />
                    <span className="bg-gradient-to-r from-[var(--app-accent)] via-purple-400 to-cyan-400 bg-clip-text text-transparent">
                        Shorts Instantly
                    </span>
                </h1>

                {/* Subtitle */}
                <p className="hero-element text-body mx-auto mb-8 md:mb-10 max-w-2xl text-base sm:text-lg md:text-xl text-[var(--app-fg-muted)] px-4">
                    The AI-powered video editor that understands context. Turn long-form content into engaging short clips with one click.
                </p>

                {/* CTA Buttons */}
                <div className="hero-element flex flex-col sm:flex-row items-center justify-center gap-4 mb-12 md:mb-20 w-full sm:w-auto px-6 sm:px-0">
                    <Link
                        to="/auth/register"
                        className="group relative inline-flex items-center justify-center overflow-hidden rounded-full bg-[var(--app-accent)] px-8 py-3 md:py-4 text-base md:text-lg font-semibold text-white shadow-lg transition-all hover:scale-105 hover:shadow-[var(--app-accent)]/25 w-full sm:w-auto"
                    >
                        <span className="relative z-10 flex items-center gap-2">
                            Start Creating <ArrowRight className="h-5 w-5" />
                        </span>
                        <div className="absolute inset-0 bg-white/20 translate-y-full transition-transform group-hover:translate-y-0" />
                    </Link>
                    <Link
                        to="/auth/login"
                        className="group inline-flex items-center justify-center gap-2 rounded-full border border-[var(--app-border)] bg-[var(--app-bg-raised)]/50 backdrop-blur-sm px-8 py-3 md:py-4 text-base md:text-lg font-medium text-[var(--app-fg)] transition-all hover:bg-[var(--app-bg-overlay)] w-full sm:w-auto"
                    >
                        <Play className="h-5 w-5 fill-current transition-transform group-hover:scale-110" />
                        Watch Demo
                    </Link>
                </div>

                {/* 3D Visual Mockup */}
                <div
                    ref={visualRef}
                    className="relative mx-auto max-w-5xl rounded-2xl border border-[var(--app-border)] bg-[var(--app-bg-raised)]/80 backdrop-blur-xl p-1.5 md:p-2 shadow-2xl transition-shadow hover:shadow-[var(--app-accent)]/10"
                    style={{ transformStyle: 'preserve-3d', perspective: '1000px' }}
                >
                    <div className="overflow-hidden rounded-xl bg-[#0a0a0b] aspect-[16/9] relative group">
                        {/* Simulated UI */}
                        <div className="absolute top-0 left-0 right-0 h-10 md:h-14 border-b border-white/10 flex items-center px-4 md:px-6 justify-between bg-white/5 backdrop-blur-md">
                            <div className="flex gap-1.5 md:gap-2">
                                <div className="w-2.5 h-2.5 md:w-3 md:h-3 rounded-full bg-red-500/80" />
                                <div className="w-2.5 h-2.5 md:w-3 md:h-3 rounded-full bg-yellow-500/80" />
                                <div className="w-2.5 h-2.5 md:w-3 md:h-3 rounded-full bg-green-500/80" />
                            </div>
                            <div className="w-1/3 h-1.5 md:h-2 rounded-full bg-white/10" />
                        </div>

                        {/* Content Area */}
                        <div className="absolute inset-0 top-10 md:top-14 flex p-3 md:p-6 gap-3 md:gap-6">
                            <div className="hidden sm:block w-1/4 h-full rounded-lg bg-white/5 border border-white/5 p-2 md:p-4 space-y-2 md:space-y-3">
                                <div className="w-3/4 h-1.5 md:h-2 rounded-full bg-white/20" />
                                <div className="w-1/2 h-1.5 md:h-2 rounded-full bg-white/10" />
                                <div className="w-full h-16 md:h-24 rounded-lg bg-white/5 mt-2 md:mt-4" />
                                <div className="w-full h-16 md:h-24 rounded-lg bg-white/5" />
                            </div>
                            <div className="flex-1 h-full rounded-lg bg-gradient-to-br from-[var(--app-accent)]/20 to-purple-500/20 border border-white/5 flex items-center justify-center relative overflow-hidden">
                                <div className="absolute inset-0 bg-[radial-gradient(circle_at_center,rgba(255,255,255,0.1),transparent)]" />
                                <Play className="h-12 w-12 md:h-20 md:w-20 text-white/80 fill-white/20 backdrop-blur-sm rounded-full p-3 md:p-4 border border-white/20 shadow-2xl" />

                                {/* Floating elements inside the 3D card */}
                                <div className="absolute top-4 right-4 md:top-10 md:right-10 bg-black/50 backdrop-blur-md border border-white/10 p-2 md:p-3 rounded-lg md:rounded-xl flex items-center gap-2 md:gap-3 animate-[float_4s_ease-in-out_infinite]">
                                    <div className="w-1.5 h-1.5 md:w-2 md:h-2 rounded-full bg-green-500 animate-pulse" />
                                    <span className="text-[10px] md:text-xs font-medium text-white">AI Processing</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    )
}
