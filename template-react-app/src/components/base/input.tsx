import type React from 'react'
import { forwardRef } from 'react'

export const Input = forwardRef(function Component(
  props: React.ComponentPropsWithoutRef<'input'>,
  ref: React.ForwardedRef<HTMLInputElement>
) {
  const { ...rest } = props
  return (
    <input
      ref={ref}
      className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-indigo-600 focus:ring-inset sm:text-sm sm:leading-6"
      {...rest}
    />
  )
})
