import { type RouteObject, createBrowserRouter, useRoutes } from 'react-router-dom'
import NotFound from '#/pages/errors/not-found'

// Application layouts
import AppLayout from '#/layouts/app-layout'
import AuthLayout from '#/layouts/auth-layout'

// Authentication pages
import SignInPage from '#/pages/auth/login'
import ForgotPasswordPage from '#/pages/auth/recovery'
import ResetPasswordPage from '#/pages/auth/resetpass'
import SignUpPage from '#/pages/auth/signup'

const route = (path: string, { ...props }: RouteObject) => ({ path, ...props })

/**
 * Using dynamic import for the pages to reduce the bundle size.
 *
 * IMPORTANT: Ensure the imported module exports both 'Component' and 'loader'.
 * These exports are required for proper routing and data loading.
 *
 * @see https://reactrouter.com/en/route/lazy#statically-defined-properties
 */
const routes: RouteObject[] = [
  route('/', {
    element: <AppLayout />,
    children: [
      {
        index: true,
        lazy: () => import('#/pages/dashboard'),
      },
      // {
      //   path: 'settings',
      //   element: <AppSettings />,
      //   children: [
      //     { index: true, lazy: () => import('#/pages/settings/general') },
      //     { path: 'account', lazy: () => import('#/pages/settings/account') },
      //   ],
      // },
    ],
  }),
  route('/auth', {
    element: <AuthLayout />,
    children: [
      { path: 'login', element: <SignInPage /> },
      { path: 'register', element: <SignUpPage /> },
      { path: 'recovery', element: <ForgotPasswordPage /> },
      { path: 'resetpass', element: <ResetPasswordPage /> },
    ],
  }),
  route('*', { element: <NotFound /> }),
]

/**
 * Creates a browser-based router instance using the provided `routes` configuration.
 * This router instance can be used with the `RouterProvider` component to render the application's routes.
 *
 * @example
 *
 * import { RouterProvider } from 'react-router-dom'
 * import { routes, browserRoutes } from './routes'
 *
 * const App = () => {
 *   return <RouterProvider router={browserRoutes} />
 * }
 *
 */
const browserRoutes = createBrowserRouter(routes)

/**
 * Renders the application's routes using the `useRoutes` hook from `react-router-dom`.
 * This component is responsible for setting up the routing structure and rendering the
 * appropriate components based on the current URL.
 *
 * @example
 *
 * import { BrowserRouter } from 'react-router-dom'
 * import AppRoutes from './routes'
 *
 * export default function App() {
 *   return (
 *     <BrowserRouter>
 *       <AppRoutes />
 *     </BrowserRouter>
 *   )
 * }
 *
 * @returns {JSX.Element} The rendered routes for the application.
 */
const AppRoutes = (): React.ReactElement | null => useRoutes(routes)

export { routes, browserRoutes }

export default AppRoutes
