import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'
import Page from '#/app/(dashboard)/page'

// Setup userEvent for interaction testing
const actor = userEvent.setup()

describe('Homepage', () => {
  it('renders navigation and content', async () => {
    render(<Page />)

    // Test navigation items
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
    expect(screen.getByText('404')).toBeInTheDocument()
    expect(screen.getByText('Sign In')).toBeInTheDocument()

    // Test cards content
    expect(screen.getByText('Zero One Starter Kit')).toBeInTheDocument()
    expect(screen.getByText('Master Next.js')).toBeInTheDocument()
    expect(screen.getByText('Star Our Repository')).toBeInTheDocument()
  })

  it('handles link interactions', async () => {
    render(<Page />)

    const learnMoreLinks = screen.getAllByText('Learn more')
    await actor.click(learnMoreLinks[0])

    expect(learnMoreLinks[0].closest('a')).toHaveAttribute('target', '_blank')
    expect(learnMoreLinks[0].closest('a')).toHaveAttribute(
      'href',
      'https://github.com/zero-one-group/monorepo'
    )
  })
})
