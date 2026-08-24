import { resolve } from 'node:path'

import tailwindcss from '@tailwindcss/vite'
import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import react from '@vitejs/plugin-react'
import { nitro } from 'nitro/vite'
import { env } from 'std-env'
import { defineConfig } from 'vite'
import devtoolsJson from 'vite-plugin-devtools-json'

// Check if the current environment is CI or test environment
const isTestOrStorybook = env.VITEST || process.argv[1]?.includes('storybook')

export default defineConfig(() => ({
  envPrefix: 'VITE_' /* Prefix for environment variables */,
  plugins: [
    tailwindcss(),
    !isTestOrStorybook &&
      tanstackStart({
        srcDirectory: 'src',
        router: {
          routesDirectory: 'routes',
        },
      }),
    !isTestOrStorybook && nitro(),
    react(),
    devtoolsJson(),
  ],
  resolve: {
    alias: {
      '#': resolve('./src'),
    },
    tsconfigPaths: true,
  },
  server: { port: 3300, host: false },
  publicDir: resolve('public'),
  optimizeDeps: {
    // Do not optimize internal workspace dependencies.
    exclude: ['@repo/shared-ui'],
  },
}))
