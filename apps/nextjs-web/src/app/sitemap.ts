import type { MetadataRoute } from 'next/types'

import ENV from '#/envar'

// @reference: https://github.com/vercel/next.js/issues/49373
export const dynamicParams = false // Prevents dynamic route parameters.
export const dynamic = 'force-dynamic' // Forces the page to be dynamically rendered on each request.
export const fetchCache = 'force-cache' // Force caching of fetch requests.
export const revalidate = 0 // Disable Incremental Static Regeneration (ISR).

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const baseUrl = ENV.NEXT_PUBLIC_BASE_URL
  const now = new Date()

  const rootPath = {
    url: baseUrl,
    lastModified: now,
    changeFrequency: 'yearly' as const,
    priority: 1,
  }

  const pages = ['/about']

  const routes = pages.map((path) => ({
    url: `${baseUrl}${path}`,
    lastModified: now,
    priority: 0.8,
  }))

  const pageRoutes = [rootPath, ...routes]

  // const allPosts = await fetch('https://.../posts')

  // if (!allPosts) return [...pageRoutes]

  // const postsRoutes = allPosts.map((post) => ({
  //   url: `${baseUrl}/blog/${post.slug}`,
  //   lastModified: post.created_at,
  //   priority: 0.5,
  // }))

  // return [...pageRoutes, ...postsRoutes]

  return [...pageRoutes]
}
