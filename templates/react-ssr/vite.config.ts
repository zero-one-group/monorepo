import { reactRouter } from '@react-router/dev/vite'
import tailwindcss from '@tailwindcss/vite'
import { resolve } from 'node:path'
import { env, isProduction } from 'std-env'
import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'

// Check if the current environment is CI or test environment
const isTestOrStorybook = env.VITEST || process.argv[1]?.includes('storybook')

export default defineConfig(({ isSsrBuild }) => ({
  envPrefix: 'VITE_' /* Prefix for environment variables */,
  plugins: [tailwindcss(), !isTestOrStorybook && reactRouter(), tsconfigPaths()],
  server: { port: {{ port_number }}, host: false },
  publicDir: resolve('public'),
  build: {
    manifest: true,
    emptyOutDir: true,
    chunkSizeWarningLimit: 1024 * 4,
    reportCompressedSize: false,
    minify: isProduction,
    rollupOptions: isSsrBuild ? { input: './server/app.ts' } : undefined,
    terserOptions: { format: { comments: false } },
  },
  esbuild: { legalComments: 'none' },
}))
