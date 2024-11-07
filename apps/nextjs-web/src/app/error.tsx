'use client' // Error components must be Client Components

import * as Lucide from 'lucide-react'
import { useEffect } from 'react'

interface ErrorProps {
  error: Error & { digest?: string }
  reset: () => void
}

export default function ErrorPage({ error, reset }: ErrorProps) {
  useEffect(() => {
    // Log the error to an error reporting service
    console.error(error)
  }, [error])

  return (
    <div className="flex h-full min-h-screen w-full flex-col bg-white px-6 md:px-0 dark:bg-gray-900">
      <main className="flex grow items-center justify-center">
        <div className="page-container">
          <section className="mx-auto flex w-full max-w-2xl flex-col items-center justify-center py-24">
            <h1 className="font-extrabold text-2xl text-gray-900 sm:text-3xl lg:text-4xl dark:invert">
              Hey, we got some problems :-(
            </h1>
            <p className="mt-8 text-center text-gray-600 text-xl leading-8 dark:text-gray-300/80">
              An unexpected error has occured. Our team has been notified and will resolve this
              error as soon as possible. We appreciate your interest, but be patient when we are
              fixing everything.
            </p>

            {error.digest && (
              <p className="mt-8 rounded bg-gray-100 px-2.5 py-2 text-center font-medium font-mono text-red-500 text-sm tracking-tight dark:bg-gray-600 dark:text-red-300">
                Error code: {error.digest}
              </p>
            )}

            <div className="mt-8 flex space-x-2 sm:mt-10">
              <button
                type="button"
                className="inline-flex items-center rounded-lg border border-gray-200 bg-gray-900 px-6 py-3 text-center font-medium text-sm text-white hover:bg-gray-700 hover:text-gray-100 focus:outline-none focus:ring-4 focus:ring-gray-100 dark:border-gray-700 dark:bg-gray-800 dark:text-white dark:focus:ring-gray-600 dark:hover:bg-gray-700"
                onClick={() => reset()} // Attempt to recover by trying to re-render the segment
              >
                <Lucide.RefreshCw className="-ml-1 mr-1 size-4" strokeWidth={1.8} />
                Try reload this page
              </button>
            </div>
          </section>
        </div>
      </main>
      <div className="relative bottom-0 w-full py-4 text-center md:absolute">
        <a
          href="/"
          className="inline-flex items-center justify-center text-center font-medium text-gray-700 text-sm hover:text-gray-900 dark:text-white dark:hover:text-gray-300"
        >
          <Lucide.ChevronsLeft className="-ml-1 mr-1 size-4" strokeWidth={1.8} />
          Back to homepage
        </a>
      </div>
    </div>
  )
}
