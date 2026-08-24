import * as ToggleGroupPrimitive from '@radix-ui/react-toggle-group'
import * as React from 'react'

import type { ToggleVariants } from '../toggle/toggle.css'
import { toggleStyles } from '../toggle/toggle.css'
import { toggleGroupStyles } from './toggle-group.css'

const ToggleGroupContext = React.createContext<ToggleVariants>({
  size: 'default',
  variant: 'default',
})

type ToggleGroupProps = React.ComponentPropsWithoutRef<'div'> & ToggleVariants

const ToggleGroup = React.forwardRef<HTMLDivElement, ToggleGroupProps>(
  ({ className, variant, size, children, ...props }, ref) => {
    const Root = ToggleGroupPrimitive.Root
    return (
      // @ts-expect-error - hoisted Radix types resolve children with React 18 ReactNode (bigint not assignable)
      <Root ref={ref} className={toggleGroupStyles({ className })} {...props}>
        <ToggleGroupContext.Provider value={{ variant, size }}>
          {children}
        </ToggleGroupContext.Provider>
      </Root>
    )
  }
)

type ToggleGroupItemProps = React.ComponentPropsWithoutRef<'button'> & ToggleVariants

const ToggleGroupItem = React.forwardRef<HTMLButtonElement, ToggleGroupItemProps>(
  ({ className, children, variant, size, ...props }, ref) => {
    const context = React.useContext(ToggleGroupContext)
    const styles = toggleStyles({
      variant: context.variant || variant,
      size: context.size || size,
      className,
    })

    const Item = ToggleGroupPrimitive.Item
    return (
      <Item ref={ref} className={styles} {...props}>
        {/* @ts-expect-error - hoisted Radix types resolve children with React 18 ReactNode (bigint not assignable) */}
        {children}
      </Item>
    )
  }
)

ToggleGroup.displayName = ToggleGroupPrimitive.Root.displayName
ToggleGroupItem.displayName = ToggleGroupPrimitive.Item.displayName

export { ToggleGroup, ToggleGroupItem }
