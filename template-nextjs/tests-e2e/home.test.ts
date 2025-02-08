import { expect, test } from '@playwright/test'

test.describe.configure({ mode: 'serial' })

test.describe('@feature-default', () => {
  test('has title', async ({ page }) => {
    await page.goto('/')
    await expect(page).toHaveTitle(/React Router App/)
  })

  test('contain title', async ({ page }) => {
    await page.goto('/')
    await expect(page.getByRole('heading', { name: 'Zero One Starter Kit' })).toBeVisible()
  })
})
