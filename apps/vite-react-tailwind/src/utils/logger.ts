import type { ConsolaInstance } from 'consola/core'
import { createConsola } from 'consola/core'
import { LOG_LEVEL } from '#/config'

/**
 * Creates a Consola instance for logging.
 * @returns {ConsolaInstance} A Consola instance for logging.

 * @see: https://unjs.io/packages/consola
 */
const logger: ConsolaInstance = createConsola({
  level: LOG_LEVEL,
})

export default logger
