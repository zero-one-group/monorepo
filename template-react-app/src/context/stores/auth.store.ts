import { persistentMap } from '@nanostores/persistent'
import type { LoginResponse } from '#/services/types/account'

type UserData = LoginResponse['data']['user']

type AuthStore = {
  loggedIn: boolean
  user: UserData | null
  role: 'admin' | 'user' | null
  accessToken: string | null
  refreshToken: string | null
}

// Default values for the AuthStore
const defaultAuthStoreValues: AuthStore = {
  loggedIn: false,
  user: null,
  role: null,
  accessToken: null,
  refreshToken: null,
}

/**
 * Configures a persistent key-value map store for the application's UI state.
 * The store is persisted to the browser's localStorage, using the 'auth:' prefix
 * for the keys. The store values are encoded and decoded using JSON.stringify
 * and JSON.parse, respectively.
 *
 * Using key-value map store. It will keep each key in separated localStorage key.
 * You can switch localStorage to any other storage for all used stores.
 * @ref: https://github.com/nanostores/persistent#persistent-engines
 */
const authStore = persistentMap<AuthStore>('auth:', defaultAuthStoreValues, {
  encode: JSON.stringify,
  decode: JSON.parse,
})

/**
 * Saves the current authentication state to the persistent store.
 * @param values - A partial object of the AuthStore type, containing the values to be updated in the store.
 */
function saveAuthState(values: Partial<AuthStore>) {
  authStore.set({ ...authStore.get(), ...values })
}

/**
 * Resets the authentication state to the default values.
 * This function can be used to log out the user and clear the authentication state.
 */
function resetAuthState() {
  authStore.set(defaultAuthStoreValues)
}

export { authStore, defaultAuthStoreValues, saveAuthState, resetAuthState }
export type { AuthStore, UserData }
