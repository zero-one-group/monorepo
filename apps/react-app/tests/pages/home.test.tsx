import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'

import { renderAt } from '../helpers'

describe('Homepage', () => {
  it('renders navigation and content', async () => {
    await renderAt('/')

    // Navigation items
    expect(await screen.findByText('Dashboard')).toBeInTheDocument()
    expect(await screen.findByText('404')).toBeInTheDocument()
    expect(await screen.findByText('Sign In')).toBeInTheDocument()

    // Cards
    expect(await screen.findByText('Zero One Starter Kit')).toBeInTheDocument()
    expect(await screen.findByText('Master TanStack Router')).toBeInTheDocument()
    expect(await screen.findByText('Star Our Repository')).toBeInTheDocument()
  })

  it('opens external cards in a new tab', async () => {
    await renderAt('/')
    const learnMoreLinks = await screen.findAllByText('Learn more')
    await userEvent.setup().click(learnMoreLinks[0])

    expect(learnMoreLinks[0].closest('a')).toHaveAttribute('target', '_blank')
    expect(learnMoreLinks[0].closest('a')).toHaveAttribute(
      'href',
      'https://github.com/zero-one-group/monorepo'
    )
  })

  it('navigates to the login page from the nav', async () => {
    const { router } = await renderAt('/')
    await userEvent.setup().click(await screen.findByRole('link', { name: 'Sign In' }))
    expect(
      await screen.findByRole('heading', { name: /sign in to your account/i })
    ).toBeInTheDocument()
    expect(router.state.location.pathname).toBe('/login')
  })
})
