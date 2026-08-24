import { createMemoryHistory, createRouter, RouterProvider } from '@tanstack/react-router'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'

import { routeTree } from '#/routeTree.gen'

// Setup userEvent for interaction testing
const actor = userEvent.setup()

const createTestRouter = () =>
  createRouter({
    routeTree,
    history: createMemoryHistory({ initialEntries: ['/'] }),
  })

describe('Homepage', () => {
  it('renders navigation and content', async () => {
    const router = createTestRouter()
    await router.load()
    render(<RouterProvider router={router} />)

    // Test navigation items
    expect(await screen.findByText('Dashboard')).toBeInTheDocument()
    expect(await screen.findByText('404')).toBeInTheDocument()
    expect(await screen.findByText('Sign In')).toBeInTheDocument()

    // Test cards content
    expect(await screen.findByText('Zero One Starter Kit')).toBeInTheDocument()
    expect(await screen.findByText('Master TanStack Router')).toBeInTheDocument()
    expect(await screen.findByText('Star Our Repository')).toBeInTheDocument()
  })

  it('handles link interactions', async () => {
    const router = createTestRouter()
    await router.load()
    render(<RouterProvider router={router} />)

    const learnMoreLinks = await screen.findAllByText('Learn more')
    await actor.click(learnMoreLinks[0])

    expect(learnMoreLinks[0].closest('a')).toHaveAttribute('target', '_blank')
    expect(learnMoreLinks[0].closest('a')).toHaveAttribute(
      'href',
      'https://github.com/zero-one-group/monorepo'
    )
  })
})
