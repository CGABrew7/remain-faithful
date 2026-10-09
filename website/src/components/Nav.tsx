'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import ShieldMark from './ShieldMark'

const navLinks = [
  { href: '/how-it-works', label: 'How It Works' },
  { href: '/app-guide', label: 'App Guide' },
  { href: '/about', label: 'About' },
  { href: '/partners', label: 'Partners' },
  { href: '/blog', label: 'Blog' },
  { href: '/about#contact', label: 'Contact' },
]

export default function Nav() {
  const [menuOpen, setMenuOpen] = useState(false)
  const pathname = usePathname()

  useEffect(() => {
    setMenuOpen(false)
  }, [pathname])

  useEffect(() => {
    document.body.style.overflow = menuOpen ? 'hidden' : ''
    return () => {
      document.body.style.overflow = ''
    }
  }, [menuOpen])

  return (
    <header className="site-nav fixed top-0 left-0 right-0 z-50 border-b border-white/10 bg-[#0C1D4C]/75 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-5 sm:px-6">
        <Link href="/" className="flex min-h-11 items-center gap-2.5">
          <ShieldMark />
          <span className="font-serif text-[1.15rem] font-semibold tracking-tight text-white">
            Remain Faithful
          </span>
        </Link>

        <nav className="hidden items-center gap-1 md:flex" aria-label="Primary">
          {navLinks.map((link) => {
            const base = link.href.split('#')[0]
            const current = !link.href.includes('#') && pathname === base
            return (
              <Link
                key={link.href}
                href={link.href}
                aria-current={current ? 'page' : undefined}
                className={`inline-flex min-h-11 items-center px-3 text-sm transition-colors duration-200 ${
                  current ? 'text-wax' : 'text-ink-soft hover:text-white'
                }`}
              >
                {link.label}
              </Link>
            )
          })}
        </nav>

        <div className="hidden md:flex">
          <Link href="/#download" className="btn-wax !px-4 !py-2">
            Get started
          </Link>
        </div>

        <button
          className="relative flex h-11 w-11 flex-col items-center justify-center gap-1.5 text-white md:hidden"
          onClick={() => setMenuOpen(!menuOpen)}
          aria-label={menuOpen ? 'Close menu' : 'Open menu'}
          aria-expanded={menuOpen}
        >
          <span
            className={`block h-px w-5 bg-current transition-[transform,opacity] duration-300 ${
              menuOpen ? 'translate-y-[3.5px] rotate-45' : ''
            }`}
          />
          <span
            className={`block h-px w-5 bg-current transition-[transform,opacity] duration-300 ${
              menuOpen ? 'scale-x-0 opacity-0' : ''
            }`}
          />
          <span
            className={`block h-px w-5 bg-current transition-[transform,opacity] duration-300 ${
              menuOpen ? '-translate-y-[3.5px] -rotate-45' : ''
            }`}
          />
        </button>
      </div>

      <div
        className={`overflow-hidden border-t border-white/10 bg-[#0C1D4C]/95 backdrop-blur-md transition-[max-height,opacity] duration-300 ease-in-out md:hidden ${
          menuOpen ? 'max-h-[28rem] opacity-100' : 'max-h-0 opacity-0 border-t-0'
        }`}
      >
        <nav className="flex flex-col px-4 pb-6 pt-2" aria-label="Mobile">
          {navLinks.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              className={`inline-flex min-h-11 items-center px-2 text-[15px] ${
                pathname === link.href ? 'text-wax' : 'text-ink-soft'
              }`}
            >
              {link.label}
            </Link>
          ))}
          <div className="mt-3 border-t border-white/10 pt-4">
            <Link href="/#download" className="btn-wax w-full">
              Get started
            </Link>
          </div>
        </nav>
      </div>
    </header>
  )
}
