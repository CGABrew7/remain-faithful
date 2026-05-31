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
    <div className="flex flex-col items-center gap-4">
      {/* One-time / Monthly toggle */}
      <div className="flex items-center gap-1 p-1 rounded-full bg-[#162235] border border-[#1E3050]">
        <button
          onClick={() => setRecurring(false)}
          className={`px-5 py-1.5 rounded-full text-sm font-semibold transition-all duration-200 ${
            !recurring
              ? 'bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] text-[#0F1B2D] shadow'
              : 'text-[#8A9BB0] hover:text-[#F0EDE8]'
          }`}
        >
          One-time
        </button>
        <button
          onClick={() => setRecurring(true)}
          className={`px-5 py-1.5 rounded-full text-sm font-semibold transition-all duration-200 ${
            recurring
              ? 'bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] text-[#0F1B2D] shadow'
              : 'text-[#8A9BB0] hover:text-[#F0EDE8]'
          }`}
        >
          Monthly
        </button>
      </div>

      {/* Amount selector */}
      <div className="flex gap-2">
        {AMOUNTS.map((a) => (
          <button
            key={a.value}
            onClick={() => setSelected(a.value)}
            className={`w-16 h-10 rounded-full text-sm font-semibold border transition-all duration-200 ${
              selected === a.value
                ? 'bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] text-[#0F1B2D] border-transparent shadow-md shadow-[#C9A84C]/30'
                : 'bg-transparent text-[#8A9BB0] border-[#1E3050] hover:border-[#C9A84C] hover:text-[#F0EDE8]'
            }`}
          >
            ${a.value}
          </button>
        ))}
      </div>

      {/* Impact label */}
      {currentImpact && (
        <p className="text-xs text-[#C9A84C]/80 text-center">
          ${selected} {currentImpact.toLowerCase()}
        </p>
      )}

      {/* Donate button */}
      <button
        onClick={handleDonate}
        disabled={loading}
        className="flex items-center gap-2 px-8 py-3 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all duration-200 shadow-lg shadow-[#C9A84C]/20 disabled:opacity-70 disabled:cursor-not-allowed"
      >
        {loading ? (
          <>
            <svg className="animate-spin w-4 h-4" viewBox="0 0 24 24" fill="none">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
            </svg>
            Redirecting...
          </>
        ) : (
          <>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
              <path d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z"/>
            </svg>
            Give ${selected}{recurring ? '/mo' : ''}
          </>
        )}
      </button>

      {error && (
        <p className="text-sm text-red-400">{error}</p>
      )}

      {/* Disclosure */}
      <div className="text-center text-xs text-[#8A9BB0]/70 max-w-sm leading-relaxed mt-1">
        Donations are made through the Woodfield Foundation Inc., a registered 501(c)(3) nonprofit organization. All donations are tax-deductible. Processed securely via Stripe.
      </div>
    </div>
  )
}
