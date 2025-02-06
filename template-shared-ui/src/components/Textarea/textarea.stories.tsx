import type { Meta, StoryObj } from '@storybook/react'
import { Textarea } from './textarea'

const meta: Meta<typeof Textarea> = {
  title: 'Basic Components/Textarea',
  component: Textarea,
  argTypes: {
    placeholder: {
      control: 'text',
      description: 'Placeholder text',
    },
    disabled: {
      control: 'boolean',
      description: 'Whether the textarea is disabled',
    },
  },
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    placeholder: 'Type your message here.',
    className: 'w-full min-w-[400px]',
  },
}

export const TextareaShowcase: Story = {
  render: () => (
    <div className="flex w-full min-w-[400px] flex-col gap-4">
      <Textarea placeholder="Default textarea" />

      <Textarea placeholder="Disabled textarea" disabled />

      <Textarea
        placeholder="With initial value"
        defaultValue="This is some default content in the textarea."
      />

      <Textarea placeholder="Custom rows" rows={8} />
    </div>
  ),
}
