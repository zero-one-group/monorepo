import { useStore } from '@nanostores/react'
import { createContext, useEffect, useState } from 'react'
import { type Theme, saveUiState, uiStore } from '../stores'

type ThemeProviderState = {
  theme: Theme
  setTheme: (theme: Theme) => void
}

const initialState: ThemeProviderState = {
  theme: 'system',
  setTheme: () => null,
}

export const ThemeProviderContext = createContext<ThemeProviderState>(initialState)

export function ThemeProvider({ children }: React.PropsWithChildren) {
  const uiState = useStore(uiStore)
  const [theme, setTheme] = useState<Theme>(() => uiState.theme)

  useEffect(() => {
    const root = document.documentElement

    // Update data-theme accordingly if user selects light or dark
    if (theme !== 'system') {
      root.dataset.theme = theme
      return
    }

    // For auto mode, we need to watch system preferences
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')

    // Set initial theme based on system preference
    root.dataset.theme = mediaQuery.matches ? 'dark' : 'light'

    // Update theme when system preference changes
    function handleChange(event: MediaQueryListEvent) {
      root.dataset.theme = event.matches ? 'dark' : 'light'
    }

    mediaQuery.addEventListener('change', handleChange)
    return () => mediaQuery.removeEventListener('change', handleChange)
  }, [theme])

  const value = {
    theme,
    setTheme: (newTheme: Theme) => {
      saveUiState({ theme: newTheme })
      setTheme(newTheme)
    },
  }

  return <ThemeProviderContext.Provider value={value}>{children}</ThemeProviderContext.Provider>
}
