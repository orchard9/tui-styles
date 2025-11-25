import { render, screen } from '@testing-library/react'

// Mock ThemeToggle to simplify tests
jest.mock('@/components/molecules/ThemeToggle', () => ({
  ThemeToggle: () => <div data-testid="theme-toggle">Theme Toggle</div>,
}))

// Import after mocks are set up
import { NavigationBar } from '../NavigationBar'

describe('NavigationBar', () => {
  it('renders brand logo with correct text', () => {
    render(<NavigationBar />)
    expect(screen.getByText('Masquerade')).toBeInTheDocument()
  })

  it('logo links to home page', () => {
    render(<NavigationBar />)
    const logo = screen.getByText('Masquerade')
    const link = logo.closest('a')
    expect(link).toHaveAttribute('href', '/')
  })

  it('renders Get Started button', () => {
    render(<NavigationBar />)
    expect(
      screen.getByRole('button', { name: /Get Started/i })
    ).toBeInTheDocument()
  })

  it('renders theme toggle', () => {
    render(<NavigationBar />)
    expect(screen.getByTestId('theme-toggle')).toBeInTheDocument()
  })

  it('has sticky navigation', () => {
    const { container } = render(<NavigationBar />)
    const nav = container.querySelector('nav')
    expect(nav).toHaveClass('sticky', 'top-0')
  })

  it('has backdrop blur styling', () => {
    const { container } = render(<NavigationBar />)
    const nav = container.querySelector('nav')
    expect(nav).toHaveClass('backdrop-blur-md')
  })

  it('has proper z-index for overlay', () => {
    const { container } = render(<NavigationBar />)
    const nav = container.querySelector('nav')
    expect(nav).toHaveClass('z-50')
  })
})
