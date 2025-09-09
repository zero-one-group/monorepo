import { Slot } from 'radix-ui'
import * as React from 'react'
import { type TextVariants, textStyles } from './text.css'

export interface TextProps extends React.HTMLAttributes<HTMLParagraphElement>, TextVariants {
  asChild?: boolean
}

const Text = React.forwardRef<HTMLParagraphElement, TextProps>(
  ({ className, size, weight, align, variant, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot.Root : 'p'
    return (
      <Comp
        ref={ref}
        className={textStyles({ size, weight, align, variant, className })}
        {...props}
      />
    )
  }
)

Text.displayName = 'Text'

export { Text }
