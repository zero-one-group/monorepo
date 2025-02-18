import './styles.css'
import { clx } from '@repo/shared-ui/utils'
import { RootProvider } from 'fumadocs-ui/provider'
import { Figtree, JetBrains_Mono } from 'next/font/google'
import type { ReactNode } from 'react'

const fontSans = Figtree({ subsets: ['latin'], variable: '--font-sans' })
const fontMono = JetBrains_Mono({ subsets: ['latin'], variable: '--font-mono' })

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" className={clx(fontSans.variable, fontMono.variable)} suppressHydrationWarning>
      <body className="flex min-h-screen flex-col">
        <RootProvider search={{ enabled: false }}>{children}</RootProvider>
      </body>
    </html>
  )
}
