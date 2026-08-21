import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { comparisonSchema } from '@/lib/structured-data'
import { COMPETITOR_PRICING_NOTE } from '@/lib/competitor-pricing'
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
        text: 'Accountable2You uses a VPN-based monitoring approach, which routes device traffic through a local VPN. This approach can increase battery drain and occasionally causes conflicts with corporate or school VPN configurations. Remain Faithful uses Apple Family Controls for always-on filtering, which has minimal battery impact.',
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
            <p className="text-wax text-sm font-semibold uppercase tracking-widest mb-4">Honest Comparison</p>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-6 leading-tight">
              Remain Faithful vs Accountable2You
            </h1>
            <p className="text-ink-soft text-lg leading-relaxed mb-4">
              Accountable2You takes a VPN-based approach to monitoring, which is meaningfully different from Remain Faithful&apos;s on-device AI model. Here is what that difference means in practice.
            </p>
            <div className="p-5 rounded-sm border border-wax/25 bg-wax/5 text-sm text-ink-soft leading-relaxed">
              <strong className="text-wax">Disclosure:</strong> We built Remain Faithful, so we are biased. We will be transparent about where Accountable2You is stronger.
            </div>
          </div>

          {/* Comparison Table */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-ink mb-6">Side-by-Side Comparison</h2>
            <p className="text-xs text-ink-soft/60 mb-4">{COMPETITOR_PRICING_NOTE}</p>
            <div className="overflow-x-auto rounded-sm border border-hairline">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-hairline">
                    <th className="text-left p-4 text-ink-soft font-semibold bg-paper-deep w-1/3">Feature</th>
                    <th className="text-center p-4 text-wax font-semibold bg-paper-deep">Remain Faithful</th>
                    <th className="text-center p-4 text-ink font-semibold bg-paper-deep">Accountable2You</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    ['Price', 'Free forever', '$121/year Personal'],
                    ['Platform', 'iOS 17+ (Android Fall 2026)', 'iOS, Android, Windows, Mac, Kindle'],
                    ['Monitoring technology', 'Family Controls + optional on-device AI', 'VPN-based traffic monitoring'],
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
                    <tr key={i} className={`border-b border-hairline ${i % 2 === 0 ? 'bg-paper-deep' : 'bg-paper'}`}>
                      <td className="p-4 text-ink-soft font-medium">{feature}</td>
                      <td className="p-4 text-center text-ink">{rf}</td>
                      <td className="p-4 text-center text-ink-soft">{a2y}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          {/* Where A2Y is stronger */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-ink mb-5">Where Accountable2You Is Stronger</h2>
            <div className="space-y-4">
              {[
                { title: 'Cross-platform coverage', body: 'Accountable2You runs on iOS, Android, Windows, Mac, and Kindle. Remain Faithful is currently iOS only. For families with mixed devices, A2Y provides unified coverage across all of them.' },
                { title: 'Detailed web activity reporting', body: 'Accountable2You logs specific web page titles and generates detailed activity reports. If your accountability partner or pastor wants more granular visibility into browsing behavior, A2Y provides more detail than Remain Faithful\'s category-level alerts.' },
                { title: 'Established track record', body: 'Accountable2You has been in operation for years with a larger installed base. Remain Faithful is in pre-launch and has not shipped yet. For users who want a proven, stable product, A2Y has an advantage in maturity.' },
              ].map((item) => (
                <div key={item.title} className="p-5 rounded-sm border border-hairline bg-paper-deep">
                  <h3 className="font-semibold text-ink mb-2">{item.title}</h3>
                  <p className="text-sm text-ink-soft leading-relaxed">{item.body}</p>
                </div>
              ))}
            </div>
          </section>

          {/* Where RF is stronger */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-ink mb-5">Where Remain Faithful Is Stronger</h2>
            <div className="space-y-4">
              {[
                { title: 'Price', body: 'Remain Faithful is free. Accountable2You publicly lists $121 per year for the Personal plan. For individuals or churches deploying accountability tools at scale, free is a meaningful difference. Verify current pricing on their site.' },
                { title: 'No VPN required: better battery life', body: 'Accountable2You routes traffic through a local VPN to monitor it. This has two downsides: it drains battery more quickly, and it can conflict with corporate or school VPN configurations. Remain Faithful uses Apple Family Controls for always-on filtering, with optional on-device AI, minimal battery impact, and no VPN conflicts.' },
                { title: 'On-device AI when you opt in', body: 'VPN-based monitoring cannot see inside encrypted HTTPS traffic or monitor what happens inside apps. Remain Faithful\'s optional Deep Scan classifies non-DRM screen frames on-device. It cannot see DRM-protected video or banking apps (Netflix, Disney+, Hulu, Prime Video, Apple TV, HBO, and banking apps render as black frames). Always-on filtering and usage monitoring do not look at screen content.' },
                { title: 'Partner privacy: category only, never page titles', body: 'Accountable2You shares specific web page titles with accountability partners, which means partners see exactly what pages were visited. Some find this level of detail helpful; others consider it overly exposing. Remain Faithful shares only the category and severity, protecting the privacy of the person being held accountable while still triggering a conversation.' },
                { title: 'Open source', body: 'Remain Faithful\'s entire codebase is publicly available for inspection. Accountable2You is proprietary. Anyone can verify what Remain Faithful does and does not transmit.' },
              ].map((item) => (
                <div key={item.title} className="p-5 rounded-sm border border-wax/25 bg-paper-deep">
                  <h3 className="font-semibold text-wax mb-2">{item.title}</h3>
                  <p className="text-sm text-ink-soft leading-relaxed">{item.body}</p>
                </div>
              ))}
            </div>
          </section>

          {/* Who should choose what */}
          <section className="mb-14">
            <div className="grid md:grid-cols-2 gap-6">
              <div className="p-6 rounded-sm border border-hairline bg-paper-deep">
                <h2 className="font-serif text-lg font-bold text-ink mb-4">Who Should Choose Accountable2You</h2>
                <ul className="space-y-2">
                  {[
                    'Need Android, Windows, Mac, or Kindle coverage',
                    'Partners want detailed web page title reports',
                    'Happy to pay a yearly subscription for a proven solution',
                    'No concerns about VPN battery drain or conflicts',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2 text-sm text-ink-soft">
                      <svg className="flex-shrink-0 mt-0.5 text-ink-soft" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><polyline points="20 6 9 17 4 12"/></svg>
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
              <div className="p-6 rounded-sm border border-wax/25 bg-paper-deep">
                <h2 className="font-serif text-lg font-bold text-ink mb-4">Who Should Choose Remain Faithful</h2>
                <ul className="space-y-2">
                  {[
                    'iPhone users who want on-device AI at zero cost',
                    'Anyone on a corporate or school VPN',
                    'Those who want battery-efficient monitoring',
                    'People who prefer category alerts over page title logs',
                    'Those who value open-source verification',
                    'Churches needing free group accountability tools',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2 text-sm text-ink-soft">
                      <svg className="flex-shrink-0 mt-0.5 text-wax" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><polyline points="20 6 9 17 4 12"/></svg>
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </section>

          {/* FAQ */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-ink mb-6">Frequently Asked Questions</h2>
            <div className="space-y-4">
              {[
                { q: 'Does Accountable2You drain the battery?', a: 'Accountable2You uses a VPN-based monitoring approach, which routes device traffic through a local VPN. This approach can increase battery drain and occasionally causes conflicts with corporate or school VPN configurations. Remain Faithful uses Apple Family Controls for always-on filtering, which has minimal battery impact.' },
                { q: 'Can I switch from Accountable2You to Remain Faithful?', a: 'Yes. Cancel your Accountable2You subscription, remove the VPN profile from your device, and download Remain Faithful. Note that Remain Faithful is currently iOS only; if your partners use Android or Windows, Accountable2You has broader platform support.' },
                { q: 'Does Remain Faithful log specific web page titles like Accountable2You?', a: 'No. Remain Faithful does not log or transmit web page titles or browsing history. Partners receive only a category label and severity level when something is flagged. Accountable2You logs specific page titles in its reports, which can expose more detail than some users want their partners to see.' },
              ].map((faq) => (
                <div key={faq.q} className="rounded-sm border border-hairline bg-paper-deep p-6">
                  <h3 className="font-semibold text-ink mb-3">{faq.q}</h3>
                  <p className="text-sm text-ink-soft leading-relaxed">{faq.a}</p>
                </div>
              ))}
            </div>
          </section>

          {/* Other comparisons */}
          <section className="mb-14">
            <h2 className="font-serif text-xl font-bold text-ink mb-4">Other Comparisons</h2>
            <div className="grid sm:grid-cols-2 gap-4">
              <Link href="/compare/covenant-eyes" className="p-4 rounded-sm border border-hairline bg-paper-deep hover:border-wax/40 transition-colors text-wax font-medium text-sm">
                Remain Faithful vs Covenant Eyes →
              </Link>
              <Link href="/compare/ever-accountable" className="p-4 rounded-sm border border-hairline bg-paper-deep hover:border-wax/40 transition-colors text-wax font-medium text-sm">
                Remain Faithful vs Ever Accountable →
              </Link>
            </div>
          </section>

          <div className="text-center p-10 rounded-sm border border-wax/25" >
            <h2 className="font-serif text-2xl font-bold text-ink mb-4">Try Remain Faithful Free</h2>
            <p className="text-ink-soft mb-6">No subscription. No VPN. No credit card. Just accountability.</p>
            <Link href="/#waitlist" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-sm font-semibold text-paper bg-wax hover:bg-wax-deep transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96]">
              Join the Waitlist
            </Link>
          </div>
        </div>
      </div>
    </>
  )
}
