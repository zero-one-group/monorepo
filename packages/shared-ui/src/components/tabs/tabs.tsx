import { Tabs as TabsPrimitive } from 'radix-ui'
import * as React from 'react'

import type { TabsVariants } from './tabs.css'
import { tabsStyles } from './tabs.css'

const Tabs = TabsPrimitive.Root as React.ComponentType<
  React.ComponentPropsWithoutRef<typeof TabsPrimitive.Root>
>

export interface TabsListProps
  extends React.ComponentPropsWithoutRef<typeof TabsPrimitive.List>, TabsVariants {}

const TabsList = React.forwardRef<React.ComponentRef<typeof TabsPrimitive.List>, TabsListProps>(
  ({ className, ...props }, ref) => {
    const styles = tabsStyles()
    return <TabsPrimitive.List ref={ref} className={styles.list({ className })} {...props} />
  }
)

const TabsTrigger: React.ComponentType<
  React.ComponentPropsWithoutRef<typeof TabsPrimitive.Trigger>
> = React.forwardRef<
  React.ComponentRef<typeof TabsPrimitive.Trigger>,
  React.ComponentPropsWithoutRef<typeof TabsPrimitive.Trigger>
>(({ className, ...props }, ref) => {
  const styles = tabsStyles()
  return <TabsPrimitive.Trigger ref={ref} className={styles.trigger({ className })} {...props} />
})

const TabsContent: React.ComponentType<
  React.ComponentPropsWithoutRef<typeof TabsPrimitive.Content>
> = React.forwardRef<
  React.ComponentRef<typeof TabsPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof TabsPrimitive.Content>
>(({ className, ...props }, ref) => {
  const styles = tabsStyles()
  return <TabsPrimitive.Content ref={ref} className={styles.content({ className })} {...props} />
})

Tabs.displayName = TabsPrimitive.Root.displayName
TabsList.displayName = TabsPrimitive.List.displayName
TabsTrigger.displayName = TabsPrimitive.Trigger.displayName
TabsContent.displayName = TabsPrimitive.Content.displayName

export { Tabs, TabsList, TabsTrigger, TabsContent }
