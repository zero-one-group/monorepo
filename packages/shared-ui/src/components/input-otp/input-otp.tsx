import { OTPInput, OTPInputContext } from 'input-otp'
import type { OTPInputProps } from 'input-otp'
import * as Lucide from 'lucide-react'
import * as React from 'react'
import { inputOtpStyles } from './input-otp.css'
import type { InputOtpVariants } from './input-otp.css'

export interface InputOTPProps extends Omit<OTPInputProps, 'children' | 'render' | 'size'> {
  size?: InputOtpVariants['size']
  children: React.ReactNode
  render?: never
}

const InputOTP = React.forwardRef<React.ComponentRef<typeof OTPInput>, InputOTPProps>(
  ({ className, containerClassName, size, ...props }, ref) => {
    const styles = inputOtpStyles({ size })
    return (
      <OTPInput
        ref={ref}
        containerClassName={styles.root({ className: containerClassName })}
        className={styles.input({ className })}
        {...props}
      />
    )
  }
)

const InputOTPGroup = React.forwardRef<
  React.ComponentRef<'div'>,
  React.ComponentPropsWithoutRef<'div'>
>(({ className, ...props }, ref) => {
  const styles = inputOtpStyles()
  return <div ref={ref} className={styles.group({ className })} {...props} />
})

interface InputOTPSlotProps extends React.ComponentPropsWithoutRef<'div'> {
  index: number
  size?: InputOtpVariants['size']
}

const InputOTPSlot = React.forwardRef<React.ComponentRef<'div'>, InputOTPSlotProps>(
  ({ index, className, size, ...props }, ref) => {
    const inputOTPContext = React.useContext(OTPInputContext)
    const { char, hasFakeCaret, isActive } = inputOTPContext.slots[index] || {}
    const styles = inputOtpStyles({ size })

    return (
      <div
        ref={ref}
        className={styles.slot({ className: [isActive && styles.slotActive(), className] })}
        {...props}
      >
        {char}
        {hasFakeCaret && (
          <div className={styles.caret()}>
            <div className={styles.caretBlink()} />
          </div>
        )}
      </div>
    )
  }
)

interface InputOTPSeparatorProps extends React.ComponentPropsWithoutRef<'div'> {
  size?: InputOtpVariants['size']
}

const InputOTPSeparator = React.forwardRef<React.ComponentRef<'div'>, InputOTPSeparatorProps>(
  ({ className, size, ...props }, ref) => {
    const styles = inputOtpStyles({ size })
    return (
      <div ref={ref} className={styles.separator({ className })} {...props}>
        <Lucide.Minus strokeWidth={2} />
      </div>
    )
  }
)

InputOTP.displayName = 'InputOTP'
InputOTPGroup.displayName = 'InputOTPGroup'
InputOTPSlot.displayName = 'InputOTPSlot'
InputOTPSeparator.displayName = 'InputOTPSeparator'

export { InputOTP, InputOTPGroup, InputOTPSlot, InputOTPSeparator }
