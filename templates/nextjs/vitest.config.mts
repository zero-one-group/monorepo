import react from '@vitejs/plugin-react'
import { isCI } from 'std-env'
import tsconfigPaths from 'vite-tsconfig-paths'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [tsconfigPaths(), react()],
  test: {
    environment: 'happy-dom',
    exclude: ['node_modules', 'tests-e2e'],
    reporters: isCI ? ['html', 'github-actions'] : ['html', 'default'],
    include: ['./tests/**/*.{test,spec}.{ts,tsx}'],
    setupFiles: ['./tests/setup-client.ts'],
    outputFile: {
      json: './tests-results/vitest-results.json',
      html: './tests-results/index.html',
    },
    coverage: {
      provider: 'v8',
      reporter: ['html-spa', 'text-summary'],
      reportsDirectory: './tests-results/coverage',
      cleanOnRerun: true,
      clean: true,
    },
    globals: true,
  },
})
