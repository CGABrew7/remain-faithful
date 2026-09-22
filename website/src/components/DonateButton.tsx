'use client'

import { useState } from 'react'
import {
  type DonateAmount,
  type DonateInterval,
  woodfieldDonateUrl,
} from '@/lib/donate-links'

const AMOUNTS: { value: DonateAmount; impact: string }[] = [
  { value: 5, impact: 'Supports one person for one week' },
  { value: 10, impact: 'Supports outreach to 5 new churches' },
  { value: 25, impact: 'Keeps the servers running for a month' },
  { value: 50, impact: 'Sponsors a full small group for a month' },
]

export default function DonateButton() {
  const [selected, setSelected] = useState<DonateAmount>(10)
  const [recurring, setRecurring] = useState(false)

  const interval: DonateInterval = recurring ? 'monthly' : 'one-time'
  const donateUrl = woodfieldDonateUrl(selected, interval)
  const currentImpact = AMOUNTS.find((a) => a.value === selected)?.impact ?? ''

  return (
    <div className="flex flex-col items-stretch sm:items-center gap-5">
      <div className="inline-flex items-center gap-1 p-1 self-center" style={{ boxShadow: 'var(--shadow-border)' }}>
        <button
          type="button"
          aria-pressed={!recurring}
          onClick={() => setRecurring(false)}
          className={`px-5 py-2 min-h-10 font-mono text-[12px] tracking-[0.06em] transition-[color,background-color,scale] duration-150 ease-out active:scale-[0.96] ${
            !recurring ? 'bg-wax text-paper' : 'text-ink-soft hover:text-ink'
          }`}
        >
          One-time
        </button>
        <button
          type="button"
          aria-pressed={recurring}
          onClick={() => setRecurring(true)}
          className={`px-5 py-2 min-h-10 font-mono text-[12px] tracking-[0.06em] transition-[color,background-color,scale] duration-150 ease-out active:scale-[0.96] ${
            recurring ? 'bg-wax text-paper' : 'text-ink-soft hover:text-ink'
          }`}
        >
          Monthly
        </button>
      </div>

      <div className="flex flex-wrap justify-center gap-2">
        {AMOUNTS.map((a) => (
          <button
            key={a.value}
            type="button"
            aria-pressed={selected === a.value}
            onClick={() => setSelected(a.value)}
            className={`min-w-16 h-10 px-3 font-mono text-sm tabular-nums transition-[color,background-color,box-shadow,scale] duration-150 ease-out active:scale-[0.96] ${
              selected === a.value
                ? 'bg-wax text-paper'
                : 'text-ink-soft hover:text-ink'
            }`}
            style={selected === a.value ? undefined : { boxShadow: 'var(--shadow-border)' }}
          >
            ${a.value}
          </button>
        ))}
      </div>

      {currentImpact && (
        <p className="font-mono text-[11px] tracking-[0.08em] uppercase text-ink-faint text-center tabular-nums">
          ${selected} {currentImpact.toLowerCase()}
        </p>
      )}

      <a
        href={donateUrl}
        target="_blank"
        rel="noopener noreferrer"
        className="btn-wax self-center px-8"
      >
        Give <span className="tabular-nums">${selected}{recurring ? '/mo' : ''}</span>
        <span className="sr-only"> (opens Stripe in a new tab)</span>
      </a>

      <div className="text-center text-sm text-ink-soft max-w-sm leading-relaxed">
        Donations are made through the Woodfield Foundation Inc., a registered 501(c)(3) nonprofit organization. All donations are tax-deductible. Processed securely via Stripe.
      </div>
    </div>
  )
}
