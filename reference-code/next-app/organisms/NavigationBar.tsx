'use client'

import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { ThemeToggle } from '@/components/molecules/ThemeToggle'
import { ArrowRight } from 'lucide-react'

export function NavigationBar() {
  return (
    <nav className="sticky top-0 z-50 border-b border-border/50 bg-background/80 backdrop-blur-md supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2">
            <div className="text-xl font-bold text-glow-primary">
              Masquerade
            </div>
          </Link>

          {/* Right Side - CTA + Theme Toggle */}
          <div className="flex items-center gap-4">
            <ThemeToggle />
            <Button size="sm" className="glow-primary-sm gap-2">
              Get Started
              <ArrowRight className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
    </nav>
  )
}
