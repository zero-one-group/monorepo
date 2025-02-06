import * as Lucide from 'lucide-react'
import * as React from 'react'
import { Link } from '#/components/base-ui'
import type { ButtonProps } from '../button/button'
import { buttonStyles } from '../button/button.css'
import { paginationStyles } from './pagination.css'

const Pagination = ({ className, ...props }: React.ComponentProps<'nav'>) => {
  const styles = paginationStyles()
  return <nav aria-label="pagination" className={styles.base({ className })} {...props} />
}

const PaginationContent = React.forwardRef<HTMLUListElement, React.ComponentProps<'ul'>>(
  ({ className, ...props }, ref) => {
    const styles = paginationStyles()
    return <ul ref={ref} className={styles.content({ className })} {...props} />
  }
)

const PaginationItem = React.forwardRef<HTMLLIElement, React.ComponentProps<'li'>>(
  ({ className, ...props }, ref) => {
    const styles = paginationStyles()
    return <li ref={ref} className={styles.item({ className })} {...props} />
  }
)

type PaginationLinkProps = {
  isActive?: boolean
} & Pick<ButtonProps, 'size'> &
  React.ComponentProps<'a'>

const PaginationLink = ({ className, isActive, size = 'icon', ...props }: PaginationLinkProps) => {
  const styles = buttonStyles({ variant: isActive ? 'outline' : 'ghost', size })
  return (
    <Link
      href={props.href || '#'}
      aria-current={isActive ? 'page' : undefined}
      className={styles.base({ className })}
      {...props}
    />
  )
}

const PaginationPrevious = ({
  className,
  ...props
}: React.ComponentProps<typeof PaginationLink>) => {
  const styles = paginationStyles()
  return (
    <PaginationLink
      className={styles.previous({ className })}
      aria-label="Go to previous page"
      size="default"
      {...props}
    >
      <Lucide.ChevronLeft className={styles.previousIcon()} />
      <span>Previous</span>
    </PaginationLink>
  )
}

const PaginationNext = ({ className, ...props }: React.ComponentProps<typeof PaginationLink>) => {
  const styles = paginationStyles()
  return (
    <PaginationLink
      className={styles.next({ className })}
      aria-label="Go to next page"
      size="default"
      {...props}
    >
      <span>Next</span>
      <Lucide.ChevronRight className={styles.nextIcon()} />
    </PaginationLink>
  )
}

const PaginationEllipsis = ({ className, ...props }: React.ComponentProps<'span'>) => {
  const styles = paginationStyles()
  return (
    <span aria-hidden className={styles.ellipsis({ className })} {...props}>
      <Lucide.MoreHorizontal className={styles.ellipsisIcon()} />
      <span className="sr-only">More pages</span>
    </span>
  )
}

Pagination.displayName = 'Pagination'
PaginationContent.displayName = 'PaginationContent'
PaginationItem.displayName = 'PaginationItem'
PaginationLink.displayName = 'PaginationLink'
PaginationPrevious.displayName = 'PaginationPrevious'
PaginationNext.displayName = 'PaginationNext'
PaginationEllipsis.displayName = 'PaginationEllipsis'

export {
  Pagination,
  PaginationContent,
  PaginationLink,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
  PaginationEllipsis,
}
