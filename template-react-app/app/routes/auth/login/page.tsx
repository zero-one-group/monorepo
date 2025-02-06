import * as Lucide from 'lucide-react'
import { Link } from '#/components/link'
import type { Route } from './+types/page'
import { LoginForm } from './login-form'

export function meta(_props: Route.MetaArgs) {
  return [{ title: 'Sign in | React Router App' }]
}

export default function Page() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8 dark:bg-gray-900">
      <div className="w-full max-w-md space-y-8">
        <div>
          <img
            className="mx-auto block h-12 w-auto dark:hidden"
            src="/images/logo-light.svg"
            alt="Logo"
          />
          <img
            className="mx-auto hidden h-12 w-auto dark:block"
            src="/images/logo-dark.svg"
            alt="Logo"
          />
          <h2 className="mt-6 text-center font-bold text-3xl text-gray-900 dark:text-white">
            Sign in to your account
          </h2>
          <p className="mt-2 text-center text-gray-600 text-sm dark:text-gray-400">
            Or{' '}
            <Link
              href="#"
              className="font-medium text-blue-600 hover:text-blue-500 dark:text-blue-400"
            >
              create a new account
            </Link>
          </p>
        </div>

        <LoginForm />

        <div className="mt-4 text-center">
          <Link
            href="/"
            className="text-gray-600 text-sm hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
          >
            <Lucide.ArrowLeft className="-ml-0.5 mr-1 inline-block size-4" />
            <span>Back to homepage</span>
          </Link>
        </div>
      </div>
    </div>
  )
}
