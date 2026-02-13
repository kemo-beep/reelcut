import { Link } from '@tanstack/react-router'
import { Github, Twitter, Linkedin, Scissors } from 'lucide-react'

export function Footer() {
    return (
        <div
            className="relative md:sticky bottom-0 z-0 h-auto md:h-[400px]" // Sticky only on desktop
            style={{ clipPath: "polygon(0% 0, 100% 0%, 100% 100%, 0 100%)" }} // Clip path mainly for desktop effect
        >
            <footer className="h-full bg-[var(--app-bg-raised)] flex flex-col justify-between py-12 px-6 border-t border-[var(--app-border)]">
                <div className="mx-auto max-w-7xl w-full flex flex-col md:flex-row justify-between items-start gap-12 mb-12 md:mb-0">
                    <div className="space-y-4">
                        <div className="flex items-center gap-2">
                            <div className="w-8 h-8 rounded-lg bg-[var(--app-accent)] flex items-center justify-center text-white">
                                <Scissors className="w-5 h-5" />
                            </div>
                            <span className="text-2xl font-bold text-[var(--app-fg)]">Reelcut</span>
                        </div>
                        <p className="text-[var(--app-fg-muted)] max-w-xs">
                            The intelligent video editor for the next generation of creators.
                        </p>
                    </div>

                    <div className="grid grid-cols-2 md:grid-cols-3 gap-8 md:gap-12 text-sm w-full md:w-auto">
                        <div className="space-y-4">
                            <h4 className="font-bold text-[var(--app-fg)]">Product</h4>
                            <div className="flex flex-col gap-2 text-[var(--app-fg-muted)]">
                                <Link to="/" className="hover:text-[var(--app-fg)] transition-colors">Features</Link>
                                <Link to="/" className="hover:text-[var(--app-fg)] transition-colors">Pricing</Link>
                                <Link to="/" className="hover:text-[var(--app-fg)] transition-colors">Changelog</Link>
                            </div>
                        </div>
                        <div className="space-y-4">
                            <h4 className="font-bold text-[var(--app-fg)]">Company</h4>
                            <div className="flex flex-col gap-2 text-[var(--app-fg-muted)]">
                                <Link to="/" className="hover:text-[var(--app-fg)] transition-colors">About</Link>
                                <Link to="/" className="hover:text-[var(--app-fg)] transition-colors">Careers</Link>
                                <Link to="/" className="hover:text-[var(--app-fg)] transition-colors">Blog</Link>
                            </div>
                        </div>
                        <div className="space-y-4">
                            <h4 className="font-bold text-[var(--app-fg)]">Legal</h4>
                            <div className="flex flex-col gap-2 text-[var(--app-fg-muted)]">
                                <Link to="/" className="hover:text-[var(--app-fg)] transition-colors">Privacy</Link>
                                <Link to="/" className="hover:text-[var(--app-fg)] transition-colors">Terms</Link>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="mx-auto max-w-7xl w-full border-t border-[var(--app-border)] pt-8 flex flex-col md:flex-row justify-between items-center gap-4">
                    <p className="text-xs text-[var(--app-fg-muted)]">
                        &copy; {new Date().getFullYear()} Reelcut Inc. All rights reserved.
                    </p>
                    <div className="flex gap-6">
                        <a href="#" className="text-[var(--app-fg-muted)] hover:text-[var(--app-fg)] transition-colors"><Twitter className="w-5 h-5" /></a>
                        <a href="#" className="text-[var(--app-fg-muted)] hover:text-[var(--app-fg)] transition-colors"><Github className="w-5 h-5" /></a>
                        <a href="#" className="text-[var(--app-fg-muted)] hover:text-[var(--app-fg)] transition-colors"><Linkedin className="w-5 h-5" /></a>
                    </div>
                </div>
            </footer>
        </div>
    )
}
