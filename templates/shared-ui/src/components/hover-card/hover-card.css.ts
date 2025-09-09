import { type VariantProps, tv } from 'tailwind-variants'

export const hoverCardStyles = tv({
  slots: {
    content: [
      'z-50 w-64 rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-none',
      // Entry animations
      'data-[state=open]:motion-safe:motion-opacity-in-0',
      'data-[state=open]:motion-safe:motion-duration-150',
      // Exit animations
      'data-[state=closed]:motion-safe:motion-opacity-out-0',
      'data-[state=closed]:motion-safe:motion-duration-100',
      // Slide animations based on position
      'data-[side=top]:motion-safe:motion-translate-y-in-2',
      'data-[side=bottom]:motion-safe:motion-translate-y-in--2',
      'data-[side=left]:motion-safe:motion-translate-x-in-2',
      'data-[side=right]:motion-safe:motion-translate-x-in--2',
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
