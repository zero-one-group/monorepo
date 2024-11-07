import type { NextRequest } from 'next/server'
import type { ServerRuntime } from 'next/types'

import { responseWithData, throwResponse } from '#/utils/response'

export const runtime: ServerRuntime = 'nodejs'

export async function GET(_req: NextRequest) {
  try {
    return responseWithData(200, 'All is well')
  } catch (error: any) {
    return error instanceof Response
      ? throwResponse(error.status, error.statusText)
      : throwResponse(500, error.message)
  }
}
