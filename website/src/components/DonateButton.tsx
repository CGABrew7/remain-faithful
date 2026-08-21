'use client'

import { useState } from 'react'

const AMOUNTS = [
  { value: 5, impact: 'Supports one person for one week' },
  { value: 10, impact: 'Supports outreach to 5 new churches' },
  { value: 25, impact: 'Keeps the servers running for a month' },
  { value: 50, impact: 'Sponsors a full small group for a month' },
]

export default function DonateButton() {
  const [selected, setSelected] = useState(10)
  const [recurring, setRecurring] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const currentImpact = AMOUNTS.find((a) => a.value === selected)?.impact ?? ''

  async function handleDonate() {
    setLoading(true)
    setError('')
    try {
      const res = await fetch('/api/donate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ amount: selected, recurring }),
      })
      if (!res.ok) throw new Error('Request failed')
      const data = await res.json()
      if (data.url) {
        window.location.href = data.url
      } else {
        throw new Error('No checkout URL returned')
      }
    } catch {
      setError('Something went wrong. Please try again.')
      setLoading(false)
    }
  }

  return (
    <div className="flex flex-col items-stretch sm:items-center gap-5">
      <div className="inline-flex items-center gap-1 p-1 self-center" style={{ boxShadow: 'var(--shadow-border)' }}>
        <button
          onClick={() => setRecurring(false)}
          className={`px-5 py-2 min-h-10 font-mono text-[12px] tracking-[0.06em] transition-[color,background-color,scale] duration-150 ease-out active:scale-[0.96] ${
            !recurring ? 'bg-wax text-paper' : 'text-ink-soft hover:text-ink'
          }`}
        >
          One-time
        </button>
        <button
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

      <button
        onClick={handleDonate}
        disabled={loading}
        className="btn-wax self-center px-8"
      >
        {loading ? (
          <>
            <svg className="animate-spin w-4 h-4" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
            </svg>
            Redirecting...
          </>
        ) : (
          <>
            Give <span className="tabular-nums">${selected}{recurring ? '/mo' : ''}</span>
          </>
        )}
      </button>

      {error && (
        <p className="text-sm text-wax">{error}</p>
      )}

      <div className="text-center text-sm text-ink-soft max-w-sm leading-relaxed">
        Donations are made through the Woodfield Foundation Inc., a registered 501(c)(3) nonprofit organization. All donations are tax-deductible. Processed securely via Stripe.
      </div>
    </div>
  )
}
