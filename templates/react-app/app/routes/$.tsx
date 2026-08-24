import { createFileRoute } from '@tanstack/react-router'

import NotFound from '#/components/errors/404'

export const Route = createFileRoute('/$')({
  component: NotFoundComponent,
})

function NotFoundComponent() {
  return <NotFound />
}
