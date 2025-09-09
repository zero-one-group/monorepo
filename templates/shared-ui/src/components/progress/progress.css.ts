import { type VariantProps, tv } from 'tailwind-variants'

export const progressStyles = tv({
  slots: {
    base: 'relative h-2 w-full overflow-hidden rounded-full bg-primary/20',
    indicator: 'size-full flex-1 bg-primary transition-all',
  },
  variants: {
    size: {
      default: {
        base: 'h-2',
      },
      sm: {
        base: 'h-1',
      },
      lg: {
        base: 'h-3',
      },
    },
  },
  defaultVariants: {
    size: 'default',
  },
})

export type ProgressVariants = VariantProps<typeof progressStyles>
