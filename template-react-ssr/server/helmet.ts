import { cspConnectSource, cspFrameSource } from './constants.js'
import { cspFontSource, cspImgSources, cspScriptSource } from './constants.js'

export function generateCspDirectives() {
  const connectSource =
    process.env.NODE_ENV === 'development' ? ['ws://localhost:*'] : cspConnectSource

  // 'unsafe-eval' required for DOMPurify
  const scriptSrc = [
    "'self'",
    "'report-sample'",
    "'unsafe-inline'",
    "'unsafe-eval'",
    ...cspScriptSource,
  ]

  // FIXME - this is a hack to get the CSP working in production
  // `unsafe-inline` allows the <LiveReload /> component to load without a nonce in the error pages
  if (process.env.NODE_ENV !== 'development') {
    // Remove unsafe-inline and add a nonce to the script-src in production
    // scriptSrc.splice(scriptSrc.indexOf("'unsafe-inline'"), 1)
    // scriptSrc.splice(scriptSrc.indexOf("'unsafe-eval'"), 1)
    // scriptSrc.push((_req, res) => `'nonce-${res.locals.nonce}'`)
  }

  return {
    'base-uri': ["'self'"],
    'child-src': ["'self'"],
    'connect-src': ["'self'", ...connectSource],
    'default-src': ["'self'"],
    'font-src': ["'self'", ...cspFontSource],
    'form-action': ["'self'"],
    'frame-ancestors': ["'none'"],
    'frame-src': ["'self'", ...cspFrameSource],
    'img-src': ["'self'", 'data:', ...cspImgSources],
    'manifest-src': ["'self'"],
    'media-src': ["'self'"],
    'object-src': ["'none'"],
    'script-src': [...scriptSrc],
    'style-src': ["'self'", "'report-sample'", "'unsafe-inline'"],
    'worker-src': ["'self'", 'blob:'],
  }
}
