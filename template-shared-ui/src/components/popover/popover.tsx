import { Popover as PopoverPrimitive } from 'radix-ui'
import * as React from 'react'
import { popoverStyles } from './popover.css'

const Popover = PopoverPrimitive.Root
const PopoverTrigger = PopoverPrimitive.Trigger
const PopoverAnchor = PopoverPrimitive.Anchor

const PopoverContent = React.forwardRef<
  React.ComponentRef<typeof PopoverPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof PopoverPrimitive.Content>
>(({ className, align = 'center', sideOffset = 4, ...props }, ref) => {
  return (
    <PopoverPrimitive.Portal>
      <PopoverPrimitive.Content
        ref={ref}
        align={align}
        sideOffset={sideOffset}
        className={popoverStyles({ className })}
        {...props}
      />
    </PopoverPrimitive.Portal>
  )
})

Popover.displayName = 'Popover'
PopoverTrigger.displayName = 'PopoverTrigger'
PopoverAnchor.displayName = 'PopoverAnchor'
PopoverContent.displayName = PopoverPrimitive.Content.displayName

export { Popover, PopoverTrigger, PopoverContent, PopoverAnchor }
