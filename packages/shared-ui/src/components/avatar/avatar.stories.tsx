import type { Meta, StoryObj } from '@storybook/react'
import { getInitials } from '../../utils'
import { Avatar, AvatarFallback, AvatarImage } from './avatar'

const meta: Meta = {
  title: 'Basic Components/Avatar',
  component: Avatar,
  argTypes: {
    className: {
      control: 'text',
      description: 'Additional classes to apply',
    },
  },
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: () => (
    <Avatar>
      <AvatarImage src="https://avatars.githubusercontent.com/u/921834?v=4" alt="@riipandi" />
      <AvatarFallback>{getInitials('Aris Ripandi')}</AvatarFallback>
    </Avatar>
  ),
}

export const Sizes: Story = {
  render: () => (
    <div className="flex items-center gap-4">
      <Avatar className="h-8 w-8">
        <AvatarImage src="https://github.com/riipandi.png" alt="@riipandi" />
        <AvatarFallback>{getInitials('Aris Ripandi')}</AvatarFallback>
      </Avatar>
      <Avatar className="h-12 w-12">
        <AvatarImage src="https://github.com/riipandi.png" alt="@riipandi" />
        <AvatarFallback>{getInitials('Aris Ripandi')}</AvatarFallback>
      </Avatar>
      <Avatar className="h-16 w-16">
        <AvatarImage src="https://github.com/riipandi.png" alt="@riipandi" />
        <AvatarFallback>{getInitials('Aris Ripandi')}</AvatarFallback>
      </Avatar>
    </div>
  ),
}

export const Rounded: Story = {
  render: () => (
    <div className="flex items-center gap-4">
      <Avatar className="rounded-sm">
        <AvatarImage src="https://github.com/riipandi.png" alt="@riipandi" />
        <AvatarFallback>{getInitials('Aris Ripandi')}</AvatarFallback>
      </Avatar>
      <Avatar className="rounded-md">
        <AvatarImage src="https://github.com/riipandi.png" alt="@riipandi" />
        <AvatarFallback>{getInitials('Aris Ripandi')}</AvatarFallback>
      </Avatar>
      <Avatar className="rounded-full">
        <AvatarImage src="https://github.com/riipandi.png" alt="@riipandi" />
        <AvatarFallback>{getInitials('Aris Ripandi')}</AvatarFallback>
      </Avatar>
    </div>
  ),
}

export const AvatarShowcase: Story = {
  render: () => (
    <div className="flex items-center gap-4">
      <Avatar>
        <AvatarImage src="https://avatars.githubusercontent.com/u/921834?v=4" alt="@riipandi" />
        <AvatarFallback>{getInitials('Aris Ripandi')}</AvatarFallback>
      </Avatar>
      <Avatar>
        <AvatarImage src="https://github.com/riipandi.png" alt="@riipandi" />
        <AvatarFallback>{getInitials('Aris Ripandi')}</AvatarFallback>
      </Avatar>
      <Avatar>
        <AvatarFallback>{getInitials('John Doe')}</AvatarFallback>
      </Avatar>
    </div>
  ),
}
