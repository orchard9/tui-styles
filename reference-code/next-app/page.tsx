'use client'

import { Button } from '@/components/ui/button'
import { ArrowRight } from 'lucide-react'

export default function HomePage() {
  return (
    <main className="min-h-screen">
      {/* Hero Section - Full viewport height */}
      <section className="relative flex min-h-screen items-center justify-center overflow-hidden px-4 py-20">
        {/* Gradient Background */}
        <div className="absolute inset-0 -z-10 bg-gradient-to-br from-background via-background to-muted" />

        {/* Glowing Orbs */}
        <div className="absolute left-1/4 top-1/4 -z-10 h-96 w-96 rounded-full bg-brand-primary opacity-20 blur-3xl" />
        <div className="absolute bottom-1/4 right-1/4 -z-10 h-96 w-96 rounded-full bg-brand-accent opacity-20 blur-3xl" />

        {/* Hero Content */}
        <div className="relative z-10 max-w-5xl text-center">
          <h1 className="mb-6 text-5xl font-bold tracking-tight text-foreground sm:text-6xl md:text-7xl lg:text-8xl">
            Create Any Identity{' '}
            <span className="text-glow-primary">For Yourself</span>
          </h1>

          <p className="mx-auto mb-12 max-w-2xl text-lg text-muted-foreground sm:text-xl md:text-2xl">
            Express freely. Transform instantly. Be anyone, anywhere.
            <br />
            <span className="text-foreground">Your identity, your rules.</span>
          </p>

          <div className="flex flex-col items-center justify-center gap-4 sm:flex-row">
            <Button
              size="lg"
              className="glow-primary h-14 cursor-pointer gap-2 px-8 text-lg font-semibold"
            >
              Get Started Free
              <ArrowRight className="h-5 w-5" />
            </Button>
            <Button
              size="lg"
              variant="outline"
              className="h-14 cursor-pointer px-8 text-lg font-semibold"
            >
              See How It Works
            </Button>
          </div>
        </div>
      </section>
    </main>
  )
}
