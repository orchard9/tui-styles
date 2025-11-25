import type { Meta, StoryObj } from '@storybook/react'
import { InputWithIcon } from './InputWithIcon'

const meta: Meta<typeof InputWithIcon> = {
  title: 'Molecules/InputWithIcon',
  component: InputWithIcon,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  decorators: [
    (Story) => (
      <div className="w-80">
        <Story />
      </div>
    ),
  ],
}

export default meta
type Story = StoryObj<typeof InputWithIcon>

export const Search: Story = {
  args: {
    icon: 'Search',
    placeholder: 'Search...',
  },
}

export const Email: Story = {
  args: {
    icon: 'Mail',
    placeholder: 'email@example.com',
    type: 'email',
  },
}

export const User: Story = {
  args: {
    icon: 'User',
    placeholder: 'Username',
  },
}

export const Lock: Story = {
  args: {
    icon: 'Lock',
    placeholder: 'Password',
    type: 'password',
  },
}

export const RightIcon: Story = {
  args: {
    icon: 'Calendar',
    iconPosition: 'right',
    placeholder: 'Select date',
  },
}

export const WithError: Story = {
  args: {
    icon: 'CircleAlert',
    placeholder: 'Error state',
    error: true,
  },
}

export const SearchWithError: Story = {
  args: {
    icon: 'Search',
    placeholder: 'Invalid search',
    error: true,
  },
}

export const Disabled: Story = {
  args: {
    icon: 'User',
    placeholder: 'Disabled',
    disabled: true,
  },
}

export const WithValue: Story = {
  args: {
    icon: 'Mail',
    value: 'user@example.com',
  },
}

export const CreditCard: Story = {
  args: {
    icon: 'CreditCard',
    placeholder: '1234 5678 9012 3456',
  },
}
