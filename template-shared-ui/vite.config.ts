import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { resolve } from 'pathe'
import { isDevelopment } from 'std-env'
import { defineConfig } from 'vite'
import dts from 'vite-plugin-dts'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    tsconfigPaths(),
    /* Generate declarations for TypeScript */
    dts({
      include: ['src'],
      rollupTypes: true,
      compilerOptions: {
        declarationMap: true,
        sourceMap: true,
      },
    }),
  ],
  server: { port: 6300, host: false },
  build: {
    emptyOutDir: true,
    chunkSizeWarningLimit: 1024 * 2,
    reportCompressedSize: false,
    copyPublicDir: false,
    minify: !isDevelopment,
    sourcemap: true,
    lib: {
      entry: {
        components: resolve('src/components/index.ts'),
        hooks: resolve('src/hooks/index.ts'),
        theme: resolve('src/theme/index.ts'),
        utils: resolve('src/utils.ts'),
      },
      formats: ['es'],
    },
    outDir: resolve('dist'),
    rollupOptions: {
      external: ['react', 'react/jsx-runtime', 'react-dom'],
      output: {
        preserveModules: false,
        preserveModulesRoot: 'src',
        entryFileNames: '[name].js',
        chunkFileNames: '[name].js',
        assetFileNames: '[name].[ext]',
      },
    },
  },
})
