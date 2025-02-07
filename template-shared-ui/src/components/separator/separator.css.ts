import { type VariantProps, tv } from 'tailwind-variants'

export const separatorStyles = tv({
  base: 'shrink-0 bg-border',
  variants: {
    orientation: {
      horizontal: 'h-[1px] w-full',
      vertical: 'h-full w-[1px]',
    },
  },
  compoundVariants: [],
  defaultVariants: {
    orientation: 'horizontal',
  },
})

export type SeparatorVariants = VariantProps<typeof separatorStyles>
