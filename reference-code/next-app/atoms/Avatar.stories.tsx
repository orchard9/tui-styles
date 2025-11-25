import type { Meta, StoryObj } from '@storybook/react'
import { Avatar } from './Avatar'

const meta: Meta<typeof Avatar> = {
  title: 'Atoms/Avatar',
  component: Avatar,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg', 'xl'],
    },
  },
}

export default meta
type Story = StoryObj<typeof Avatar>

export const WithImage: Story = {
  args: {
    src: 'https://avatars.githubusercontent.com/u/1?v=4',
    alt: 'User avatar',
  },
}

export const WithFallback: Story = {
  args: {
    fallback: 'JD',
  },
}

export const WithIcon: Story = {
  args: {
    // No src or fallback, shows icon
  },
}

export const Small: Story = {
  args: {
    size: 'sm',
    fallback: 'SM',
  },
}

export const Medium: Story = {
  args: {
    size: 'md',
    fallback: 'MD',
  },
}

export const Large: Story = {
  args: {
    size: 'lg',
    fallback: 'LG',
  },
}

export const ExtraLarge: Story = {
  args: {
    size: 'xl',
    fallback: 'XL',
  },
}

export const SmallWithImage: Story = {
  args: {
    size: 'sm',
    src: 'https://avatars.githubusercontent.com/u/1?v=4',
    alt: 'Small avatar',
  },
}

export const LargeWithImage: Story = {
  args: {
    size: 'lg',
    src: 'https://avatars.githubusercontent.com/u/1?v=4',
    alt: 'Large avatar',
  },
}
