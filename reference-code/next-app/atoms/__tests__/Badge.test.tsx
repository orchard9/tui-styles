import { render, screen } from '@testing-library/react'
import { Badge } from '../Badge'

describe('Badge', () => {
  it('renders children correctly', () => {
    render(<Badge>New</Badge>)
    expect(screen.getByText('New')).toBeInTheDocument()
  })

  it('applies default variant styles', () => {
    const { container } = render(<Badge>Default</Badge>)
    const badge = container.firstChild as HTMLElement
    expect(badge).toHaveClass('bg-primary')
  })

  it('applies success variant styles', () => {
    const { container } = render(<Badge variant="success">Success</Badge>)
    const badge = container.firstChild as HTMLElement
    expect(badge).toHaveClass('bg-[hsl(var(--semantic-success))]')
  })

  it('applies warning variant styles', () => {
    const { container } = render(<Badge variant="warning">Warning</Badge>)
    const badge = container.firstChild as HTMLElement
    expect(badge).toHaveClass('bg-[hsl(var(--semantic-warning))]')
  })

  it('applies error variant styles', () => {
    const { container } = render(<Badge variant="error">Error</Badge>)
    const badge = container.firstChild as HTMLElement
    expect(badge).toHaveClass('bg-[hsl(var(--semantic-error))]')
  })

  it('applies info variant styles', () => {
    const { container } = render(<Badge variant="info">Info</Badge>)
    const badge = container.firstChild as HTMLElement
    expect(badge).toHaveClass('bg-[hsl(var(--semantic-info))]')
  })

  it('applies small size', () => {
    const { container } = render(<Badge size="sm">Small</Badge>)
    const badge = container.firstChild as HTMLElement
    expect(badge).toHaveClass('px-2')
  })

  it('applies large size', () => {
    const { container } = render(<Badge size="lg">Large</Badge>)
    const badge = container.firstChild as HTMLElement
    expect(badge).toHaveClass('px-3')
  })

  it('applies custom className', () => {
    const { container } = render(<Badge className="custom">Custom</Badge>)
    const badge = container.firstChild as HTMLElement
    expect(badge).toHaveClass('custom')
  })
})
