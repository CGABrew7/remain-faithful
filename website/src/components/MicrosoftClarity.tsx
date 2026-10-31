'use client'

import { useEffect, useState } from 'react'
import { usePathname } from 'next/navigation'
import Script from 'next/script'
import { ANALYTICS_CONSENT_EVENT, readAnalyticsConsent } from '@/lib/analytics-consent'
import { shouldLoadClarity } from '@/lib/clarity-paths'

const CLARITY_PROJECT_ID = 'z6e76ynt8r'

type ClarityWindow = Window & {
  clarity?: (...args: unknown[]) => void
}

export default function MicrosoftClarity() {
  const pathname = usePathname() || ''
  const [consent, setConsent] = useState<string | null>(null)
  const enabled = shouldLoadClarity(pathname, consent)

  useEffect(() => {
    function sync() {
      setConsent(readAnalyticsConsent(window.localStorage))
    }
    sync()
    window.addEventListener('storage', sync)
    window.addEventListener(ANALYTICS_CONSENT_EVENT, sync)
    return () => {
      window.removeEventListener('storage', sync)
      window.removeEventListener(ANALYTICS_CONSENT_EVENT, sync)
    }
  }, [])

  useEffect(() => {
    const clarity = (window as ClarityWindow).clarity
    if (typeof clarity !== 'function') return
    clarity(enabled ? 'start' : 'stop')
  }, [enabled])

  if (!enabled) return null

  return (
    <Script id="microsoft-clarity" strategy="afterInteractive">
      {`(function(c,l,a,r,i,t,y){c[a]=c[a]||function(){(c[a].q=c[a].q||[]).push(arguments)};t=l.createElement(r);t.async=1;t.src="https://www.clarity.ms/tag/"+i;y=l.getElementsByTagName(r)[0];y.parentNode.insertBefore(t,y);})(window, document, "clarity", "script", "${CLARITY_PROJECT_ID}");`}
    </Script>
  )
}
