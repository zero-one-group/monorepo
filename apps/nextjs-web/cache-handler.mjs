import { CacheHandler } from '@neshca/cache-handler'
import createLocalHandler from '@neshca/cache-handler/local-lru'

CacheHandler.onCreation(async () => {
  const localHandler = createLocalHandler({
    maxItemsNumber: 10000,
    maxItemSizeBytes: 1024 * 1024 * 500,
  })

  return { handlers: [localHandler] }
})

export default CacheHandler
