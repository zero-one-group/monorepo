import type { Route } from './+types/page'
import { Welcome } from './welcome'

export function meta(_props: Route.MetaArgs) {
  return [
    { title: 'React Router App' },
    { name: 'description', content: 'Welcome to React Router!' },
  ]
}

export default function Page() {
  return <Welcome />
}
