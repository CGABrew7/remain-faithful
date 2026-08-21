'use client'

import { useSearchParams } from 'next/navigation'
import { useEffect, useState } from 'react'

export default function DonationSuccessBanner() {
  const searchParams = useSearchParams()
  const [show, setShow] = useState(false)

  useEffect(() => {
    if (searchParams.get('donated') === 'true') {
      setShow(true)
      const url = new URL(window.location.href)
      url.searchParams.delete('donated')
      window.history.replaceState({}, '', url.toString())
    }
  }, [searchParams])

  if (!show) return null

  return (
    <div className="fixed top-20 left-1/2 -translate-x-1/2 z-50 w-full max-w-md px-4">
      <div
        className="p-5 flex items-start gap-4 bg-paper-raised"
        style={{ boxShadow: 'var(--shadow-border)' }}
      >
        <div className="flex-1">
          <p className="font-serif text-lg text-ink mb-1">Thank you for your donation.</p>
          <p className="text-sm text-ink-soft leading-relaxed">
            Your support keeps Remain Faithful free for everyone. We&apos;re deeply grateful.
          </p>
        </div>
        <button
          onClick={() => setShow(false)}
          className="relative text-ink-faint hover:text-ink transition-colors duration-200 shrink-0 w-10 h-10 -mr-2 -mt-2 inline-flex items-center justify-center"
          aria-label="Dismiss"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
    </div>
  )
}
