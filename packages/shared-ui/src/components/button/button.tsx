import { Slot } from '@radix-ui/react-slot'
import * as Lucide from 'lucide-react'
import * as React from 'react'
import { type ButtonVariants, buttonStyles } from './button.css'

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement>, ButtonVariants {
  asChild?: boolean
  isLoading?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ variant, size, className, isLoading, disabled, children, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : 'button'
    const isDisabled = disabled || isLoading
    const styles = buttonStyles({ variant, size, isLoading })

    // Wrap children in fragment when loading to ensure single child
    const content = isLoading ? (
      <>
        <Lucide.Loader2 strokeWidth={2} />
        {children}
      </>
    ) : (
      children
    )

    return (
      <Comp
        ref={ref}
        className={styles.base({ className })}
        data-loading={isLoading}
        disabled={isDisabled}
        {...props}
      >
        {content}
      </Comp>
    )
  }
)

Button.displayName = 'Button'

export { Button }
