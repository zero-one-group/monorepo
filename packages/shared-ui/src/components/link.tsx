/**
 * TODO: Update this component to use your client-side framework's link component.
 */

import type React from 'react'
import { forwardRef } from 'react'

interface CustomLinkProps {
  href: string
  newTab?: boolean
}

export const Link = forwardRef(function Component(
  props: CustomLinkProps & React.ComponentPropsWithoutRef<'a'>,
  ref: React.ForwardedRef<HTMLAnchorElement>
) {
  const { href, newTab, ...rest } = props
  return (
    <a
      ref={ref}
      href={href}
      className="font-medium text-primary-600 hover:underline dark:text-primary-500"
      target={newTab ? '_blank' : '_self'}
      rel={newTab ? 'noopener noreferrer' : undefined}
      {...rest}
    />
  )
})
