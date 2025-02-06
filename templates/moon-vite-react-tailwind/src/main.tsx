import './assets/styles/main.css'

import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import AppProvider from '#/context/providers/app-provider'
import { browserRoutes } from './routes'

const rootElement = document.getElementById('root')

if (!rootElement) {
  throw new Error(
    "Root element not found. Check if it's in your index.html or if the id is correct."
  )
}

// When you use Strict Mode, React renders each component twice to help you find unexpected side effects.
// @ref: https://react.dev/blog/2022/03/08/react-18-upgrade-guide#react
ReactDOM.createRoot(rootElement).render(
  <React.StrictMode>
    <AppProvider debugScreenSize={import.meta.env.DEV}>
      <RouterProvider router={browserRoutes} />
    </AppProvider>
  </React.StrictMode>
)
