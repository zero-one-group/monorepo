import * as Lucide from 'lucide-react'
import type { Metadata } from 'next'
import Image from 'next/image'
import Link from '#/app/link'

export const metadata: Metadata = { title: 'Welcome to Next.js App' }

const navItems = [
  { name: 'Dashboard', href: '/' },
  { name: '404', href: '/404' },
  { name: 'Sign In', href: '/login' },
]

const cards = [
  {
    title: 'Zero One Starter Kit',
    description:
      'Launch your next project in minutes with our battle-tested monorepo template and development tools.',
    href: 'https://github.com/zero-one-group/monorepo',
  },
  {
    title: 'Master Next.js',
    description: 'Build a full-stack app with Next.js, TypeScript, Tailwind CSS, and more.',
    href: 'https://nextjs.org/docs',
  },
  {
    title: 'Star Our Repository',
    description:
      'Support our work by starring our GitHub repository and stay updated with latest features.',
    href: 'https://github.com/zero-one-group/monorepo',
  },
]

export default function Page() {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Navbar */}
      <nav className="bg-white shadow-sm dark:bg-gray-800">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex h-16 justify-between">
            <div className="flex items-center">
              <Image
                src="/next.svg"
                className="block h-8 w-auto dark:invert"
                alt="Next.js logo"
                width={180}
                height={38}
                priority
              />
            </div>
            <div className="flex items-center gap-4">
              {navItems.map((item) => (
                <Link
                  key={item.name}
                  href={item.href}
                  className="text-gray-700 hover:text-gray-900 dark:text-gray-200 dark:hover:text-white"
                >
                  {item.name}
                </Link>
              ))}
            </div>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {cards.map((card) => (
            <div key={card.title} className="rounded-lg bg-white p-6 shadow dark:bg-gray-800">
              <h3 className="font-medium text-gray-900 text-lg dark:text-white">{card.title}</h3>
              <p className="mt-2 text-gray-600 dark:text-gray-300">{card.description}</p>
              <Link
                href={card.href}
                className="mt-4 inline-flex items-center text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
                newTab
              >
                <span>Learn more</span>
                <Lucide.ChevronsRight className="ml-1 size-5" />
              </Link>
            </div>
          ))}
        </div>
      </main>
    </div>
  )
}
