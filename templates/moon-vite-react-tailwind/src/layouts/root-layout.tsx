import React, { Suspense } from 'react'
import { Toaster } from 'sonner'
import PageLoader from '#/components/page-loader'
import { cn } from '#/utils/helper'

// Used for health check
// import { useErrorBoundary } from 'react-error-boundary'
// import { useApiClient } from '#/context/hooks/api-client'
// import type { ErrorResponse } from '#/services'
// import logger from '#/utils/logger'

interface RootLayoutProps {
  children: React.ReactNode
  className?: string
}

export default function RootLayout({ children, className }: RootLayoutProps) {
  // const { showBoundary } = useErrorBoundary()
  // const api = useRef(useApiClient()).current

  // useEffect(() => {
  //   const doHealthCheck = async () => {
  //     try {
  //       const result = await api._healthCheck()
  //       if (result.status !== 200) {
  //         showBoundary({ message: result.message || 'API is not healthy' })
  //       }
  //       logger.info('[RESULT] doHealthCheck', result.status)
  //     } catch (error: any) {
  //       logger.error('[ERROR] doHealthCheck', error)
  //       if (error instanceof Error) {
  //         showBoundary({ message: error.message })
  //       } else if (typeof error.error === 'object' && 'error' in error) {
  //         const errorResponse = error as ErrorResponse
  //         showBoundary({
  //           message: errorResponse.error?.reason || 'Unknown error occurred',
  //           code: errorResponse.error?.code,
  //         })
  //       } else {
  //         showBoundary({ message: 'An unexpected error occurred' })
  //       }
  //     }
  //   }

  //   doHealthCheck()
  // }, [api, showBoundary])

  return (
    <React.Fragment>
      <Suspense fallback={<PageLoader />}>
        <div className={cn('bg-neutral-50 dark:bg-neutral-950', className)}>{children}</div>
      </Suspense>
      <Toaster richColors theme="system" />
    </React.Fragment>
  )
}
