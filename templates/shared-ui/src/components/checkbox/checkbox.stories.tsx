import type { Meta, StoryObj } from '@storybook/react'
import { Label } from '../label/label'
import { Checkbox } from './checkbox'

const meta: Meta<typeof Checkbox> = {
  title: 'Basic Components/Checkbox',
  component: Checkbox,
  argTypes: {
    checked: {
      control: 'boolean',
      description: 'The controlled checked state of the checkbox',
    },
    disabled: {
      control: 'boolean',
      description: 'When true, prevents the user from interacting with the checkbox',
    },
  },
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: () => (
    <div className="flex items-center space-x-2">
      <Checkbox id="terms" />
      <Label htmlFor="terms">Accept terms and conditions</Label>
    </div>
  ),
}

export const CheckboxShowcase: Story = {
  render: () => (
    <div className="flex flex-col gap-4">
      <div className="flex items-center space-x-2">
        <Checkbox id="default" />
        <Label htmlFor="default">Default Checkbox</Label>
      </div>

      <div className="flex items-center space-x-2">
        <Checkbox id="checked" defaultChecked />
        <Label htmlFor="checked">Checked Checkbox</Label>
      </div>

      <div className="flex items-center space-x-2">
        <Checkbox id="disabled" disabled />
        <Label htmlFor="disabled">Disabled Checkbox</Label>
      </div>

      <div className="flex items-center space-x-2">
        <Checkbox id="disabled-checked" disabled defaultChecked />
        <Label htmlFor="disabled-checked">Disabled Checked</Label>
      </div>
    </div>
  ),
}
