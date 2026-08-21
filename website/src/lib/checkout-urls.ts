/** Canonical production site for Stripe Checkout return URLs. */
export const PRODUCTION_SITE_ORIGIN = 'https://www.remainfaithful.com'

export const PRODUCTION_SUCCESS_URL = `${PRODUCTION_SITE_ORIGIN}/?donated=true`
export const PRODUCTION_CANCEL_URL = `${PRODUCTION_SITE_ORIGIN}/#donate`

/**
 * Stripe success/cancel URLs for a donate checkout session.
 *
 * Production (and any non-local Origin, including Vercel previews and the
 * dead remainfaithful.app domain) always returns to www.remainfaithful.com.
 * Localhost Origins are kept so local Stripe test checkouts can return to
 * the running Next.js app.
 */
export function checkoutReturnUrls(requestOrigin: string | null | undefined): {
  successUrl: string
  cancelUrl: string
} {
  const origin = checkoutOrigin(requestOrigin)
  return {
    successUrl: `${origin}/?donated=true`,
    cancelUrl: `${origin}/#donate`,
  }
}

export function checkoutOrigin(requestOrigin: string | null | undefined): string {
  if (requestOrigin && isLocalDevOrigin(requestOrigin)) {
    return requestOrigin.replace(/\/$/, '')
  }
  return PRODUCTION_SITE_ORIGIN
}

function isLocalDevOrigin(origin: string): boolean {
  try {
    const url = new URL(origin)
    if (url.protocol !== 'http:' && url.protocol !== 'https:') {
      return false
    }
    return isLocalhostHostname(url.hostname)
  } catch {
    return false
  }
}

function isLocalhostHostname(hostname: string): boolean {
  return hostname === 'localhost' || hostname === '127.0.0.1' || hostname === '[::1]'
}
