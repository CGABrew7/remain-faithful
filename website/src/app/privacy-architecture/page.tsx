import type { Metadata } from 'next'
import Link from 'next/link'

export const metadata: Metadata = {
  title: 'Privacy Architecture',
  description:
    'A detailed technical explanation of how Remain Faithful protects your privacy through on-device AI, encrypted data flows, and open-source transparency.',
}

const GITHUB_URL = 'https://github.com/CGABrew7/remain-faithful'

export default function PrivacyArchitecturePage() {
  const jsonLd = {
    '@context': 'https://schema.org',
    '@type': 'TechArticle',
    headline: 'How Your Privacy is Protected — Remain Faithful Privacy Architecture',
    description:
      'A detailed look at how Remain Faithful processes content without compromising your privacy.',
    author: { '@type': 'Person', name: 'Jeff Brewer' },
    publisher: {
      '@type': 'Organization',
      name: 'Remain Faithful',
      url: 'https://remainfaithful.com',
    },
  }

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />

      {/* Hero */}
      <section className="pt-32 pb-20 border-b border-[#1E3050]">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-[#C9A84C]/10 border border-[#C9A84C]/20 mb-6">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2" strokeLinecap="round">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
          </div>
          <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-4">Privacy Architecture</p>
          <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-6">
            How Your Privacy is Protected
          </h1>
          <p className="text-[#8A9BB0] text-lg max-w-2xl mx-auto">
            A detailed look at how Remain Faithful processes content without compromising your privacy — from device to partner notification.
          </p>
        </div>
      </section>

      {/* Section 1: Three-Tier Pipeline */}
      <section className="py-20">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">
              The Three-Tier Classification Pipeline
            </h2>
            <p className="text-[#8A9BB0] max-w-2xl mx-auto">
              Every screen frame passes through a tiered pipeline designed to keep your content on your device in the vast majority of cases.
            </p>
          </div>

          {/* Pipeline diagram */}
          <div className="relative">
            {/* Vertical connector line (desktop) */}
            <div className="hidden md:block absolute left-1/2 top-0 bottom-0 w-px bg-gradient-to-b from-[#C9A84C]/40 via-[#C9A84C]/20 to-[#C9A84C]/40 -translate-x-1/2" />

            <div className="space-y-6">
              {/* Device */}
              <PipelineStep
                side="center"
                badge="Start"
                badgeColor="#162235"
                title="Your Device Captures a Screen Frame"
                desc="Apple's ReplayKit creates a sandboxed broadcast extension process. This process cannot make network requests — it is architecturally isolated from the internet."
                icon={
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2" strokeLinecap="round">
                    <rect x="5" y="2" width="14" height="20" rx="2"/>
                    <line x1="12" y1="18" x2="12.01" y2="18"/>
                  </svg>
                }
              />

              {/* Tier 1 */}
              <PipelineStep
                side="center"
                badge="Tier 1 · 70% of cases"
                badgeColor="#C9A84C"
                title="Rules: URL Blocklist + Keyword Matching"
                desc="Known adult domains are checked against a local blocklist. Visible text is pattern-matched against regex rules. Fast, deterministic, 100% on-device. No AI required."
                icon={
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2" strokeLinecap="round">
                    <polyline points="9 11 12 14 22 4"/>
                    <path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/>
                  </svg>
                }
              />

              {/* Tier 2 */}
              <PipelineStep
                side="center"
                badge="Tier 2 · 25% of cases"
                badgeColor="#C9A84C"
                title="On-Device AI: Apple SensitiveContentAnalysis + Text Classifier"
                desc="Apple Vision OCR extracts text; SensitiveContentAnalysis detects explicit imagery. Both run on the device's Neural Engine — the dedicated AI chip in modern iPhones. No server involved."
                icon={
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2" strokeLinecap="round">
                    <circle cx="12" cy="12" r="3"/>
                    <path d="M12 1v4M12 19v4M4.22 4.22l2.83 2.83M16.95 16.95l2.83 2.83M1 12h4M19 12h4M4.22 19.78l2.83-2.83M16.95 7.05l2.83-2.83"/>
                  </svg>
                }
              />

              {/* Tier 3 */}
              <PipelineStep
                side="center"
                badge="Tier 3 · <5% of cases"
                badgeColor="#8A9BB0"
                title="Rare Cloud Fallback: Text-Only Category Query"
                desc="Only when both Tier 1 and Tier 2 are uncertain, an anonymized category query is sent to our secure classification server. This query contains no screenshots, no URLs, no personal information — only the anonymized text category."
                icon={
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2" strokeLinecap="round">
                    <path d="M18 10h-1.26A8 8 0 109 20h9a5 5 0 000-10z"/>
                  </svg>
                }
              />

              {/* Alert */}
              <PipelineStep
                side="center"
                badge="Result"
                badgeColor="#162235"
                title="Discreet Alert Delivered to Partners"
                desc="Partners receive: category label (e.g., 'Adult Content'), severity level (Low/Medium/High), and timestamp. Never a screenshot. Never your browsing history. Never raw content."
                icon={
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2" strokeLinecap="round">
                    <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9M13.73 21a2 2 0 01-3.46 0"/>
                  </svg>
                }
              />
            </div>
          </div>

          <div
            className="mt-10 p-5 rounded-2xl border border-[#C9A84C]/20 text-sm text-[#8A9BB0] text-center"
            style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}
          >
            In the rare Tier 3 case, the server receives only an anonymized text category query — never screenshots, never identifying information. Full encryption in transit via TLS 1.3.
          </div>
        </div>
      </section>

      {/* Section 2: What We Can/Cannot See */}
      <section className="py-20 bg-[#0A1420] border-y border-[#1E3050]">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">
              What We Can See vs. What We Cannot See
            </h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
              The architecture enforces these limits, not just our policies.
            </p>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr>
                  <th className="text-left pb-4 text-xs uppercase tracking-widest text-[#C9A84C] font-semibold w-1/3">Data type</th>
                  <th className="text-center pb-4 text-xs uppercase tracking-widest text-[#C9A84C] font-semibold w-1/3">Remain Faithful server</th>
                  <th className="text-center pb-4 text-xs uppercase tracking-widest text-[#C9A84C] font-semibold w-1/3">Your partners</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#1E3050]">
                {[
                  ['Screenshots / screen frames', '✗ Never', '✗ Never'],
                  ['Raw screen content or text', '✗ Never', '✗ Never'],
                  ['Browsing history or URLs', '✗ Never', '✗ Never'],
                  ['App usage details', '✗ Never', '✗ Never'],
                  ['Passwords or financial data', '✗ Never', '✗ Never'],
                  ['Message content', '✗ Never', '✗ Never'],
                  ['Photos and videos', '✗ Never', '✗ Never'],
                  ['Alert category (e.g. "Adult Content")', '✓ Encrypted metadata', '✓ Yes'],
                  ['Severity level (Low / Medium / High)', '✓ Encrypted metadata', '✓ Yes'],
                  ['Timestamp', '✓ Encrypted metadata', '✓ Yes'],
                  ['System-generated description', '✓ Encrypted metadata', '✓ Yes'],
                  ['Your name and email (account info)', '✓ Encrypted at rest', '✗ No'],
                ].map(([item, server, partners]) => (
                  <tr key={item} className="hover:bg-[#162235]/50 transition-colors">
                    <td className="py-3.5 text-[#F0EDE8]">{item}</td>
                    <td className={`py-3.5 text-center font-medium ${server.startsWith('✗') ? 'text-red-400' : 'text-green-400'}`}>{server}</td>
                    <td className={`py-3.5 text-center font-medium ${partners.startsWith('✗') ? 'text-red-400' : 'text-green-400'}`}>{partners}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </section>

      {/* Section 3: Data Flow */}
      <section className="py-20">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">Data Flow Diagram</h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
              How a flagged event travels from your device to your partner&apos;s notification — with encryption at every step.
            </p>
          </div>

          <div className="grid md:grid-cols-5 gap-3 items-center">
            {[
              { label: 'Your Device', sub: 'Screen frame classified locally', icon: '📱' },
              { label: 'Alert Metadata', sub: 'Category + severity only', icon: '📋', connector: true },
              { label: 'RF Server', sub: 'Encrypted at rest (AES-256)', icon: '🔐', connector: true },
              { label: 'APNs', sub: 'Apple Push (TLS 1.3)', icon: '📡', connector: true },
              { label: "Partner's Device", sub: 'Notification received', icon: '🔔', connector: true },
            ].map((node, i) => (
              <div key={i} className="flex md:flex-col items-center gap-3">
                {node.connector && (
                  <div className="md:hidden flex items-center justify-center text-[#C9A84C]">→</div>
                )}
                {node.connector && (
                  <div className="hidden md:block w-full h-px bg-gradient-to-r from-[#C9A84C]/40 to-[#C9A84C]/20 -ml-3 mb-2" />
                )}
                <div
                  className="w-full rounded-2xl p-5 border border-[#1E3050] bg-[#162235] text-center"
                >
                  <div className="text-3xl mb-2">{node.icon}</div>
                  <p className="font-semibold text-[#F0EDE8] text-sm mb-1">{node.label}</p>
                  <p className="text-xs text-[#8A9BB0]">{node.sub}</p>
                </div>
              </div>
            ))}
          </div>

          <p className="text-center text-sm text-[#8A9BB0] mt-8">
            All communication between the app and server uses TLS 1.3. Data at rest is AES-256 encrypted. The ReplayKit broadcast extension is sandboxed and cannot make any network requests directly.
          </p>
        </div>
      </section>

      {/* Section 4: Threat Model */}
      <section className="py-20 bg-[#0A1420] border-y border-[#1E3050]">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">Threat Model</h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
              What happens in the worst-case scenarios? We&apos;ve thought through them.
            </p>
          </div>

          <div className="grid md:grid-cols-2 gap-6">
            {[
              {
                threat: 'What if your servers are hacked?',
                answer: 'We do not store screenshots or browsing content. The database contains only encrypted alert metadata (category, severity, timestamp) and account information (name, email, bcrypt-hashed password). A breach would expose metadata, not your screen content.',
              },
              {
                threat: 'What if data is intercepted in transit?',
                answer: 'All communication between the app, server, and Apple Push Notification Service uses TLS 1.3 with certificate pinning. Interception would yield only encrypted ciphertext with no practical path to decryption.',
              },
              {
                threat: 'What if a partner is malicious?',
                answer: 'Partners only see alert categories and timestamps — never raw content, screenshots, or browsing history. A malicious partner has nothing to expose. You can remove a partner instantly at any time.',
              },
              {
                threat: 'What if the app itself is compromised?',
                answer: 'The entire codebase is open source and auditable by anyone. We run pre-commit secret scanning on every contribution. The ReplayKit sandbox architecture means the broadcast extension physically cannot exfiltrate screen content over the network.',
              },
            ].map((t) => (
              <div key={t.threat} className="rounded-2xl p-7 border border-[#1E3050] bg-[#162235]">
                <div className="flex items-start gap-3 mb-3">
                  <div className="w-6 h-6 rounded-full bg-[#C9A84C]/20 flex items-center justify-center flex-shrink-0 mt-0.5">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2.5" strokeLinecap="round">
                      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                    </svg>
                  </div>
                  <h3 className="font-semibold text-[#F0EDE8] leading-snug">{t.threat}</h3>
                </div>
                <p className="text-sm text-[#8A9BB0] leading-relaxed">{t.answer}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Section 5: Open Source Commitment */}
      <section className="py-20">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-5">Open Source Commitment</h2>
              <p className="text-[#8A9BB0] leading-relaxed mb-4">
                The entire Remain Faithful codebase — iOS app, Go backend, and this website — is publicly available on GitHub. This is not optional for an app that handles sensitive behavioral data.
              </p>
              <p className="text-[#8A9BB0] leading-relaxed mb-4">
                Our privacy architecture is not a policy claim. It is verifiable in the code. Anyone can confirm that the broadcast extension cannot make network requests, that classification happens on-device, and that partner alerts contain only metadata.
              </p>
              <p className="text-[#8A9BB0] leading-relaxed mb-6">
                Security researchers and privacy advocates are invited to review, test, and report findings. We take responsible disclosure seriously.
              </p>
              <a
                href={GITHUB_URL}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 px-6 py-3 rounded-xl border border-[#1E3050] bg-[#162235] text-[#F0EDE8] text-sm font-semibold hover:border-[#C9A84C]/50 transition-colors"
              >
                <svg width="18" height="18" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/>
                </svg>
                View Source on GitHub
              </a>
            </div>
            <div
              className="rounded-2xl p-7 border border-[#C9A84C]/20"
              style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}
            >
              <h3 className="font-semibold text-[#F0EDE8] mb-4">Why open source matters for trust</h3>
              <ul className="space-y-3">
                {[
                  'Anyone can verify our privacy claims by reading the code',
                  'Security researchers can find and report vulnerabilities',
                  'The community can audit every update before it ships',
                  'No "trust us" black boxes when handling intimate behavioral data',
                  'Pre-commit secret scanning prevents credential leaks',
                ].map((item) => (
                  <li key={item} className="flex items-start gap-2 text-sm text-[#8A9BB0]">
                    <svg className="flex-shrink-0 mt-0.5" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2.5" strokeLinecap="round">
                      <polyline points="20 6 9 17 4 12"/>
                    </svg>
                    {item}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      </section>

      {/* Section 6: Competitor Comparison */}
      <section className="py-20 bg-[#0A1420] border-y border-[#1E3050]">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">
              How We Compare to Other Tools
            </h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
              Privacy dimensions compared across the most common accountability apps.
            </p>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr>
                  <th className="text-left pb-4 text-xs uppercase tracking-widest text-[#C9A84C] font-semibold">Privacy Dimension</th>
                  <th className="text-center pb-4 text-xs uppercase tracking-widest text-[#C9A84C] font-semibold">Remain Faithful</th>
                  <th className="text-center pb-4 text-xs uppercase tracking-widest text-[#8A9BB0] font-semibold">Covenant Eyes</th>
                  <th className="text-center pb-4 text-xs uppercase tracking-widest text-[#8A9BB0] font-semibold">Ever Accountable</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#1E3050]">
                {[
                  ['On-device AI processing', '✓ Yes', '✗ No (cloud)', '✗ No (cloud)'],
                  ['Open source codebase', '✓ Yes', '✗ No', '✗ No'],
                  ['Screenshots stored on server', '✗ Never', '✓ Yes', '✓ Yes'],
                  ['Partners see raw content', '✗ Never', '✓ Yes', '✓ Yes'],
                  ['Cloud dependency for classification', '< 5% of events', 'Always', 'Always'],
                  ['Cost', '100% Free', 'Paid subscription', 'Paid subscription'],
                  ['Auditable by security researchers', '✓ Yes', '✗ No', '✗ No'],
                ].map(([dim, rf, ce, ea]) => (
                  <tr key={dim} className="hover:bg-[#162235]/50 transition-colors">
                    <td className="py-3.5 text-[#F0EDE8]">{dim}</td>
                    <td className="py-3.5 text-center font-medium text-green-400">{rf}</td>
                    <td className="py-3.5 text-center font-medium text-[#8A9BB0]">{ce}</td>
                    <td className="py-3.5 text-center font-medium text-[#8A9BB0]">{ea}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <p className="text-xs text-[#8A9BB0]/60 mt-6 text-center">
            Competitor information based on publicly available documentation. All claims are verifiable via our open-source codebase.
          </p>
        </div>
      </section>

      {/* CTA */}
      <section className="py-20">
        <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">Questions About Our Privacy Model?</h2>
          <p className="text-[#8A9BB0] mb-8">
            Read the source code, open a GitHub issue, or contact us directly. Transparency is not just a commitment — it is a practice.
          </p>
          <div className="flex flex-wrap justify-center gap-4">
            <a
              href={GITHUB_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-6 py-3 rounded-full border border-[#1E3050] bg-[#162235] text-[#F0EDE8] text-sm font-semibold hover:border-[#C9A84C]/50 transition-colors"
            >
              View on GitHub
            </a>
            <Link
              href="/about#contact"
              className="inline-flex items-center gap-2 px-6 py-3 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all duration-200 text-sm"
            >
              Contact Us
            </Link>
          </div>
        </div>
      </section>
    </>
  )
}

function PipelineStep({
  badge,
  badgeColor,
  title,
  desc,
  icon,
}: {
  side?: string
  badge: string
  badgeColor: string
  title: string
  desc: string
  icon: React.ReactNode
}) {
  return (
    <div className="relative flex justify-center">
      <div
        className="w-full max-w-2xl rounded-2xl p-6 border border-[#1E3050] bg-[#162235]"
      >
        <div className="flex items-start gap-4">
          <div className="w-10 h-10 rounded-xl bg-[#C9A84C]/10 border border-[#C9A84C]/20 flex items-center justify-center flex-shrink-0">
            {icon}
          </div>
          <div className="flex-1">
            <div className="flex items-center gap-3 mb-2 flex-wrap">
              <span
                className="text-xs font-bold uppercase tracking-wider px-2.5 py-1 rounded-full"
                style={{
                  background: badgeColor === '#C9A84C' ? 'rgba(201,168,76,0.15)' : 'rgba(138,155,176,0.15)',
                  color: badgeColor === '#C9A84C' ? '#C9A84C' : '#8A9BB0',
                }}
              >
                {badge}
              </span>
            </div>
            <h3 className="font-semibold text-[#F0EDE8] mb-2">{title}</h3>
            <p className="text-sm text-[#8A9BB0] leading-relaxed">{desc}</p>
          </div>
        </div>
      </div>
    </div>
  )
}
