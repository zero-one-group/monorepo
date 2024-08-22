import type React from 'react'
import { forwardRef } from 'react'

export const Label = forwardRef(function Component(
  props: React.ComponentPropsWithoutRef<'label'>,
  ref: React.ForwardedRef<HTMLLabelElement>
) {
  const { ...rest } = props
  return (
    <label ref={ref} className="mb-1 block font-medium text-gray-900 text-sm leading-6" {...rest} />
  )
})
