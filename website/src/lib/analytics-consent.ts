/** Same browser choice as the Allow analytics control. */
export const ANALYTICS_CONSENT_KEY = 'rf-analytics'

/** Fired in this tab when the visitor taps Allow analytics or No thanks. */
export const ANALYTICS_CONSENT_EVENT = 'rf-analytics-consent'

export function readAnalyticsConsent(
  storage: { getItem(key: string): string | null } | null | undefined,
): 'allow' | 'deny' | null {
  if (!storage) return null
  const saved = storage.getItem(ANALYTICS_CONSENT_KEY)
  if (saved === 'allow' || saved === 'deny') return saved
  return null
}
