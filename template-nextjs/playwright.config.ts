/**
 * Reference: https://playwright.dev/docs/test-configuration
 * Example: https://dev.to/this-is-learning/playwright-lets-start-2mdj
 */

import 'dotenv/config'
import { defineConfig, devices } from '@playwright/test'
import { resolve } from 'pathe'
import { env, isCI } from 'std-env'

export const STORAGE_STATE = resolve('.playwright/user.json')

export default defineConfig({
  quiet: !!isCI,
  testDir: resolve('tests-e2e'),
  outputDir: resolve('tests-results/e2e'),
  fullyParallel: true,
  forbidOnly: !!isCI,
  retries: isCI ? 2 : 0,
  workers: isCI ? 1 : undefined,
  reporter: [['html', { open: 'never', outputFolder: './tests-results/e2e' }], ['list']],
  use: {
    baseURL: env.URL || 'http://localhost:{{ port_number }}',
    ...devices['Desktop Chrome'],
    defaultBrowserType: 'chromium',
    colorScheme: 'no-preference',
    locale: 'en-US',
    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',
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
        command: 'moon {{ package_name | kebab_case }}:build && moon {{ package_name | kebab_case }}:start-node',
        reuseExistingServer: !isCI,
        timeout: 30_000,
        port: {{ port_number }},
      },
})
