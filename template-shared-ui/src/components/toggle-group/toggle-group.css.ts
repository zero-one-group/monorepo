import { type VariantProps, tv } from 'tailwind-variants'

export const toggleGroupStyles = tv({
  base: 'flex items-center justify-center gap-1',
})

export type ToggleGroupVariants = VariantProps<typeof toggleGroupStyles>
