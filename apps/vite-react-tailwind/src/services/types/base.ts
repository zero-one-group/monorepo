import type { LogLevel } from 'consola'

export interface ApiClientOptions {
  /** Base URL for API requests */
  baseUrl?: string
  /** Custom headers for all requests */
  headers?: { [key: string]: string }
  /** Enable debug mode or provide a custom logging function */
  logLevel?: LogLevel
}

export interface BaseApiResponse {
  status: number
  success: boolean
  message?: string
}

export interface ErrorResponse extends BaseApiResponse {
  error?: {
    code?: number
    reason?: string
  }
}

export interface SuccessResponse<T> extends BaseApiResponse {
  data: T
}

export interface HealthCheckResponse extends BaseApiResponse {
  data: {
    uptime: number
    timestamp: number
  }
}
