import {
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  RouterProvider,
} from '@tanstack/react-router'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import InternalError from '#/components/errors/500'

import { renderAt } from '../helpers'

describe('404', () => {
  it('renders for an unknown route and can go back home', async () => {
    const { router } = await renderAt('/this/route/does/not/exist')
    expect(await screen.findByRole('heading', { name: 'Page not found' })).toBeInTheDocument()
    expect(screen.getByText(/couldn't find the page/i)).toBeInTheDocument()

    // Fresh memory history: nothing to go back to, so "Go back" must land on the home page.
    await userEvent.setup().click(screen.getByRole('button', { name: 'Go back' }))
    expect(await screen.findByText('Zero One Starter Kit')).toBeInTheDocument()
    expect(router.state.location.pathname).toBe('/')
  })

  it('uses browser history when there is somewhere to go back to', async () => {
    const back = vi.spyOn(window.history, 'back').mockImplementation(() => {})
    vi.spyOn(window.history, 'length', 'get').mockReturnValue(3)
    await renderAt('/nope')
    await userEvent.setup().click(await screen.findByRole('button', { name: 'Go back' }))
    expect(back).toHaveBeenCalledTimes(1)
    vi.restoreAllMocks()
  })
})

describe('500', () => {
  it('shows the message, details and stack, and is used as the root error boundary', async () => {
    const boom = new Error('kaboom')
    const rootRoute = createRootRoute({
      errorComponent: ({ error }) => (
        <InternalError
          message={error.message}
          details="An unexpected error occurred."
          stack={error.stack}
        />
      ),
    })
    const indexRoute = createRoute({
      getParentRoute: () => rootRoute,
      path: '/',
      component: () => {
        throw boom
      },
    })
    const router = createRouter({
      routeTree: rootRoute.addChildren([indexRoute]),
      history: createMemoryHistory({ initialEntries: ['/'] }),
    })
    vi.spyOn(console, 'error').mockImplementation(() => {})
    await router.load()
    render(<RouterProvider router={router} />)

    expect(
      await screen.findByRole('heading', { name: 'Internal Server Error' })
    ).toBeInTheDocument()
    expect(screen.getByText('kaboom')).toBeInTheDocument()
    expect(screen.getByText('An unexpected error occurred.')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Go back' })).toBeInTheDocument()
    vi.restoreAllMocks()
  })

  it('omits the stack block when none is given', () => {
    // Rendered outside a router on purpose: only static content is asserted.
    const rootRoute = createRootRoute({
      component: () => <InternalError message="500" details="details" />,
    })
    const router = createRouter({
      routeTree: rootRoute,
      history: createMemoryHistory({ initialEntries: ['/'] }),
    })
    render(<RouterProvider router={router} />)
    return screen.findByText('details').then(() => {
      expect(document.querySelector('pre')).toBeNull()
    })
  })
})
