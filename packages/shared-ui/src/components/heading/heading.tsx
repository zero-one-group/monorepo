import { Slot } from '@radix-ui/react-slot'
import * as React from 'react'
import { type HeadingVariants, headingStyles } from './heading.css'

export interface HeadingProps extends React.HTMLAttributes<HTMLHeadingElement>, HeadingVariants {
  asChild?: boolean
}

const Heading = React.forwardRef<HTMLHeadingElement, HeadingProps>(
  ({ className, level, weight, align, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : 'h2'
    return (
      <Comp ref={ref} className={headingStyles({ level, weight, align, className })} {...props} />
    )
  }
)

Heading.displayName = 'Heading'

export { Heading }
