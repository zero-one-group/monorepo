import { type VariantProps, tv } from 'tailwind-variants'

export const resizableStyles = tv({
  slots: {
    group: 'flex size-full data-[panel-group-direction=vertical]:flex-col',
    groupFixed: 'fixed inset-0',
    handle:
      'group relative flex w-px items-center justify-center bg-border transition-colors hover:bg-foreground/20 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1',
    handleVertical:
      'data-[panel-group-direction=vertical]:h-px data-[panel-group-direction=vertical]:w-full',
    handleAfter: 'after:-translate-x-1/2 after:absolute after:inset-y-0 after:left-1/2 after:w-1',
    handleAfterVertical:
      'data-[panel-group-direction=vertical]:after:left-0 data-[panel-group-direction=vertical]:after:h-1 data-[panel-group-direction=vertical]:after:w-full data-[panel-group-direction=vertical]:after:translate-x-0',
    handleRotate: '[&[data-panel-group-direction=vertical]>div]:rotate-90',
    handleButton:
      'z-10 flex h-4 w-3 items-center justify-center rounded-xs border bg-border opacity-0 transition-opacity group-hover:opacity-100',
    handleIcon: 'size-2.5',
  },
})

export type ResizableVariants = VariantProps<typeof resizableStyles>
