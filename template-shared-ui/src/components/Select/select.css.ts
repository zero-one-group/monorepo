import { type VariantProps, tv } from 'tailwind-variants'

export const selectStyles = tv({
  slots: {
    trigger: [
      'flex h-9 w-full items-center justify-between whitespace-nowrap rounded-md border border-input',
      'bg-transparent px-3 py-2 text-sm shadow-xs ring-offset-background placeholder:text-muted-foreground',
      'focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50',
      '[&>span]:line-clamp-1',
    ],
    icon: 'size-4 opacity-50',
    content: [
      'relative z-50 max-h-96 min-w-[8rem] overflow-hidden rounded-md border bg-popover text-popover-foreground shadow-md',
      'data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 data-[state=open]:animate-in',
      'data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=closed]:animate-out',
      'data-[side=top]:slide-in-from-bottom-2 data-[side=bottom]:slide-in-from-top-2',
      'data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2',
    ],
    viewport: 'p-1',
    viewportPopper:
      'h-[var(--radix-select-trigger-height)] w-full min-w-[var(--radix-select-trigger-width)]',
    scrollButton: 'flex cursor-default items-center justify-center py-1',
    label: 'px-2 py-1.5 font-semibold text-sm',
    item: [
      'relative flex w-full cursor-default select-none items-center rounded-xs py-1.5 pr-8 pl-2',
      'text-sm outline-none focus:bg-accent focus:text-accent-foreground',
      'data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
    ],
    itemIndicator: 'absolute right-2 flex size-3.5 items-center justify-center',
    itemIndicatorIcon: 'size-4',
    separator: '-mx-1 my-1 h-px bg-muted',
  },
})

export type SelectVariants = VariantProps<typeof selectStyles>
