import type { Metadata } from 'next/types'

import ENV from '#/envar'
import { cn } from '#/utils/helpers'

import '#/styles/globals.css'

export const metadata: Metadata = {
  title: {
    default: 'Zero One Group',
    template: '%s - Zero One Group',
  },
  applicationName: 'Zero One Group',
  description: 'Leading integrated technology services company based in Indonesia.',
  keywords: ['nextjs', 'website', 'agency', 'indonesia', 'stealth startup'],
  robots: {
    index: ENV.NEXT_PUBLIC_ALLOW_INDEXING,
    follow: ENV.NEXT_PUBLIC_ALLOW_INDEXING,
  },
  manifest: '/manifest.webmanifest',
  icons: [
    { rel: 'icon', type: 'image/x-icon', url: '/favicon.ico' },
    { rel: 'icon', type: 'image/svg+xml', url: '/favicon.svg' },
    { rel: 'apple-touch-icon', url: '/favicon.png' },
  ],
  metadataBase: new URL(ENV.NEXT_PUBLIC_BASE_URL),
  openGraph: {
    type: 'website',
    url: new URL(ENV.NEXT_PUBLIC_BASE_URL),
    title: 'Zero One Group',
    description: 'Leading integrated technology services company based in Indonesia.',
    siteName: 'Zero One Group',
    images: [{ url: `${ENV.NEXT_PUBLIC_BASE_URL}/images/og-image.png` }],
  },
}

export default function RootLayout({ children }: React.PropsWithChildren) {
  return (
    <html lang="en">
      <body
        className={cn(ENV.NODE_ENV === 'development' && 'debug-breakpoints')}
        suppressHydrationWarning={true}
      >
        {children}
      </body>
    </html>
  )
}
