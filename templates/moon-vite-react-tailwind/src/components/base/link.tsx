import type React from 'react'
import { forwardRef } from 'react'
import { type LinkProps, Link as RouterLink } from 'react-router-dom'

interface CustomLinkProps extends Omit<LinkProps, 'to'> {
  href: string
  newTab?: boolean
}

/**
 * Custom Link component that wraps React Router's Link component.
 * This component replaces the `to` prop with `href` for consistency with HTML anchor elements.
 *
 * @param props - The properties for the Link component.
 * @param ref - The forwarded ref for the anchor element.
 * @returns A React element that renders a link.
 *
 * Example usage:
 * ```tsx
 * <Link href="/path" className="custom-class">Link Text</Link>
 * ```
 */
export const Link = forwardRef(function Component(
  props: CustomLinkProps & React.ComponentPropsWithoutRef<'a'>,
  ref: React.ForwardedRef<HTMLAnchorElement>
) {
  const { href, newTab, ...rest } = props
  return (
    <RouterLink
      ref={ref}
      to={href}
      className="font-medium text-primary-600 hover:underline dark:text-primary-500"
      target={newTab ? '_blank' : '_self'}
      rel={newTab ? 'noopener noreferrer' : undefined}
      {...rest}
    />
  )
})
