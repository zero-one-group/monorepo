import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { loginAction } from '#/routes/login/-action'

const formWith = (fields: Record<string, string>) => {
  const fd = new FormData()
  for (const [k, v] of Object.entries(fields)) fd.set(k, v)
  return fd
}

describe('loginAction', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.spyOn(console, 'log').mockImplementation(() => {})
    vi.spyOn(console, 'error').mockImplementation(() => {})
  })
  afterEach(() => {
    vi.useRealTimers()
    vi.restoreAllMocks()
  })

  it('rejects a submission missing the email or the password', async () => {
    const incomplete: Record<string, string>[] = [
      { email: 'jane@example.com' },
      { password: 'secret' },
      {},
    ]
    for (const fields of incomplete) {
      const pending = loginAction({}, formWith(fields))
      await vi.runAllTimersAsync()
      expect(await pending).toEqual({ error: 'Email and password are required!' })
    }
  })

  it('returns a success state echoing the email with a timestamp', async () => {
    vi.setSystemTime(new Date('2026-08-26T10:00:00Z'))
    const pending = loginAction({}, formWith({ email: 'jane@example.com', password: 'secret' }))
    await vi.runAllTimersAsync()
    const state = await pending
    expect(state.success).toBe(true)
    expect(state.error).toBeUndefined()
    expect(state.data?.email).toBe('jane@example.com')
    // The action sleeps 1s before answering, so the stamp is taken after the fake clock advanced.
    expect(state.data?.timestamp).toBe('2026-08-26T10:00:01.000Z')
  })
})
