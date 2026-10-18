import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { comparisonSchema } from '@/lib/structured-data'
import { COMPARE_PRICING_NOTE } from '@/lib/competitor-pricing'
import { Breadcrumbs } from '@/components/Breadcrumbs'
import { UpdatedStamp } from '@/components/UpdatedStamp'

export const metadata: Metadata = {
  title: 'Remain Faithful vs Ever Accountable: Honest Comparison (2026)',
  description: 'A detailed comparison of Remain Faithful and Ever Accountable for Christian accountability. See how they differ on price, privacy, monitoring, and approach.',
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
        text: 'Remain Faithful is free and classifies on your iPhone. Screenshots never leave the device. Ever Accountable states 700,000+ installs on its public pages. A June 6, 2026 post on their site describes an AI on the iPhone that looks at the screen. Their how-it-works page says partners have access to browsing history. Which fits depends on what a partner should see, and on price.',
      },
    },
    {
      '@type': 'Question',
      name: 'Does Ever Accountable also do on-device processing?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Ever Accountable’s June 6, 2026 post says an AI on the iPhone looks at the screen, using the phone’s screen-recording feature rather than making a video anyone watches. That is their description, not a claim that the products match. Remain Faithful is free, never transmits screenshots, and is open source. Their plans page lists $14.99 per month or $129 per year for one reporting user. Their how-it-works page says partners have access to browsing history.',
      },
    },
    {
      '@type': 'Question',
      name: 'Can I switch from Ever Accountable to Remain Faithful?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes. Cancel Ever Accountable and set up Remain Faithful on iPhone (iOS 17+). Invite your partners and walk through the covenant setup. Remain Faithful is available for iPhone only; if your partners use Android, Windows, or Mac, Ever Accountable has broader coverage.',
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
            <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
              <p className="text-wax text-sm font-semibold uppercase tracking-widest">Honest Comparison</p>
              <UpdatedStamp />
            </div>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-6 leading-tight text-balance">
              Remain Faithful vs Ever Accountable
            </h1>
            <p className="text-ink-soft text-lg leading-relaxed mb-4 text-pretty">
              Remain Faithful is free. It classifies on your iPhone and never sends screenshots to a partner. Ever Accountable&apos;s public plans page lists $14.99 per month or $129 per year for one reporting user. Their how-it-works page says partners have access to browsing history. A June 6, 2026 post on their site describes an AI on the iPhone that looks at the screen. Those are different sharing models.
            </p>
            <div className="p-5 rounded-sm border border-wax/25 bg-wax/5 text-sm text-ink-soft leading-relaxed">
              <strong className="text-wax">Disclosure:</strong> We built Remain Faithful, so we are biased. We will be transparent about where Ever Accountable is stronger.
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
                    <th className="text-center p-4 text-ink font-semibold bg-paper-deep">Ever Accountable</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    ['Price', 'Free forever', '$14.99/month or $129/year ($10.75/month)'],
                    ['Who the price covers', 'Free forever', 'One reporting user; unlimited devices and partners'],
                    ['Platform', 'iOS 17+ (Android planned)', 'Android, iOS, Windows, Mac, Chromebook, and Kindle'],
                    ['Monitoring approach', 'On-device AI. Screenshots never leave the phone.', 'How-it-works lists Screenshot Accountability and browsing-history reports. A June 6, 2026 post describes on-phone AI for iPhone.'],
                    ['Screenshots shared with partners?', 'Never. See the privacy page.', 'How-it-works says partners have browsing history and will see blacklisted or questionable content. The device table lists Screenshot Accountability on each platform they name.'],
                    ['Where AI runs', 'Entirely on your device', 'Their June 6, 2026 post says an AI on the iPhone looks at the screen.'],
                    ['Open source', 'Yes, full codebase on GitHub', 'No'],
                    ['Accountability model', 'Covenant-based partnership or small group', 'Partner-based report sharing'],
                    ['Professional coach option', 'No', 'Optional. Essential $99/month or Premium $197/month on the plans page.'],
                    ['Installs they publish', 'Available for iPhone (iOS 17+)', '700,000+ installs, stated on their site'],
                    ['DRM streaming monitoring', 'No (iOS limitation)', 'Not stated on the public pages we checked'],
                  ].map(([feature, rf, ea], i) => (
                    <tr key={i} className={`border-b border-hairline ${i % 2 === 0 ? 'bg-paper-deep' : 'bg-paper'}`}>
                      <td className="p-4 text-ink-soft font-medium">{feature}</td>
                      <td className="p-4 text-center text-ink">{rf}</td>
                      <td className="p-4 text-center text-ink-soft">{ea}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          {/* Where EA is stronger */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-ink mb-5">Where Ever Accountable Is Stronger</h2>
            <div className="space-y-4">
              {[
                { title: 'Install base they publish', body: 'Ever Accountable states 700,000+ installs on its public pages. A June 6, 2026 post says they rebuilt the iPhone app four times over six years. Remain Faithful is newer and available for iPhone (iOS 17+).' },
                { title: 'Multi-platform coverage', body: 'Ever Accountable lists Android, iOS, Windows, Mac, Chromebook, and Kindle. Remain Faithful is iPhone only right now. If your accountability partner uses Android, Windows, or Mac, Ever Accountable serves them. Remain Faithful does not yet.' },
                { title: 'Professional accountability coaching', body: 'Their plans page lists optional coaching: Essential at $99 per month (30 minutes) and Premium at $197 per month (2 hours; the 30-minute line is struck through). Remain Faithful is peer-only. It does not offer a coached plan.' },
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
                { title: 'Price', body: 'Remain Faithful is free. On 2026-10-13, Ever Accountable’s public plans page listed $14.99 per month or $129 per year for one reporting user. The annual price is also shown as $10.75 per month, with a 28% discount label. The same page lists optional coaching at $99 per month (Essential) and $197 per month (Premium). Membership includes a 14-day trial link. Verify current amounts on their site.' },
                { title: 'Open source transparency', body: 'Remain Faithful is fully open source. Anyone can read the code and verify that screenshots, OCR text, and raw screen content are never transmitted, and the privacy page says what a partner can receive. Ever Accountable is proprietary. Both say they protect your privacy; only one lets you verify it.' },
                { title: 'What partners receive', body: 'Ever Accountable’s how-it-works page says partners have access to browsing history on phones, computers, and tablets, and that partners will see blacklisted or questionable content. Their device table lists Screenshot Accountability. What a Remain Faithful partner can receive is on the privacy page. Partners never see screenshots, OCR text, or raw screen content.' },
                { title: 'Covenant model', body: 'The covenant framework grounds Remain Faithful\'s accountability in a theological understanding of relationship. Partners agree to respond with grace before gaining access. This is designed to shape how alerts are received.' },
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
                <h2 className="font-serif text-lg font-bold text-ink mb-4">Who Should Choose Ever Accountable</h2>
                <ul className="space-y-2">
                  {[
                    'Need Android, Windows, Mac, Chromebook, or Kindle coverage today',
                    'Want the install base they publish (700,000+)',
                    'Want the optional coaching plans listed on their site',
                    'Are ready to pay the public monthly or yearly price',
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
                    'iPhone users who want free, privacy-first accountability',
                    'Anyone who values open-source code verification',
                    'Churches and small groups needing free group tools',
                    'People who believe partners should not see screenshots',
                    'Those drawn to a covenant-based theological framework',
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
                { q: 'Is Remain Faithful as effective as Ever Accountable?', a: 'Remain Faithful is free and classifies on your iPhone. Screenshots never leave the device. Ever Accountable states 700,000+ installs on its public pages. A June 6, 2026 post on their site describes an AI on the iPhone that looks at the screen. Their how-it-works page says partners have access to browsing history. Which fits depends on what a partner should see, and on price.' },
                { q: 'Does Ever Accountable also do on-device processing?', a: 'Ever Accountable’s June 6, 2026 post says an AI on the iPhone looks at the screen, using the phone’s screen-recording feature rather than making a video anyone watches. That is their description, not a claim that the products match. Remain Faithful is free, never transmits screenshots, and is open source. Their plans page lists $14.99 per month or $129 per year for one reporting user. Their how-it-works page says partners have access to browsing history.' },
                { q: 'Can I switch from Ever Accountable to Remain Faithful?', a: 'Yes. Cancel Ever Accountable and set up Remain Faithful on iPhone (iOS 17+). Invite your partners and walk through the covenant setup. Remain Faithful is available for iPhone only; if your partners use Android, Windows, or Mac, Ever Accountable has broader coverage.' },
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
              <Link href="/compare/accountable2you" className="p-4 rounded-sm border border-hairline bg-paper-deep hover:border-wax/40 transition-colors text-wax font-medium text-sm">
                Remain Faithful vs Accountable2You →
              </Link>
            </div>
          </section>

          <div className="text-center p-10 rounded-sm border border-wax/25" >
            <h2 className="font-serif text-2xl font-bold text-ink mb-4">Try Remain Faithful Free</h2>
            <p className="text-ink-soft mb-6">No subscription. No credit card. Just accountability.</p>
            <Link href="/#download" className="btn-wax">
              Get started
            </Link>
          </div>
        </div>
      </div>
    </>
  )
}
