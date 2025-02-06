import type { Meta, StoryObj } from '@storybook/react'
import * as Lucide from 'lucide-react'
import { Button } from '../button/button'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from './collapsible'

const meta: Meta<typeof Collapsible> = {
  title: 'Basic Components/Collapsible',
  component: Collapsible,
  argTypes: {
    open: {
      control: 'boolean',
      description: 'The controlled open state of the collapsible',
    },
    defaultOpen: {
      control: 'boolean',
      description: 'The default open state when initially rendered',
    },
  },
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: () => (
    <Collapsible className="w-[350px] space-y-2">
      <div className="flex w-full items-center justify-between space-x-4 pr-0 pl-4">
        <h4 className="font-semibold text-sm">@riipandi starred 3 repositories</h4>
        <CollapsibleTrigger asChild>
          <Button variant="ghost" size="sm" className="w-9 p-0">
            <Lucide.ChevronsUpDown className="size-4" />
            <span className="sr-only">Toggle</span>
          </Button>
        </CollapsibleTrigger>
      </div>
      <div className="rounded-md border px-4 py-3 font-mono text-sm">@radix-ui/primitives</div>
      <CollapsibleContent className="space-y-2">
        <div className="rounded-md border px-4 py-3 font-mono text-sm">@radix-ui/colors</div>
        <div className="rounded-md border px-4 py-3 font-mono text-sm">@stitches/react</div>
      </CollapsibleContent>
    </Collapsible>
  ),
}

export const CollapsibleShowcase: Story = {
  render: () => (
    <div className="space-y-8">
      <div className="space-y-4">
        <h4 className="font-medium text-sm">Basic Collapsible</h4>
        <Collapsible className="w-[350px] space-y-2">
          <div className="flex items-center justify-between space-x-4">
            <h4 className="font-semibold text-sm">Show More</h4>
            <CollapsibleTrigger asChild>
              <Button variant="ghost" size="sm">
                <Lucide.ChevronDown className="size-4" />
              </Button>
            </CollapsibleTrigger>
          </div>
          <CollapsibleContent className="space-y-2">
            <div className="rounded-md border px-4 py-2">Hidden Content 1</div>
            <div className="rounded-md border px-4 py-2">Hidden Content 2</div>
          </CollapsibleContent>
        </Collapsible>
      </div>
    </div>
  ),
}
