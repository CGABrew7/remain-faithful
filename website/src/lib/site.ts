/** Public site origin. Production serves this host. */
export const SITE_ORIGIN = 'https://www.remainfaithful.com'

export function siteUrl(path = '/'): string {
  if (path === '/' || path === '') return SITE_ORIGIN
  return `${SITE_ORIGIN}${path.startsWith('/') ? path : `/${path}`}`
}
