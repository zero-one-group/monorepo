import { join, resolve } from 'node:path'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import inspect from 'vite-plugin-inspect'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  plugins: [
    react({ include: /\.(mdx|js|jsx|ts|tsx)$/ }),
    inspect({ build: false, open: false }),
    tsconfigPaths(),
  ],
  define: { 'import.meta.env.APP_VERSION': `"${process.env.npm_package_version}"` },
  envDir: join(__dirname, '..'),
  envPrefix: ['VITE_'],
  clearScreen: true,
  server: {
    port: 8000,
    strictPort: true,
  },
  base: '/',
  root: resolve(__dirname),
  optimizeDeps: {
    // Do not optimize internal workspace dependencies.
    exclude: ['@myorg/shared-ui'],
  },
  resolve: {
    alias: [
      {
        // Configure an alias for the package. So, we don't have to restart
        // the Vite server each time when the former is performed.
        find: '@myorg/shared-ui',
        replacement: resolve(__dirname, '../../packages/shared-ui/src/index.ts'),
      },
    ],
  },
  build: {
    emptyOutDir: true,
    chunkSizeWarningLimit: 1024,
    reportCompressedSize: false,
    outDir: resolve(__dirname, 'dist'),
    rollupOptions: {
      output: {
        // Generate output with hash in filename.
        entryFileNames: `assets/[name]-[hash].js`,
        chunkFileNames: `assets/[name]-[hash].js`,
        assetFileNames: `assets/[name]-[hash].[ext]`,
      },
    },
  },
})
