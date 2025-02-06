import clsx, { type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

/**
 * Combines multiple CSS class names using the `clsx` and `tailwind-merge` libraries.
 *
 * @param inputs - An array of CSS class names to be combined.
 * @returns The combined CSS class names.
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
