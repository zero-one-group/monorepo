import * as React from 'react'
import { type SkeletonVariants, skeletonStyles } from './skeleton.css'

export interface SkeletonProps extends React.ComponentPropsWithoutRef<'div'>, SkeletonVariants {}

const Skeleton = React.forwardRef<HTMLDivElement, SkeletonProps>(({ className, ...props }, ref) => {
  return <div ref={ref} className={skeletonStyles({ className })} {...props} />
})

Skeleton.displayName = 'Skeleton'

export { Skeleton }
