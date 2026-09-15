import { COMPARE_UPDATED_ON } from '@/lib/competitor-pricing'

/** Visible recency stamp for buying-intent compare pages. */
export function UpdatedStamp() {
  return (
    <p className="text-xs text-ink-soft tabular-nums">
      Updated on {COMPARE_UPDATED_ON}
    </p>
  )
}
