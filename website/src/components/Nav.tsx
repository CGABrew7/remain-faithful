'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

const navLinks = [
  { href: '/how-it-works', label: 'How It Works' },
  { href: '/about', label: 'About' },
  { href: '/partners', label: 'Partners' },
  { href: '/blog', label: 'Blog' },
  { href: '/about#contact', label: 'Contact' },
]

export default function Nav() {
  const [menuOpen, setMenuOpen] = useState(false)
  const [scrolled, setScrolled] = useState(false)
  const pathname = usePathname()

  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 20)
    window.addEventListener('scroll', handleScroll, { passive: true })
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  useEffect(() => {
    setMenuOpen(false)
  }, [pathname])

  useEffect(() => {
    if (menuOpen) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
    return () => {
      document.body.style.overflow = ''
    }
  }, [menuOpen])

  const overDusk = pathname === '/' && !scrolled && !menuOpen
  const inkNav = !overDusk

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 transition-[background-color,box-shadow] duration-300 ${
        inkNav ? 'bg-paper/95 shadow-[0_1px_0_var(--hairline)]' : 'bg-transparent'
      }`}
    >
      <div className="max-w-5xl mx-auto px-4 sm:px-6">
        <div className="flex items-center justify-between h-16">
          <Link href="/" className="flex items-center gap-2.5 group min-h-10">
            <WaxMark />
            <span className={`font-serif text-[1.15rem] tracking-tight transition-colors duration-200 ${inkNav ? 'text-ink' : 'text-paper'}`}>
              Remain Faithful
            </span>
          </Link>

          <nav className="hidden md:flex items-center gap-1">
            {navLinks.map((link) => {
              const current =
                pathname === link.href ||
                (link.href !== '/' &&
                  pathname?.startsWith(link.href.split('#')[0]) &&
                  link.href.split('#')[0] === pathname)
              return (
                <Link
                  key={link.href}
                  href={link.href}
                  className={`relative px-3 py-2 min-h-10 inline-flex items-center font-mono text-[12px] tracking-[0.06em] transition-colors duration-200 ${
                    current
                      ? inkNav
                        ? 'text-wax'
                        : 'text-paper'
                      : inkNav
                        ? 'text-ink-soft hover:text-ink'
                        : 'text-paper/70 hover:text-paper'
                  }`}
                >
                  {link.label}
                </Link>
              )
            })}
          </nav>

          <div className="hidden md:flex items-center">
            <Link
              href="/#waitlist"
              className={inkNav ? 'btn-wax !py-2 !px-4' : 'btn-ghost !py-2 !px-4 !text-paper !shadow-[0_0_0_1px_rgba(255,255,255,0.18)] hover:!bg-white/5'}
            >
              Join the Waitlist
            </Link>
          </div>

          <button
            className={`md:hidden relative w-10 h-10 flex flex-col items-center justify-center gap-1.5 ${
              inkNav ? 'text-ink' : 'text-paper'
            }`}
            onClick={() => setMenuOpen(!menuOpen)}
            aria-label={menuOpen ? 'Close menu' : 'Open menu'}
            aria-expanded={menuOpen}
          >
            <span
              className={`block w-5 h-px bg-current transition-[transform,opacity] duration-300 ${
                menuOpen ? 'rotate-45 translate-y-[3.5px]' : ''
              }`}
            />
            <span
              className={`block w-5 h-px bg-current transition-[transform,opacity] duration-300 ${
                menuOpen ? 'opacity-0 scale-x-0' : ''
              }`}
            />
            <span
              className={`block w-5 h-px bg-current transition-[transform,opacity] duration-300 ${
                menuOpen ? '-rotate-45 -translate-y-[3.5px]' : ''
              }`}
            />
          </button>
        </div>
      </div>

      <div
        className={`md:hidden overflow-hidden transition-[max-height,opacity] duration-300 ease-in-out ${
          menuOpen ? 'max-h-[500px] opacity-100' : 'max-h-0 opacity-0'
        }`}
      >
        <div className="bg-paper px-4 pt-2 pb-6 shadow-[0_1px_0_var(--hairline)]">
          <nav className="flex flex-col">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={`px-2 py-3 min-h-10 font-mono text-[13px] tracking-[0.06em] ${
                  pathname === link.href ? 'text-wax' : 'text-ink-soft hover:text-ink'
                }`}
              >
                {link.label}
              </Link>
            ))}
            <div className="mt-3 pt-4 border-t border-hairline">
              <Link href="/#waitlist" className="btn-wax w-full">
                Join the Waitlist
              </Link>
            </div>
          </nav>
        </div>
      </div>
    </header>
  )
}

function WaxMark() {
  return (
    <svg
      width="22"
      height="24"
      viewBox="0 0 32 36"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden="true"
    >
      <path
        d="M16 0L2 6V18C2 26.284 8.268 33.916 16 36C23.732 33.916 30 26.284 30 18V6L16 0Z"
        fill="#7a2c28"
      />
      <path
        d="M11 18L14.5 21.5L21 14"
        stroke="#ece6d8"
        strokeWidth="2.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}
