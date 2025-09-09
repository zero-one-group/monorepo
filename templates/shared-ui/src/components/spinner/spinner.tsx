import { useEffect, useState } from 'react'

export function LoadingSpinner() {
  return (
    <svg
      className="motion-preset-spin motion-duration-1000 size-3.5"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
      <path
        className="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      />
    </svg>
  )
}

interface Props {
  delay?: number
  fallback?: React.ReactNode
  element?: React.ReactNode
}

export function LazyLoadingSpinner(props: Props) {
  const { delay = 500 } = props
  const [show, setShow] = useState(false)

  useEffect(() => {
    const timeout = setTimeout(() => {
      setShow(true)
    }, delay)

    return () => {
      clearTimeout(timeout)
    }
  }, [delay])

  return show ? (props.element ?? <LoadingSpinner />) : (props.fallback ?? null)
}
