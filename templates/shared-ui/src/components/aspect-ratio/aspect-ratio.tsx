import { AspectRatio as AspectRatioPrimitive } from 'radix-ui'
import * as React from 'react'
import { type AspectRatioVariants, aspectRatioStyles } from './aspect-ratio.css'

const AspectRatio = React.forwardRef<
  React.ComponentRef<typeof AspectRatioPrimitive.Root>,
  React.ComponentPropsWithoutRef<typeof AspectRatioPrimitive.Root> & AspectRatioVariants
>(({ className, ...props }, ref) => {
  return (
    <AspectRatioPrimitive.Root ref={ref} className={aspectRatioStyles({ className })} {...props} />
  )
})

AspectRatio.displayName = 'AspectRatio'

export { AspectRatio }
