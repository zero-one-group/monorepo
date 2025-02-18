import { createMDX } from 'fumadocs-mdx/next'
import { isProduction } from 'std-env'

const withMDX = createMDX()

/** @type {import('next').NextConfig} */
const nextConfig = {
  cleanDistDir: true,
  reactStrictMode: true,
  poweredByHeader: false,
  productionBrowserSourceMaps: false,
  images: { remotePatterns: [{ protocol: 'https', hostname: '**' }] },
  eslint: { ignoreDuringBuilds: isProduction },
  typescript: { ignoreBuildErrors: isProduction },
  logging: { fetches: { fullUrl: true } },
}

export default withMDX(nextConfig)
