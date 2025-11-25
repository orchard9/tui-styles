import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { InputWithIcon } from '../InputWithIcon'

// eslint-disable-next-line max-lines-per-function
describe('InputWithIcon', () => {
  it('renders input with placeholder', () => {
    render(<InputWithIcon placeholder="Search..." />)
    const input = screen.getByRole('textbox')
    expect(input).toBeInTheDocument()
    expect(input).toHaveAttribute('placeholder', 'Search...')
  })

  it('renders icon when provided', () => {
    const { container } = render(<InputWithIcon icon="Search" />)
    const svg = container.querySelector('svg')
    expect(svg).toBeInTheDocument()
  })

  it('applies error styles when error prop is true', () => {
    render(<InputWithIcon error placeholder="Error input" />)
    const input = screen.getByRole('textbox')
    expect(input).toHaveClass('border-destructive')
    expect(input).toHaveAttribute('aria-invalid', 'true')
  })

  it('handles user input', async () => {
    const handleChange = jest.fn()
    render(<InputWithIcon onChange={handleChange} />)

    const input = screen.getByRole('textbox')
    await userEvent.type(input, 'test')

    expect(handleChange).toHaveBeenCalled()
    expect(input).toHaveValue('test')
  })

  it('adds left padding when icon present', () => {
    render(<InputWithIcon icon="Mail" />)
    const input = screen.getByRole('textbox')
    expect(input).toHaveClass('pl-10')
  })

  it('adds right padding when icon on right', () => {
    render(<InputWithIcon icon="Mail" iconPosition="right" />)
    const input = screen.getByRole('textbox')
    expect(input).toHaveClass('pr-10')
  })

  it('renders without icon padding when no icon', () => {
    render(<InputWithIcon />)
    const input = screen.getByRole('textbox')
    expect(input).not.toHaveClass('pl-10')
    expect(input).not.toHaveClass('pr-10')
  })

  it('applies custom className', () => {
    render(<InputWithIcon className="custom-class" />)
    const input = screen.getByRole('textbox')
    expect(input).toHaveClass('custom-class')
  })

  it('forwards ref correctly', () => {
    const ref = jest.fn()
    render(<InputWithIcon ref={ref} />)
    expect(ref).toHaveBeenCalled()
  })
})
