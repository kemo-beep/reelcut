import { useRef } from 'react'
import { gsap } from 'gsap'
import { useGSAP } from '@gsap/react'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { Brain, Type, Layers, Zap } from 'lucide-react'

gsap.registerPlugin(ScrollTrigger)

const features = [
    {
        icon: <Brain className="h-8 w-8 md:h-10 md:w-10 text-[var(--app-accent)]" />,
        title: "Context-Aware Clipping",
        description: "Our AI understands the narrative of your video, not just the pixels. It finds the hooks, the jokes, and the emotional peaks.",
        // Using gradients that work in both modes or relying on opacity
        gradient: "from-[var(--app-accent)]/10 to-purple-500/10",
        border: "border-[var(--app-accent)]/20"
    },
    {
        icon: <Type className="h-8 w-8 md:h-10 md:w-10 text-blue-400" />,
        title: "Dynamic Captions",
        description: "Generate 99.9% accurate captions. Animate them with our preset styles or create your own brand kit.",
        gradient: "from-blue-500/10 to-cyan-500/10",
        border: "border-blue-500/20"
    },
    {
        icon: <Layers className="h-8 w-8 md:h-10 md:w-10 text-green-400" />,
        title: "Multi-Format Export",
        description: "One click to render for TikTok, Reels, Shorts, and LinkedIn. Smart cropping keeps the action in focus.",
        gradient: "from-green-500/10 to-emerald-500/10",
        border: "border-green-500/20"
    },
    {
        icon: <Zap className="h-8 w-8 md:h-10 md:w-10 text-yellow-400" />,
        title: "Cloud Rendering",
        description: "No more heating up your laptop. All processing happens on our dedicated GPU clusters.",
        gradient: "from-yellow-500/10 to-orange-500/10",
        border: "border-yellow-500/20"
    }
]

export function Features() {
    const containerRef = useRef<HTMLDivElement>(null)
    const scrollContainerRef = useRef<HTMLDivElement>(null)

    useGSAP(() => {
        const scrollContainer = scrollContainerRef.current
        if (!scrollContainer) return

        const sections = gsap.utils.toArray('.feature-card')

        ScrollTrigger.matchMedia({
            // Desktop: Horizontal Scroll
            "(min-width: 1024px)": function () {
                gsap.to(sections, {
                    xPercent: -100 * (sections.length - 1),
                    ease: "none",
                    scrollTrigger: {
                        trigger: containerRef.current,
                        pin: true,
                        scrub: 1,
                        end: () => "+=" + scrollContainer.offsetWidth,
                        snap: 1 / (sections.length - 1),
                        invalidateOnRefresh: true
                    }
                })
            },

            // Mobile: Simple Vertical Stacking (Cleanup)
            "(max-width: 1023px)": function () {
                // Reset styles if coming from desktop resize
                gsap.set(sections, { xPercent: 0, clearProps: "all" })
                // Optional: Add simple fade-in or slide-up for mobile
                sections.forEach((section: any) => {
                    gsap.from(section, {
                        y: 50,
                        opacity: 0,
                        duration: 0.8,
                        scrollTrigger: {
                            trigger: section,
                            start: "top 85%"
                        }
                    })
                })
            }
        })

    }, { scope: containerRef })

    return (
        <section ref={containerRef} className="relative min-h-screen bg-[var(--app-bg-raised)] overflow-hidden flex flex-col items-center">
            {/* Background Accents */}
            <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_right,var(--app-accent-muted),transparent)] pointer-events-none opacity-50" />

            <div className="lg:absolute lg:top-10 lg:left-20 z-10 w-full max-w-7xl px-6 pt-20 lg:pt-0 text-center lg:text-left">
                <h2 className="text-display text-4xl md:text-5xl font-bold text-[var(--app-fg)] mb-4">
                    Why Top Creators <br />
                    <span className="text-[var(--app-fg-muted)]">Choose Reelcut</span>
                </h2>
            </div>

            {/* 
                Desktop: Flex row for horizontal scroll 
                Mobile: Flex col for vertical stacking
            */}
            <div
                ref={scrollContainerRef}
                className="flex flex-col lg:flex-row h-auto lg:h-full w-full lg:w-[400%]"
            >
                {features.map((feature, index) => (
                    <div
                        key={index}
                        className="feature-card h-auto lg:h-screen w-full lg:w-screen flex items-center justify-center p-6 lg:p-20 box-border lg:border-r border-[var(--app-border)] last:border-r-0 relative pt-10 pb-20 lg:pt-0 lg:pb-0"
                    >
                        {/* Background Card */}
                        <div className={`relative w-full max-w-2xl aspect-auto lg:aspect-[4/3] rounded-3xl border ${feature.border} bg-gradient-to-br ${feature.gradient} backdrop-blur-sm p-6 md:p-10 flex flex-col justify-end transition-all hover:scale-[1.01] shadow-2xl overflow-hidden`}>
                            <div className="absolute top-6 left-6 md:top-10 md:left-10 p-3 md:p-4 bg-[var(--app-bg)]/80 rounded-2xl border border-[var(--app-border)] backdrop-blur-md shadow-sm">
                                {feature.icon}
                            </div>

                            <div className="mt-20 lg:mt-0">
                                <h3 className="text-display text-2xl md:text-4xl font-bold text-[var(--app-fg)] mb-4 md:mb-6">
                                    {feature.title}
                                </h3>
                                <p className="text-lg md:text-xl text-[var(--app-fg-muted)] leading-relaxed max-w-lg">
                                    {feature.description}
                                </p>
                            </div>

                            {/* Decorative Elements */}
                            <div className="absolute top-1/2 right-10 w-32 h-32 bg-[var(--app-accent)]/10 rounded-full blur-2xl pointer-events-none" />
                            <div className="absolute bottom-4 right-6 md:bottom-10 md:right-10 text-[100px] md:text-[200px] font-bold text-[var(--app-fg)]/5 pointer-events-none leading-none select-none">
                                {index + 1}
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </section>
    )
}
