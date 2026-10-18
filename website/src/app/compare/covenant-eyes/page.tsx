import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { comparisonSchema } from '@/lib/structured-data'
import { COMPARE_PRICING_NOTE } from '@/lib/competitor-pricing'
import { Breadcrumbs } from '@/components/Breadcrumbs'
import { UpdatedStamp } from '@/components/UpdatedStamp'

export const metadata: Metadata = {
  title: 'Remain Faithful vs Covenant Eyes: Honest Comparison (2026)',
  description: 'A detailed comparison of Remain Faithful and Covenant Eyes for Christian accountability. See how they differ on price, privacy, monitoring, and approach.',
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
        text: 'Remain Faithful is free and classifies on your iPhone. Screenshots never leave the device. Covenant Eyes, on its public Victory page, describes screenshot-based Screen Accountability on iOS, Android, Mac, and Windows, and says the mission started 25 years ago. The products take different approaches. Which one fits depends on the devices you need and whether a partner should see screenshots.',
      },
    },
    {
      '@type': 'Question',
      name: 'Can I switch from Covenant Eyes to Remain Faithful?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes. Cancel Covenant Eyes and set up Remain Faithful on iPhone (iOS 17+), then invite your partners. Remain Faithful is available for iPhone only. If your partners use Android, Windows, or Mac, Covenant Eyes covers those platforms and Remain Faithful does not yet.',
      },
    },
    {
      '@type': 'Question',
      name: 'Does Remain Faithful work on the same devices as Covenant Eyes?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'No. Covenant Eyes lists Screen Accountability on iOS 15.1+, Android 11+, Mac 11+, and Windows 10+. Their page says Screen Accountability is not available on Chromebooks, Kindle Fires, or Smart TVs. Remain Faithful is available for iPhone (iOS 17+). Android is planned, with no launch date. If you need those other platforms today, Covenant Eyes covers them and Remain Faithful does not.',
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
            <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
              <p className="text-wax text-sm font-semibold uppercase tracking-widest">Honest Comparison</p>
              <UpdatedStamp />
            </div>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-6 leading-tight text-balance">
              Remain Faithful vs Covenant Eyes
            </h1>
            <p className="text-ink-soft text-lg leading-relaxed mb-4 text-pretty">
              Remain Faithful is free. It classifies on your iPhone and never sends screenshots to a partner. What a partner can receive is on the privacy page. Covenant Eyes, on its public Victory page, describes screenshot-based Screen Accountability sent to an ally. This page is about that difference. It is not a claim that the products match.
            </p>
            <div className="p-5 rounded-sm border border-wax/25 bg-wax/5 text-sm text-ink-soft leading-relaxed">
              <strong className="text-wax">Disclosure:</strong> We built Remain Faithful, so we are biased. We will be transparent about where Covenant Eyes is stronger.
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
                    <th className="text-center p-4 text-ink font-semibold bg-paper-deep">Covenant Eyes</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    ['Price', 'Free forever', '$18/month, $198/year, or $950 lifetime'],
                    ['Who the price covers', '$0', 'One account: unlimited devices for up to 10 users in a family'],
                    ['Platform', 'iOS 17+ (Android planned)', 'iOS 15.1+, Android 11+, Mac 11+, Windows 10+'],
                    ['Monitoring approach', 'On-device AI classification', 'Screen Accountability: screenshots sent to an ally'],
                    ['Screenshots shared with partners?', 'Never. See the privacy page.', 'Yes. They say screenshots are sent to an ally, blurred, distorted, and encrypted, then deleted after 30 days. On iOS they list Safari only.'],
                    ['Where classification runs', 'Entirely on your device', 'Screenshots are captured and sent to an ally'],
                    ['Open source', 'Yes, full codebase on GitHub', 'No'],
                    ['Accountability model', 'Covenant-based partnership or small group', 'Ally-based reports'],
                    ['Household coverage', 'Built-in group mode (up to 12)', 'Public plans cover up to 10 users on one account'],
                    ['Content blocking', 'Yes, always-on Family Controls (adult self-monitor)', 'Yes. They list blocking explicit domains and custom block/allow lists.'],
                    ['DRM streaming monitoring', 'No (iOS limitation)', 'Not stated on their public Victory page'],
                    ['Established since', '2026', 'Their site says the mission started 25 years ago'],
                    ['User base', 'Available for iPhone (iOS 17+)', '1.7 million people, as stated on their site'],
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
                { title: 'Platform coverage', body: 'Covenant Eyes lists Screen Accountability and blocking on iOS 15.1+, Mac 11+, Android 11+, and Windows 10+. Their page says Screen Accountability is not available on Chromebooks, Kindle Fires, or Smart TVs. Remain Faithful is iPhone only (iOS 17+) right now. If you or your partners use those other platforms, Covenant Eyes covers them and Remain Faithful does not.' },
                { title: 'Track record', body: 'Their site says the mission launched 25 years ago and cites 1.7 million people who have used Covenant Eyes. Remain Faithful is newer and available for iPhone (iOS 17+).' },
                { title: 'Content blocking', body: 'Their device list includes blocking explicit domains and custom block/allow lists. On iOS they also say to couple that with Screen Time Limits and App Restrictions. Remain Faithful uses Apple Family Controls in the adult self-monitor (.individual) authorization to block chosen apps and categories.' },
                { title: 'One account for a household', body: 'The public Victory pricing page says every plan covers unlimited devices for up to 10 individual users within a family. Remain Faithful includes group mode for up to 12 people and does not charge for it. That page does not list a separate enterprise price.' },
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
                { title: 'Price', body: 'Remain Faithful is free. On 2026-10-13, the public Victory pricing page listed $18 per month, $198 per year (marked Popular), and a $950 one-time Lifetime payment (marked “Save 50% over 10 years”). The same page says all plans cover unlimited devices for up to 10 individual users within a family. That is an account price, not a per-person fee. The page also says to try it for 30 days, and a note excludes some in-app purchases from the money-back guarantee. Verify current amounts on their site.' },
                { title: 'Privacy: screenshots stay on your iPhone', body: 'Remain Faithful classifies on-device. Screenshots, OCR text, and raw screen content are never transmitted. What a partner can receive is on the privacy page. Covenant Eyes says Victory Shield takes screenshots and sends them to an ally. On iOS they list Safari only. They also say those images are blurred, distorted, and encrypted with AES 256-bit security, then deleted after 30 days.' },
                { title: 'Partner protection', body: 'Covenant Eyes still sends screenshots to an ally, even when those images are described as blurred or distorted. What a Remain Faithful partner can receive is on the privacy page. Partners never see screenshots, OCR text, or raw screen content.' },
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
                    'Want one paid account for up to 10 users in a family',
                    'Want the multi-platform coverage and the 25-year history they state on their site',
                    'Partners who prefer to see screenshots',
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
                { q: 'Is Remain Faithful as effective as Covenant Eyes?', a: 'Remain Faithful is free and classifies on your iPhone. Screenshots never leave the device. Covenant Eyes, on its public Victory page, describes screenshot-based Screen Accountability on iOS, Android, Mac, and Windows, and says the mission started 25 years ago. The products take different approaches. Which one fits depends on the devices you need and whether a partner should see screenshots.' },
                { q: 'Can I switch from Covenant Eyes to Remain Faithful?', a: 'Yes. Cancel Covenant Eyes and set up Remain Faithful on iPhone (iOS 17+), then invite your partners. Remain Faithful is available for iPhone only. If your partners use Android, Windows, or Mac, Covenant Eyes covers those platforms and Remain Faithful does not yet.' },
                { q: 'Does Remain Faithful work on the same devices as Covenant Eyes?', a: 'No. Covenant Eyes lists Screen Accountability on iOS 15.1+, Android 11+, Mac 11+, and Windows 10+. Their page says Screen Accountability is not available on Chromebooks, Kindle Fires, or Smart TVs. Remain Faithful is available for iPhone (iOS 17+). Android is planned, with no launch date. If you need those other platforms today, Covenant Eyes covers them and Remain Faithful does not.' },
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
            <Link href="/#download" className="btn-wax">
              Get started
            </Link>
          </div>
        </div>
      </div>
    </>
  )
}
