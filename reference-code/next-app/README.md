# Masquerade Creator Studio Web

Web interface where artists create and monetize unique face avatars for gaming, streaming, and content creation.

## Tech Stack

- **Framework**: Next.js 16 (App Router)
- **Language**: TypeScript (strict mode)
- **UI Library**: React 19.2
- **Styling**: Tailwind CSS v4 + shadcn/ui components
- **State Management**: React Context (theme), React Hook Form (forms)
- **Testing**: Jest + React Testing Library, Vitest (Storybook)
- **Documentation**: Storybook 10
- **Quality**: ESLint 9, Prettier, lint-staged
- **Build Tools**: SWC (fast compilation), next/font optimization

## Getting Started

### Prerequisites

- Node.js 18+ (LTS recommended)
- npm 9+
- Git or Perforce workspace configured

### Installation

1. Clone or sync the repository:

```bash
# Git
git clone <repository-url>
cd apps/creator-studio-web

# Perforce
p4 sync
cd apps/creator-studio-web
```

2. Install dependencies:

```bash
npm install
```

3. Create environment file:

```bash
cp .env.local.example .env.local
```

4. Edit `.env.local` and set your values:

```env
NEXT_PUBLIC_APP_URL=http://localhost:34072
NEXT_PUBLIC_API_URL=http://localhost:34075
```

5. Start development server:

```bash
npm run dev
```

6. Open http://localhost:34072 in your browser

### First Time Setup

After installation, the project is ready for development. Run `npm run quality` to verify all tools are working correctly.

## Development

### Available Scripts

| Script                    | Description                                                                | Example                   |
| ------------------------- | -------------------------------------------------------------------------- | ------------------------- |
| `npm run dev`             | Start development server (port 34072)                                      | `npm run dev`             |
| `npm run build`           | Build production bundle (includes env check)                               | `npm run build`           |
| `npm start`               | Start production server                                                    | `npm start`               |
| `npm run lint`            | Run ESLint (fail on warnings)                                              | `npm run lint`            |
| `npm run lint:fix`        | Auto-fix linting issues                                                    | `npm run lint:fix`        |
| `npm run format`          | Format code with Prettier                                                  | `npm run format`          |
| `npm run format:check`    | Check formatting (CI)                                                      | `npm run format:check`    |
| `npm run typecheck`       | Run TypeScript compiler                                                    | `npm run typecheck`       |
| `npm run find-dead-code`  | Detect unused code with Knip                                               | `npm run find-dead-code`  |
| `npm run complexity`      | Check code complexity via ESLint                                           | `npm run complexity`      |
| `npm run check-circular`  | Detect circular dependencies                                               | `npm run check-circular`  |
| `npm run env:check`       | Validate environment variables                                             | `npm run env:check`       |
| `npm run quality`         | Run all quality checks (typecheck, lint, format, dead code, circular deps) | `npm run quality`         |
| `npm test`                | Run Jest tests                                                             | `npm test`                |
| `npm run test:watch`      | Run tests in watch mode                                                    | `npm run test:watch`      |
| `npm run test:coverage`   | Generate coverage report                                                   | `npm run test:coverage`   |
| `npm run storybook`       | Start Storybook (port 6006)                                                | `npm run storybook`       |
| `npm run build-storybook` | Build Storybook for deployment                                             | `npm run build-storybook` |

### Development Workflow

