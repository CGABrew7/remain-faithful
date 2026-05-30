import Link from 'next/link'

const footerLinks = {
  product: {
    title: 'Product',
    links: [
      { label: 'How It Works', href: '/how-it-works' },
      { label: 'Download', href: '/#download' },
      { label: 'Donate', href: '/#donate' },
    ],
  },
  ministry: {
    title: 'Ministry',
    links: [
      { label: 'Partners', href: '/partners' },
      { label: 'Group Setup', href: '/partners#group-setup' },
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
    <footer className="bg-[#0A1420] border-t border-[#1E3050]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-16 pb-8">
        {/* Top grid */}
        <div className="grid grid-cols-2 md:grid-cols-6 gap-8 pb-12 border-b border-[#1E3050]">
          {/* Logo + tagline */}
          <div className="col-span-2 md:col-span-2">
            <Link href="/" className="flex items-center gap-2.5 mb-4">
              <svg
                width="28"
                height="32"
                viewBox="0 0 32 36"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M16 0L2 6V18C2 26.284 8.268 33.916 16 36C23.732 33.916 30 26.284 30 18V6L16 0Z"
                  fill="url(#footerShieldGrad)"
                />
                <path
                  d="M11 18L14.5 21.5L21 14"
                  stroke="white"
                  strokeWidth="2.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
                <defs>
                  <linearGradient id="footerShieldGrad" x1="0" y1="0" x2="32" y2="36" gradientUnits="userSpaceOnUse">
                    <stop stopColor="#C9A84C" />
                    <stop offset="1" stopColor="#E8C87A" />
                  </linearGradient>
                </defs>
              </svg>
              <span className="font-serif font-semibold text-[#F0EDE8]">Remain Faithful</span>
            </Link>
            <p className="text-sm text-[#8A9BB0] leading-relaxed max-w-[240px]">
              Accountability That Works. Free forever, privacy-first, built for men who are serious about purity.
            </p>
            {/* Social */}
            <div className="flex items-center gap-3 mt-6">
              <a
                href="https://github.com/CGABrew7/remain-faithful"
                target="_blank"
                rel="noopener noreferrer"
                className="w-9 h-9 rounded-lg bg-[#162235] border border-[#1E3050] flex items-center justify-center text-[#8A9BB0] hover:text-[#F0EDE8] hover:border-[#C9A84C] transition-colors"
                aria-label="GitHub"
              >
                <svg width="16" height="16" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/>
                </svg>
              </a>
            </div>
          </div>

          {/* Link columns */}
          {Object.values(footerLinks).map((col) => (
            <div key={col.title} className="col-span-1">
              <h3 className="text-xs font-semibold uppercase tracking-wider text-[#C9A84C] mb-4">
                {col.title}
              </h3>
              <ul className="space-y-2.5">
                {col.links.map((link) => (
                  <li key={link.label}>
                    <Link
                      href={link.href}
                      className="text-sm text-[#8A9BB0] hover:text-[#F0EDE8] transition-colors"
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        {/* Bottom bar */}
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4 pt-8">
          <p className="text-sm text-[#8A9BB0]">
            © 2025 Remain Faithful. Free forever. Open source.
          </p>
          <p className="text-xs text-[#8A9BB0]/60">
            Built with integrity, for the men who need it most.
          </p>
        </div>
      </div>
    </footer>
  )
}
