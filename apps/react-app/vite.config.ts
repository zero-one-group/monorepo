import { resolve } from 'node:path'

import tailwindcss from '@tailwindcss/vite'
import { tanstackRouter } from '@tanstack/router-plugin/vite'
import react from '@vitejs/plugin-react'
import { env, isProduction } from 'std-env'
import { defineConfig } from 'vite'
import devtoolsJson from 'vite-plugin-devtools-json'
import tsconfigPaths from 'vite-tsconfig-paths'

// Check if the current environment is CI or test environment
const isTestOrStorybook = env.VITEST || process.argv[1]?.includes('storybook')

export default defineConfig({
  envPrefix: 'VITE_' /* Prefix for environment variables */,
  plugins: [
    tailwindcss(),
    !isTestOrStorybook &&
      tanstackRouter({
        target: 'react',
        autoCodeSplitting: true,
        routesDirectory: './app/routes',
        generatedRouteTree: './app/routeTree.gen.ts',
      }),
    react(),
    tsconfigPaths(),
    devtoolsJson(),
  ],
  server: { port: 3000, host: false },
  preview: { host: '127.0.0.1' },
  publicDir: resolve('public'),
  optimizeDeps: {
    // Do not optimize internal workspace dependencies.
    exclude: ['@repo/shared-ui'],
  },
  build: {
    outDir: 'build/client',
    emptyOutDir: true,
    ssr: false,
    minify: isProduction,
    cssMinify: isProduction,
    chunkSizeWarningLimit: 1024 * 2,
    reportCompressedSize: false,
    manifest: true,
    terserOptions: { format: { comments: false } },
  },
})
