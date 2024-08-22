import { ErrorBoundary } from 'react-error-boundary'
import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '#/context/hooks/use-auth'
import RootLayout from '#/layouts/root-layout'
import InternalError from '#/pages/errors/internal-error'
import logger from '#/utils/logger'

export default function AuthLayout() {
  const { loggedIn } = useAuth()

  if (loggedIn) {
    return <Navigate to="/" replace />
  }

  return (
    <ErrorBoundary
      FallbackComponent={InternalError}
      onReset={() => logger.info('ErrorBoundary', 'reset errror state - auth')}
      resetKeys={['auth']}
    >
      <RootLayout className="flex h-full min-h-screen items-center">
        <Outlet />
      </RootLayout>
    </ErrorBoundary>
  )
}
