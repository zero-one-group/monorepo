/**
 * A custom Link component that wraps TanStack Router's Link. It supports both
 * internal routes (via `href`) and external/hash links (rendered as plain anchors).
 *
 * @param props.href - Internal path (e.g. "/login") or external/hash URL.
 * @param props.newTab - Whether to open the link in a new tab.
 */

import { Link as RouterLink } from '@tanstack/react-router'
import * as React from 'react'

// TanStack Link expects a typed `to`; we accept a plain string `href` instead.
const RouterLinkCmp = RouterLink as unknown as React.ComponentType<
  React.ComponentPropsWithRef<'a'> & { to: string }
>

interface LinkProps extends Omit<React.ComponentPropsWithoutRef<'a'>, 'href'> {
  href: string
  newTab?: boolean
}

const Link = React.forwardRef<HTMLAnchorElement, LinkProps>(function Component(props, ref) {
  const { className, newTab, href, ...rest } = props

  const isExternal = /^https?:\/\//.test(href) || href.startsWith('mailto:')
  const isHash = href.startsWith('#')
  const target = newTab ? '_blank' : undefined
  const rel = newTab ? 'noopener noreferrer' : undefined

  if (isExternal || isHash) {
    return <a href={href} className={className} target={target} rel={rel} ref={ref} {...rest} />
  }

  return (
    <RouterLinkCmp to={href} className={className} target={target} rel={rel} ref={ref} {...rest} />
  )
})

export { Link }
