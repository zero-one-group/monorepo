import { ScrollArea as ScrollAreaPrimitive } from 'radix-ui'
import * as React from 'react'
import { scrollAreaStyles } from './scroll-area.css'

const ScrollArea = React.forwardRef<
  React.ComponentRef<typeof ScrollAreaPrimitive.Root>,
  React.ComponentPropsWithoutRef<typeof ScrollAreaPrimitive.Root>
>(({ className, children, ...props }, ref) => {
  const styles = scrollAreaStyles()
  return (
    <ScrollAreaPrimitive.Root ref={ref} className={styles.root({ className })} {...props}>
      <ScrollAreaPrimitive.Viewport className={styles.viewport()}>
        {children}
      </ScrollAreaPrimitive.Viewport>
      <ScrollBar />
      <ScrollAreaPrimitive.Corner />
    </ScrollAreaPrimitive.Root>
  )
})

const ScrollBar = React.forwardRef<
  React.ComponentRef<typeof ScrollAreaPrimitive.ScrollAreaScrollbar>,
  React.ComponentPropsWithoutRef<typeof ScrollAreaPrimitive.ScrollAreaScrollbar>
>(({ className, orientation = 'vertical', ...props }, ref) => {
  const styles = scrollAreaStyles({ orientation })
  return (
    <ScrollAreaPrimitive.ScrollAreaScrollbar
      ref={ref}
      orientation={orientation}
      className={styles.scrollbar({ className })}
      {...props}
    >
      <ScrollAreaPrimitive.ScrollAreaThumb className={styles.thumb()} />
    </ScrollAreaPrimitive.ScrollAreaScrollbar>
  )
})

ScrollArea.displayName = ScrollAreaPrimitive.Root.displayName
ScrollBar.displayName = ScrollAreaPrimitive.ScrollAreaScrollbar.displayName

export { ScrollArea, ScrollBar }
