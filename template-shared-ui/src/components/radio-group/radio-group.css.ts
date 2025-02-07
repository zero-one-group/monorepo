import { type VariantProps, tv } from 'tailwind-variants'

export const radioGroupStyles = tv({
  slots: {
    root: 'grid gap-2',
    item: 'aspect-square size-4 rounded-full border border-primary text-primary shadow-sm focus:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50',
    indicator: 'flex items-center justify-center',
    icon: 'size-3.5',
  },
  variants: {
    orientation: {
      horizontal: {
        root: 'flex gap-4',
      },
      vertical: {
        root: 'grid gap-2',
      },
    },
  },
  defaultVariants: {
    orientation: 'vertical',
  },
})

export type RadioGroupVariants = VariantProps<typeof radioGroupStyles>
