/** Website donation section. Hidden unless this is exactly "true". */
export function donationsEnabled(): boolean {
  return process.env.NEXT_PUBLIC_DONATIONS_ENABLED === 'true'
}

/** Server donate routes. A Stripe key does not turn these on. */
export function donationsApiEnabled(): boolean {
  return process.env.DONATIONS_ENABLED === 'true'
}
