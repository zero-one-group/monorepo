import { useMemo } from 'react'
import { API_BASE_URL, LOG_LEVEL } from '#/config'
import { ApiClient } from '#/services'

export function useApiClient() {
  return useMemo(
    () =>
      ApiClient.getInstance({
        baseUrl: API_BASE_URL,
        logLevel: LOG_LEVEL,
      }),
    []
  )
}
