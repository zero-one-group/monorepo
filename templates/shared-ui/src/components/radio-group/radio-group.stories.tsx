import type { Meta, StoryObj } from '@storybook/react'
import { Label } from '../label/label'
import { RadioGroup, RadioGroupItem } from './radio-group'
import type { RadioGroupVariants } from './radio-group.css'

const orientationOptions: NonNullable<RadioGroupVariants['orientation']>[] = [
  'horizontal',
  'vertical',
]

const meta: Meta<typeof RadioGroup> = {
  title: 'Basic Components/RadioGroup',
  component: RadioGroup,
  argTypes: {
    orientation: {
      control: 'radio',
      options: orientationOptions,
      table: {
        defaultValue: { summary: 'vertical' },
        type: { summary: 'RadioGroupVariants["orientation"]' },
      },
    },
  },
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: (args) => (
    <RadioGroup defaultValue="option-1" {...args}>
      <div className="flex items-center space-x-2">
        <RadioGroupItem id="option-1" value="option-1" />
        <Label htmlFor="option-1">Option 1</Label>
      </div>
      <div className="flex items-center space-x-2">
        <RadioGroupItem id="option-2" value="option-2" />
        <Label htmlFor="option-2">Option 2</Label>
      </div>
      <div className="flex items-center space-x-2">
        <RadioGroupItem id="option-3" value="option-3" />
        <Label htmlFor="option-3">Option 3</Label>
      </div>
    </RadioGroup>
  ),
}
