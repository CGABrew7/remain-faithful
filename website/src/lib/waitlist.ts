// Website-side check for waitlist emails. The backend runs the full check again.
export function plausibleEmail(value: unknown): value is string {
  if (typeof value !== 'string') return false
  const email = value.trim()
  return email.length > 0 && email.length <= 254 && /^[^\s@]+@[^\s@.]+(\.[^\s@.]+)+$/.test(email)
}
