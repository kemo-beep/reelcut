import { useRef } from 'react'
import { gsap } from 'gsap'
import { useGSAP } from '@gsap/react'

const testimonials = [
    {
        name: "Alex R.",
        role: "YouTuber (2.5M)",
        quote: "Reelcut saved me 20 hours a week. The AI knows exactly what to keep.",
        avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Alex",
    },
    {
        name: "Sarah J.",
        role: "TikTok Strategist",
        quote: "The viral score prediction is scary accurate. My views are up 300%.",
        avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Sarah",
    },
    {
        name: "Mike T.",
        role: "Podcast Host",
        quote: "Finally, a tool that handles long-form context properly.",
        avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Mike",
    },
    {
        name: "Emily W.",
        role: "Content Agency",
        quote: "We use Reelcut for all our clients. It's a game changer.",
        avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Emily",
    },
    {
        name: "David K.",
        role: "Streamer",
        quote: "From Twitch VOD to TikTok banger in minutes.",
        avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=David",
    }
]

export function Testimonials() {
    const containerRef = useRef<HTMLDivElement>(null)
    const row1Ref = useRef<HTMLDivElement>(null)
    const row2Ref = useRef<HTMLDivElement>(null)

    useGSAP(() => {
        const row1 = row1Ref.current
        const row2 = row2Ref.current
        if (!row1 || !row2) return

        // Duplicate content for seamless loop
        const content1 = Array.from(row1.children)
        content1.forEach(item => row1.appendChild(item.cloneNode(true)))

        const content2 = Array.from(row2.children)
        content2.forEach(item => row2.appendChild(item.cloneNode(true)))

        // Row 1 - Left
        gsap.to(row1, {
            x: "-50%",
            duration: 20,
            ease: "none",
            repeat: -1
        })

        // Row 2 - Right
        gsap.fromTo(row2,
            { x: "-50%" },
            {
                x: "0%",
                duration: 25,
                ease: "none",
                repeat: -1
            }
        )
    }, { scope: containerRef })

    return (
        <section ref={containerRef} className="py-20 md:py-32 bg-[var(--app-bg)] overflow-hidden">
            <div className="text-center mb-12 md:mb-16 px-6">
                <h2 className="text-display text-4xl md:text-5xl font-bold text-[var(--app-fg)]">
                    Creators Love Us
                </h2>
            </div>

            <div className="flex flex-col gap-6 md:gap-8 mask-gradient-x">
                {/* Row 1 */}
                <div ref={row1Ref} className="flex gap-4 md:gap-6 w-max">
                    {testimonials.map((t, i) => (
                        <div key={i} className="w-[280px] md:w-[350px] p-6 md:p-8 rounded-2xl bg-[var(--app-bg-raised)] border border-[var(--app-border)] shadow-sm hover:border-[var(--app-accent)] transition-colors">
                            <p className="text-base md:text-lg text-[var(--app-fg)] mb-4 md:mb-6 leading-relaxed">"{t.quote}"</p>
                            <div className="flex items-center gap-3 md:gap-4">
                                <img src={t.avatar} alt={t.name} className="w-8 h-8 md:w-10 md:h-10 rounded-full bg-[var(--app-bg-overlay)]" />
                                <div>
                                    <div className="font-bold text-[var(--app-fg)] text-sm md:text-base">{t.name}</div>
                                    <div className="text-xs md:text-sm text-[var(--app-fg-muted)]">{t.role}</div>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Row 2 */}
                <div ref={row2Ref} className="flex gap-4 md:gap-6 w-max">
                    {testimonials.map((t, i) => (
                        <div key={i} className="w-[280px] md:w-[350px] p-6 md:p-8 rounded-2xl bg-[var(--app-bg-raised)] border border-[var(--app-border)] shadow-sm hover:border-[var(--app-accent)] transition-colors">
                            <p className="text-base md:text-lg text-[var(--app-fg)] mb-4 md:mb-6 leading-relaxed">"{t.quote}"</p>
                            <div className="flex items-center gap-3 md:gap-4">
                                <img src={t.avatar} alt={t.name} className="w-8 h-8 md:w-10 md:h-10 rounded-full bg-[var(--app-bg-overlay)]" />
                                <div>
                                    <div className="font-bold text-[var(--app-fg)] text-sm md:text-base">{t.name}</div>
                                    <div className="text-xs md:text-sm text-[var(--app-fg-muted)]">{t.role}</div>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    )
}
