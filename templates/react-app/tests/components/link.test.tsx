import {
  createMemoryHistory,
  createRootRoute,
  createRouter,
  RouterProvider,
} from '@tanstack/react-router'
import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { Link } from '#/components/link'

async function renderInRouter(ui: React.ReactNode) {
  const rootRoute = createRootRoute({ component: () => <div>{ui}</div> })
  const router = createRouter({
    routeTree: rootRoute,
    history: createMemoryHistory({ initialEntries: ['/'] }),
  })
  await router.load()
  render(<RouterProvider router={router} />)
}

describe('<Link>', () => {
  it('renders internal paths through the router', async () => {
    await renderInRouter(<Link href="/login">Sign in</Link>)
    const a = screen.getByRole('link', { name: 'Sign in' })
    expect(a).toHaveAttribute('href', '/login')
    expect(a).not.toHaveAttribute('target')
  })

  it('renders external and mailto links as plain anchors', async () => {
    await renderInRouter(
      <>
        <Link href="https://example.com">Site</Link>
        <Link href="mailto:hi@example.com">Mail</Link>
      </>
    )
    expect(screen.getByRole('link', { name: 'Site' })).toHaveAttribute('href', 'https://example.com')
    expect(screen.getByRole('link', { name: 'Mail' })).toHaveAttribute(
      'href',
      'mailto:hi@example.com'
    )
  })

  it('opens in a new tab safely when asked', async () => {
    await renderInRouter(
      <Link href="https://example.com" newTab>
        Site
      </Link>
    )
    const a = screen.getByRole('link', { name: 'Site' })
    expect(a).toHaveAttribute('target', '_blank')
    expect(a).toHaveAttribute('rel', 'noopener noreferrer')
  })

  it('treats hash links as plain anchors and forwards className', async () => {
    await renderInRouter(
      <Link href="#top" className="x">
        Top
      </Link>
    )
    const a = screen.getByRole('link', { name: 'Top' })
    expect(a).toHaveAttribute('href', '#top')
    expect(a).toHaveClass('x')
  })
})
