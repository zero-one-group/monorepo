import { type NextRequest, NextResponse } from 'next/server'

/**
 * HTTP response status codes.
 * @reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
 */
export const HttpStatusCode: { [key: number]: string } = {
  200: 'The request was successful',
  400: 'Bad Request',
  401: 'Unauthorized',
  403: 'Forbidden',
  404: 'Resource Not Found',
  405: 'Method Not Allowed',
  422: 'Unprocessable Entity',
  429: 'Too Many Requests',
  500: 'Internal Server Error',
  502: 'Bad Gateway',
  503: 'Service Unavailable',
}

/**
 * Handles or throws a REST API error.
 *
 * @param {number} status - The HTTP status code for the response.
 * @param {string} message - An optional message to include in the response body.
 * @param {string|object|Error} reasons - An optional explanation or reasons for the response.
 * @returns {Response} - An HTTP response object.
 */
export function throwResponse(
  status: number,
  message: string | undefined = undefined,
  reasons: string | object | Error | undefined = undefined
): Response {
  const headers = new Headers()
  const authUrl = status === 401 ? 'http://localhost:3000/auth/login' : undefined
  const result = { code: status, message: message || HttpStatusCode[status], reasons, authUrl }

  headers.set('Content-Type', 'application/json')

  const responseInit: ResponseInit = { status: status, headers: headers }

  return new Response(JSON.stringify(result), responseInit)
}

/**
 * Creates an HTTP response with the specified status code, optional message, and optional data.
 *
 * @template T - The type of the data payload.
 * @param {number} status - The HTTP status code for the response.
 * @param {string} message - An optional message to include in the response body.
 * @param {T[] | T | undefined} data - The payload to include in the response. Can be an array of objects or a single object.
 * @returns {Response} - An HTTP response object.
 */
export function responseWithData<T>(
  status: number,
  message: string | undefined = undefined,
  data: T[] | T | undefined = undefined
): Response {
  const headers = new Headers()
  const payloadInfo = { code: status, message: message || HttpStatusCode[status] }

  const result = Array.isArray(data)
    ? { ...payloadInfo, data: [...data] }
    : { ...payloadInfo, data }

  headers.set('Content-Type', 'application/json')

  const responseInit: ResponseInit = { status: status, headers: headers }

  return new Response(JSON.stringify(result), responseInit)
}

/**
 * Redirects an HTTP request to the specified URL with an optional status code.
 *
 * @param {NextRequest} request - The incoming HTTP request object.
 * @param {string|URL} url - The target URL to redirect to, as a string or URL object.
 * @param {301|302|307|undefined} statusCode - An optional HTTP status code for the redirection.
 *   Defaults to 302 (Found) if not provided.
 * @returns {NextResponse} - A Next.js HTTP response object representing the redirection.
 */
export function redirect(
  request: NextRequest,
  url: string | URL,
  statusCode: 301 | 302 | 307 | undefined = 302
): NextResponse {
  const target = typeof url === 'string' ? new URL(url, request.nextUrl) : url
  return NextResponse.redirect(target, statusCode)
}
