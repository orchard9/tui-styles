import { render } from '@testing-library/react'
import { ThemeToggle } from '../ThemeToggle'

describe('ThemeToggle', () => {
  it('renders null (dark mode is forced)', () => {
    const { container } = render(<ThemeToggle />)
    expect(container.firstChild).toBeNull()
  })

  it('does not render any button', () => {
    const { container } = render(<ThemeToggle />)
    const button = container.querySelector('button')
    expect(button).toBeNull()
  })

  it('has no DOM output', () => {
    const { container } = render(<ThemeToggle />)
    expect(container.innerHTML).toBe('')
  })
})
