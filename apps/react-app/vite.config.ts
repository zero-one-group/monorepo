import { reactRouter } from '@react-router/dev/vite'
import tailwindcss from '@tailwindcss/vite'
import { resolve } from 'pathe'
import { env, isProduction } from 'std-env'
import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'

// Check if the current environment is CI or test environment
const isTestOrStorybook = env.VITEST || process.argv[1]?.includes('storybook')

export default defineConfig({
  envPrefix: 'VITE_' /* Prefix for environment variables */,
  plugins: [tailwindcss(), !isTestOrStorybook && reactRouter(), tsconfigPaths()],
  server: { port: 3000, host: false },
  publicDir: resolve('public'),
  optimizeDeps: {
    // Do not optimize internal workspace dependencies.
    exclude: ['@repo/shared-ui'],
  },
  build: {
    ssr: false,
    minify: isProduction,
    cssMinify: isProduction,
    chunkSizeWarningLimit: 1024 * 2,
    reportCompressedSize: false,
    emptyOutDir: true,
    manifest: true,
    terserOptions: { format: { comments: false } },
  },
  esbuild: { legalComments: 'none' },
})
