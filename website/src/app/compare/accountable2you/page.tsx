import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { comparisonSchema } from '@/lib/structured-data'
import { Breadcrumbs } from '@/components/Breadcrumbs'

export const metadata: Metadata = {
  title: 'Remain Faithful vs Accountable2You: Honest Comparison (2026)',
  description: 'A detailed, honest comparison of Remain Faithful and Accountable2You for Christian accountability. See how they differ on price, privacy, monitoring approach, and battery impact.',
  alternates: { canonical: 'https://remainfaithful.com/compare/accountable2you' },
}

const pageFaqSchema = {
  '@context': 'https://schema.org',
  '@type': 'FAQPage',
  mainEntity: [
    {
      '@type': 'Question',
      name: 'Does Accountable2You drain the battery?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Accountable2You uses a VPN-based monitoring approach, which routes device traffic through a local VPN. This approach can increase battery drain and occasionally causes conflicts with corporate or school VPN configurations. Remain Faithful uses Apple\'s Screen Time framework for always-on monitoring, which has minimal battery impact.',
      },
    },
    {
      '@type': 'Question',
      name: 'Can I switch from Accountable2You to Remain Faithful?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes. Cancel your Accountable2You subscription, remove the VPN profile from your device, and download Remain Faithful. Note that Remain Faithful is currently iOS only; if your partners use Android or Windows, Accountable2You has broader platform support.',
      },
    },
    {
      '@type': 'Question',
      name: 'Does Remain Faithful log specific web page titles like Accountable2You?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'No. Remain Faithful does not log or transmit web page titles or browsing history. Partners receive only a category label and severity level when something is flagged. Accountable2You logs specific page titles in its reports, which can expose more detail than some users want their partners to see.',
      },
    },
  ],
}

