import type { Meta, StoryObj } from '@storybook/react'
import * as Lucide from 'lucide-react'
import { ToggleGroup, ToggleGroupItem } from './toggle-group'

const meta = {
  title: 'Basic Components/ToggleGroup',
  component: ToggleGroup,
  parameters: {
    layout: 'centered',
  },
} satisfies Meta<typeof ToggleGroup>

export default meta

export const Default: StoryObj = {
  render: () => (
    <ToggleGroup type="multiple">
      <ToggleGroupItem value="bold" aria-label="Toggle bold">
        <Lucide.Bold className="size-4" />
      </ToggleGroupItem>
      <ToggleGroupItem value="italic" aria-label="Toggle italic">
        <Lucide.Italic className="size-4" />
      </ToggleGroupItem>
      <ToggleGroupItem value="underline" aria-label="Toggle underline">
        <Lucide.Underline className="size-4" />
      </ToggleGroupItem>
    </ToggleGroup>
  ),
}

export const SingleSelect: StoryObj = {
  render: () => (
    <ToggleGroup type="single">
      <ToggleGroupItem value="left" aria-label="Align left">
        <Lucide.AlignLeft className="size-4" />
      </ToggleGroupItem>
      <ToggleGroupItem value="center" aria-label="Align center">
        <Lucide.AlignCenter className="size-4" />
      </ToggleGroupItem>
      <ToggleGroupItem value="right" aria-label="Align right">
        <Lucide.AlignRight className="size-4" />
      </ToggleGroupItem>
    </ToggleGroup>
  ),
}

export const Variants: StoryObj = {
  render: () => (
    <div className="flex flex-col gap-4">
      <ToggleGroup type="multiple" variant="default">
        <ToggleGroupItem value="bold">
          <Lucide.Bold className="size-4" />
        </ToggleGroupItem>
        <ToggleGroupItem value="italic">
          <Lucide.Italic className="size-4" />
        </ToggleGroupItem>
      </ToggleGroup>

      <ToggleGroup type="multiple" variant="outline">
        <ToggleGroupItem value="bold">
          <Lucide.Bold className="size-4" />
        </ToggleGroupItem>
        <ToggleGroupItem value="italic">
          <Lucide.Italic className="size-4" />
        </ToggleGroupItem>
      </ToggleGroup>
    </div>
  ),
}
