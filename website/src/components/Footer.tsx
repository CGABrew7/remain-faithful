import Link from 'next/link'
import ShieldMark from './ShieldMark'
import WaitlistForm from './WaitlistForm'

const columns = [
  {
    title: 'Product',
    links: [
      { label: 'How It Works', href: '/how-it-works' },
      { label: 'App Guide', href: '/app-guide' },
      { label: 'Privacy Architecture', href: '/privacy-architecture' },
      { label: 'Updates', href: '/#waitlist' },
      { label: 'Donate', href: '/#donate' },
    ],
  },
  {
    title: 'Ministry',
    links: [
      { label: 'Partners', href: '/partners' },
      { label: 'Group Setup', href: '/partners#group-setup' },
      { label: 'Group Setup Guide', href: '/group-setup-guide' },
    ],
  },
  {
    title: 'Company',
    links: [
      { label: 'About', href: '/about' },
      { label: 'Blog', href: '/blog' },
      { label: 'Contact', href: '/about#contact' },
      { label: 'Suggest an Idea', href: '/feedback' },
      { label: 'Privacy Policy', href: '/privacy' },
      { label: 'Terms of Service', href: '/terms' },
    ],
  },
]

export default function Footer() {
  return (
    <footer className="border-t border-white/10 bg-[#0C1D4C]">
      <div className="mx-auto max-w-6xl px-5 pb-10 pt-16 sm:px-6">
        <div className="grid gap-12 border-b border-white/10 pb-12 md:grid-cols-[minmax(0,1.2fr)_minmax(0,2fr)]">
          <div>
            <Link href="/" className="inline-flex min-h-11 items-center gap-2.5">
              <ShieldMark />
              <span className="font-serif text-lg font-semibold text-white">Remain Faithful</span>
            </Link>
            <p className="mt-4 max-w-xs text-[1.02rem] leading-relaxed text-ink-soft">
              Free forever, privacy-first, built for believers serious about purity.
            </p>
            <div className="mt-6 flex flex-col items-start gap-4">
              <Link href="/#download" className="btn-wax">
                Get started
              </Link>
              <a
                href="https://github.com/CGABrew7/remain-faithful"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex min-h-11 items-center gap-2 text-sm text-ink-faint transition-colors duration-200 hover:text-white"
              >
              <svg width="16" height="16" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/>
              </svg>
                GitHub
              </a>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-8 sm:grid-cols-3">
            {columns.map((col) => (
              <div key={col.title}>
                <h2 className="kicker mb-4">{col.title}</h2>
                <ul className="space-y-1">
                  {col.links.map((link) => (
                    <li key={link.label}>
                      <Link
                        href={link.href}
                        className="inline-flex min-h-11 items-center text-[0.98rem] text-ink-soft transition-colors duration-200 hover:text-white"
                      >
                        {link.label}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </div>

        <div className="grid items-end gap-6 border-b border-white/10 py-10 sm:grid-cols-2">
          <div>
            <h2 className="font-serif text-xl font-medium text-white">
              Join the accountability newsletter
            </h2>
            <p className="mt-2 max-w-md text-ink-soft">
              Monthly encouragement, guides, and updates from Remain Faithful.
            </p>
          </div>
          <WaitlistForm variant="footer" buttonText="Subscribe" />
        </div>

        <div className="flex flex-col gap-2 pt-8 text-sm text-ink-faint">
          <p>© 2026 Remain Faithful — Free forever. Open source.</p>
          <p>Woodfield Foundation Inc. is a 501(c)(3) private foundation, EIN 39-2184435. All donations are tax-deductible.</p>
        </div>
      </div>
    </footer>
  )
}
