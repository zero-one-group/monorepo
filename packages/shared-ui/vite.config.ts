import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { resolve } from 'pathe'
import { env, isDevelopment } from 'std-env'
import { defineConfig } from 'vite'
import dts from 'vite-plugin-dts'
import tsconfigPaths from 'vite-tsconfig-paths'
import pkg from './package.json' assert { type: 'json' }

// Check if the current environment is CI or test environment
const isTestOrStorybook = env.VITEST || process.argv[1]?.includes('storybook')

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    tsconfigPaths(),
    /* Generate declarations for TypeScript */
    !isTestOrStorybook &&
      dts({
        include: ['src'],
        rollupTypes: true,
        exclude: ['**/*.stories.@(ts|tsx)'],
        tsconfigPath: resolve('tsconfig.json'),
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
        theme: resolve('src/theme.tsx'),
        utils: resolve('src/utils.ts'),
      },
      formats: ['es'],
    },
    outDir: resolve('dist'),
    rollupOptions: {
      watch: {
        include: ['src/**'],
        exclude: ['src/**/*.stories.@(ts|tsx)'],
      },
      treeshake: true,
      output: {
        format: 'es',
        exports: 'named',
        entryFileNames: '[name]/index.js',
        chunkFileNames: '_chunks/[name].js',
        assetFileNames: 'assets/[name].[ext]',
        manualChunks: undefined, // Let rollup handle chunking automatically
        preserveModules: false, // To enable tree-shaking set to `true`
        reexportProtoFromExternal: false,
        preserveModulesRoot: 'src',
        banner: "'use client';",
      },
      external: [...Object.keys(pkg.peerDependencies || {})],
    },
  },
})
