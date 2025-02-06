import * as SwitchPrimitives from '@radix-ui/react-switch'
import * as React from 'react'
import { switchStyles } from './switch.css.ts'

const Switch = React.forwardRef<
  React.ComponentRef<typeof SwitchPrimitives.Root>,
  React.ComponentPropsWithoutRef<typeof SwitchPrimitives.Root>
>(({ className, ...props }, ref) => {
  const styles = switchStyles()
  return (
    <SwitchPrimitives.Root className={styles.base({ className })} {...props} ref={ref}>
      <SwitchPrimitives.Thumb className={styles.thumb()} />
    </SwitchPrimitives.Root>
  )
})

Switch.displayName = SwitchPrimitives.Root.displayName

export { Switch }