export default function Accountable2YouCompare() {
  return (
    <>
      <JsonLd data={comparisonSchema('Accountable2You')} />
      <JsonLd data={pageFaqSchema} />

      <div className="pt-24 pb-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <Breadcrumbs items={[
            { name: 'Compare', url: 'https://remainfaithful.com/compare/accountable2you' },
            { name: 'Accountable2You', url: 'https://remainfaithful.com/compare/accountable2you' },
          ]} />

          <div className="mb-10">
            <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-4">Honest Comparison</p>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-6 leading-tight">
              Remain Faithful vs Accountable2You
            </h1>
            <p className="text-[#8A9BB0] text-lg leading-relaxed mb-4">
              Accountable2You takes a VPN-based approach to monitoring, which is meaningfully different from Remain Faithful&apos;s on-device AI model. Here is what that difference means in practice.
            </p>
            <div className="p-5 rounded-xl border border-[#C9A84C]/20 bg-[#C9A84C]/5 text-sm text-[#8A9BB0] leading-relaxed">
              <strong className="text-[#C9A84C]">Disclosure:</strong> We built Remain Faithful, so we are biased. We will be transparent about where Accountable2You is stronger.
            </div>
          </div>

          {/* Comparison Table */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-6">Side-by-Side Comparison</h2>
            <p className="text-xs text-[#8A9BB0]/60 mb-4">Pricing as of June 2026. Check each provider&apos;s website for current pricing.</p>
            <div className="overflow-x-auto rounded-2xl border border-[#1E3050]">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-[#1E3050]">
                    <th className="text-left p-4 text-[#8A9BB0] font-semibold bg-[#0A1420] w-1/3">Feature</th>
                    <th className="text-center p-4 text-[#C9A84C] font-semibold bg-[#0A1420]">Remain Faithful</th>
                    <th className="text-center p-4 text-[#F0EDE8] font-semibold bg-[#0A1420]">Accountable2You</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    ['Price', 'Free forever', '~$80/year'],
                    ['Platform', 'iOS 17+ (Android Fall 2026)', 'iOS, Android, Windows, Mac, Kindle'],
                    ['Monitoring technology', 'On-device AI + Screen Time framework', 'VPN-based traffic monitoring'],
                    ['Web page title logging', 'No. Category only.', 'Yes. Specific page titles logged.'],
                    ['Battery impact', 'Minimal', 'Higher (VPN runs continuously)'],
                    ['VPN required', 'No', 'Yes'],
                    ['Where AI runs', 'Entirely on your device', 'Traffic routed through local VPN'],
                    ['Open source', 'Yes, full codebase on GitHub', 'No'],
                    ['Partner data shared', 'Category + severity label only', 'Detailed activity reports with titles'],
                    ['Accountability model', 'Covenant-based partnership', 'Report-based accountability'],
                    ['Group mode', 'Yes, up to 12 members', 'Yes'],
                    ['DRM streaming monitoring', 'No (iOS limitation)', 'No'],
                  ].map(([feature, rf, a2y], i) => (
                    <tr key={i} className={`border-b border-[#1E3050] ${i % 2 === 0 ? 'bg-[#162235]' : 'bg-[#0F1B2D]'}`}>
                      <td className="p-4 text-[#8A9BB0] font-medium">{feature}</td>
                      <td className="p-4 text-center text-[#F0EDE8]">{rf}</td>
                      <td className="p-4 text-center text-[#8A9BB0]">{a2y}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          {/* Where A2Y is stronger */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Where Accountable2You Is Stronger</h2>
            <div className="space-y-4">
              {[
                { title: 'Cross-platform coverage', body: 'Accountable2You runs on iOS, Android, Windows, Mac, and Kindle. Remain Faithful is currently iOS only. For families with mixed devices, A2Y provides unified coverage across all of them.' },
                { title: 'Detailed web activity reporting', body: 'Accountable2You logs specific web page titles and generates detailed activity reports. If your accountability partner or pastor wants more granular visibility into browsing behavior, A2Y provides more detail than Remain Faithful\'s category-level alerts.' },
                { title: 'Established track record', body: 'Accountable2You has been in operation for years with a larger installed base. Remain Faithful is in beta. For users who want a proven, stable product, A2Y has an advantage in maturity.' },
              ].map((item) => (
                <div key={item.title} className="p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
                  <h3 className="font-semibold text-[#F0EDE8] mb-2">{item.title}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.body}</p>
                </div>
              ))}
            </div>
          </section>

          {/* Where RF is stronger */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Where Remain Faithful Is Stronger</h2>
            <div className="space-y-4">
              {[
                { title: 'Price', body: 'Remain Faithful is free. Accountable2You costs approximately $80 per year. That savings compounds: $400 over five years, $800 over ten. For individuals or churches deploying accountability tools at scale, free is a meaningful difference.' },
                { title: 'No VPN required: better battery life', body: 'Accountable2You routes traffic through a local VPN to monitor it. This has two downsides: it drains battery more quickly, and it can conflict with corporate or school VPN configurations. Remain Faithful uses Apple\'s native Screen Time framework and on-device AI, with minimal battery impact and no VPN conflicts.' },
                { title: 'On-device AI is more thorough', body: 'VPN-based monitoring cannot see inside encrypted HTTPS traffic or monitor what happens inside apps. Remain Faithful\'s on-device approach can classify screen content across all apps because the classification happens at the display level, not the network level.' },
                { title: 'Partner privacy: category only, never page titles', body: 'Accountable2You shares specific web page titles with accountability partners, which means partners see exactly what pages were visited. Some find this level of detail helpful; others consider it overly exposing. Remain Faithful shares only the category and severity, protecting the privacy of the person being held accountable while still triggering a conversation.' },
                { title: 'Open source', body: 'Remain Faithful\'s entire codebase is publicly available for inspection. Accountable2You is proprietary. Anyone can verify what Remain Faithful does and does not transmit.' },
              ].map((item) => (
                <div key={item.title} className="p-5 rounded-xl border border-[#C9A84C]/20 bg-[#162235]">
                  <h3 className="font-semibold text-[#C9A84C] mb-2">{item.title}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.body}</p>
                </div>
              ))}
            </div>
          </section>

          {/* Who should choose what */}
          <section className="mb-14">
            <div className="grid md:grid-cols-2 gap-6">
              <div className="p-6 rounded-2xl border border-[#1E3050] bg-[#162235]">
                <h2 className="font-serif text-lg font-bold text-[#F0EDE8] mb-4">Who Should Choose Accountable2You</h2>
                <ul className="space-y-2">
                  {[
                    'Need Android, Windows, Mac, or Kindle coverage',
                    'Partners want detailed web page title reports',
                    'Happy to pay ~$80/year for a proven solution',
                    'No concerns about VPN battery drain or conflicts',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2 text-sm text-[#8A9BB0]">
                      <svg className="flex-shrink-0 mt-0.5 text-[#8A9BB0]" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><polyline points="20 6 9 17 4 12"/></svg>
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
              <div className="p-6 rounded-2xl border border-[#C9A84C]/20 bg-[#162235]">
                <h2 className="font-serif text-lg font-bold text-[#F0EDE8] mb-4">Who Should Choose Remain Faithful</h2>
                <ul className="space-y-2">
                  {[
                    'iPhone users who want on-device AI at zero cost',
                    'Anyone on a corporate or school VPN',
                    'Those who want battery-efficient monitoring',
                    'People who prefer category alerts over page title logs',
                    'Those who value open-source verification',
                    'Churches needing free group accountability tools',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2 text-sm text-[#8A9BB0]">
                      <svg className="flex-shrink-0 mt-0.5 text-[#C9A84C]" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><polyline points="20 6 9 17 4 12"/></svg>
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </section>

          {/* FAQ */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-6">Frequently Asked Questions</h2>
            <div className="space-y-4">
              {[
                { q: 'Does Accountable2You drain the battery?', a: 'Accountable2You uses a VPN-based monitoring approach, which routes device traffic through a local VPN. This approach can increase battery drain and occasionally causes conflicts with corporate or school VPN configurations. Remain Faithful uses Apple\'s Screen Time framework for always-on monitoring, which has minimal battery impact.' },
                { q: 'Can I switch from Accountable2You to Remain Faithful?', a: 'Yes. Cancel your Accountable2You subscription, remove the VPN profile from your device, and download Remain Faithful. Note that Remain Faithful is currently iOS only; if your partners use Android or Windows, Accountable2You has broader platform support.' },
                { q: 'Does Remain Faithful log specific web page titles like Accountable2You?', a: 'No. Remain Faithful does not log or transmit web page titles or browsing history. Partners receive only a category label and severity level when something is flagged. Accountable2You logs specific page titles in its reports, which can expose more detail than some users want their partners to see.' },
              ].map((faq) => (
                <div key={faq.q} className="rounded-2xl border border-[#1E3050] bg-[#162235] p-6">
                  <h3 className="font-semibold text-[#F0EDE8] mb-3">{faq.q}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{faq.a}</p>
                </div>
              ))}
            </div>
          </section>

          {/* Other comparisons */}
          <section className="mb-14">
            <h2 className="font-serif text-xl font-bold text-[#F0EDE8] mb-4">Other Comparisons</h2>
            <div className="grid sm:grid-cols-2 gap-4">
              <Link href="/compare/covenant-eyes" className="p-4 rounded-xl border border-[#1E3050] bg-[#162235] hover:border-[#C9A84C]/40 transition-colors text-[#C9A84C] font-medium text-sm">
                Remain Faithful vs Covenant Eyes →
              </Link>
              <Link href="/compare/ever-accountable" className="p-4 rounded-xl border border-[#1E3050] bg-[#162235] hover:border-[#C9A84C]/40 transition-colors text-[#C9A84C] font-medium text-sm">
                Remain Faithful vs Ever Accountable →
              </Link>
            </div>
          </section>

          <div className="text-center p-10 rounded-3xl border border-[#C9A84C]/20" style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}>
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-4">Try Remain Faithful Free</h2>
            <p className="text-[#8A9BB0] mb-6">No subscription. No VPN. No credit card. Just accountability.</p>
            <Link href="/#waitlist" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96]">
              Join the Waitlist
            </Link>
          </div>
        </div>
      </div>
    </>
  )
}
