import * as Lucide from 'lucide-react'
import { Select as SelectPrimitive } from 'radix-ui'
import * as React from 'react'
import { selectStyles } from './select.css'

const Select = SelectPrimitive.Root
const SelectGroup = SelectPrimitive.Group
const SelectValue = SelectPrimitive.Value

const SelectTrigger = React.forwardRef<
  React.ComponentRef<typeof SelectPrimitive.Trigger>,
  React.ComponentPropsWithoutRef<typeof SelectPrimitive.Trigger>
>(({ className, children, ...props }, ref) => {
  const styles = selectStyles()
  return (
    <SelectPrimitive.Trigger ref={ref} className={styles.trigger({ className })} {...props}>
      {children}
      <SelectPrimitive.Icon asChild>
        <Lucide.ChevronsUpDown className={styles.icon()} strokeWidth={2} />
      </SelectPrimitive.Icon>
    </SelectPrimitive.Trigger>
  )
})

const SelectScrollUpButton = React.forwardRef<
  React.ComponentRef<typeof SelectPrimitive.ScrollUpButton>,
  React.ComponentPropsWithoutRef<typeof SelectPrimitive.ScrollUpButton>
>(({ className, ...props }, ref) => {
  const styles = selectStyles()
  return (
    <SelectPrimitive.ScrollUpButton
      ref={ref}
      className={styles.scrollButton({ className })}
      {...props}
    >
      <Lucide.ChevronUp strokeWidth={2} />
    </SelectPrimitive.ScrollUpButton>
  )
})

const SelectScrollDownButton = React.forwardRef<
  React.ComponentRef<typeof SelectPrimitive.ScrollDownButton>,
  React.ComponentPropsWithoutRef<typeof SelectPrimitive.ScrollDownButton>
>(({ className, ...props }, ref) => {
  const styles = selectStyles()
  return (
    <SelectPrimitive.ScrollDownButton
      ref={ref}
      className={styles.scrollButton({ className })}
      {...props}
    >
      <Lucide.ChevronDown strokeWidth={2} />
    </SelectPrimitive.ScrollDownButton>
  )
})

const SelectContent = React.forwardRef<
  React.ComponentRef<typeof SelectPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof SelectPrimitive.Content>
>(({ className, children, position = 'popper', ...props }, ref) => {
  const styles = selectStyles()
  return (
    <SelectPrimitive.Portal>
      <SelectPrimitive.Content
        ref={ref}
        className={styles.content({ className })}
        position={position}
        {...props}
      >
        <SelectScrollUpButton />
        <SelectPrimitive.Viewport
          className={position === 'popper' ? styles.viewportPopper() : styles.viewport()}
        >
          {children}
        </SelectPrimitive.Viewport>
        <SelectScrollDownButton />
      </SelectPrimitive.Content>
    </SelectPrimitive.Portal>
  )
})

const SelectLabel = React.forwardRef<
  React.ComponentRef<typeof SelectPrimitive.Label>,
  React.ComponentPropsWithoutRef<typeof SelectPrimitive.Label>
>(({ className, ...props }, ref) => {
  const styles = selectStyles()
  return <SelectPrimitive.Label ref={ref} className={styles.label({ className })} {...props} />
})

const SelectItem = React.forwardRef<
  React.ComponentRef<typeof SelectPrimitive.Item>,
  React.ComponentPropsWithoutRef<typeof SelectPrimitive.Item>
>(({ className, children, ...props }, ref) => {
  const styles = selectStyles()
  return (
    <SelectPrimitive.Item ref={ref} className={styles.item({ className })} {...props}>
      <SelectPrimitive.ItemIndicator className={styles.itemIndicator()}>
        <Lucide.Check className={styles.itemIndicatorIcon()} strokeWidth={2} />
      </SelectPrimitive.ItemIndicator>
      <SelectPrimitive.ItemText>{children}</SelectPrimitive.ItemText>
    </SelectPrimitive.Item>
  )
})

const SelectSeparator = React.forwardRef<
  React.ComponentRef<typeof SelectPrimitive.Separator>,
  React.ComponentPropsWithoutRef<typeof SelectPrimitive.Separator>
>(({ className, ...props }, ref) => {
  const styles = selectStyles()
  return (
    <SelectPrimitive.Separator ref={ref} className={styles.separator({ className })} {...props} />
  )
})

Select.displayName = SelectPrimitive.Root.displayName
SelectGroup.displayName = SelectPrimitive.Group.displayName
SelectValue.displayName = SelectPrimitive.Value.displayName
SelectTrigger.displayName = SelectPrimitive.Trigger.displayName
SelectContent.displayName = SelectPrimitive.Content.displayName
SelectLabel.displayName = SelectPrimitive.Label.displayName
SelectItem.displayName = SelectPrimitive.Item.displayName
SelectSeparator.displayName = SelectPrimitive.Separator.displayName
SelectScrollUpButton.displayName = SelectPrimitive.ScrollUpButton.displayName
SelectScrollDownButton.displayName = SelectPrimitive.ScrollDownButton.displayName

export {
  Select,
  SelectGroup,
  SelectValue,
  SelectTrigger,
  SelectContent,
  SelectLabel,
  SelectItem,
  SelectSeparator,
  SelectScrollUpButton,
  SelectScrollDownButton,
}
