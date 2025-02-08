import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

/**
 * Combines multiple CSS class values using the `clsx` and `tailwind-merge` libraries.
 *
 * @param args - An array of CSS class values to be combined.
 * @returns The combined CSS class value.
 */
export function clx(...args: ClassValue[]) {
  return twMerge(clsx(...args))
}
