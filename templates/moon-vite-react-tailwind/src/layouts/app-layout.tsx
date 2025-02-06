import { ErrorBoundary } from 'react-error-boundary'
import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useAuth } from '#/context/hooks/use-auth'
import InternalError from '#/pages/errors/internal-error'
import logger from '#/utils/logger'
import RootLayout from './root-layout'

export default function AppLayout() {
  const { loggedIn } = useAuth()
  const { pathname } = useLocation()

  if (!loggedIn) {
    return <Navigate to={`/auth/login?returnTo=${pathname}`} replace />
  }

  return (
    <ErrorBoundary
      FallbackComponent={InternalError}
      onReset={() => logger.info('ErrorBoundary', 'reset errror state - app')}
      resetKeys={['app']}
    >
      <RootLayout className="h-full min-h-screen p-10">
        <Outlet />
      </RootLayout>
    </ErrorBoundary>
  )
}
