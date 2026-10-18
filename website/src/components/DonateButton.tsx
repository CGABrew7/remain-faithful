'use client'

import { useState } from 'react'
import {
  type DonateAmount,
  type DonateInterval,
  woodfieldDonateUrl,
} from '@/lib/donate-links'

const AMOUNTS: DonateAmount[] = [5, 10, 25, 50]

export default function DonateButton() {
  const [selected, setSelected] = useState<DonateAmount>(10)
  const [recurring, setRecurring] = useState(false)

  const interval: DonateInterval = recurring ? 'monthly' : 'one-time'
  const donateUrl = woodfieldDonateUrl(selected, interval)

  return (
    <div className="flex flex-col items-stretch sm:items-center gap-5">
      <div className="inline-flex items-center gap-1 self-center rounded-[calc(var(--radius-control)+4px)] border border-white/10 bg-[#0C1D4C]/40 p-1">
        <button
          type="button"
          aria-pressed={!recurring}
          onClick={() => setRecurring(false)}
          className={`min-h-11 rounded-[var(--radius-control)] px-5 py-2 text-sm transition-[color,background-color,transform] duration-150 ease-out active:scale-[0.96] ${
            !recurring ? 'bg-wax text-paper' : 'text-ink-soft hover:text-white'
          }`}
        >
          One-time
        </button>
        <button
          type="button"
          aria-pressed={recurring}
          onClick={() => setRecurring(true)}
          className={`min-h-11 rounded-[var(--radius-control)] px-5 py-2 text-sm transition-[color,background-color,transform] duration-150 ease-out active:scale-[0.96] ${
            recurring ? 'bg-wax text-paper' : 'text-ink-soft hover:text-white'
          }`}
        >
          Monthly
        </button>
      </div>

      <div className="flex flex-wrap justify-center gap-2">
        {AMOUNTS.map((amount) => (
          <button
            key={amount}
            type="button"
            aria-pressed={selected === amount}
            onClick={() => setSelected(amount)}
            className={`h-11 min-w-16 rounded-[var(--radius-control)] px-3 text-sm tabular-nums transition-[color,background-color,box-shadow,transform] duration-150 ease-out active:scale-[0.96] ${
              selected === amount
                ? 'bg-wax text-paper'
                : 'border border-white/10 text-ink-soft hover:text-white'
            }`}
          >
            ${amount}
          </button>
        ))}
      </div>

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
