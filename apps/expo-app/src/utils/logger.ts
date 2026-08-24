type LogLevel = 'debug' | 'info' | 'warn' | 'error'

const prefix = '[expo-app]'

export function log(level: LogLevel, message: string, meta?: unknown) {
  if (meta === undefined) {
    console[level](`${prefix} ${message}`)
    return
  }

  console[level](`${prefix} ${message}`, meta)
}
