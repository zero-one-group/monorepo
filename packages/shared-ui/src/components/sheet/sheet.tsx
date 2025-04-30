import * as Lucide from 'lucide-react'
import { Dialog as SheetPrimitive } from 'radix-ui'
import * as React from 'react'
import { sheetStyles } from './sheet.css'
import type { SheetVariants } from './sheet.css'

const Sheet = SheetPrimitive.Root
const SheetTrigger = SheetPrimitive.Trigger
const SheetClose = SheetPrimitive.Close
const SheetPortal = SheetPrimitive.Portal

const SheetOverlay = React.forwardRef<
  React.ComponentRef<typeof SheetPrimitive.Overlay>,
  React.ComponentPropsWithoutRef<typeof SheetPrimitive.Overlay>
>(({ className, ...props }, ref) => {
  const styles = sheetStyles()
  return <SheetPrimitive.Overlay className={styles.overlay({ className })} {...props} ref={ref} />
})

interface SheetContentProps
  extends React.ComponentPropsWithoutRef<typeof SheetPrimitive.Content>,
    SheetVariants {}

const SheetContent = React.forwardRef<
  React.ComponentRef<typeof SheetPrimitive.Content>,
  SheetContentProps
>(({ side = 'right', className, children, ...props }, ref) => {
  const styles = sheetStyles({ side })
  return (
    <SheetPortal>
      <SheetOverlay />
      <SheetPrimitive.Content ref={ref} className={styles.base({ className })} {...props}>
        <SheetPrimitive.Close className={styles.contentCloseWrapper()}>
          <Lucide.XIcon className={styles.contentCloseIcon()} strokeWidth={2} />
          <span className="sr-only">Close</span>
        </SheetPrimitive.Close>
        {children}
      </SheetPrimitive.Content>
    </SheetPortal>
  )
})

const SheetHeader = ({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) => {
  const styles = sheetStyles()
  return <div className={styles.header({ className })} {...props} />
}

const SheetFooter = ({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) => {
  const styles = sheetStyles()
  return <div className={styles.footer({ className })} {...props} />
}

const SheetTitle = React.forwardRef<
  React.ComponentRef<typeof SheetPrimitive.Title>,
  React.ComponentPropsWithoutRef<typeof SheetPrimitive.Title>
>(({ className, ...props }, ref) => {
  const styles = sheetStyles()
  return <SheetPrimitive.Title ref={ref} className={styles.title({ className })} {...props} />
})

const SheetDescription = React.forwardRef<
  React.ComponentRef<typeof SheetPrimitive.Description>,
  React.ComponentPropsWithoutRef<typeof SheetPrimitive.Description>
>(({ className, ...props }, ref) => {
  const styles = sheetStyles()
  return (
    <SheetPrimitive.Description
      className={styles.description({ className })}
      ref={ref}
      {...props}
    />
  )
})

SheetOverlay.displayName = SheetPrimitive.Overlay.displayName
SheetContent.displayName = SheetPrimitive.Content.displayName
SheetHeader.displayName = 'SheetHeader'
SheetFooter.displayName = 'SheetFooter'
SheetTitle.displayName = SheetPrimitive.Title.displayName
SheetDescription.displayName = SheetPrimitive.Description.displayName

export {
  Sheet,
  SheetPortal,
  SheetOverlay,
  SheetTrigger,
  SheetClose,
  SheetContent,
  SheetHeader,
  SheetFooter,
  SheetTitle,
  SheetDescription,
}
