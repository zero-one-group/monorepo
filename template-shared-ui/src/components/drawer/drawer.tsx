import * as React from 'react'
import { Drawer as DrawerPrimitive } from 'vaul'
import { drawerStyles } from './drawer.css'

const Drawer = ({
  shouldScaleBackground = true,
  ...props
}: React.ComponentProps<typeof DrawerPrimitive.Root>) => (
  <DrawerPrimitive.Root shouldScaleBackground={shouldScaleBackground} {...props} />
)

const DrawerTrigger = DrawerPrimitive.Trigger
const DrawerPortal = DrawerPrimitive.Portal
const DrawerClose = DrawerPrimitive.Close

const DrawerOverlay = React.forwardRef<
  React.ComponentRef<typeof DrawerPrimitive.Overlay>,
  React.ComponentPropsWithoutRef<typeof DrawerPrimitive.Overlay>
>(({ className, ...props }, ref) => {
  const styles = drawerStyles()
  return <DrawerPrimitive.Overlay ref={ref} className={styles.overlay({ className })} {...props} />
})

const DrawerContent = React.forwardRef<
  React.ComponentRef<typeof DrawerPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof DrawerPrimitive.Content>
>(({ className, children, ...props }, ref) => {
  const styles = drawerStyles()
  return (
    <DrawerPortal>
      <DrawerOverlay />
      <DrawerPrimitive.Content ref={ref} className={styles.content({ className })} {...props}>
        <div className={styles.handle()} />
        {children}
      </DrawerPrimitive.Content>
    </DrawerPortal>
  )
})

const DrawerHeader = ({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) => {
  const styles = drawerStyles()
  return <div className={styles.header({ className })} {...props} />
}

const DrawerFooter = ({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) => {
  const styles = drawerStyles()
  return <div className={styles.footer({ className })} {...props} />
}

const DrawerTitle = React.forwardRef<
  React.ComponentRef<typeof DrawerPrimitive.Title>,
  React.ComponentPropsWithoutRef<typeof DrawerPrimitive.Title>
>(({ className, ...props }, ref) => {
  const styles = drawerStyles()
  return <DrawerPrimitive.Title ref={ref} className={styles.title({ className })} {...props} />
})

const DrawerDescription = React.forwardRef<
  React.ComponentRef<typeof DrawerPrimitive.Description>,
  React.ComponentPropsWithoutRef<typeof DrawerPrimitive.Description>
>(({ className, ...props }, ref) => {
  const styles = drawerStyles()
  return (
    <DrawerPrimitive.Description
      ref={ref}
      className={styles.description({ className })}
      {...props}
    />
  )
})

Drawer.displayName = 'Drawer'
DrawerTrigger.displayName = 'DrawerTrigger'
DrawerClose.displayName = 'DrawerClose'
DrawerOverlay.displayName = DrawerPrimitive.Overlay.displayName
DrawerContent.displayName = 'DrawerContent'
DrawerHeader.displayName = 'DrawerHeader'
DrawerFooter.displayName = 'DrawerFooter'
DrawerTitle.displayName = DrawerPrimitive.Title.displayName
DrawerDescription.displayName = DrawerPrimitive.Description.displayName

export {
  Drawer,
  DrawerPortal,
  DrawerOverlay,
  DrawerTrigger,
  DrawerClose,
  DrawerContent,
  DrawerHeader,
  DrawerFooter,
  DrawerTitle,
  DrawerDescription,
}
