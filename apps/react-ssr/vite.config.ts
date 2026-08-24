import { resolve } from 'node:path'

import { reactRouter } from '@react-router/dev/vite'
import tailwindcss from '@tailwindcss/vite'
import { env, isProduction } from 'std-env'
import { defineConfig } from 'vite'
import devtoolsJson from 'vite-plugin-devtools-json'

// Check if the current environment is CI or test environment
const isTestOrStorybook = env.VITEST || process.argv[1]?.includes('storybook')

export default defineConfig(({ isSsrBuild }) => ({
  envPrefix: 'VITE_' /* Prefix for environment variables */,
  plugins: [tailwindcss(), !isTestOrStorybook && reactRouter(), devtoolsJson()],
  server: { port: 3100, host: false },
  publicDir: resolve('public'),
  optimizeDeps: {
    // Do not optimize internal workspace dependencies.
    exclude: ['@repo/shared-ui'],
  },
  resolve: {
    alias: {
      '#': resolve('./app'),
    },
    tsconfigPaths: true,
  },
  build: {
    manifest: true,
    emptyOutDir: true,
    chunkSizeWarningLimit: 1024 * 4,
    reportCompressedSize: false,
    minify: isProduction,
    rolldownOptions: isSsrBuild ? { input: './server/app.ts' } : undefined,
  },
}))
