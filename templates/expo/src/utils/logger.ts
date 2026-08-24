type LogLevel = 'debug' | 'info' | 'warn' | 'error'

const prefix = '[{{ package_name | kebab_case }}]'

export function log(level: LogLevel, message: string, meta?: unknown) {
  if (meta === undefined) {
    console[level](`${prefix} ${message}`)
    return
  }

  console[level](`${prefix} ${message}`, meta)
}
