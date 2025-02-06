import { type VariantProps, tv } from 'tailwind-variants'

export const skeletonStyles = tv({
  base: 'animate-pulse rounded-md bg-foreground/10',
})

export type SkeletonVariants = VariantProps<typeof skeletonStyles>
