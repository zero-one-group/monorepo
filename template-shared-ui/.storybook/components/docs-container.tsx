import { DocsContainer as BaseContainer } from '@storybook/blocks'
import * as React from 'react'
import { dark, light, listenToColorScheme } from '../themes'

const themes = { light, dark }

/**
 * Switch color scheme based on the global types or system preferences
 */
export const DocsContainer: typeof BaseContainer = ({ children, context }) => {
  const [theme, setTheme] = React.useState('light')

  React.useEffect(() => listenToColorScheme(context.channel, setTheme), [context.channel])

  React.useEffect(() => {
    document.documentElement.dataset.theme = theme
  }, [theme])

  return (
    <BaseContainer context={context} theme={themes[theme]}>
      {children}
    </BaseContainer>
  )
}
