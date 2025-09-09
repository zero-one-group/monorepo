import { persistentMap } from '@nanostores/persistent'
import type { MapStore } from 'nanostores'
import { storeDecode, storeEncode } from './utils'

export type Theme = 'dark' | 'light' | 'system'

type UIStore = {
  sidebar: 'expanded' | 'collapsed'
  theme: Theme
}

/**
 * The default values for the UI store, which includes the initial state of the sidebar.
 */
const defaultUIStoreValues: UIStore = {
  sidebar: 'expanded',
  theme: 'system',
}

/**
 * A persistent map store for the UI state, with the default values for the sidebar state.
 * Using key-value map store. It will keep each key in separated localStorage key.
 * You can switch localStorage to any other storage for all used stores.
 * @ref: https://github.com/nanostores/persistent#persistent-engines
 */
const uiStore: MapStore<UIStore> = persistentMap<UIStore>('ui:', defaultUIStoreValues, {
  encode: storeEncode,
  decode: storeDecode,
})

/**
 * Saves the current UI state by merging the provided partial UI store values with the existing values.
 * @param values - A partial object of the UI store values to be merged with the existing state.
 */
function saveUiState(values: Partial<UIStore>) {
  uiStore.set({ ...uiStore.get(), ...values })
}

/**
 * Resets the UI store to its default values, which includes setting the sidebar state to 'expanded'.
 */
function resetUiState() {
  uiStore.set(defaultUIStoreValues)
}

export { uiStore, defaultUIStoreValues, saveUiState, resetUiState }
export type { UIStore }
