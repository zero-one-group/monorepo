import '../styles/global.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import type * as React from 'react'
import { clx } from '#/libs/utils'

// Load the optimized fonts from Google Fonts
const fontSans = Inter({ variable: '--font-inter', subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Next.js App',
  description: 'A simple Next.js app with Tailwind CSS and TypeScript',
}

export default function RootLayout({ children }: React.PropsWithChildren) {
  return (
    <html lang="en" data-theme="light">
      <body className={clx(fontSans.variable, 'antialiased')}>{children}</body>
    </html>
  )
}
