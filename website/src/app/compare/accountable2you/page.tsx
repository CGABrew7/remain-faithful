import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { comparisonSchema } from '@/lib/structured-data'
import { COMPARE_PRICING_NOTE } from '@/lib/competitor-pricing'
import { Breadcrumbs } from '@/components/Breadcrumbs'
import { UpdatedStamp } from '@/components/UpdatedStamp'

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
        text: 'The Accountable2You pricing, features, and iOS pages checked on 2026-10-13 do not state a battery figure. Their iOS page says VPN technology records internet activity in Safari, Chrome, or any browser. Remain Faithful does not use a VPN. It uses Apple Family Controls for always-on filtering, with optional on-device classification.',
      },
    },
    {
      '@type': 'Question',
      name: 'Can I switch from Accountable2You to Remain Faithful?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes. Cancel Accountable2You and set up Remain Faithful on iPhone (iOS 17+). Remain Faithful is iPhone only. Accountable2You lists iPhone, iPad, Android, Windows, Mac, Linux, and Chromebook. Their public pages checked on 2026-10-13 do not describe a required VPN-profile removal step.',
      },
    },
    {
      '@type': 'Question',
      name: 'Does Remain Faithful log specific web page titles like Accountable2You?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'No. Remain Faithful does not log or transmit web page titles or browsing history. Partners may receive a short system-generated summary in addition to category, severity, and timestamp when something is flagged — not page titles. Accountable2You reports browsing history and activity (including Incognito, per their site), which can expose more detail than some users want their partners to see.',
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
            <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
              <p className="text-wax text-sm font-semibold uppercase tracking-widest">Honest Comparison</p>
              <UpdatedStamp />
            </div>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-6 leading-tight text-balance">
              Remain Faithful vs Accountable2You
            </h1>
            <p className="text-ink-soft text-lg leading-relaxed mb-4 text-pretty">
              Remain Faithful is free. It classifies on your iPhone and never sends screenshots to a partner. On 2026-10-13, Accountable2You&apos;s public pricing page listed Personal at $11 per month or $121 per year ($10.08 per month, billed annually) and Family at $16 per month or $176 per year ($14.67 per month, billed annually). They say they do not capture screenshots. They do report browsing history, including Incognito, and full web page titles.
            </p>
            <div className="p-5 rounded-sm border border-wax/25 bg-wax/5 text-sm text-ink-soft leading-relaxed">
              <strong className="text-wax">Disclosure:</strong> We built Remain Faithful, so we are biased. We will be transparent about where Accountable2You is stronger.
            </div>
          </div>

          {/* Comparison Table */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-ink mb-6">Side-by-Side Comparison</h2>
            <p className="text-xs text-ink-faint mb-4">{COMPARE_PRICING_NOTE}</p>
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
                    ['Price', 'Free forever', 'Personal: $11/month or $121/year ($10.08/month billed annually). Family: $16/month or $176/year ($14.67/month billed annually).'],
                    ['Platform', 'iOS 17+ (Android planned)', 'iPhone, iPad, Android, Windows, Mac, Linux, Chromebook'],
                    ['Monitoring technology', 'Family Controls + optional on-device AI', 'iOS page: VPN technology records browser activity. Features page: browsing history, app usage, and page titles.'],
                    ['Browsing history reported', 'No. No page titles or browsing history.', 'Yes. Browsing history, including Incognito, and full web page titles.'],
                    ['Screenshots shared?', 'Never. Category, severity, timestamp, and a short summary.', 'They say they do not capture screenshots.'],
                    ['VPN', 'No', 'Named on their iOS page for recording browser activity.'],
                    ['Where AI runs', 'Entirely on your device', 'Not described as on-device AI on the public pages we checked'],
                    ['Open source', 'Yes, full codebase on GitHub', 'No'],
                    ['Partner data shared', 'Category, severity, timestamp, and a short system-generated summary', 'Activity reports and alerts. They say no screenshots. Partners can be set to full reports or alerts only.'],
                    ['Accountability model', 'Covenant-based partnership', 'Report-based accountability'],
                    ['Group / family', 'Yes, up to 12 members', 'Family: parent and child accounts, up to 20 installs. Personal: up to 6 installs.'],
                    ['Content blocking', 'Yes, always-on Family Controls', 'Their pricing FAQ says they do not block websites or filter content.'],
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
                { title: 'Cross-platform coverage', body: 'Accountable2You lists iPhone, iPad, Android, Windows, Mac, Linux, and Chromebook. Remain Faithful is iPhone only right now. For mixed devices, Accountable2You covers platforms Remain Faithful does not ship on yet.' },
                { title: 'Detailed web activity reporting', body: 'Their features page says reports include full web page titles, and that Incognito and private browsing are visible to partners. Their monitoring page says iOS app usage is limited compared with other platforms. If a partner wants that browsing detail, Accountable2You describes more of it than Remain Faithful’s category, severity, timestamp, and short summary.' },
                { title: 'Family and personal install limits', body: 'The public pricing page lists a Family plan with up to 20 active installs, parent and child accounts, and time-monitoring alerts, and a Personal plan with up to 6 active installs and unlimited partners. Remain Faithful does not publish a paid family tier. Its group mode is included at no charge, up to 12 members.' },
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
                { title: 'Price', body: 'Remain Faithful is free. On 2026-10-13, Accountable2You’s public pricing page listed Personal at $11 per month or $121 per year ($10.08 USD per month, billed annually) and Family at $16 per month or $176 per year ($14.67 USD per month, billed annually). The annual cards say “GET 1 MONTH FREE.” The page offers a 10-day free trial. Verify current amounts on their site.' },
                { title: 'On-device classification, no VPN in Remain Faithful', body: 'Remain Faithful uses Apple Family Controls for always-on filtering, with optional on-device AI. Screenshots stay on the phone. Accountable2You’s iOS page says their VPN technology records internet activity in Safari, Chrome, or any browser. Their public pages checked that day do not state a battery figure. Their pricing FAQ says they do not block or filter content, and that they do not capture screenshots.' },
                { title: 'On-device AI when you opt in', body: 'Remain Faithful’s optional Deep Scan classifies non-DRM screen frames on-device. It cannot see DRM-protected video or banking apps (Netflix, Disney+, Hulu, Prime Video, Apple TV, HBO, and banking apps render as black frames). Always-on filtering and usage monitoring do not look at screen content. Accountable2You’s public pages describe browsing history, page titles, and app usage. They do not describe on-device screen classification.' },
                { title: 'What partners receive', body: 'Accountable2You’s features page says partners can see full web page titles and Incognito browsing, and their pricing FAQ says they do not capture screenshots. You can set a partner to full activity reports or alerts only. Remain Faithful may share a short system-generated summary in addition to category, severity, and timestamp. Partners do not receive page titles, screenshots, or raw screen content.' },
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
                    'Need iPad, Android, Windows, Mac, Linux, or Chromebook coverage',
                    'Partners want browsing-history reports and page titles',
                    'The public Personal or Family price fits the household',
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
                    'iPhone users who want on-device classification at no cost',
                    'People who do not want browsing history or page titles sent to a partner',
                    'People who do not want screenshots sent to a partner',
                    'Those who value open-source verification',
                    'Churches using the included group mode (up to 12)',
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
                { q: 'Does Accountable2You drain the battery?', a: 'The Accountable2You pricing, features, and iOS pages checked on 2026-10-13 do not state a battery figure. Their iOS page says VPN technology records internet activity in Safari, Chrome, or any browser. Remain Faithful does not use a VPN. It uses Apple Family Controls for always-on filtering, with optional on-device classification.' },
                { q: 'Can I switch from Accountable2You to Remain Faithful?', a: 'Yes. Cancel Accountable2You and set up Remain Faithful on iPhone (iOS 17+). Remain Faithful is iPhone only. Accountable2You lists iPhone, iPad, Android, Windows, Mac, Linux, and Chromebook. Their public pages checked on 2026-10-13 do not describe a required VPN-profile removal step.' },
                { q: 'Does Remain Faithful log specific web page titles like Accountable2You?', a: 'No. Remain Faithful does not log or transmit web page titles or browsing history. Partners may receive a short system-generated summary in addition to category, severity, and timestamp when something is flagged — not page titles. Accountable2You reports browsing history and activity (including Incognito, per their site), which can expose more detail than some users want their partners to see.' },
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
            <Link href="/#download" className="btn-wax">
              Get started
            </Link>
          </div>
        </div>
      </div>
    </>
  )
}
