import { ofetch } from 'ofetch'
import { Await, json, useLoaderData } from 'react-router'

import { Suspense } from 'react'
import viteLogo from '/vite.svg'
import { useAuth } from '#/context/hooks/use-auth'
import reactLogo from '../assets/images/react.svg'

type Todo = {
  limit: number
  skip: number
  total: number
  todos: Array<{
    id: number
    todo: string
    completed: boolean
    userId: number
  }>
}

export async function loader() {
  const response = await ofetch<Todo>('https://dummyjson.com/todos?limit=9')
  return json<Todo['todos']>(response.todos)
}

export function Component() {
  const { logout } = useAuth()
  const data = useLoaderData() as Todo['todos']

  return (
    <>
      <button
        type="button"
        className="absolute top-4 right-12 font-semibold text-sm uppercase hover:underline"
        onClick={logout}
      >
        Logout
      </button>

      <div className="flex flex-row items-center justify-center gap-8">
        <a href="https://vitejs.dev" target="_blank" rel="noreferrer">
          <img src={viteLogo} className="size-16" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank" rel="noreferrer">
          <img src={reactLogo} className="size-16" alt="React logo" />
        </a>
      </div>

      <Suspense fallback={<div>Loading data...</div>}>
        <div className="mx-auto mt-8 grid max-w-6xl grid-cols-1 gap-4 sm:mt-14 sm:grid-cols-2 md:grid-cols-3">
          <Await resolve={data}>
            {data.map((item) => (
              <div key={item.id} className="rounded-md border bg-white p-4 shadow-sm">
                <h2 className="font-bold text-lg">#{item.id}</h2>
                <p>
                  <strong>Task:</strong> {item.todo}
                </p>
                <p>
                  <strong>Completed:</strong> {item.completed}
                </p>
                <p>
                  <strong>User ID:</strong> {item.userId}
                </p>
              </div>
            ))}
          </Await>
        </div>
      </Suspense>
    </>
  )
}
