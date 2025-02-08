import { randomBytes } from '@noble/hashes/utils'

/**
 * Generates a cryptographically secure random key with specified length
 * Uses uncrypto which provides isomorphic crypto API
 */
interface RandomStringOptions {
  size?: number
  pattern?: string
  prefix?: string
  digitsOnly?: boolean
  includeLower?: boolean
  includeUpper?: boolean
  includeSpecial?: boolean
}

export function generateRandomStr(config: RandomStringOptions = {}): string {
  const digits = '0123456789'
  const lowerChars = 'abcdefghijklmnopqrstuvwxyz'
  const upperChars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
  const specialChars = '!@#$%^&*()_+-=[]{}|;:,.<>?'

  let allowedChars = lowerChars + upperChars + digits

  if (config.pattern) allowedChars = config.pattern
  if (config.digitsOnly) allowedChars = digits
  if (config.includeLower === false) allowedChars = allowedChars.replace(lowerChars, '')
  if (config.includeUpper === false) allowedChars = allowedChars.replace(upperChars, '')
  if (config.includeSpecial) allowedChars += specialChars

  const size = config.size || 10
  const bytes = randomBytes(size)

  let result = ''
  for (let i = 0; i < size; i++) {
    const byte = bytes[i] ?? 0
    result += allowedChars[byte % allowedChars.length] ?? allowedChars[0]
  }

  return config.prefix ? `${config.prefix}${result}` : result
}
