import { useStore } from '@nanostores/react'
import * as React from 'react'
import { useIsMobile } from '#/context/hooks/use-mobile'
import { saveUiState, uiStore } from '#/context/stores/ui.store'
import { TooltipProvider } from '../tooltip/tooltip'
import { sidebarStyles } from './sidebar.css'

type SidebarContext = {
  state: 'expanded' | 'collapsed'
  open: boolean
  setOpen: (open: boolean) => void
  openMobile: boolean
  setOpenMobile: (open: boolean) => void
  isMobile: boolean
  toggleSidebar: () => void
}

const SidebarContext = React.createContext<SidebarContext | null>(null)

function useSidebar() {
  const context = React.useContext(SidebarContext)
  if (!context) {
    throw new Error('useSidebar must be used within a SidebarProvider.')
  }

  return context
}

const SidebarProvider = React.forwardRef<
  HTMLDivElement,
  React.ComponentProps<'div'> & {
    defaultOpen?: boolean
    open?: boolean
    onOpenChange?: (open: boolean) => void
    shortcutKey?: string
  }
>(
  (
    {
      defaultOpen = true,
      open: openProp,
      onOpenChange: setOpenProp,
      shortcutKey = 'e',
      className,
      style,
      children,
      ...props
    },
    ref
  ) => {
    const isMobile = useIsMobile()
    const [openMobile, setOpenMobile] = React.useState(false)
    const uiState = useStore(uiStore)
    const styles = sidebarStyles()

    // This is the internal state of the sidebar.
    // We use openProp and setOpenProp for control from outside the component.
    const open = openProp ?? uiState.sidebar === 'expanded'

    const setOpen = React.useCallback(
      (value: boolean | ((value: boolean) => boolean)) => {
        const openState = typeof value === 'function' ? value(open) : value
        if (setOpenProp) {
          setOpenProp(openState)
        } else {
          saveUiState({ sidebar: openState ? 'expanded' : 'collapsed' })
        }

        // Save the state to local storage to persist the state between page loads.
        saveUiState({ sidebar: openState ? 'expanded' : 'collapsed' })
      },
      [setOpenProp, open]
    )

    // Helper to toggle the sidebar.
    const toggleSidebar = React.useCallback(() => {
      return isMobile ? setOpenMobile((open) => !open) : setOpen((open) => !open)
    }, [isMobile, setOpen])

    // Adds a keyboard shortcut to toggle the sidebar.
    React.useEffect(() => {
      const handleKeyDown = (event: KeyboardEvent) => {
        if (event.key === shortcutKey && (event.metaKey || event.ctrlKey) && event.shiftKey) {
          event.preventDefault()
          toggleSidebar()
        }
      }

      window.addEventListener('keydown', handleKeyDown)
      return () => window.removeEventListener('keydown', handleKeyDown)
    }, [toggleSidebar, shortcutKey])

    // We add a state so that we can do data-state="expanded" or "collapsed".
    // This makes it easier to style the sidebar with Tailwind classes.
    const state = open ? 'expanded' : 'collapsed'

    const contextValue = React.useMemo<SidebarContext>(
      () => ({
        state,
        open,
        setOpen,
        isMobile,
        openMobile,
        setOpenMobile,
        toggleSidebar,
      }),
      [state, open, setOpen, isMobile, openMobile, toggleSidebar]
    )

    return (
      <SidebarContext.Provider value={contextValue}>
        <TooltipProvider delayDuration={0}>
          <div
            style={{ ...style } as React.CSSProperties}
            className={styles.sidebarProvider({ className })}
            ref={ref}
            {...props}
          >
            {children}
          </div>
        </TooltipProvider>
      </SidebarContext.Provider>
    )
  }
)

SidebarProvider.displayName = 'SidebarProvider'

export { SidebarProvider, useSidebar }