1. **Pick a task** from the roadmap (`roadmap/milestone-1/`)
2. **Create a branch** with descriptive name (Git: `git checkout -b feature/task-006-shadcn-ui`)
3. **Update task status** to `in_progress` (rename file: `*_ready.md` → `*_in_progress.md`)
4. **Implement changes** following acceptance criteria
5. **Write tests** for new functionality (>80% coverage)
6. **Run quality checks**: `npm run quality`
7. **Test manually** in browser (http://localhost:34072)
8. **Submit changes** to Perforce after quality checks pass
9. **Create review** with task reference and description
10. **Mark task complete** when approved (rename file to `*_complete.md`)

### Project Structure

```
apps/creator-studio-web/
├── app/                      # Next.js 16 App Router
│   ├── layout.tsx           # Root layout with providers
│   ├── page.tsx             # Landing page (/)
│   ├── providers.tsx        # Client providers (ThemeProvider)
│   ├── globals.css          # Global styles + Tailwind directives
│   ├── dashboard/           # Dashboard page (/dashboard)
│   └── test-components/     # Test components page (dev)
├── components/
│   ├── ui/                  # shadcn/ui base components (Button, Input, Card, etc.)
│   ├── atoms/               # Simple components (Icon, Badge, Avatar)
│   ├── molecules/           # Composed components (InputWithIcon, ThemeToggle)
│   ├── organisms/           # Complex components (NavigationBar, Footer)
│   └── templates/           # Page layouts (CreatorStudioLayout)
├── lib/
│   ├── env.ts               # Environment variable validation (Zod)
│   └── utils.ts             # Utility functions (cn, etc.)
├── public/                  # Static assets (images, fonts)
├── .storybook/              # Storybook configuration
├── .vscode/                 # VS Code settings
├── eslint.config.mjs        # ESLint 9 flat config
├── jest.config.ts           # Jest configuration
├── jest.setup.ts            # Jest setup (testing-library)
├── vitest.config.ts         # Vitest config (Storybook tests)
├── next.config.ts           # Next.js configuration
├── postcss.config.mjs       # PostCSS config (Tailwind v4)
├── tsconfig.json            # TypeScript configuration
├── knip.json                # Knip config (dead code detection)
├── components.json          # shadcn/ui config
└── package.json
```

### Component Development

Follow atomic design principles:

1. **Atoms**: Basic building blocks (Button, Input, Icon, Badge, Avatar)
2. **Molecules**: Simple combinations (InputWithIcon, ThemeToggle)
3. **Organisms**: Complex sections (NavigationBar, Footer)
4. **Templates**: Page layouts (CreatorStudioLayout)

Example component pattern with CVA (Class Variance Authority):

```typescript
// components/atoms/Example.tsx
'use client' // Only if using hooks/interactivity

import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const exampleVariants = cva(
  'base-classes transition-colors',
  {
    variants: {
      variant: {
        default: 'bg-primary text-primary-foreground',
        secondary: 'bg-secondary text-secondary-foreground',
      },
      size: {
        sm: 'text-sm px-2 py-1',
        md: 'text-base px-4 py-2',
        lg: 'text-lg px-6 py-3',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'md',
    },
  }
)

export interface ExampleProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof exampleVariants> {
  // Add custom props here
}

export function Example({
  className,
  variant,
  size,
  ...props
}: ExampleProps) {
  return (
    <div
      className={cn(exampleVariants({ variant, size }), className)}
      {...props}
    />
  )
}
```

### Server vs Client Components

**Server Components** (default - no 'use client'):

- Render on server, send HTML to client
- Can fetch data directly (no API calls needed)
- Reduce bundle size (no JavaScript sent to client)
- Better for SEO

**Client Components** ('use client' directive):

- Required for interactivity (onClick, onChange)
- Required for hooks (useState, useEffect, useTheme)
- Required for browser APIs (localStorage, window)

**Decision**: Use server components by default, convert to client only when needed.

## Testing

### Running Tests

```bash
# Run all tests
npm test

# Watch mode (re-run on changes)
npm run test:watch

# Coverage report (minimum 80% required)
npm run test:coverage
```

### Writing Tests

Place tests next to components in `__tests__` directory:

```typescript
// components/atoms/__tests__/Example.test.tsx
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Example } from '../Example'

describe('Example', () => {
  it('renders correctly', () => {
    render(<Example>Test Content</Example>)
    expect(screen.getByText('Test Content')).toBeInTheDocument()
  })

  it('handles user interaction', async () => {
    const handleClick = jest.fn()
    render(<Example onClick={handleClick}>Click me</Example>)

    await userEvent.click(screen.getByText('Click me'))
    expect(handleClick).toHaveBeenCalledTimes(1)
  })

  it('applies variant styles', () => {
    const { container } = render(
      <Example variant="secondary">Secondary</Example>
    )
    expect(container.firstChild).toHaveClass('bg-secondary')
  })
})
```

**Testing Guidelines**:

- Test behavior, not implementation
- Use user-centric queries (getByRole, getByLabelText)
- Avoid testing library internals
- Coverage threshold: 80% (lines, branches, functions, statements)
- Test accessibility (use jest-axe if needed)

## Storybook

View and develop components in isolation:

```bash
npm run storybook
```

Open http://localhost:6006

Create stories for new components:

```typescript
// components/atoms/Example.stories.tsx
import type { Meta, StoryObj } from '@storybook/react'
import { Example } from './Example'

const meta: Meta<typeof Example> = {
  title: 'Atoms/Example',
  component: Example,
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'secondary'],
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
    },
  },
}

export default meta
type Story = StoryObj<typeof Example>

export const Default: Story = {
  args: {
    children: 'Example Component',
  },
}

export const Secondary: Story = {
  args: {
    children: 'Secondary Example',
    variant: 'secondary',
  },
}

export const Large: Story = {
  args: {
    children: 'Large Example',
    size: 'lg',
  },
}
```

**Storybook Features**:

- **Autodocs**: Automatic documentation from component props
- **Accessibility addon**: WCAG 2.1 AA compliance checks
- **Vitest addon**: Run tests in Storybook with Playwright
- **Chromatic integration**: Visual regression testing (future)

## Quality Standards

All code must pass:

- ✅ ESLint (zero warnings)
- ✅ Prettier formatting
- ✅ TypeScript compilation (strict mode)
- ✅ >80% test coverage
- ✅ No dead code (Knip)
- ✅ No circular dependencies (Madge)
- ✅ WCAG 2.1 AA accessibility

Run `npm run quality` before submitting to Perforce to enforce these standards.

## Environment Variables

Environment variables are validated at build time using Zod (see `lib/env.ts`).

### Available Variables

```env
# Required: Application URL
NEXT_PUBLIC_APP_URL=http://localhost:34072

# Required: Backend API URL
NEXT_PUBLIC_API_URL=http://localhost:34075

# Auto-set by Next.js
NODE_ENV=development
```

### Adding New Variables

1. Update `lib/env.ts` schema:

```typescript
const envSchema = z.object({
  NEXT_PUBLIC_NEW_VAR: z.string().min(1),
})
```

2. Add to `.env.local.example`
3. Document in this README

**Client-side variables**: Must be prefixed with `NEXT_PUBLIC_`
**Server-side variables**: No prefix (not exposed to browser)

## Deployment

### Production Build

```bash
# Build application (includes env validation)
npm run build

# Start production server
npm start
```

### Environment Variables (Production)

Required for production:

```env
NEXT_PUBLIC_APP_URL=https://studio.masquerade.com
NEXT_PUBLIC_API_URL=https://api.masquerade.com
NODE_ENV=production
```

### Docker Deployment (Future)

See `Dockerfile` when created in future milestone.

## Performance Optimizations

- **Code Splitting**: Automatic route-based splitting via Next.js
- **Font Optimization**: next/font for Geist Sans and Geist Mono
- **Image Optimization**: next/image for responsive images (when added)
- **SWC Compilation**: Faster than Babel/Webpack
- **Server Components**: Reduce client bundle size
- **CSS**: Tailwind v4 with Lightning CSS (faster than PostCSS)

## Troubleshooting

See [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) for common issues and solutions.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for development guidelines and workflow.

## Architecture

See [../../docs/creator-studio-architecture.md](../../docs/creator-studio-architecture.md) for detailed architecture documentation.

## Resources

- [Next.js Documentation](https://nextjs.org/docs)
- [React Documentation](https://react.dev/)
- [shadcn/ui Components](https://ui.shadcn.com/)
- [Tailwind CSS v4](https://tailwindcss.com/docs)
- [Storybook](https://storybook.js.org/)
- [Testing Library](https://testing-library.com/)
- [Masquerade Design System](../../DESIGN_SYSTEM.md) (when created)
- [Project Roadmap](../../roadmap/milestone-1/)

## License

Proprietary - All rights reserved
