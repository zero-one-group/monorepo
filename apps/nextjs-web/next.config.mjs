import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import NextBundleAnalyzer from '@next/bundle-analyzer'

const isVercel = process.env.IS_VERCEL_ENV === 'true'
const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

const withBundleAnalyzer = NextBundleAnalyzer({
  enabled: process.env.ANALYZE === 'true',
  openAnalyzer: true,
})

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  cleanDistDir: true,
  reactStrictMode: true,
  poweredByHeader: false,
  productionBrowserSourceMaps: false,
  images: { remotePatterns: [{ protocol: 'https', hostname: '**' }] },

  // @ref: https://nextjs.org/blog/next-14-1#improved-self-hosting
  cacheMaxMemorySize: 0, // disable default in-memory caching
  cacheHandler:
    process.env.NODE_ENV === 'production' && !isVercel
      ? resolve(__dirname, 'cache-handler.mjs')
      : undefined,

  eslint: { ignoreDuringBuilds: true },
  typescript: { ignoreBuildErrors: true },
  logging: { fetches: { fullUrl: true } },
  transpilePackages: ['shared-ui'],

  experimental: {
    // This is required for the experimental feature of
    // pre-populating the cache with the initial data.
    instrumentationHook: true,
  },
}

export default withBundleAnalyzer(nextConfig)
