import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { comparisonSchema } from '@/lib/structured-data'
import { Breadcrumbs } from '@/components/Breadcrumbs'

export const metadata: Metadata = {
  title: 'Remain Faithful vs Ever Accountable: Honest Comparison (2026)',
  description: 'A detailed, honest comparison of Remain Faithful and Ever Accountable for Christian accountability. See how they differ on price, privacy, monitoring, and approach.',
  alternates: { canonical: 'https://remainfaithful.com/compare/ever-accountable' },
}

const pageFaqSchema = {
  '@context': 'https://schema.org',
  '@type': 'FAQPage',
  mainEntity: [
    {
      '@type': 'Question',
      name: 'Is Remain Faithful as effective as Ever Accountable?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Both apps use on-device AI, which is a meaningful shared strength. Ever Accountable has a more mature iOS implementation with four generations of development. Remain Faithful is newer but is entirely free, open source, and built on a covenant accountability model. Which is more effective depends on your specific situation.',
      },
    },
    {
      '@type': 'Question',
      name: 'Does Ever Accountable also do on-device processing?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes. Ever Accountable also performs on-device AI analysis. This is a shared strength with Remain Faithful. The key differences are price (Remain Faithful is free, Ever Accountable is $99/year), open source transparency (Remain Faithful is fully open source, Ever Accountable is not), and the covenant-based model that Remain Faithful uses.',
      },
    },
    {
      '@type': 'Question',
      name: 'Can I switch from Ever Accountable to Remain Faithful?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes. Cancel your Ever Accountable subscription and download Remain Faithful. Invite your accountability partners and walk through the covenant setup. Note that Remain Faithful is currently iOS only; if your partners use Android, Windows, or Mac, Ever Accountable has broader coverage.',
      },
    },
  ],
}

