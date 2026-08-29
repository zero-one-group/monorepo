import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import { renderAt } from '../helpers'

vi.mock('#/routes/login/-action', async (importOriginal) => {
  const mod = await importOriginal<typeof import('#/routes/login/-action')>()
  return {
    ...mod,
    // The template action sleeps 1s; the form tests are about the form, not the delay.
    loginAction: vi.fn(async (_prev: unknown, fd: FormData) => {
      const email = fd.get('email')
      const password = fd.get('password')
      if (!email || !password) return { error: 'Email and password are required!' }
      return { success: true, data: { email: String(email), timestamp: 'now' } }
    }),
  }
})

describe('/login', () => {
  it('renders the sign-in form with labelled fields and links', async () => {
    await renderAt('/login')
    expect(
      await screen.findByRole('heading', { name: /sign in to your account/i })
    ).toBeInTheDocument()
    expect(screen.getByLabelText('Email address')).toHaveAttribute('type', 'email')
    expect(screen.getByLabelText('Password')).toHaveAttribute('type', 'password')
    expect(screen.getByLabelText('Remember me')).not.toBeChecked()
    expect(screen.getByRole('button', { name: 'Sign in' })).toBeEnabled()
    expect(screen.getByRole('link', { name: /back to homepage/i })).toHaveAttribute('href', '/')
  })

  it('shows the action error when submitted empty', async () => {
    await renderAt('/login')
    await userEvent.setup().click(await screen.findByRole('button', { name: 'Sign in' }))
    expect(await screen.findByText('Email and password are required!')).toBeInTheDocument()
  })

  it('greets the user after a successful submission', async () => {
    const user = userEvent.setup()
    await renderAt('/login')
    await user.type(await screen.findByLabelText('Email address'), 'jane@example.com')
    await user.type(screen.getByLabelText('Password'), 'secret')
    await user.click(screen.getByRole('button', { name: 'Sign in' }))
    await waitFor(() => expect(screen.getByText('Hello jane@example.com')).toBeInTheDocument())
    expect(screen.queryByText('Email and password are required!')).not.toBeInTheDocument()
  })
})
