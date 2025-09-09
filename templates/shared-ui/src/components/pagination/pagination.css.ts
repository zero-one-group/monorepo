import { type VariantProps, tv } from 'tailwind-variants'

export const paginationStyles = tv({
  slots: {
    base: 'mx-auto flex w-full justify-center',
    content: 'flex flex-row items-center gap-1',
    item: 'relative hover:no-underline',
    previous: 'gap-1 pl-2.5 hover:no-underline',
    previousIcon: 'size-4',
    next: 'gap-1 pr-2.5 hover:no-underline',
    nextIcon: 'size-4',
    ellipsis: 'flex size-9 items-center justify-center',
    ellipsisIcon: 'size-4',
  },
})

export type PaginationVariants = VariantProps<typeof paginationStyles>
