import type { Meta, StoryObj } from '@storybook/react'
import { Slider } from './slider'

const meta: Meta<typeof Slider> = {
  title: 'Basic Components/Slider',
  component: Slider,
  argTypes: {
    defaultValue: {
      control: 'object',
      description: 'The default value of the slider',
    },
    max: {
      control: 'number',
      description: 'The maximum value for the slider',
    },
    step: {
      control: 'number',
      description: 'The stepping interval',
    },
    disabled: {
      control: 'boolean',
      description: 'When true, prevents the user from interacting with the slider',
    },
  },
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: () => (
    <div className="w-full min-w-[600px] p-6">
      <Slider defaultValue={[50]} max={100} step={1} className="w-[60%]" />
    </div>
  ),
}

export const SliderShowcase: Story = {
  render: () => (
    <div className="flex w-full min-w-[600px] flex-col gap-8">
      <div className="space-y-4">
        <h4 className="font-medium text-sm">Default Slider</h4>
        <Slider defaultValue={[50]} max={100} step={1} />
      </div>

      <div className="space-y-4">
        <h4 className="font-medium text-sm">Range Slider</h4>
        <Slider defaultValue={[25, 75]} max={100} step={1} />
      </div>

      <div className="space-y-4">
        <h4 className="font-medium text-sm">Disabled Slider</h4>
        <Slider defaultValue={[40]} max={100} step={1} disabled />
      </div>

      <div className="space-y-4">
        <h4 className="font-medium text-sm">Step Slider</h4>
        <Slider defaultValue={[20]} max={100} step={20} />
      </div>
    </div>
  ),
}
