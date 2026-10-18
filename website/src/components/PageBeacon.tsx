'use client'

import { useEffect } from 'react'
import { usePathname, useSearchParams } from 'next/navigation'
import { trackEvent } from '@/lib/track'

export default function PageBeacon() {
  const pathname = usePathname()
  const search = useSearchParams()

  useEffect(() => {
    trackEvent('pageview', pathname)
  }, [pathname, search])

  return null
}
