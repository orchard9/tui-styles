import type { Meta, StoryObj } from '@storybook/react'

const meta: Meta = {
  title: 'Design System/Tokens',
  parameters: {
    layout: 'padded',
  },
  tags: ['autodocs'],
}

export default meta

export const DesignTokens: StoryObj = {
  // eslint-disable-next-line max-lines-per-function
  render: () => (
    <div className="max-w-4xl mx-auto space-y-12">
      <section>
        <h1 className="text-4xl font-bold mb-4">Design Tokens</h1>
        <p className="text-muted-foreground mb-8">
          Design tokens are the visual design atoms of the design system —
          specifically, they are named entities that store visual design
          attributes.
        </p>
      </section>

      <section>
        <h2 className="text-2xl font-semibold mb-4">Colors</h2>

        <h3 className="text-xl font-medium mb-3 mt-6">Brand Colors</h3>
        <div className="grid grid-cols-3 gap-4 mb-6">
          <div className="space-y-2">
            <div className="h-24 rounded-lg bg-primary"></div>
            <p className="font-medium">Primary</p>
            <p className="text-sm text-muted-foreground">Vibrant purple</p>
          </div>
          <div className="space-y-2">
            <div className="h-24 rounded-lg bg-secondary"></div>
            <p className="font-medium">Secondary</p>
            <p className="text-sm text-muted-foreground">Electric blue</p>
          </div>
          <div className="space-y-2">
            <div className="h-24 rounded-lg bg-accent"></div>
            <p className="font-medium">Accent</p>
            <p className="text-sm text-muted-foreground">Hot pink</p>
          </div>
        </div>

        <h3 className="text-xl font-medium mb-3 mt-6">Semantic Colors</h3>
        <div className="grid grid-cols-4 gap-4 mb-6">
          <div className="space-y-2">
            <div className="h-16 rounded-lg bg-green-500"></div>
            <p className="font-medium">Success</p>
          </div>
          <div className="space-y-2">
            <div className="h-16 rounded-lg bg-yellow-500"></div>
            <p className="font-medium">Warning</p>
          </div>
          <div className="space-y-2">
            <div className="h-16 rounded-lg bg-red-500"></div>
            <p className="font-medium">Error</p>
          </div>
          <div className="space-y-2">
            <div className="h-16 rounded-lg bg-blue-500"></div>
            <p className="font-medium">Info</p>
          </div>
        </div>
      </section>

      <section>
        <h2 className="text-2xl font-semibold mb-4">Typography</h2>

        <h3 className="text-xl font-medium mb-3">Font Families</h3>
        <div className="space-y-2 mb-6">
          <p className="font-sans">Sans: Inter, system-ui, sans-serif</p>
          <p className="font-mono">Mono: JetBrains Mono, Consolas, monospace</p>
        </div>

        <h3 className="text-xl font-medium mb-3">Font Sizes</h3>
        <div className="space-y-2">
          <p className="text-xs">Extra Small (xs): 0.64rem</p>
          <p className="text-sm">Small (sm): 0.8rem</p>
          <p className="text-base">Base: 1rem</p>
          <p className="text-lg">Large (lg): 1.25rem</p>
          <p className="text-xl">Extra Large (xl): 1.563rem</p>
          <p className="text-2xl">2XL: 1.953rem</p>
          <p className="text-3xl">3XL: 2.441rem</p>
        </div>
      </section>

      <section>
        <h2 className="text-2xl font-semibold mb-4">Spacing</h2>
        <p className="text-muted-foreground mb-4">Based on 4px baseline grid</p>
        <div className="space-y-2">
          <div className="flex items-center gap-4">
            <div className="w-1 h-8 bg-primary"></div>
            <span>1 (4px) - Micro spacing</span>
          </div>
          <div className="flex items-center gap-4">
            <div className="w-2 h-8 bg-primary"></div>
            <span>2 (8px) - Component internal</span>
          </div>
          <div className="flex items-center gap-4">
            <div className="w-4 h-8 bg-primary"></div>
            <span>4 (16px) - Component internal</span>
          </div>
          <div className="flex items-center gap-4">
            <div className="w-6 h-8 bg-primary"></div>
            <span>6 (24px) - Component external</span>
          </div>
          <div className="flex items-center gap-4">
            <div className="w-8 h-8 bg-primary"></div>
            <span>8 (32px) - Component external</span>
          </div>
          <div className="flex items-center gap-4">
            <div className="w-12 h-8 bg-primary"></div>
            <span>12 (48px) - Section spacing</span>
          </div>
        </div>
      </section>

      <section>
        <h2 className="text-2xl font-semibold mb-4">Border Radius</h2>
        <div className="grid grid-cols-4 gap-4">
          <div className="space-y-2">
            <div className="h-16 bg-primary/20 rounded-sm"></div>
            <p className="text-sm">sm (2px)</p>
          </div>
          <div className="space-y-2">
            <div className="h-16 bg-primary/20 rounded-md"></div>
            <p className="text-sm">md (6px)</p>
          </div>
          <div className="space-y-2">
            <div className="h-16 bg-primary/20 rounded-lg"></div>
            <p className="text-sm">lg (8px)</p>
          </div>
          <div className="space-y-2">
            <div className="h-16 bg-primary/20 rounded-xl"></div>
            <p className="text-sm">xl (12px)</p>
          </div>
        </div>
      </section>

      <section>
        <h2 className="text-2xl font-semibold mb-4">Shadows</h2>
        <div className="grid grid-cols-3 gap-4">
          <div className="space-y-2">
            <div className="h-16 bg-background border rounded-lg shadow-xs"></div>
            <p className="text-sm">xs</p>
          </div>
          <div className="space-y-2">
            <div className="h-16 bg-background border rounded-lg shadow-sm"></div>
            <p className="text-sm">sm</p>
          </div>
          <div className="space-y-2">
            <div className="h-16 bg-background border rounded-lg shadow-md"></div>
            <p className="text-sm">md</p>
          </div>
        </div>
      </section>

      <section>
        <h2 className="text-2xl font-semibold mb-4">Usage Guidelines</h2>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground">
          <li>Always use design tokens instead of hard-coded values</li>
          <li>
            Use semantic color names (success, warning) over literal values
          </li>
          <li>Stick to the 4px baseline grid for spacing</li>
          <li>Use the defined type scale for consistent hierarchy</li>
          <li>Ensure sufficient color contrast (WCAG AA minimum)</li>
        </ul>
      </section>
    </div>
  ),
}
