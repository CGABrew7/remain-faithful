'use client'

import { useEffect, useState } from 'react'
import { usePathname } from 'next/navigation'
import Script from 'next/script'

const KEY = 'rf-analytics'

export default function AnalyticsConsent({ gaId }: { gaId?: string }) {
  const pathname = usePathname()
  const [choice, setChoice] = useState<'unknown' | 'allow' | 'deny'>('unknown')

  useEffect(() => {
    const saved = window.localStorage.getItem(KEY)
    if (saved === 'allow' || saved === 'deny') setChoice(saved)
  }, [])

  if (!gaId || pathname === '/admin' || pathname.startsWith('/admin/')) return null

  function choose(next: 'allow' | 'deny') {
    window.localStorage.setItem(KEY, next)
    setChoice(next)
  }

  return (
    <>
      {choice === 'allow' && (
        <>
          <Script src={`https://www.googletagmanager.com/gtag/js?id=${gaId}`} strategy="afterInteractive" />
          <Script id="google-analytics" strategy="afterInteractive">
            {`
              window.dataLayer = window.dataLayer || [];
              function gtag(){dataLayer.push(arguments);}
              gtag('js', new Date());
              gtag('config', '${gaId}');
            `}
          </Script>
        </>
      )}
      {choice === 'unknown' && (
        <div
          role="dialog"
          aria-label="Website analytics"
          className="fixed bottom-4 left-4 right-4 z-50 mx-auto max-w-xl rounded-[14px] border border-white/10 bg-[#0C1D4C] p-4 shadow-lg"
        >
          <p className="text-sm leading-relaxed text-ink-soft">
            Allow analytics so we can see which pages people open. It stays off until you say yes. The iPhone app does not include analytics.
          </p>
          <div className="mt-3 flex flex-col gap-2 sm:flex-row">
            <button type="button" className="btn-wax" onClick={() => choose('allow')}>
              Allow analytics
            </button>
            <button type="button" className="btn-ghost" onClick={() => choose('deny')}>
              No thanks
            </button>
          </div>
        </div>
      )}
    </>
  )
}
