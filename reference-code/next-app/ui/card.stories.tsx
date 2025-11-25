import type { Meta, StoryObj } from '@storybook/react'
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
  CardAction,
} from './card'
import { Button } from './button'

const meta: Meta<typeof Card> = {
  title: 'UI/Card',
  component: Card,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  decorators: [
    (Story) => (
      <div className="w-96">
        <Story />
      </div>
    ),
  ],
}

export default meta
type Story = StoryObj<typeof Card>

export const Default: Story = {
  render: () => (
    <Card>
      <CardHeader>
        <CardTitle>Card Title</CardTitle>
        <CardDescription>Card description goes here</CardDescription>
      </CardHeader>
      <CardContent>
        <p>Card content with some text to demonstrate the layout.</p>
      </CardContent>
    </Card>
  ),
}

export const Elevated: Story = {
  render: () => (
    <Card variant="elevated">
      <CardHeader>
        <CardTitle>Elevated Card</CardTitle>
        <CardDescription>This card has elevated shadow styling</CardDescription>
      </CardHeader>
      <CardContent>
        <p>Hover over this card to see the shadow effect.</p>
      </CardContent>
    </Card>
  ),
}

export const Interactive: Story = {
  render: () => (
    <Card variant="interactive">
      <CardHeader>
        <CardTitle>Interactive Card</CardTitle>
        <CardDescription>This card changes on hover</CardDescription>
      </CardHeader>
      <CardContent>
        <p>Hover over this card to see the interactive effect.</p>
      </CardContent>
    </Card>
  ),
}

export const WithFooter: Story = {
  render: () => (
    <Card>
      <CardHeader>
        <CardTitle>Card with Footer</CardTitle>
        <CardDescription>This card includes a footer section</CardDescription>
      </CardHeader>
      <CardContent>
        <p>Main card content goes here.</p>
      </CardContent>
      <CardFooter className="border-t">
        <Button variant="outline" className="mr-2">
          Cancel
        </Button>
        <Button>Save</Button>
      </CardFooter>
    </Card>
  ),
}

export const WithAction: Story = {
  render: () => (
    <Card>
      <CardHeader>
        <CardTitle>Card with Action</CardTitle>
        <CardDescription>
          This card has an action button in the header
        </CardDescription>
        <CardAction>
          <Button size="sm" variant="outline">
            Edit
          </Button>
        </CardAction>
      </CardHeader>
      <CardContent>
        <p>The action button is positioned in the top-right corner.</p>
      </CardContent>
    </Card>
  ),
}

export const Complete: Story = {
  render: () => (
    <Card variant="elevated">
      <CardHeader>
        <CardTitle>Complete Card Example</CardTitle>
        <CardDescription>
          Demonstrating all card components together
        </CardDescription>
        <CardAction>
          <Button size="sm" variant="ghost">
            ...
          </Button>
        </CardAction>
      </CardHeader>
      <CardContent>
        <p className="mb-4">
          This card demonstrates all available subcomponents working together.
        </p>
        <ul className="list-disc list-inside space-y-2">
          <li>Header with title and description</li>
          <li>Action button in header</li>
          <li>Content area with rich content</li>
          <li>Footer with action buttons</li>
        </ul>
      </CardContent>
      <CardFooter className="border-t gap-2">
        <Button variant="outline" className="flex-1">
          Cancel
        </Button>
        <Button className="flex-1">Confirm</Button>
      </CardFooter>
    </Card>
  ),
}
