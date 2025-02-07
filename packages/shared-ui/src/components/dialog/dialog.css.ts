import { type VariantProps, tv } from 'tailwind-variants'

export const dialogStyles = tv({
  slots: {
    overlay: [
      'fixed inset-0 z-50 bg-black/80',
      'data-[state=open]:motion-preset-fade data-[state=open]:motion-duration-300',
      'data-[state=closed]:motion-opacity-out-0 motion-duration-[0.35s]/opacity',
    ],
    content: [
      'fixed top-[50%] left-[50%] z-50 grid w-full max-w-lg translate-x-[-50%]',
      'translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 sm:rounded-lg',
      'data-[state=open]:motion-scale-in-[0.5] data-[state=open]:motion-opacity-in-[0%]',
      'data-[state=open]:motion-blur-in-[5px] data-[state=open]:motion-duration-[0.00s]',
      'data-[state=closed]:motion-scale-out-[0.5] data-[state=closed]:motion-opacity-out-0',
      'data-[state=closed]:motion-blur-out-[5px] data-[state=closed]:motion-duration-[0.00s]',
      'motion-duration-[0.39s]/scale motion-duration-[0.35s]/opacity motion-duration-[0.35s]/blur',
    ],
    close: [
      'absolute top-4 right-4 rounded-xs opacity-70 ring-offset-background transition-opacity hover:opacity-100',
      'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none',
      'cursor-pointer data-[state=open]:bg-accent data-[state=open]:text-muted-foreground',
    ],
    closeIcon: 'size-4',
    header: 'flex flex-col space-y-1.5 text-center sm:text-left',
    footer: 'flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2',
    title: 'font-semibold text-lg leading-none tracking-tight',
    description: 'text-muted-foreground text-sm',
  },
})

export type DialogVariants = VariantProps<typeof dialogStyles>
