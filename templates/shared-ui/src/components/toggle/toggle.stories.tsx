import type { Meta, StoryObj } from '@storybook/react'
import * as Lucide from 'lucide-react'
import { Toggle } from './toggle'
import type { ToggleVariants } from './toggle.css'

const variantOptions: NonNullable<ToggleVariants['variant']>[] = ['default', 'outline']

const sizeOptions: NonNullable<ToggleVariants['size']>[] = ['sm', 'default', 'lg']

const meta: Meta<typeof Toggle> = {
  title: 'Basic Components/Toggle',
  component: Toggle,
  argTypes: {
    variant: {
      control: { type: 'radio' },
      options: variantOptions,
      table: {
        defaultValue: { summary: 'default' },
        type: { summary: 'ToggleVariants["variant"]' },
      },
    },
    size: {
      control: { type: 'inline-radio' },
      options: sizeOptions,
      table: {
        type: { summary: 'ToggleVariants["size"]' },
      },
    },
    pressed: {
      control: 'boolean',
    },
    disabled: {
      control: 'boolean',
    },
  },
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    children: <Lucide.Bold className="size-4" />,
    'aria-label': 'Toggle bold',
  },
}

export const VariantShowcase: Story = {
  parameters: {
    controls: { exclude: ['variant'] },
  },
  render: (args) => (
    <div className="flex items-center gap-4">
      <Toggle {...args}>
        <Lucide.Bold className="size-4" />
      </Toggle>
      <Toggle {...args} variant="outline">
        <Lucide.Italic className="size-4" />
      </Toggle>
    </div>
  ),
}

export const SizeShowcase: Story = {
  parameters: {
    controls: { exclude: ['size'] },
  },
  render: (args) => (
    <div className="flex items-center gap-4">
      <Toggle {...args} size="sm">
        <Lucide.Bold className="size-4" />
      </Toggle>
      <Toggle {...args} size="default">
        <Lucide.Bold className="size-4" />
      </Toggle>
      <Toggle {...args} size="lg">
        <Lucide.Bold className="size-4" />
      </Toggle>
    </div>
  ),
}

export const WithText: Story = {
  render: (args) => (
    <div className="flex items-center gap-4">
      <Toggle {...args}>
        <Lucide.Bold className="size-4" />
        Bold
      </Toggle>
      <Toggle {...args}>
        <Lucide.Italic className="size-4" />
        Italic
      </Toggle>
    </div>
  ),
}
