import { type VariantProps, tv } from 'tailwind-variants'

export const dropdownMenuStyles = tv({
  slots: {
    trigger: [
      'flex cursor-default select-none items-center rounded-xs px-2 py-1.5 text-sm outline-none focus:bg-accent',
      ' data-[state=open]:bg-accent',
    ],
    triggerIcon: 'ml-auto size-4',
    content: [
      'min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md',
      'data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 data-[state=open]:animate-in',
      'data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=closed]:animate-out',
      'data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 z-50',
      'data-[side=bottom]:slide-in-from-top-2 data-[side=top]:slide-in-from-bottom-2',
    ],
    item: [
      'relative flex cursor-default select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none',
      'transition-colors focus:bg-accent focus:text-accent-foreground [&>svg]:size-4 [&>svg]:shrink-0',
      'data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
    ],
    checkboxItem: [
      'relative flex cursor-default select-none items-center rounded-xs py-1.5 pr-2 pl-8 text-sm',
      'outline-none transition-colors focus:bg-accent focus:text-accent-foreground',
      'data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
    ],
    checkboxItemIndicator: 'absolute left-2 flex size-3.5 items-center justify-center',
    checkboxItemIcon: 'size-4',
    radioItem: [
      'relative flex cursor-default select-none items-center rounded-xs py-1.5 pr-2 pl-8 text-sm',
      'outline-none transition-colors focus:bg-accent focus:text-accent-foreground',
      'data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
    ],
    radioItemIndicator: 'absolute left-2 flex size-3.5 items-center justify-center',
    radioItemIcon: 'size-4 fill-current',
    label: 'px-2 py-1.5 font-semibold text-sm',
    separator: '-mx-1 my-1 h-px bg-muted',
    shortcut: 'ml-auto text-xs tracking-widest opacity-60',
    subTrigger: [
      'flex cursor-default select-none items-center rounded-xs px-2 py-1.5 text-sm outline-none focus:bg-accent',
      'data-[state=open]:bg-accent',
    ],
    subContent: [
      'z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-lg',
      'data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 data-[state=open]:animate-in',
      'data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=closed]:animate-out',
      'data-[side=bottom]:slide-in-from-top-2 data-[side=top]:slide-in-from-bottom-2',
      'data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2',
    ],
  },
  variants: {
    inset: {
      true: {
        item: 'pl-8',
        label: 'pl-8',
        subTrigger: 'pl-8',
      },
    },
  },
})

export type DropdownMenuVariants = VariantProps<typeof dropdownMenuStyles>
