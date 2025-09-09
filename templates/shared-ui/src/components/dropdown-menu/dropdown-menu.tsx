import * as Lucide from 'lucide-react'
import { DropdownMenu as DropdownMenuPrimitive } from 'radix-ui'
import * as React from 'react'
import { type DropdownMenuVariants, dropdownMenuStyles } from './dropdown-menu.css'

const DropdownMenu = DropdownMenuPrimitive.Root
const DropdownMenuTrigger = DropdownMenuPrimitive.Trigger
const DropdownMenuGroup = DropdownMenuPrimitive.Group
const DropdownMenuPortal = DropdownMenuPrimitive.Portal
const DropdownMenuSub = DropdownMenuPrimitive.Sub
const DropdownMenuRadioGroup = DropdownMenuPrimitive.RadioGroup

const DropdownMenuSubTrigger = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.SubTrigger>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.SubTrigger> & DropdownMenuVariants
>(({ className, inset, children, ...props }, ref) => {
  const styles = dropdownMenuStyles({ inset })
  return (
    <DropdownMenuPrimitive.SubTrigger
      ref={ref}
      className={styles.subTrigger({ className })}
      {...props}
    >
      {children}
      <Lucide.ChevronRight className={styles.triggerIcon()} strokeWidth={2} />
    </DropdownMenuPrimitive.SubTrigger>
  )
})

const DropdownMenuSubContent = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.SubContent>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.SubContent>
>(({ className, ...props }, ref) => {
  const styles = dropdownMenuStyles()
  return (
    <DropdownMenuPrimitive.SubContent
      ref={ref}
      className={styles.subContent({ className })}
      {...props}
    />
  )
})

const DropdownMenuContent = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Content>
>(({ className, sideOffset = 4, ...props }, ref) => {
  const styles = dropdownMenuStyles()
  return (
    <DropdownMenuPrimitive.Portal>
      <DropdownMenuPrimitive.Content
        ref={ref}
        sideOffset={sideOffset}
        className={styles.content({ className })}
        {...props}
      />
    </DropdownMenuPrimitive.Portal>
  )
})

const DropdownMenuItem = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.Item>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Item> & DropdownMenuVariants
>(({ className, inset, ...props }, ref) => {
  const styles = dropdownMenuStyles({ inset })
  return <DropdownMenuPrimitive.Item ref={ref} className={styles.item({ className })} {...props} />
})

const DropdownMenuCheckboxItem = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.CheckboxItem>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.CheckboxItem>
>(({ className, children, checked, ...props }, ref) => {
  const styles = dropdownMenuStyles()
  return (
    <DropdownMenuPrimitive.CheckboxItem
      ref={ref}
      className={styles.checkboxItem({ className })}
      checked={checked}
      {...props}
    >
      <DropdownMenuPrimitive.ItemIndicator className={styles.checkboxItemIndicator()}>
        <Lucide.Check className={styles.checkboxItemIcon()} strokeWidth={2} />
      </DropdownMenuPrimitive.ItemIndicator>
      {children}
    </DropdownMenuPrimitive.CheckboxItem>
  )
})

const DropdownMenuRadioItem = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.RadioItem>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.RadioItem>
>(({ className, children, ...props }, ref) => {
  const styles = dropdownMenuStyles()
  return (
    <DropdownMenuPrimitive.RadioItem
      ref={ref}
      className={styles.radioItem({ className })}
      {...props}
    >
      <DropdownMenuPrimitive.ItemIndicator className={styles.radioItemIndicator()}>
        <Lucide.Dot className={styles.radioItemIcon()} strokeWidth={2} />
      </DropdownMenuPrimitive.ItemIndicator>
      {children}
    </DropdownMenuPrimitive.RadioItem>
  )
})

const DropdownMenuLabel = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.Label>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Label> & DropdownMenuVariants
>(({ className, inset, ...props }, ref) => {
  const styles = dropdownMenuStyles({ inset })
  return (
    <DropdownMenuPrimitive.Label ref={ref} className={styles.label({ className })} {...props} />
  )
})

const DropdownMenuSeparator = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.Separator>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Separator>
>(({ className, ...props }, ref) => {
  const styles = dropdownMenuStyles()
  return (
    <DropdownMenuPrimitive.Separator
      ref={ref}
      className={styles.separator({ className })}
      {...props}
    />
  )
})

const DropdownMenuShortcut = ({ className, ...props }: React.HTMLAttributes<HTMLSpanElement>) => {
  const styles = dropdownMenuStyles()
  return <span className={styles.shortcut({ className })} {...props} />
}

DropdownMenu.displayName = 'DropdownMenu'
DropdownMenuTrigger.displayName = 'DropdownMenuTrigger'
DropdownMenuGroup.displayName = 'DropdownMenuGroup'
DropdownMenuPortal.displayName = 'DropdownMenuPortal'
DropdownMenuSub.displayName = 'DropdownMenuSub'
DropdownMenuRadioGroup.displayName = 'DropdownMenuRadioGroup'
DropdownMenuSubTrigger.displayName = DropdownMenuPrimitive.SubTrigger.displayName
DropdownMenuSubContent.displayName = DropdownMenuPrimitive.SubContent.displayName
DropdownMenuContent.displayName = DropdownMenuPrimitive.Content.displayName
DropdownMenuItem.displayName = DropdownMenuPrimitive.Item.displayName
DropdownMenuCheckboxItem.displayName = DropdownMenuPrimitive.CheckboxItem.displayName
DropdownMenuRadioItem.displayName = DropdownMenuPrimitive.RadioItem.displayName
DropdownMenuLabel.displayName = DropdownMenuPrimitive.Label.displayName
DropdownMenuSeparator.displayName = DropdownMenuPrimitive.Separator.displayName
DropdownMenuShortcut.displayName = 'DropdownMenuShortcut'

export {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuCheckboxItem,
  DropdownMenuRadioItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuGroup,
  DropdownMenuPortal,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuRadioGroup,
}
