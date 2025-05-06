import { defineCommand } from 'citty'
import { consola } from 'consola'
import { generateRandomStr } from '../utils/string'

function generateBase64Key(length = 24): string {
  const randomStr = generateRandomStr({ size: length })
  return Buffer.from(randomStr).toString('base64')
}

function generateAppKeys(count = 4): string {
  return Array.from({ length: count })
    .map(() => generateBase64Key())
    .join(',')
}

export default defineCommand({
  meta: {
    name: 'make:strapi-keys',
    description: 'Create application keys for Strapi',
  },
  args: {
    help: {
      type: 'boolean',
      description: 'Print information about the command',
      default: false,
    },
  },
  run() {
    try {
      const appKeys = generateAppKeys()
      const apiTokenSalt = generateBase64Key()
      const adminJwtSecret = generateBase64Key()
      const transferTokenSalt = generateBase64Key()
      const jwtSecret = generateBase64Key()

      const output = `\nAPP_KEYS=${appKeys}
API_TOKEN_SALT=${apiTokenSalt}
ADMIN_JWT_SECRET=${adminJwtSecret}
TRANSFER_TOKEN_SALT=${transferTokenSalt}
JWT_SECRET=${jwtSecret}\n`

      consola.log(output)
    } catch (error) {
      consola.error(error instanceof Error ? error.message : 'Unknown error occurred')
      process.exit(1)
    }
  },
})
