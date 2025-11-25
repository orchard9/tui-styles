import type { Metadata, Viewport } from 'next'
import { Geist, Geist_Mono } from 'next/font/google'
import './globals.css'
import { Providers } from './providers'
import { Footer } from '@/components/organisms/Footer'

// Optimize font loading
const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin'],
  display: 'swap',
  preload: true,
})

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
  display: 'swap',
  preload: true,
})

export const metadata: Metadata = {
  title: {
    default: 'Masquerade - Create Any Identity',
    template: '%s | Masquerade',
  },
  description:
    'Transform yourself. Express freely. Be anyone, anywhere. Create and transform your identity instantly with AI-powered face transformation.',
  keywords: [
    'identity',
    'transformation',
    'face swap',
    'avatar',
    'gaming',
    'streaming',
    'self-expression',
    'AI transformation',
  ],
  authors: [{ name: 'Masquerade' }],
  creator: 'Masquerade',
  publisher: 'Masquerade',
  robots: {
    index: true,
    follow: true,
  },
  openGraph: {
    type: 'website',
    locale: 'en_US',
    url: 'https://masquerade.com',
    title: 'Masquerade - Create Any Identity For Yourself',
    description: 'Transform yourself. Express freely. Be anyone, anywhere.',
    siteName: 'Masquerade',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Masquerade - Create Any Identity',
    description: 'Transform yourself. Express freely. Be anyone, anywhere.',
    creator: '@masquerade',
  },
}

export const viewport: Viewport = {
  width: 'device-width',
  initialScale: 1,
  maximumScale: 5,
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <Providers>
          {children}
          <Footer />
        </Providers>
      </body>
    </html>
  )
}
