import type { MetadataRoute } from 'next'

import ENV from '#/envar'

export const dynamicParams = false // Prevents dynamic route parameters.
export const dynamic = 'force-static' // Forces the page to be statically rendered.
export const fetchCache = 'force-cache' // Force caching of fetch requests.
export const revalidate = 0 // Disable Incremental Static Regeneration (ISR).

export default function robots(): MetadataRoute.Robots {
  const DISALLOW_LIST: string[] = ['/api', '/_next']

  return {
    rules: {
      userAgent: '*',
      disallow: ENV.NEXT_PUBLIC_ALLOW_INDEXING ? DISALLOW_LIST : '*',
    },
    sitemap: ENV.NEXT_PUBLIC_ALLOW_INDEXING ? `${ENV.NEXT_PUBLIC_BASE_URL}/sitemap.xml` : undefined,
  }
}
