import { useRef } from 'react'
import { gsap } from 'gsap'
import { useGSAP } from '@gsap/react'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { Upload, Scissors, Share2 } from 'lucide-react'

gsap.registerPlugin(ScrollTrigger)

const steps = [
    {
        id: 'step-1',
        title: 'Upload Content',
        description: 'Drop your long-form video link or file. We handle 4K resolution and all major formats.',
        icon: <Upload className="h-6 w-6 text-white" />,
        color: 'bg-blue-600',
        gradient: 'from-blue-600 to-indigo-600'
    },
    {
        id: 'step-2',
        title: 'AI Processing',
        description: 'Our neural networks analyze audio, visual cues, and pacing to identify viral moments.',
        icon: <Scissors className="h-6 w-6 text-white" />,
        color: 'bg-purple-600',
        gradient: 'from-purple-600 to-pink-600'
    },
    {
        id: 'step-3',
        title: 'Review & Export',
        description: 'Fine-tune the results with our precision editor. Export to all platforms instantly.',
        icon: <Share2 className="h-6 w-6 text-white" />,
        color: 'bg-green-600',
        gradient: 'from-green-600 to-emerald-600'
    },
]

export function HowItWorks() {
    const containerRef = useRef<HTMLDivElement>(null)
    const visualsRef = useRef<(HTMLDivElement | null)[]>([])

    useGSAP(() => {
        const visuals = visualsRef.current

        ScrollTrigger.matchMedia({
            // Desktop: Parallax & Complex Animations
            "(min-width: 768px)": function () {
                steps.forEach((step, index) => {
                    const trigger = document.getElementById(step.id)
                    if (!trigger || !visuals[index]) return

                    // Parallax effect on the visual card
                    gsap.to(visuals[index], {
                        scrollTrigger: {
                            trigger: trigger,
                            start: "top center",
                            end: "bottom center",
                            scrub: 1
                        },
                        y: -50,
                        rotateZ: index % 2 === 0 ? 5 : -5,
                        scale: 1.05
                    })

                    // Reveal opacity based on scroll
                    gsap.fromTo(visuals[index],
                        { opacity: 0.3, filter: 'blur(10px)', scale: 0.9 },
                        {
                            opacity: 1,
                            filter: 'blur(0px)',
                            scale: 1,
                            scrollTrigger: {
                                trigger: trigger,
                                start: "top bottom",
                                end: "center center",
                                scrub: 1
                            }
                        }
                    )
                })
            },
            // Mobile: Simple Fade In
            "(max-width: 767px)": function () {
                steps.forEach((step, index) => {
                    if (!visuals[index]) return
                    gsap.from(visuals[index], {
                        opacity: 0,
                        y: 30,
                        duration: 0.8,
                        scrollTrigger: {
                            trigger: visuals[index],
                            start: "top 85%"
                        }
                    })
                })
            }
        })

    }, { scope: containerRef })

    return (
        <section ref={containerRef} className="relative bg-[var(--app-bg-raised)] py-20 md:py-32 overflow-hidden">
            <div className="mx-auto max-w-7xl px-6">
                <div className="mb-16 md:mb-24 text-center">
                    <h2 className="text-display text-4xl md:text-5xl font-bold text-[var(--app-fg)] leading-[1.1]">
                        Streamlined Workflow
                    </h2>
                    <p className="text-lg md:text-xl text-[var(--app-fg-muted)] mt-4">Average processing time: 45 seconds</p>
                </div>

                <div className="relative">
                    {/* Connecting Line - Desktop Only */}
                    <div className="hidden md:block absolute left-1/2 top-0 bottom-0 w-0.5 bg-[var(--app-border)] -translate-x-1/2" />

                    <div className="space-y-20 md:space-y-32">
                        {steps.map((step, index) => (
                            <div
                                key={step.id}
                                id={step.id}
                                className={`flex flex-col md:flex-row gap-8 md:gap-12 items-center ${index % 2 === 0 ? 'md:flex-row' : 'md:flex-row-reverse'
                                    }`}
                            >
                                {/* Text Content */}
                                <div className={`flex-1 w-full md:w-1/2 relative md:text-right ${index % 2 !== 0 ? 'md:text-left' : ''
                                    }`}>
                                    {/* Timeline Node - Desktop Only */}
                                    <div className={`hidden md:flex absolute ${index % 2 === 0 ? 'right-[-38px]' : 'left-[-38px]'} top-0 w-14 h-14 rounded-full border-4 border-[var(--app-bg-raised)] bg-gradient-to-br ${step.gradient} items-center justify-center shadow-lg z-10`}>
                                        <div className="text-white">
                                            {step.icon}
                                        </div>
                                    </div>

                                    {/* Mobile Icon */}
                                    <div className="md:hidden mb-4 inline-flex p-3 rounded-xl bg-gradient-to-br from-[var(--app-accent-muted)] to-transparent border border-[var(--app-border)]">
                                        {/* Clone icon to change color if needed or use same */}
                                        <div className="text-[var(--app-fg)]">{step.icon}</div>
                                    </div>

                                    <div className={`pr-0 ${index % 2 === 0 ? 'md:pr-16' : 'md:pl-16'}`}>
                                        <h3 className="text-2xl md:text-3xl font-bold text-[var(--app-fg)] mb-3 md:mb-4">{step.title}</h3>
                                        <p className="text-base md:text-lg text-[var(--app-fg-muted)] leading-relaxed">
                                            {step.description}
                                        </p>
                                    </div>
                                </div>

                                {/* Visual Card */}
                                <div className="flex-1 w-full md:w-1/2 perspective-1000">
                                    <div
                                        ref={el => { visualsRef.current[index] = el }}
                                        className="relative aspect-video rounded-2xl bg-[var(--app-bg)] border border-[var(--app-border)] shadow-xl overflow-hidden group"
                                    >
                                        <div className={`absolute inset-0 opacity-10 bg-gradient-to-br ${step.gradient}`} />

                                        {/* Mock UI Content */}
                                        <div className="absolute inset-0 p-4 md:p-6 flex flex-col justify-between">
                                            <div className="h-6 md:h-8 bg-[var(--app-bg-overlay)] rounded-md w-1/3" />
                                            <div className="flex gap-3 md:gap-4">
                                                <div className="w-2/3 space-y-2 md:space-y-3">
                                                    <div className="w-full h-16 md:h-20 bg-[var(--app-bg-overlay)] rounded-lg" />
                                                    <div className="w-full h-16 md:h-20 bg-[var(--app-bg-overlay)] rounded-lg opacity-50" />
                                                </div>
                                                <div className="w-1/3 bg-[var(--app-bg-overlay)] rounded-lg" />
                                            </div>
                                        </div>

                                        {/* Floating Badge */}
                                        <div className="absolute bottom-4 right-4 md:bottom-6 md:right-6 px-3 py-1.5 md:px-4 md:py-2 bg-black/80 backdrop-blur-md rounded-full border border-white/10 text-[10px] md:text-xs font-mono text-white">
                                            Executing Step {index + 1}...
                                        </div>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </section>
    )
}
