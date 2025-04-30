import { Slot } from 'radix-ui'
import * as React from 'react'
import { clx } from '../../utils'
import { Skeleton } from '../skeleton/skeleton'
import { Tooltip, TooltipContent, TooltipTrigger } from '../tooltip/tooltip'
import { useSidebar } from './sidebar-provider'
import { type SidebarMenuButtonVariants, sidebarMenuButtonStyles } from './sidebar.css'
import { sidebarStyles } from './sidebar.css'

const SidebarMenu = React.forwardRef<HTMLUListElement, React.ComponentProps<'ul'>>(
  ({ className, ...props }, ref) => {
    const styles = sidebarStyles()
    return (
      <ul ref={ref} data-sidebar="menu" className={styles.sidebarMenu({ className })} {...props} />
    )
  }
)

const SidebarMenuItem = React.forwardRef<HTMLLIElement, React.ComponentProps<'li'>>(
  ({ className, ...props }, ref) => {
    const styles = sidebarStyles()
    return (
      <li
        ref={ref}
        data-sidebar="menu-item"
        className={styles.sidebarMenuItem({ className })}
        {...props}
      />
    )
  }
)

const SidebarMenuButton = React.forwardRef<
  HTMLButtonElement,
  React.ComponentProps<'button'> & {
    asChild?: boolean
    isActive?: boolean
    tooltip?: string | React.ComponentProps<typeof TooltipContent>
  } & SidebarMenuButtonVariants
>(
  (
    {
      asChild = false,
      isActive = false,
      variant = 'default',
      size = 'default',
      tooltip,
      className,
      ...props
    },
    ref
  ) => {
    const Comp = asChild ? Slot.Root : 'button'
    const { isMobile, state } = useSidebar()

    const button = (
      <Comp
        ref={ref}
        data-sidebar="menu-button"
        data-size={size}
        data-active={isActive}
        className={sidebarMenuButtonStyles({ variant, size, className })}
        {...props}
      />
    )

    if (!tooltip) {
      return button
    }

    if (typeof tooltip === 'string') {
      tooltip = {
        children: tooltip,
      }
    }

    return (
      <Tooltip>
        <TooltipTrigger asChild>{button}</TooltipTrigger>
        <TooltipContent
          side="right"
          align="center"
          hidden={state !== 'collapsed' || isMobile}
          {...tooltip}
        />
      </Tooltip>
    )
  }
)

const SidebarMenuAction = React.forwardRef<
  HTMLButtonElement,
  React.ComponentProps<'button'> & {
    asChild?: boolean
    showOnHover?: boolean
  }
>(({ className, asChild = false, showOnHover = false, ...props }, ref) => {
  const Comp = asChild ? Slot.Root : 'button'
  const styles = sidebarStyles({ showOnHover })

  return (
    <Comp
      ref={ref}
      data-sidebar="menu-action"
      className={styles.sidebarMenuAction({ className })}
      {...props}
    />
  )
})

const SidebarMenuBadge = React.forwardRef<HTMLDivElement, React.ComponentProps<'div'>>(
  ({ className, ...props }, ref) => {
    const styles = sidebarStyles()

    return (
      <div
        ref={ref}
        data-sidebar="menu-badge"
        className={styles.sidebarMenuBadge({ className })}
        {...props}
      />
    )
  }
)

const SidebarMenuSkeleton = React.forwardRef<
  HTMLDivElement,
  React.ComponentProps<'div'> & {
    showIcon?: boolean
  }
>(({ className, showIcon = false, ...props }, ref) => {
  // Random width between 50 to 90%.
  const width = React.useMemo(() => {
    return `${Math.floor(Math.random() * 40) + 50}%`
  }, [])

  const styles = sidebarStyles()

  return (
    <div
      ref={ref}
      data-sidebar="menu-skeleton"
      className={styles.sidebarMenuSkeletonWrapper({ className })}
      {...props}
    >
      {showIcon && (
        <Skeleton className={styles.sidebarMenuSkeletonIcon()} data-sidebar="menu-skeleton-icon" />
      )}
      <Skeleton
        className={styles.sidebarMenuSkeletonText()}
        style={% raw %}{{ '--skeleton-width': width } as React.CSSProperties}{% endraw %}
        data-sidebar="menu-skeleton-text"
      />
    </div>
  )
})

const SidebarMenuSub = React.forwardRef<HTMLUListElement, React.ComponentProps<'ul'>>(
  ({ className, ...props }, ref) => {
    const styles = sidebarStyles()
    return (
      <ul
        ref={ref}
        data-sidebar="menu-sub"
        className={styles.sidebarMenuSub({ className })}
        {...props}
      />
    )
  }
)

const SidebarMenuSubItem = React.forwardRef<HTMLLIElement, React.ComponentProps<'li'>>(
  ({ ...props }, ref) => <li ref={ref} {...props} />
)

const SidebarMenuSubButton = React.forwardRef<
  HTMLAnchorElement,
  React.ComponentProps<'a'> & {
    asChild?: boolean
    size?: 'sm' | 'md'
    isActive?: boolean
  }
>(({ asChild = false, size = 'md', isActive, className, ...props }, ref) => {
  const Comp = asChild ? Slot.Root : 'a'
  // const styles = sidebarStyles()

  return (
    <Comp
      ref={ref}
      data-sidebar="menu-sub-button"
      data-size={size}
      data-active={isActive}
      className={clx(
        '-translate-x-px flex h-7 min-w-0 items-center gap-2 overflow-hidden rounded-md px-2 text-sidebar-foreground outline-none ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus:ring-0 focus-visible:ring-1 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:text-sidebar-accent-foreground',
        'data-[active=true]:bg-sidebar-accent data-[active=true]:text-sidebar-accent-foreground',
        size === 'sm' && 'text-xs',
        size === 'md' && 'text-sm',
        'group-data-[collapsible=icon]:hidden',
        className
      )}
      {...props}
    />
  )
})

SidebarMenu.displayName = 'SidebarMenu'
SidebarMenuItem.displayName = 'SidebarMenuItem'
SidebarMenuButton.displayName = 'SidebarMenuButton'
SidebarMenuAction.displayName = 'SidebarMenuAction'
SidebarMenuBadge.displayName = 'SidebarMenuBadge'
SidebarMenuSkeleton.displayName = 'SidebarMenuSkeleton'
SidebarMenuSub.displayName = 'SidebarMenuSub'
SidebarMenuSubItem.displayName = 'SidebarMenuSubItem'
SidebarMenuSubButton.displayName = 'SidebarMenuSubButton'

export {
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSkeleton,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
}
