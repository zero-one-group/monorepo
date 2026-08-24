import { createRootRoute, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'

import NotFound from '#/components/errors/404'
import InternalError from '#/components/errors/500'

function RootErrorComponent({ error }: { error: Error }) {
  const message = error?.message || 'Error'
  const details = 'An unexpected error occurred.'
  const stack = error?.stack

  return <InternalError message={message} details={details} stack={stack} />
}

export const Route = createRootRoute({
  component: RootComponent,
  errorComponent: RootErrorComponent,
  notFoundComponent: () => <NotFound />,
})

function RootComponent() {
  return (
    <>
      <Outlet />
      {import.meta.env.DEV && <TanStackRouterDevtools position="bottom-right" />}
    </>
  )
}
