import type { MetadataRoute } from 'next'

import ENV from '#/envar'

export const dynamicParams = false // Prevents dynamic route parameters.
export const dynamic = 'force-static' // Forces the page to be statically rendered.
export const fetchCache = 'force-cache' // Force caching of fetch requests.
export const revalidate = 0 // Disable Incremental Static Regeneration (ISR).

export default function manifest(): MetadataRoute.Manifest {
  return {
    lang: 'en',
    dir: 'ltr',
    name: 'Zero One Group',
    short_name: 'zero-one-group',
    description: 'Leading integrated technology services company based in Indonesia',
    theme_color: '#121314',
    background_color: '#ffffff',
    start_url: `${ENV.NEXT_PUBLIC_BASE_URL}/?source=pwa`,
    id: `${ENV.NEXT_PUBLIC_BASE_URL}/?source=pwa`,
    icons: [
      {
        src: '/favicon.svg',
        sizes: '36x36',
        type: 'image/svg+xml',
      },
      {
        src: '/favicon.svg',
        sizes: '48x48',
        type: 'image/svg+xml',
      },
      {
        src: '/favicon.svg',
        sizes: '72x72',
        type: 'image/svg+xml',
      },
      {
        src: '/favicon.svg',
        sizes: '96x96',
        type: 'image/svg+xml',
      },
      {
        src: '/favicon.svg',
        sizes: '144x144',
        type: 'image/svg+xml',
      },
      {
        src: '/favicon.svg',
        sizes: '192x192',
        type: 'image/svg+xml',
      },
    ],
    display: 'standalone',
    orientation: 'natural',
  }
}
