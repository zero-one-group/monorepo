import * as ToggleGroupPrimitive from '@radix-ui/react-toggle-group'
import * as React from 'react'

import type { ToggleVariants } from '../toggle/toggle.css'
import { toggleStyles } from '../toggle/toggle.css'
import { toggleGroupStyles } from './toggle-group.css'

const ToggleGroupContext = React.createContext<ToggleVariants>({
  size: 'default',
  variant: 'default',
})

/**
 * Radix ToggleGroup props, redeclared locally instead of derived via
 * `React.ComponentPropsWithoutRef<typeof ToggleGroupPrimitive.Root>`.
 *
 * Deriving them triggers TS2742 when emitting declarations: the inferred type
 * references `@types/react` hoisted under `@radix-ui/*` (React 18), which is not
 * portable from this package. Keep this union in sync with
 * `ToggleGroupSingleProps` / `ToggleGroupMultipleProps` from
 * `@radix-ui/react-toggle-group`.
 */
type ToggleGroupValueProps =
  | {
      type: 'single'
      value?: string
      defaultValue?: string
      onValueChange?: (value: string) => void
    }
  | {
      type: 'multiple'
      value?: string[]
      defaultValue?: string[]
      onValueChange?: (value: string[]) => void
    }

type ToggleGroupProps = Omit<
  React.ComponentPropsWithoutRef<'div'>,
  'defaultValue' | 'dir'
> &
  ToggleGroupValueProps &
  ToggleVariants & {
    disabled?: boolean
    rovingFocus?: boolean
    loop?: boolean
    orientation?: 'horizontal' | 'vertical'
    dir?: 'ltr' | 'rtl'
  }

const ToggleGroup = React.forwardRef<HTMLDivElement, ToggleGroupProps>(
  ({ className, variant, size, children, ...props }, ref) => {
    const Root = ToggleGroupPrimitive.Root
    return (
      <Root ref={ref} className={toggleGroupStyles({ className })} {...props}>
        <ToggleGroupContext.Provider value={{ variant, size }}>
          {children}
        </ToggleGroupContext.Provider>
      </Root>
    )
  }
)

type ToggleGroupItemProps = React.ComponentPropsWithoutRef<'button'> &
  ToggleVariants & {
    value: string
  }

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
export type { ToggleGroupItemProps, ToggleGroupProps }
