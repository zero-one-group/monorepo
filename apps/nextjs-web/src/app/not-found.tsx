import * as Lucide from 'lucide-react'
import type { Metadata } from 'next/types'

import Link from '#/components/link'

export const metadata: Metadata = { title: '404 Page not found' }

export default function NotFound() {
  return (
    <main className="flex grow items-center justify-center">
      <div className="page-container">
        <section className="mx-auto flex w-full max-w-2xl flex-col items-center justify-center py-24">
          <h1 className="font-extrabold text-2xl text-gray-900 sm:text-3xl lg:text-4xl dark:invert">
            404 - Page not found
          </h1>
          <p className="mt-8 text-center text-gray-600 text-xl leading-8 dark:text-gray-300/80">
            Sorry, we can&rsquo;t find that page. Check that you typed the address correctly, or try
            using our site search to find something specific.
          </p>
          <div className="mt-8 sm:mt-12">
            <Link
              href="/"
              className="inline-flex items-center rounded-lg border border-gray-200 bg-gray-900 px-6 py-3 text-center font-medium text-sm text-white hover:bg-gray-700 hover:text-gray-100 focus:outline-none focus:ring-4 focus:ring-gray-100 dark:border-gray-700 dark:bg-gray-800 dark:text-white dark:focus:ring-gray-600 dark:hover:bg-gray-700"
            >
              <Lucide.ChevronsLeft className="-ml-1 mr-1 size-4" strokeWidth={1.8} />
              Back to homepage
            </Link>
          </div>
        </section>
      </div>
    </main>
  )
}
