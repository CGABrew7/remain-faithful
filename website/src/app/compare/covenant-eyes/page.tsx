import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { comparisonSchema } from '@/lib/structured-data'
import { COMPETITOR_PRICING_NOTE } from '@/lib/competitor-pricing'
import { Breadcrumbs } from '@/components/Breadcrumbs'

export const metadata: Metadata = {
  title: 'Remain Faithful vs Covenant Eyes: Honest Comparison (2026)',
  description: 'A detailed, honest comparison of Remain Faithful and Covenant Eyes for Christian accountability. See how they differ on price, privacy, monitoring, and approach.',
  alternates: { canonical: 'https://remainfaithful.com/compare/covenant-eyes' },
}

const pageFaqSchema = {
  '@context': 'https://schema.org',
  '@type': 'FAQPage',
  mainEntity: [
    {
      '@type': 'Question',
      name: 'Is Remain Faithful as effective as Covenant Eyes?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Remain Faithful and Covenant Eyes take different approaches. Covenant Eyes has 25 years of proven history and broad multi-platform coverage. Remain Faithful is newer but introduces on-device AI that keeps screen content entirely private, a covenant-based accountability model, and costs nothing. Effectiveness depends on what your situation requires.',
      },
    },
    {
      '@type': 'Question',
      name: 'Can I switch from Covenant Eyes to Remain Faithful?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes. Switching is straightforward: cancel your Covenant Eyes subscription, download Remain Faithful, and invite your accountability partners. Setup takes about 15 minutes. Note that Remain Faithful is currently iOS only; if your partners use Android, Windows, or Mac, Covenant Eyes may be the better fit for now.',
      },
    },
    {
      '@type': 'Question',
      name: 'Does Remain Faithful work on the same devices as Covenant Eyes?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'No. Covenant Eyes supports iOS, Android, Windows, Mac, and Chromebook. Remain Faithful currently supports iPhone (iOS 17+). Android support is planned for Fall 2026. If you need multi-platform coverage today, Covenant Eyes has the advantage.',
      },
    },
  ],
}

