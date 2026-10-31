/** Routes that are not public website pages. Clarity stays off these. */
const PRIVATE_PREFIXES = [
  '/admin',
  '/reset-password',
  '/verify-email',
  '/invite',
  '/accept-invite',
  '/join',
]

export function isPublicClarityPath(pathname: string): boolean {
  if (!pathname || pathname.startsWith('/api/')) return false
  return !PRIVATE_PREFIXES.some((prefix) => pathname === prefix || pathname.startsWith(`${prefix}/`))
}

/** Clarity's tag (and its _clck / _clsk cookies) only run after Allow analytics. */
export function shouldLoadClarity(pathname: string, consent: string | null | undefined): boolean {
  return consent === 'allow' && isPublicClarityPath(pathname)
}
