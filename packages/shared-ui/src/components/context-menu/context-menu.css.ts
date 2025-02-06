import { type VariantProps, tv } from 'tailwind-variants'

export const contextMenuStyles = tv({
  slots: {
    trigger: [
      'flex cursor-default select-none items-center rounded-xs px-2 py-1.5 text-sm',
      'outline-none focus:bg-accent focus:text-accent-foreground',
      'data-[state=open]:bg-accent data-[state=open]:text-accent-foreground',
    ],
    content: [
      'z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md',
      'data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 data-[state=open]:animate-in',
      'data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=closed]:animate-out',
      'data-[side=bottom]:slide-in-from-top-2 data-[side=top]:slide-in-from-bottom-2',
      'data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2',
    ],
    item: [
      'relative flex cursor-default select-none items-center rounded-xs px-2 py-1.5 text-sm',
      'outline-none focus:bg-accent focus:text-accent-foreground',
      'data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
    ],
    checkboxItem: [
      'relative flex cursor-default select-none items-center rounded-xs py-1.5 pr-2 pl-8 text-sm',
      'outline-none focus:bg-accent focus:text-accent-foreground',
      'data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
    ],
    radioItem: [
      'relative flex cursor-default select-none items-center rounded-xs py-1.5 pr-2 pl-8 text-sm',
      'outline-none focus:bg-accent focus:text-accent-foreground',
      'data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
    ],
    label: 'px-2 py-1.5 font-semibold text-foreground text-sm',
    separator: '-mx-1 my-1 h-px bg-border',
    shortcut: 'ml-auto text-muted-foreground text-xs tracking-widest',
    subTrigger: [
      'flex cursor-default select-none items-center rounded-xs px-2 py-1.5 text-sm',
      'outline-none focus:bg-accent focus:text-accent-foreground',
      'data-[state=open]:bg-accent data-[state=open]:text-accent-foreground',
    ],
    subTriggerIcon: 'ml-auto size-4',
    subContent: [
      'z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-lg',
      'data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 data-[state=open]:animate-in',
      'data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=closed]:animate-out',
      'data-[side=bottom]:slide-in-from-top-2 data-[side=top]:slide-in-from-bottom-2',
      'data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2',
    ],
    itemIndicator: 'absolute left-2 flex h-3.5 w-3.5 items-center justify-center',
    itemIndicatorIcon: 'size-4',
    radioItemIcon: 'size-4 fill-current',
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

export type ContextMenuVariants = VariantProps<typeof contextMenuStyles>
