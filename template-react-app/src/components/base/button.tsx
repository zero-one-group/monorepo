import type React from 'react'
import { forwardRef } from 'react'

export const Button = forwardRef(function Component(
  props: React.ComponentPropsWithoutRef<'button'>,
  ref: React.ForwardedRef<HTMLButtonElement>
) {
  const { ...rest } = props
  return (
    <button
      ref={ref}
      className="rounded-md bg-gray-900 px-3 py-2 font-semibold text-sm text-white shadow-sm ring-1 ring-gray-800 ring-inset hover:bg-gray-800"
      {...rest}
    />
  )
})
