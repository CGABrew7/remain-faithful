'use client'

import { trackEvent } from '@/lib/track'

export const APP_STORE_URL = 'https://apps.apple.com/app/id6774615426'

export default function AppStoreLink({
  className = 'btn-wax',
  children = 'Download on the App Store',
}: {
  className?: string
  children?: React.ReactNode
}) {
  return (
    <a
      href={APP_STORE_URL}
      className={className}
      onClick={() => trackEvent('app_store_click')}
    >
      {children}
    </a>
  )
}
