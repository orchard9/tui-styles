import type { Meta, StoryObj } from '@storybook/react'
import { Icon } from './Icon'

const meta: Meta<typeof Icon> = {
  title: 'Atoms/Icon',
  component: Icon,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    size: {
      control: 'number',
    },
  },
}

export default meta
type Story = StoryObj<typeof Icon>

export const User: Story = {
  args: {
    name: 'User',
    size: 24,
  },
}

export const Heart: Story = {
  args: {
    name: 'Heart',
    size: 24,
  },
}

export const Settings: Story = {
  args: {
    name: 'Settings',
    size: 24,
  },
}

export const Star: Story = {
  args: {
    name: 'Star',
    size: 24,
  },
}

export const Search: Story = {
  args: {
    name: 'Search',
    size: 24,
  },
}

export const Bell: Story = {
  args: {
    name: 'Bell',
    size: 24,
  },
}

export const Mail: Story = {
  args: {
    name: 'Mail',
    size: 24,
  },
}

export const Upload: Story = {
  args: {
    name: 'Upload',
    size: 24,
  },
}

export const Download: Story = {
  args: {
    name: 'Download',
    size: 24,
  },
}

export const Colored: Story = {
  args: {
    name: 'Star',
    size: 32,
    className: 'text-primary',
  },
}

export const Large: Story = {
  args: {
    name: 'Heart',
    size: 48,
  },
}

export const Small: Story = {
  args: {
    name: 'User',
    size: 16,
  },
}

export const IconShowcase: Story = {
  render: () => (
    <div className="grid grid-cols-5 gap-4">
      <div className="flex flex-col items-center gap-2">
        <Icon name="User" size={32} />
        <span className="text-xs">User</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <Icon name="Heart" size={32} />
        <span className="text-xs">Heart</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <Icon name="Settings" size={32} />
        <span className="text-xs">Settings</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <Icon name="Star" size={32} />
        <span className="text-xs">Star</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <Icon name="Search" size={32} />
        <span className="text-xs">Search</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <Icon name="Bell" size={32} />
        <span className="text-xs">Bell</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <Icon name="Mail" size={32} />
        <span className="text-xs">Mail</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <Icon name="Upload" size={32} />
        <span className="text-xs">Upload</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <Icon name="Download" size={32} />
        <span className="text-xs">Download</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <Icon name="Trash" size={32} />
        <span className="text-xs">Trash</span>
      </div>
    </div>
  ),
}
