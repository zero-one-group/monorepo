import { useCallback, useState } from 'react'
import logger from '#/utils/logger'

export const useLocalStorage = <T>(
  keyName: string,
  defaultValue: T
): [T, (newValue: T) => void] => {
  const isClient = typeof window !== 'undefined'

  const [storedValue, setStoredValue] = useState<T>(() => {
    if (!isClient) {
      return defaultValue
    }

    try {
      const value = window.localStorage.getItem(keyName)
      if (!value) {
        window.localStorage.setItem(keyName, JSON.stringify(defaultValue))
        return defaultValue
      }
      return JSON.parse(value)
    } catch (err) {
      logger.error('useLocalStorage', err)
      return defaultValue
    }
  })

  const setValue = useCallback(
    (newValue: T) => {
      if (!isClient) {
        return
      }

      try {
        window.localStorage.setItem(keyName, JSON.stringify(newValue))
        setStoredValue(newValue)
      } catch (err) {
        logger.error('setValue', err)
      }
    },
    [keyName, isClient]
  )

  return [storedValue, setValue]
}
