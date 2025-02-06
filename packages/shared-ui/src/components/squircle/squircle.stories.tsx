import type { Meta, StoryObj } from '@storybook/react'
import * as Lucide from 'lucide-react'
import { LargeSquircle, SquircleIconContainer } from './squircle'

const meta: Meta<typeof SquircleIconContainer> = {
  title: 'Basic Components/Squircle',
  component: SquircleIconContainer,
  argTypes: {
    size: {
      control: { type: 'radio' },
      options: ['small', 'large'],
      table: {
        defaultValue: { summary: 'large' },
      },
    },
  },
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: (args) => (
    <SquircleIconContainer {...args}>
      <Lucide.Star className="relative z-10 size-5" />
    </SquircleIconContainer>
  ),
}

export const SquircleShowcase: Story = {
  render: () => (
    <div className="flex items-center gap-4">
      <SquircleIconContainer size="small" className="bg-blue-500/10 text-blue-500">
        <Lucide.Star className="relative z-10 size-4" />
      </SquircleIconContainer>

      <SquircleIconContainer size="large" className="bg-green-500/10 text-green-500">
        <Lucide.Heart className="relative z-10 size-5" />
      </SquircleIconContainer>

      <SquircleIconContainer className="bg-purple-500/10 text-purple-500">
        <Lucide.Bell className="relative z-10 size-5" />
      </SquircleIconContainer>
    </div>
  ),
}
