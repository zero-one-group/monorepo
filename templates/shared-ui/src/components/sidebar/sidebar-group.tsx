import { Slot } from 'radix-ui'
import * as React from 'react'
import { sidebarStyles } from './sidebar.css'

const SidebarGroup = React.forwardRef<HTMLDivElement, React.ComponentProps<'div'>>(
  ({ className, ...props }, ref) => {
    const styles = sidebarStyles()

    return (
      <div
        ref={ref}
        data-sidebar="group"
        className={styles.sidebarGroup({ className })}
        {...props}
      />
    )
  }
)

const SidebarGroupLabel = React.forwardRef<
  HTMLDivElement,
  React.ComponentProps<'div'> & { asChild?: boolean }
>(({ className, asChild = false, ...props }, ref) => {
  const Comp = asChild ? Slot.Root : 'div'
  const styles = sidebarStyles()

  return (
    <Comp
      ref={ref}
      data-sidebar="group-label"
      className={styles.sidebarGroupLabel({ className })}
      {...props}
    />
  )
})

const SidebarGroupAction = React.forwardRef<
  HTMLButtonElement,
  React.ComponentProps<'button'> & { asChild?: boolean }
>(({ className, asChild = false, ...props }, ref) => {
  const Comp = asChild ? Slot.Root : 'button'
  const styles = sidebarStyles()

  return (
    <Comp
      ref={ref}
      data-sidebar="group-action"
      className={styles.sidebarGroupAction({ className })}
      {...props}
    />
  )
})

const SidebarGroupContent = React.forwardRef<HTMLDivElement, React.ComponentProps<'div'>>(
  ({ className, ...props }, ref) => {
    const styles = sidebarStyles()

    return (
      <div
        ref={ref}
        data-sidebar="group-content"
        className={styles.sidebarGroupContent({ className })}
        {...props}
      />
    )
  }
)

SidebarGroup.displayName = 'SidebarGroup'
SidebarGroupLabel.displayName = 'SidebarGroupLabel'
SidebarGroupAction.displayName = 'SidebarGroupAction'
SidebarGroupContent.displayName = 'SidebarGroupContent'

export { SidebarGroup, SidebarGroupAction, SidebarGroupContent, SidebarGroupLabel }