export default function CovenantEyesCompare() {
  return (
    <>
      <JsonLd data={comparisonSchema('Covenant Eyes')} />
      <JsonLd data={pageFaqSchema} />

      <div className="pt-24 pb-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <Breadcrumbs items={[
            { name: 'Compare', url: 'https://remainfaithful.com/compare/covenant-eyes' },
            { name: 'Covenant Eyes', url: 'https://remainfaithful.com/compare/covenant-eyes' },
          ]} />

          <div className="mb-10">
            <p className="text-wax text-sm font-semibold uppercase tracking-widest mb-4">Honest Comparison</p>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-6 leading-tight">
              Remain Faithful vs Covenant Eyes
            </h1>
            <p className="text-ink-soft text-lg leading-relaxed mb-4">
              Both Remain Faithful and Covenant Eyes exist to help Christians pursue purity through accountability. They take fundamentally different approaches. This page is an honest comparison so you can choose the right tool for your situation.
            </p>
            <div className="p-5 rounded-sm border border-wax/25 bg-wax/5 text-sm text-ink-soft leading-relaxed">
              <strong className="text-wax">Disclosure:</strong> We built Remain Faithful, so we are biased. We will be transparent about where Covenant Eyes is stronger.
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
                    <th className="text-center p-4 text-ink font-semibold bg-paper-deep">Covenant Eyes</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    ['Price', 'Free forever', '$18/mo or $198/yr (account, not per-user)'],
                    ['Annual cost', '$0', '$198/yr (up to 10 members on one account)'],
                    ['Platform', 'iOS 17+ (Android Fall 2026)', 'iOS, Android, Windows, Mac, Chromebook'],
                    ['Monitoring approach', 'On-device AI classification', 'Screen capture + cloud AI analysis'],
                    ['Screenshots shared with partners?', 'Never. Category + severity only.', 'Yes. Screenshots sent to partner.'],
                    ['Where AI runs', 'Entirely on your device', 'Cloud servers'],
                    ['Open source', 'Yes, full codebase on GitHub', 'No'],
                    ['Accountability model', 'Covenant-based partnership or small group', 'Ally-based reports'],
                    ['Church/group tools', 'Built-in group mode (up to 12)', 'Available with enterprise plan'],
                    ['Content blocking', 'Yes, always-on app blocking via Screen Time', 'Yes, customizable filtering'],
                    ['DRM streaming monitoring', 'No (iOS limitation)', 'No (iOS limitation)'],
                    ['Established since', '2026', '2000'],
                    ['User base', 'Pre-launch (waitlist)', '1.5M+ users'],
                  ].map(([feature, rf, ce], i) => (
                    <tr key={i} className={`border-b border-hairline ${i % 2 === 0 ? 'bg-paper-deep' : 'bg-paper'}`}>
                      <td className="p-4 text-ink-soft font-medium">{feature}</td>
                      <td className="p-4 text-center text-ink">{rf}</td>
                      <td className="p-4 text-center text-ink-soft">{ce}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          {/* Where CE is stronger */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-ink mb-5">Where Covenant Eyes Is Stronger</h2>
            <div className="space-y-4">
              {[
                { title: 'Platform coverage', body: 'Covenant Eyes supports 5 platforms: iOS, Android, Windows, Mac, and Chromebook. Remain Faithful is iOS only right now. If you or your partners use non-Apple devices, CE has the clear advantage.' },
                { title: 'Track record', body: 'Covenant Eyes has been in operation since 2000 with over 1.5 million users. That is 25 years of proven history. Remain Faithful is in pre-launch and has not shipped yet. Maturity matters for trust.' },
                { title: 'Content filtering', body: 'Covenant Eyes offers robust website and app blocking with granular category controls. Remain Faithful also blocks apps via Screen Time integration, but CE has a more mature filtering system overall.' },
                { title: 'Enterprise and church administration', body: 'Covenant Eyes has mature organizational management tools designed for large-scale church or ministry deployments. For denominations or large organizations with IT requirements, CE has purpose-built infrastructure.' },
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
                { title: 'Price', body: 'Remain Faithful is free. Covenant Eyes publicly lists $18/month or $198/year for an account (up to 10 members) — not a per-user fee, and not multiplied by household size. There is no version of Covenant Eyes that is free; Remain Faithful is free without limitations. Verify current pricing on their site.' },
                { title: 'Privacy: no screenshots ever leave your device', body: 'Covenant Eyes sends screenshots of your screen to cloud servers for AI analysis. Remain Faithful does all classification on your device using Apple\'s Neural Engine. No screen content, no screenshots, and no OCR text is ever transmitted. Your screen stays on your phone.' },
                { title: 'Partner protection', body: 'Covenant Eyes partners receive actual screenshots of flagged content, which means partners are exposed to the harmful material. Remain Faithful partners receive only a category label and severity level. Partners never see the content that triggered the alert.' },
                { title: 'Open source transparency', body: 'Anyone can read every line of Remain Faithful\'s code on GitHub and verify exactly what is and is not transmitted. Covenant Eyes is proprietary. For an app handling sensitive behavioral data, open source is a meaningful trust advantage.' },
                { title: 'Covenant model', body: 'The covenant framework grounds accountability in a biblical understanding of relationship rather than a surveillance model. Partners agree to respond with grace before they gain any access. This shapes how alerts are received and processed.' },
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
                <h2 className="font-serif text-lg font-bold text-ink mb-4">Who Should Choose Covenant Eyes</h2>
                <ul className="space-y-2">
                  {[
                    'Need Windows, Mac, or Android coverage today',
                    'Organizations with mature IT admin requirements',
                    'Families who want a proven 25-year track record',
                    'Users whose partners prefer to see screenshots',
                    'Large ministries needing enterprise management',
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
                    'iPhone users who want privacy-first accountability at zero cost',
                    'Anyone who believes partners should not see harmful content',
                    'Those who value open-source code transparency',
                    'Churches and small groups needing free tools',
                    'People who want a covenant-based theological framework',
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
                { q: 'Is Remain Faithful as effective as Covenant Eyes?', a: 'Remain Faithful and Covenant Eyes take different approaches. Covenant Eyes has 25 years of proven history and broad multi-platform coverage. Remain Faithful is newer but introduces on-device AI that keeps screen content entirely private, a covenant-based accountability model, and costs nothing. Effectiveness depends on what your situation requires.' },
                { q: 'Can I switch from Covenant Eyes to Remain Faithful?', a: 'Yes. Switching is straightforward: cancel your Covenant Eyes subscription, download Remain Faithful, and invite your accountability partners. Setup takes about 15 minutes. Note that Remain Faithful is currently iOS only; if your partners use Android, Windows, or Mac, Covenant Eyes may be the better fit for now.' },
                { q: 'Does Remain Faithful work on the same devices as Covenant Eyes?', a: 'No. Covenant Eyes supports iOS, Android, Windows, Mac, and Chromebook. Remain Faithful currently supports iPhone (iOS 17+). Android support is planned for Fall 2026. If you need multi-platform coverage today, Covenant Eyes has the advantage.' },
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
              <Link href="/compare/ever-accountable" className="p-4 rounded-sm border border-hairline bg-paper-deep hover:border-wax/40 transition-colors text-wax font-medium text-sm">
                Remain Faithful vs Ever Accountable →
              </Link>
              <Link href="/compare/accountable2you" className="p-4 rounded-sm border border-hairline bg-paper-deep hover:border-wax/40 transition-colors text-wax font-medium text-sm">
                Remain Faithful vs Accountable2You →
              </Link>
            </div>
          </section>

          {/* CTA */}
          <div className="text-center p-10 rounded-sm border border-wax/25" >
            <h2 className="font-serif text-2xl font-bold text-ink mb-4">Try Remain Faithful Free</h2>
            <p className="text-ink-soft mb-6">No subscription. No credit card. Just accountability.</p>
            <Link href="/#waitlist" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-sm font-semibold text-paper bg-wax hover:bg-wax-deep transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96]">
              Join the Waitlist
            </Link>
          </div>
        </div>
      </div>
    </>
  )
}