export default function EverAccountableCompare() {
  return (
    <>
      <JsonLd data={comparisonSchema('Ever Accountable')} />
      <JsonLd data={pageFaqSchema} />

      <div className="pt-24 pb-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <Breadcrumbs items={[
            { name: 'Compare', url: 'https://remainfaithful.com/compare/ever-accountable' },
            { name: 'Ever Accountable', url: 'https://remainfaithful.com/compare/ever-accountable' },
          ]} />

          <div className="mb-10">
            <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-4">Honest Comparison</p>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-6 leading-tight">
              Remain Faithful vs Ever Accountable
            </h1>
            <p className="text-[#8A9BB0] text-lg leading-relaxed mb-4">
              Both apps take a privacy-first approach with on-device AI. This is a closer comparison than most. Here is where they actually differ.
            </p>
            <div className="p-5 rounded-xl border border-[#C9A84C]/20 bg-[#C9A84C]/5 text-sm text-[#8A9BB0] leading-relaxed">
              <strong className="text-[#C9A84C]">Disclosure:</strong> We built Remain Faithful, so we are biased. We will be transparent about where Ever Accountable is stronger.
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
                    <th className="text-center p-4 text-[#F0EDE8] font-semibold bg-[#0A1420]">Ever Accountable</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    ['Price', 'Free forever', '$99/year'],
                    ['Platform', 'iOS 17+ (Android Fall 2026)', 'iOS, Android, Windows, Mac'],
                    ['Monitoring approach', 'On-device AI classification', 'On-device AI classification'],
                    ['Screenshots shared with partners?', 'Never. Category + severity only.', 'Partners see screenshot reports'],
                    ['Where AI runs', 'Entirely on your device', 'On-device'],
                    ['Open source', 'Yes, full codebase on GitHub', 'No'],
                    ['Accountability model', 'Covenant-based partnership or small group', 'Partner-based report sharing'],
                    ['Professional coach option', 'No', 'Yes, paid add-on'],
                    ['App maturity', 'New (beta, 2026)', '4th-generation iOS app, 700K+ installs'],
                    ['DRM streaming monitoring', 'No (iOS limitation)', 'No (iOS limitation)'],
                  ].map(([feature, rf, ea], i) => (
                    <tr key={i} className={`border-b border-[#1E3050] ${i % 2 === 0 ? 'bg-[#162235]' : 'bg-[#0F1B2D]'}`}>
                      <td className="p-4 text-[#8A9BB0] font-medium">{feature}</td>
                      <td className="p-4 text-center text-[#F0EDE8]">{rf}</td>
                      <td className="p-4 text-center text-[#8A9BB0]">{ea}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          {/* Where EA is stronger */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Where Ever Accountable Is Stronger</h2>
            <div className="space-y-4">
              {[
                { title: 'More mature iOS implementation', body: 'Ever Accountable is on its fourth generation of iOS development with over 700,000 installs. That level of real-world testing produces a more polished, battle-tested app. Remain Faithful is in beta.' },
                { title: 'Multi-platform coverage', body: 'Ever Accountable supports iOS, Android, Windows, and Mac. Remain Faithful is iOS only right now. If your accountability partner uses Android, Windows, or Mac, Ever Accountable serves them. Remain Faithful does not yet.' },
                { title: 'Professional accountability coach option', body: 'Ever Accountable offers a paid add-on to connect with a professional accountability coach. Remain Faithful is peer-only; it does not offer a coached accountability option.' },
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
                { title: 'Price', body: 'Remain Faithful is completely free. Ever Accountable costs $99 per year. Over five years, that is $495 in savings. For ministries deploying accountability tools to dozens of people, free is a significant difference.' },
                { title: 'Open source transparency', body: 'Remain Faithful is fully open source. Anyone can read the code and verify that screen content is never transmitted. Ever Accountable is proprietary. Both say they protect your privacy; only one lets you verify it.' },
                { title: 'Partner protection', body: 'Ever Accountable shares screenshot reports with accountability partners, exposing them to the content that was flagged. Remain Faithful partners receive only a category and severity label. Partners never see the material that triggered the alert.' },
                { title: 'Covenant model', body: 'The covenant framework grounds Remain Faithful\'s accountability in a theological understanding of relationship. Partners agree to respond with grace before gaining access. This is designed to shape how alerts are received.' },
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
                <h2 className="font-serif text-lg font-bold text-[#F0EDE8] mb-4">Who Should Choose Ever Accountable</h2>
                <ul className="space-y-2">
                  {[
                    'Need Android, Windows, or Mac coverage today',
                    'Want a more mature, battle-tested iOS app',
                    'Interested in professional accountability coaching',
                    'Comfortable paying $99/year for a proven solution',
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
                    'iPhone users who want free, privacy-first accountability',
                    'Anyone who values open-source code verification',
                    'Churches and small groups needing free group tools',
                    'People who believe partners should not see screenshots',
                    'Those drawn to a covenant-based theological framework',
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
                { q: 'Is Remain Faithful as effective as Ever Accountable?', a: 'Both apps use on-device AI, which is a meaningful shared strength. Ever Accountable has a more mature iOS implementation with four generations of development. Remain Faithful is newer but is entirely free, open source, and built on a covenant accountability model. Which is more effective depends on your specific situation.' },
                { q: 'Does Ever Accountable also do on-device processing?', a: 'Yes. Ever Accountable also performs on-device AI analysis. This is a shared strength with Remain Faithful. The key differences are price (Remain Faithful is free, Ever Accountable is $99/year), open source transparency (Remain Faithful is fully open source, Ever Accountable is not), and the covenant-based model that Remain Faithful uses.' },
                { q: 'Can I switch from Ever Accountable to Remain Faithful?', a: 'Yes. Cancel your Ever Accountable subscription and download Remain Faithful. Invite your accountability partners and walk through the covenant setup. Note that Remain Faithful is currently iOS only; if your partners use Android, Windows, or Mac, Ever Accountable has broader coverage.' },
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
              <Link href="/compare/accountable2you" className="p-4 rounded-xl border border-[#1E3050] bg-[#162235] hover:border-[#C9A84C]/40 transition-colors text-[#C9A84C] font-medium text-sm">
                Remain Faithful vs Accountable2You →
              </Link>
            </div>
          </section>

          <div className="text-center p-10 rounded-3xl border border-[#C9A84C]/20" style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}>
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-4">Try Remain Faithful Free</h2>
            <p className="text-[#8A9BB0] mb-6">No subscription. No credit card. Just accountability.</p>
            <Link href="/#waitlist" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96]">
              Join the Waitlist
            </Link>
          </div>
        </div>
      </div>
    </>
  )
}
