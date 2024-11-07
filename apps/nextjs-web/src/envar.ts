import { createEnv } from '@t3-oss/env-core'
import { z } from 'zod'

export default createEnv({
  server: {
    // DATABASE_URL: z.string().url(),
  },

  client: {},

  shared: {
    NODE_ENV: z.enum(['development', 'test', 'production']).default('development'),
    NEXT_PUBLIC_BASE_URL: z.string().url().optional().default('http://localhost:3000'),
    NEXT_PUBLIC_ALLOW_INDEXING: z.coerce.boolean().default(false),
  },

  /**
   * The prefix that client-side variables must have.
   * This is enforced both at a type-level and at runtime.
   */
  clientPrefix: 'NEXT_PUBLIC_',

  /**
   * What object holds the environment variables at runtime.
   * This is usually `process.env` or `import.meta.env`.
   */
  runtimeEnv: process.env,

  /**
   * By default, this library will feed the environment variables directly to
   * the Zod validator.
   *
   * This means that if you have an empty string for a value that is supposed
   * to be a number (e.g. `PORT=` in a ".env" file), Zod will incorrectly flag
   * it as a type mismatch violation. Additionally, if you have an empty string
   * for a value that is supposed to be a string with a default value (e.g.
   * `DOMAIN=` in an ".env" file), the default value will never be applied.
   *
   * In order to solve these issues, we recommend that all new projects
   * explicitly specify this option as true.
   */
  emptyStringAsUndefined: true,

  // Tell the library to skip validation if condition is true.
  skipValidation: false,
})
