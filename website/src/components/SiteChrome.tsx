'use client'

import { usePathname } from 'next/navigation'
import Nav from '@/components/Nav'
import Footer from '@/components/Footer'

export default function SiteChrome({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  if (pathname === '/admin' || pathname.startsWith('/admin/')) {
    return (
      <main id="main" className="min-h-screen">
        {children}
      </main>
    )
  }
  return (
    <>
      <Nav />
      <main id="main" className="min-h-screen">
        {children}
      </main>
      <Footer />
    </>
  )
}
