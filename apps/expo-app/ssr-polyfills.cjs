// Polyfills for browser globals that some dependencies (e.g. react-native-reanimated)
// reference at runtime but are undefined in the Node.js SSR environment used by
// `expo export -p web`. Defined lazily so they never override a real browser global.
if (typeof globalThis.requestAnimationFrame === 'undefined') {
  globalThis.requestAnimationFrame = (callback) => setTimeout(() => callback(Date.now()), 0)
  globalThis.cancelAnimationFrame = (id) => clearTimeout(id)
}
