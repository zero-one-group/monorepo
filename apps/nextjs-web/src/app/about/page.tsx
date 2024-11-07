import type { Metadata } from 'next'

import AboutPage from './about'

export const metadata: Metadata = { title: 'About' }

export default function Page() {
  return <AboutPage />
}
