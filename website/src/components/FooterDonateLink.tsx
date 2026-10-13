'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

export default function FooterDonateLink() {
  const path = usePathname() || ''
  if (path.startsWith('/join') || path.startsWith('/app-guide')) return null
  return (
    <Link
      href="/#donate"
      className="inline-flex min-h-11 items-center text-[0.98rem] text-ink-soft transition-colors duration-200 hover:text-white"
    >
      Donate
    </Link>
  )
}
