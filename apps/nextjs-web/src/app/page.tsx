import type { Metadata } from 'next'

import { Button } from 'shared-ui'
import Link from '#/components/link'

export const metadata: Metadata = {
  title: {
    absolute: 'Zero One Group',
  },
}

export default function Page() {
  return (
    <>
      <main className="mx-auto flex w-full max-w-4xl grow flex-col justify-center px-4 sm:px-6 lg:px-8">
        <div className="py-16">
          <div className="text-center">
            <p className="font-semibold text-lg sm:text-2xl">Howdy, fellas!</p>
            <h1 className="mt-6 font-bold text-3xl tracking-tight sm:text-4xl lg:text-5xl">
              Welcome to Zero One Group
            </h1>
            <div className="mx-auto mt-8 max-w-3xl">
              <p className="text-lg leading-7">
                We are a leading integrated technology services company, providing personalised
                solutions to even the most challenging business problems. We work closely with our
                clients at any point of their business journey to achieve impactful and
                results‑driven digital transformation.
              </p>
            </div>
            <div className="mt-12 flex items-center justify-center space-x-4">
              <Button variant="default" size="lg" asChild>
                <Link href="/about">About</Link>
              </Button>
              <Button variant="default" size="lg" asChild>
                <Link href="https://github.com/zero-one-group" newTab>
                  GitHub
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </>
  )
}
