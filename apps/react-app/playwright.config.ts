/**
 * Reference: https://playwright.dev/docs/test-configuration
 * Example: https://dev.to/this-is-learning/playwright-lets-start-2mdj
 */

import 'dotenv/config'
import { defineConfig, devices } from '@playwright/test'
import { resolve } from 'pathe'
import { env, isCI } from 'std-env'

export const STORAGE_STATE = resolve('.playwright/user.json')
const APP_PORT = env.PORT || 3000

export default defineConfig({
  quiet: !!isCI,
  testDir: resolve('tests-e2e'),
  outputDir: resolve('tests-results/e2e'),
  fullyParallel: true,
  forbidOnly: !!isCI,
  retries: isCI ? 2 : 0,
  workers: isCI ? 1 : undefined,
  reporter: [['html', { open: 'never', outputFolder: './tests-results/e2e-html' }], ['list']],
  preserveOutput: 'always',
  use: {
    ...devices['Desktop Chrome'],
    baseURL: env.URL || `http://localhost:${APP_PORT}`,
    defaultBrowserType: 'chromium',
    colorScheme: 'no-preference',
    ignoreHTTPSErrors: true,
    locale: 'en-US',
    trace: 'on-all-retries' /* @see https://playwright.dev/docs/trace-viewer */,
    video: { mode: isCI ? 'off' : 'on' },
    contextOptions: {
      recordVideo: { dir: resolve('tests-results/e2e-videos') },
      reducedMotion: 'reduce',
    },
    screenshot: isCI ? 'off' : 'only-on-failure',
    launchOptions: {
      slowMo: 1000, // a 1000 milliseconds pause before each operation. Useful for slow systems.
    },
  },
  /* Configure projects for major browsers */
  projects: [
    { name: 'Chromium', use: { ...devices['Desktop Chrome'] } },
    { name: 'Firefox', use: { ...devices['Desktop Firefox'] } },
    { name: 'Safari', use: { ...devices['Desktop Safari'] } },
    { name: 'Mobile Chrome', use: { ...devices['Pixel 5'] } },
  ],
  /* Run your local dev server before starting the tests */
  webServer: env.URL
    ? undefined
    : {
        command: 'moon react-app:build && moon react-app:start',
        reuseExistingServer: !isCI,
        port: Number(APP_PORT),
        timeout: 30_000,
      },
})
