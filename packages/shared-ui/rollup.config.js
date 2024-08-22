import { createRequire } from 'node:module'
import typescript from '@rollup/plugin-typescript'

const isProduction = process.env.NODE_ENV === 'production'
const require = createRequire(import.meta.url)
const pkgJson = require('./package.json')

/** @type {import('rollup').RollupOptions} */
export default {
  plugins: [
    typescript({
      cacheDir: 'node_modules/.cache/rollup-plugin-typescript',
      outputToFilesystem: false,
      tsconfig: 'tsconfig.json',
    }),
    isProduction && (await import('@rollup/plugin-terser')).default(),
  ],
  input: 'src/index.ts',
  treeshake: true,
  output: [
    {
      format: 'commonjs',
      file: pkgJson.main,
      sourcemap: !isProduction,
      name: pkgJson.name,
    },
    {
      format: 'esm',
      file: pkgJson.module,
      sourcemap: !isProduction,
    },
  ],
  external: [/\.css$/, ...Object.keys(pkgJson.dependencies), ...Object.keys(pkgJson.dependencies)],
}
