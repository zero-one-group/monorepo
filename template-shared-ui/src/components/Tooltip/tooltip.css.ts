import { type VariantProps, tv } from 'tailwind-variants'

export const tooltipStyles = tv({
  slots: {
    base: [
      'z-50 overflow-hidden rounded-lg bg-black px-3 py-1.5 text-white text-xs',
      'fade-in-0 zoom-in-95 animate-in duration-150 ease-in-out',
      'data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=closed]:animate-out',
      'data-[side=bottom]:slide-in-from-top-2 data-[side=top]:slide-in-from-bottom-2',
      'data-[side=right]:slide-in-from-left-2 data-[side=left]:slide-in-from-right-2',
    ],
    arrow: 'transition-transform duration-150 ease-in-out',
  },
})

export type TooltipVariants = VariantProps<typeof tooltipStyles>
