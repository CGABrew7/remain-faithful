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
