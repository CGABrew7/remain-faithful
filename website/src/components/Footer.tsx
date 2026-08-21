import Link from 'next/link'
import WaitlistForm from './WaitlistForm'

const footerLinks = {
  product: {
    title: 'Product',
    links: [
      { label: 'How It Works', href: '/how-it-works' },
      { label: 'Privacy Architecture', href: '/privacy-architecture' },
      { label: 'Waitlist', href: '/#waitlist' },
      { label: 'Donate', href: '/#donate' },
    ],
  },
  ministry: {
    title: 'Ministry',
    links: [
      { label: 'Partners', href: '/partners' },
      { label: 'Group Setup', href: '/partners#group-setup' },
      { label: 'Group Setup Guide', href: '/group-setup-guide' },
    ],
  },
  company: {
    title: 'Company',
    links: [
      { label: 'About', href: '/about' },
      { label: 'Blog', href: '/blog' },
      { label: 'Contact', href: '/about#contact' },
    ],
  },
  legal: {
    title: 'Legal',
    links: [
      { label: 'Privacy Policy', href: '/privacy' },
      { label: 'Terms of Service', href: '/terms' },
    ],
  },
}

export default function Footer() {
  return (
    <footer className="border-t border-hairline bg-paper-deep">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 pt-16 pb-8">
        <div className="pb-12 mb-12 border-b border-hairline grid sm:grid-cols-2 gap-6 items-end">
          <div>
            <p className="kicker mb-3">Letter</p>
            <h3 className="font-serif text-xl font-medium text-ink mb-2">
              Join the accountability newsletter
            </h3>
            <p className="text-ink-soft text-[1.05rem] leading-relaxed">
              Monthly encouragement, guides, and updates from Remain Faithful.
            </p>
          </div>
          <WaitlistForm variant="footer" buttonText="Subscribe" />
        </div>

        <div className="grid grid-cols-2 md:grid-cols-6 gap-8 pb-12 border-b border-hairline">
          <div className="col-span-2 md:col-span-2">
            <Link href="/" className="flex items-center gap-2.5 mb-4 min-h-10">
              <svg width="22" height="24" viewBox="0 0 32 36" fill="none" aria-hidden="true">
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
              <span className="font-serif text-ink">Remain Faithful</span>
            </Link>
            <p className="text-ink-soft leading-relaxed max-w-[240px]">
              Free forever, privacy-first, built for believers serious about purity.
            </p>
            <div className="flex items-center gap-3 mt-6">
              <a
                href="https://github.com/CGABrew7/remain-faithful"
                target="_blank"
                rel="noopener noreferrer"
                className="relative w-10 h-10 inline-flex items-center justify-center text-ink-soft hover:text-ink transition-colors duration-200"
                aria-label="GitHub"
              >
                <svg width="16" height="16" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/>
                </svg>
              </a>
            </div>
          </div>

          {Object.values(footerLinks).map((col) => (
            <div key={col.title} className="col-span-1">
              <h3 className="kicker mb-4">{col.title}</h3>
              <ul className="space-y-2.5">
                {col.links.map((link) => (
                  <li key={link.label}>
                    <Link
                      href={link.href}
                      className="text-ink-soft hover:text-ink transition-colors duration-200"
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="flex flex-col sm:flex-row items-center justify-between gap-4 pt-8">
          <p className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint">
            © 2026 Remain Faithful. Free forever. Open source.
          </p>
          <p className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint">
            Fort Wayne
          </p>
        </div>
      </div>
    </footer>
  )
}
