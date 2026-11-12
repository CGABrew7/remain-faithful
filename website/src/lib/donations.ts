/** Website donation section. Hidden unless this is exactly "true". */
export function donationsEnabled(): boolean {
  return process.env.NEXT_PUBLIC_DONATIONS_ENABLED === 'true'
}
