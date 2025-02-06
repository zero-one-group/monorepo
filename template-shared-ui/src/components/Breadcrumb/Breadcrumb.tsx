import { Slot } from '@radix-ui/react-slot'
import * as Lucide from 'lucide-react'
import * as React from 'react'
import { breadcrumbStyles } from './breadcrumb.css'

const Breadcrumb = React.forwardRef<
  HTMLElement,
  React.ComponentPropsWithoutRef<'nav'> & {
    separator?: React.ReactNode
  }
>(({ ...props }, ref) => <nav ref={ref} aria-label="breadcrumb" {...props} />)

const BreadcrumbList = React.forwardRef<HTMLOListElement, React.ComponentPropsWithoutRef<'ol'>>(
  ({ className, ...props }, ref) => {
    const styles = breadcrumbStyles()
    return <ol ref={ref} className={styles.list({ className })} {...props} />
  }
)

const BreadcrumbItem = React.forwardRef<HTMLLIElement, React.ComponentPropsWithoutRef<'li'>>(
  ({ className, ...props }, ref) => {
    const styles = breadcrumbStyles()
    return <li ref={ref} className={styles.item({ className })} {...props} />
  }
)

const BreadcrumbLink = React.forwardRef<
  HTMLAnchorElement,
  React.ComponentPropsWithoutRef<'a'> & {
    asChild?: boolean
  }
>(({ asChild, className, ...props }, ref) => {
  const Comp = asChild ? Slot : 'a'
  const styles = breadcrumbStyles()
  return <Comp ref={ref} className={styles.link({ className })} {...props} />
})

const BreadcrumbPage = React.forwardRef<HTMLSpanElement, React.ComponentPropsWithoutRef<'span'>>(
  ({ className, ...props }, ref) => {
    const styles = breadcrumbStyles()
    return (
      <span
        ref={ref}
        aria-disabled="true"
        aria-current="page"
        className={styles.page({ className })}
        {...props}
      />
    )
  }
)

const BreadcrumbSeparator = ({ children, className, ...props }: React.ComponentProps<'span'>) => {
  const styles = breadcrumbStyles()
  return (
    <span
      role="presentation"
      aria-hidden="true"
      className={styles.separator({ className })}
      {...props}
    >
      {children ?? <Lucide.ChevronRight strokeWidth={2} />}
    </span>
  )
}

const BreadcrumbEllipsis = ({ className, ...props }: React.ComponentProps<'span'>) => {
  const styles = breadcrumbStyles()
  return (
    <span
      role="presentation"
      aria-hidden="true"
      className={styles.ellipsis({ className })}
      {...props}
    >
      <Lucide.Ellipsis className={styles.ellipsisIcon()} strokeWidth={2} />
      <span className="sr-only">More</span>
    </span>
  )
}

Breadcrumb.displayName = 'Breadcrumb'
BreadcrumbList.displayName = 'BreadcrumbList'
BreadcrumbItem.displayName = 'BreadcrumbItem'
BreadcrumbLink.displayName = 'BreadcrumbLink'
BreadcrumbPage.displayName = 'BreadcrumbPage'
BreadcrumbSeparator.displayName = 'BreadcrumbSeparator'
BreadcrumbEllipsis.displayName = 'BreadcrumbEllipsis'

export {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbPage,
  BreadcrumbSeparator,
  BreadcrumbEllipsis,
}
