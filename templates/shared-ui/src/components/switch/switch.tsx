import { Switch as SwitchPrimitives } from 'radix-ui'
import * as React from 'react'
import { switchStyles } from './switch.css'

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
