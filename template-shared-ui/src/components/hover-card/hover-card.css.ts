import { type VariantProps, tv } from 'tailwind-variants'

export const hoverCardStyles = tv({
  slots: {
    content: [
      'z-50 w-64 rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-none',
      'data-[state=open]:zoom-in-95 data-[state=open]:fade-in-0 data-[state=open]:animate-in',
      'data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=closed]:animate-out',
      'data-[side=top]:slide-in-from-bottom-2 data-[side=bottom]:slide-in-from-top-2',
      'data-[side=right]:slide-in-from-left-2 data-[side=left]:slide-in-from-right-2',
    ],
  },
  variants: {
    align: {
      start: {},
      center: {},
      end: {},
    },
  },
  defaultVariants: {
    align: 'center',
  },
})

export type HoverCardVariants = VariantProps<typeof hoverCardStyles>
