import { reactRouter } from '@react-router/dev/vite'
import tailwindcss from '@tailwindcss/vite'
import { resolve } from 'pathe'
import { isProduction } from 'std-env'
import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  envPrefix: 'VITE_' /* Prefix for environment variables */,
  plugins: [tailwindcss(), reactRouter(), tsconfigPaths()],
  server: { port: 3000, host: true },
  publicDir: resolve('public'),
  build: {
    ssr: false,
    minify: isProduction,
    cssMinify: isProduction,
    chunkSizeWarningLimit: 1024 * 2,
    reportCompressedSize: false,
    emptyOutDir: true,
    manifest: true,
  },
})
