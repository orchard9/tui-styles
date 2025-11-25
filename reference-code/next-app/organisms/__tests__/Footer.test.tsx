import { render, screen } from '@testing-library/react'
import { Footer } from '../Footer'

describe('Footer', () => {
  it('renders brand name', () => {
    render(<Footer />)
    expect(screen.getByText('Masquerade')).toBeInTheDocument()
  })

  it('renders brand tagline', () => {
    render(<Footer />)
    expect(
      screen.getByText(
        /Transform yourself. Express freely. Be anyone, anywhere./i
      )
    ).toBeInTheDocument()
  })

  it('renders Sparkles icon', () => {
    render(<Footer />)
    const footer = screen.getByRole('contentinfo')
    const sparklesIcon = footer.querySelector('svg')
    expect(sparklesIcon).toBeInTheDocument()
  })

  it('renders Privacy Policy link', () => {
    render(<Footer />)
    const privacyLink = screen.getByText('Privacy Policy')
    expect(privacyLink).toBeInTheDocument()
    expect(privacyLink).toHaveAttribute('href', '/privacy')
  })

  it('renders Terms of Service link', () => {
    render(<Footer />)
    const termsLink = screen.getByText('Terms of Service')
    expect(termsLink).toBeInTheDocument()
    expect(termsLink).toHaveAttribute('href', '/terms')
  })

  it('renders copyright with current year', () => {
    render(<Footer />)
    const currentYear = new Date().getFullYear()
    expect(
      screen.getByText(
        new RegExp(`© ${currentYear} Masquerade. All rights reserved.`)
      )
    ).toBeInTheDocument()
  })
})
