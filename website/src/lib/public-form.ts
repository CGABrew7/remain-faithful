export const MIN_SUBMIT_MS = 2000
export const MAX_SUBMIT_MS = 24 * 60 * 60 * 1000

const MIN_LOADED_AT = 1_000_000_000_000

export function honeypotFilled(company: unknown) {
  return typeof company === 'string' && company.trim() !== ''
}

// submittedTooFast is true when the form was sent too quickly, the timestamp
// is missing, or the timestamp is older than a day.
export function submittedTooFast(loadedAt: unknown, now = Date.now()) {
  if (typeof loadedAt !== 'number' || !Number.isFinite(loadedAt) || loadedAt < MIN_LOADED_AT) {
    return true
  }
  const elapsed = now - loadedAt
  return elapsed < MIN_SUBMIT_MS || elapsed > MAX_SUBMIT_MS
}

export function emailLimitKey(email: string) {
  return `email:${email.trim().toLowerCase()}`
}
