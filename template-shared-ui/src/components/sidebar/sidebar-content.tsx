import * as React from 'react'
import { clx } from '../../utils'
import { Input } from '../input/input'
import { ScrollArea } from '../scroll-area/scroll-area'
import { Separator } from '../separator/separator'
import { useSidebar } from './sidebar-provider'

const SidebarRail = React.forwardRef<HTMLButtonElement, React.ComponentProps<'button'>>(
  ({ className, ...props }, ref) => {
    const { toggleSidebar } = useSidebar()

    return (
      <button
        ref={ref}
        data-sidebar="rail"
        aria-label="Toggle Sidebar"
        tabIndex={-1}
        onClick={toggleSidebar}
        title="Toggle Sidebar"
        className={clx(
          '-translate-x-1/2 group-data-[side=left]:-right-4 absolute inset-y-0 z-20 hidden w-4 transition-all ease-linear after:absolute after:inset-y-0 after:left-1/2 after:w-[2px] hover:after:bg-sidebar-border group-data-[side=right]:left-0 sm:flex',
          '[[data-side=left]_&]:cursor-w-resize [[data-side=right]_&]:cursor-e-resize',
          '[[data-side=left][data-state=collapsed]_&]:cursor-e-resize [[data-side=right][data-state=collapsed]_&]:cursor-w-resize',
          'group-data-[collapsible=offcanvas]:translate-x-0 group-data-[collapsible=offcanvas]:hover:bg-sidebar group-data-[collapsible=offcanvas]:after:left-full',
          '[[data-side=left][data-collapsible=offcanvas]_&]:-right-2',
          '[[data-side=right][data-collapsible=offcanvas]_&]:-left-2',
          className
        )}
        {...props}
      />
    )
  }
)

const SidebarInset = React.forwardRef<HTMLDivElement, React.ComponentProps<'main'>>(
  ({ className, ...props }, ref) => {
    return (
      <main
        ref={ref}
        className={clx(
          'relative flex min-h-svh flex-1 flex-col bg-background',
          'peer-data-[variant=inset]:min-h-[calc(100svh-theme(spacing.4))] md:peer-data-[state=collapsed]:peer-data-[variant=inset]:ml-2 md:peer-data-[variant=inset]:m-2 md:peer-data-[variant=inset]:ml-0 md:peer-data-[variant=inset]:rounded-xl md:peer-data-[variant=inset]:shadow',
          className
        )}
        {...props}
      />
    )
  }
)

const SidebarInput = React.forwardRef<
  React.ComponentRef<typeof Input>,
  React.ComponentProps<typeof Input>
>(({ className, ...props }, ref) => {
  return (
    <Input
      ref={ref}
      data-sidebar="input"
      className={clx(
        'h-8 w-full bg-background shadow-none focus:ring-0 focus-visible:ring-1 focus-visible:ring-primary/50',
        className
      )}
      {...props}
    />
  )
})

const SidebarHeader = React.forwardRef<HTMLDivElement, React.ComponentProps<'div'>>(
  ({ className, ...props }, ref) => {
    return (
      <div
        ref={ref}
        data-sidebar="header"
        className={clx('flex flex-col gap-2 p-2', className)}
        {...props}
      />
    )
  }
)

const SidebarFooter = React.forwardRef<HTMLDivElement, React.ComponentProps<'div'>>(
  ({ className, ...props }, ref) => {
    return (
      <div
        ref={ref}
        data-sidebar="footer"
        className={clx('flex flex-col gap-2 p-2', className)}
        {...props}
      />
    )
  }
)

const SidebarSeparator: React.ForwardRefExoticComponent<
  React.ComponentPropsWithRef<typeof Separator>
> = React.forwardRef<React.ComponentRef<typeof Separator>, React.ComponentProps<typeof Separator>>(
  ({ className, ...props }, ref) => {
    return (
      <Separator
        ref={ref}
        data-sidebar="separator"
        className={clx('mx-2 w-auto bg-sidebar-border', className)}
        {...props}
      />
    )
  }
)

const SidebarContent = React.forwardRef<HTMLDivElement, React.ComponentProps<'div'>>(
  ({ className, ...props }, ref) => {
    return (
      <ScrollArea
        className={clx(
          'flex min-h-0 flex-1 flex-col gap-2 overflow-auto group-data-[collapsible=icon]:overflow-hidden',
          className
        )}
      >
        <div ref={ref} data-sidebar="content" {...props} />
      </ScrollArea>
    )
  }
)

SidebarRail.displayName = 'SidebarRail'
SidebarInset.displayName = 'SidebarInset'
SidebarInput.displayName = 'SidebarInput'
SidebarHeader.displayName = 'SidebarHeader'
SidebarFooter.displayName = 'SidebarFooter'
SidebarSeparator.displayName = 'SidebarSeparator'
SidebarContent.displayName = 'SidebarContent'

export {
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarInput,
  SidebarInset,
  SidebarRail,
  SidebarSeparator,
}
