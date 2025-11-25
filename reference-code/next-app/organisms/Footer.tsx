import Link from 'next/link'
import { Sparkles } from 'lucide-react'

export function Footer() {
  const currentYear = new Date().getFullYear()

  return (
    <footer className="border-t border-border bg-background">
      <div className="container mx-auto px-4 py-12">
        <div className="flex flex-col items-center gap-8">
          {/* Brand */}
          <div className="text-center">
            <div className="mb-4 flex items-center justify-center gap-2">
              <Sparkles className="h-6 w-6 text-primary" />
              <span className="text-xl font-bold">Masquerade</span>
            </div>
            <p className="text-sm text-muted-foreground">
              Transform yourself. Express freely. Be anyone, anywhere.
            </p>
          </div>

          {/* Legal Links */}
          <div className="flex w-full flex-col items-center gap-4 border-t border-border pt-8 md:flex-row md:justify-between">
            <p className="text-sm text-muted-foreground">
              © {currentYear} Masquerade. All rights reserved.
            </p>
            <div className="flex gap-6 text-sm">
              <Link
                href="/privacy"
                className="text-muted-foreground transition-colors hover:text-foreground"
              >
                Privacy Policy
              </Link>
              <Link
                href="/terms"
                className="text-muted-foreground transition-colors hover:text-foreground"
              >
                Terms of Service
              </Link>
            </div>
          </div>
        </div>
      </div>
    </footer>
  )
}
