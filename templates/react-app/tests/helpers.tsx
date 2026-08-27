import { createMemoryHistory, createRouter, RouterProvider } from '@tanstack/react-router'
import { render } from '@testing-library/react'

import { routeTree } from '#/routeTree.gen'

/** Renders the real app router at `path`, exactly as main.tsx wires it. */
export async function renderAt(path: string) {
  const router = createRouter({
    routeTree,
    history: createMemoryHistory({ initialEntries: [path] }),
  })
  await router.load()
  const utils = render(<RouterProvider router={router} />)
  return { router, ...utils }
}
