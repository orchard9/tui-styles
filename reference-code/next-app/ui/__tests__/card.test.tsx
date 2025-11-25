import { render } from '@testing-library/react'
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from '../card'

describe('Card', () => {
  it('renders children correctly', () => {
    const { container } = render(<Card>Card content</Card>)
    expect(container.textContent).toContain('Card content')
  })

  it('applies default variant', () => {
    const { container } = render(<Card>Default</Card>)
    const card = container.querySelector('[data-slot="card"]')
    expect(card).toHaveClass('bg-card')
  })

  it('applies elevated variant', () => {
    const { container } = render(<Card variant="elevated">Elevated</Card>)
    const card = container.querySelector('[data-slot="card"]')
    expect(card).toHaveClass('shadow-md')
  })

  it('applies interactive variant', () => {
    const { container } = render(<Card variant="interactive">Interactive</Card>)
    const card = container.querySelector('[data-slot="card"]')
    expect(card).toHaveClass('cursor-pointer')
  })

  it('applies custom className', () => {
    const { container } = render(<Card className="custom-class">Custom</Card>)
    const card = container.querySelector('[data-slot="card"]')
    expect(card).toHaveClass('custom-class')
  })
})

describe('CardHeader', () => {
  it('renders children correctly', () => {
    const { container } = render(<CardHeader>Header</CardHeader>)
    expect(container.textContent).toContain('Header')
  })

  it('has correct data-slot attribute', () => {
    const { container } = render(<CardHeader>Header</CardHeader>)
    expect(
      container.querySelector('[data-slot="card-header"]')
    ).toBeInTheDocument()
  })

  it('applies custom className', () => {
    const { container } = render(
      <CardHeader className="custom">Header</CardHeader>
    )
    const header = container.querySelector('[data-slot="card-header"]')
    expect(header).toHaveClass('custom')
  })
})

describe('CardTitle', () => {
  it('renders children correctly', () => {
    const { container } = render(<CardTitle>Title</CardTitle>)
    expect(container.textContent).toContain('Title')
  })

  it('has correct data-slot attribute', () => {
    const { container } = render(<CardTitle>Title</CardTitle>)
    expect(
      container.querySelector('[data-slot="card-title"]')
    ).toBeInTheDocument()
  })
})

describe('CardDescription', () => {
  it('renders children correctly', () => {
    const { container } = render(<CardDescription>Description</CardDescription>)
    expect(container.textContent).toContain('Description')
  })

  it('has correct data-slot attribute', () => {
    const { container } = render(<CardDescription>Desc</CardDescription>)
    expect(
      container.querySelector('[data-slot="card-description"]')
    ).toBeInTheDocument()
  })
})

describe('CardContent', () => {
  it('renders children correctly', () => {
    const { container } = render(<CardContent>Content</CardContent>)
    expect(container.textContent).toContain('Content')
  })

  it('has correct data-slot attribute', () => {
    const { container } = render(<CardContent>Content</CardContent>)
    expect(
      container.querySelector('[data-slot="card-content"]')
    ).toBeInTheDocument()
  })
})

describe('CardFooter', () => {
  it('renders children correctly', () => {
    const { container } = render(<CardFooter>Footer</CardFooter>)
    expect(container.textContent).toContain('Footer')
  })

  it('has correct data-slot attribute', () => {
    const { container } = render(<CardFooter>Footer</CardFooter>)
    expect(
      container.querySelector('[data-slot="card-footer"]')
    ).toBeInTheDocument()
  })
})

describe('Card composition', () => {
  it('renders complete card structure', () => {
    const { container } = render(
      <Card>
        <CardHeader>
          <CardTitle>Test Title</CardTitle>
          <CardDescription>Test Description</CardDescription>
        </CardHeader>
        <CardContent>Test Content</CardContent>
        <CardFooter>Test Footer</CardFooter>
      </Card>
    )

    expect(container.textContent).toContain('Test Title')
    expect(container.textContent).toContain('Test Description')
    expect(container.textContent).toContain('Test Content')
    expect(container.textContent).toContain('Test Footer')
  })
})
