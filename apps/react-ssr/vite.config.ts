import { reactRouter } from '@react-router/dev/vite'
import tailwindcss from '@tailwindcss/vite'
import { resolve } from 'pathe'
import { isProduction } from 'std-env'
import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig(({ isSsrBuild }) => ({
  envPrefix: 'VITE_' /* Prefix for environment variables */,
  plugins: [tailwindcss(), reactRouter(), tsconfigPaths()],
  server: { port: 3001, host: false },
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
  esbuild: { legalComments: 'inline' },
}))
