import type { ApiClientOptions } from './types/base'

const DEFAULT_OPTIONS: Omit<Required<ApiClientOptions>, 'headers'> = {
  baseUrl: import.meta.env.VITE_API_URL || '/api',
  logLevel: Number(import.meta.env.LOG_LEVEL) || (process.env.NODE_ENV === 'production' ? 3 : 4),
}

export default DEFAULT_OPTIONS
