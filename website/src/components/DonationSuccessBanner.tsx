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
    <div className="fixed top-20 left-1/2 -translate-x-1/2 z-50 w-full max-w-md px-4 animate-in fade-in slide-in-from-top-2 duration-300">
      <div
        className="rounded-2xl p-5 border border-[#C9A84C]/30 shadow-xl shadow-black/40 flex items-start gap-4"
        style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}
      >
        <div className="text-2xl flex-shrink-0">❤️</div>
        <div className="flex-1">
          <p className="font-semibold text-[#F0EDE8] mb-1">Thank you for your donation!</p>
          <p className="text-sm text-[#8A9BB0] leading-relaxed">
            Your support keeps Remain Faithful free for everyone. We&apos;re deeply grateful.
          </p>
        </div>
        <button
          onClick={() => setShow(false)}
          className="text-[#8A9BB0] hover:text-[#F0EDE8] transition-colors flex-shrink-0 mt-0.5"
          aria-label="Dismiss"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
    </div>
  )
}
