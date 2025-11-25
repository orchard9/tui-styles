import { render } from '@testing-library/react'
import { Icon, type IconName } from '../Icon'

describe('Icon', () => {
  it('renders lucide icon correctly', () => {
    const { container } = render(<Icon name="User" />)
    const svg = container.querySelector('svg')
    expect(svg).toBeInTheDocument()
  })

  it('applies custom size', () => {
    const { container } = render(<Icon name="Heart" size={32} />)
    const svg = container.querySelector('svg')
    expect(svg).toHaveAttribute('width', '32')
    expect(svg).toHaveAttribute('height', '32')
  })

  it('applies default size 24', () => {
    const { container } = render(<Icon name="Star" />)
    const svg = container.querySelector('svg')
    expect(svg).toHaveAttribute('width', '24')
    expect(svg).toHaveAttribute('height', '24')
  })

  it('applies custom className', () => {
    const { container } = render(<Icon name="Star" className="text-primary" />)
    const svg = container.querySelector('svg')
    expect(svg).toHaveClass('text-primary')
  })

  it('logs warning for invalid icon name', () => {
    const consoleWarnSpy = jest.spyOn(console, 'warn').mockImplementation()
    render(<Icon name={'InvalidIcon' as IconName} />)
    expect(consoleWarnSpy).toHaveBeenCalledWith(
      'Icon "InvalidIcon" not found in Lucide icons'
    )
    consoleWarnSpy.mockRestore()
  })

  it('returns null for invalid icon name', () => {
    const consoleWarnSpy = jest.spyOn(console, 'warn').mockImplementation()
    const { container } = render(<Icon name={'InvalidIcon' as IconName} />)
    expect(container.firstChild).toBeNull()
    consoleWarnSpy.mockRestore()
  })
})
