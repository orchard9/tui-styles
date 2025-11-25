import { render, screen } from '@testing-library/react'
import { Avatar } from '../Avatar'

// Mock Next.js Image component
jest.mock('next/image', () => ({
  __esModule: true,
  default: (props: React.ImgHTMLAttributes<HTMLImageElement>) => {
    // eslint-disable-next-line @next/next/no-img-element, jsx-a11y/alt-text
    return <img {...props} />
  },
}))

describe('Avatar', () => {
  it('renders image when src provided', () => {
    render(<Avatar src="/test.jpg" alt="Test user" />)
    const img = screen.getByRole('img', { name: /test user/i })
    expect(img).toBeInTheDocument()
    expect(img).toHaveAttribute('src', '/test.jpg')
  })

  it('renders fallback initials when no src', () => {
    render(<Avatar fallback="JD" />)
    expect(screen.getByText('JD')).toBeInTheDocument()
  })

  it('renders user icon when no src or fallback', () => {
    const { container } = render(<Avatar />)
    const svg = container.querySelector('svg')
    expect(svg).toBeInTheDocument()
  })

  it('applies small size class', () => {
    const { container } = render(<Avatar size="sm" fallback="S" />)
    const avatar = container.firstChild as HTMLElement
    expect(avatar).toHaveClass('h-8', 'w-8')
  })

  it('applies medium size class', () => {
    const { container } = render(<Avatar size="md" fallback="M" />)
    const avatar = container.firstChild as HTMLElement
    expect(avatar).toHaveClass('h-10', 'w-10')
  })

  it('applies large size class', () => {
    const { container } = render(<Avatar size="lg" fallback="L" />)
    const avatar = container.firstChild as HTMLElement
    expect(avatar).toHaveClass('h-14', 'w-14')
  })

  it('applies xl size class', () => {
    const { container } = render(<Avatar size="xl" fallback="XL" />)
    const avatar = container.firstChild as HTMLElement
    expect(avatar).toHaveClass('h-20', 'w-20')
  })

  it('applies custom className', () => {
    const { container } = render(<Avatar className="custom" fallback="C" />)
    const avatar = container.firstChild as HTMLElement
    expect(avatar).toHaveClass('custom')
  })
})
