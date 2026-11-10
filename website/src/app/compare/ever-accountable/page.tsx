import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { comparisonSchema } from '@/lib/structured-data'
import { COMPARE_PRICING_NOTE } from '@/lib/competitor-pricing'
import { Breadcrumbs } from '@/components/Breadcrumbs'
import { UpdatedStamp } from '@/components/UpdatedStamp'
import { CheckedSources, WhatRemainFaithfulDoes } from '../what-rf-does'

export const metadata: Metadata = {
  title: 'Remain Faithful vs Ever Accountable: Direct Comparison (2026)',
  description: 'Remain Faithful is free on iPhone and keeps screenshots on the device. Ever Accountable lists $14.99/month or $129/year for one reporting user. Rechecked 2026-11-10.',
  alternates: { canonical: 'https://www.remainfaithful.com/compare/ever-accountable' },
}

const faqs = [
  {
    q: 'What is the main difference?',
    a: 'Remain Faithful is free on iPhone. Optional Deep Scan classifies on the phone, and screenshots stay on the device. A partner gets a short summary, not a screenshot and not a browsing history. Ever Accountable’s plans page lists $14.99 a month or $129 a year for one reporting user. Their how-it-works page says partners have access to browsing history. Their plans table checks Screenshot Accountability on every platform they list. A June 6, 2026 post says an AI on the iPhone looks at the screen. Those statements describe different sharing, and this page does not merge them into one claim.',
  },
  {
    q: 'Does Ever Accountable also classify on the iPhone?',
    a: 'Their June 6, 2026 post says they use the iPhone’s screen-recording feature, not to make a video anyone watches, but so an AI on the phone can look at the screen. Their iOS explainer says App Unlock uses an AI to watch the screen for sexualized content, and that it is not recording a video. Their plans table still checks Screenshot Accountability for iPhone and iPad. Their how-it-works page still says partners have browsing history. Remain Faithful is free, does not transmit screenshots, and publishes the iOS app, backend, and website source.',
  },
  {
    q: 'Can I switch from Ever Accountable to Remain Faithful?',
    a: 'You can set up Remain Faithful on an iPhone running iOS 17 or later and invite your partners. Remain Faithful does not cover Android, Windows, Mac, Chromebook, or Kindle. If those devices still need a report, Ever Accountable lists them and this app does not. End the Ever Accountable plan from their own account.',
  },
]

const pageFaqSchema = {
  '@context': 'https://schema.org',
  '@type': 'FAQPage',
  mainEntity: faqs.map((faq) => ({
    '@type': 'Question',
    name: faq.q,
    acceptedAnswer: { '@type': 'Answer', text: faq.a },
  })),
}

const rows = [
  ['Price', 'Free. No subscription.', '$14.99/month or $129/year. The annual card says $10.75/month and a 28% discount.'],
  ['Who that covers', 'You, plus the partners you invite. A group holds up to 12.', 'The plan card says one reporting user, unlimited monitored devices, and unlimited accountability partners.'],
  ['Devices', 'iPhone (iOS 17+)', 'Plans table: Android, iPhone/iPad, Windows, Mac, ChromeOS, and Kindle. The FAQ also names Chromebook.'],
  ['How it watches', 'Apple Family Controls for apps and categories you choose. Optional Deep Scan classifies on the phone.', 'How-it-works says partners have browsing history. The plans table checks Screenshot Accountability on each listed platform. A June 6, 2026 post describes on-phone AI for iPhone.'],
  ['What a partner receives', 'Time, category, severity, and a short summary. Not a screenshot and not a browsing history.', 'How-it-works says partners see blacklisted or questionable content and have browsing history. The plans table checks Screenshot Accountability. The June 6 post says the iPhone AI is not a video anyone watches.'],
  ['Blocking', 'Chosen apps and categories, with Apple’s Screen Time shield.', 'Plans table: checked for Android, iPhone/iPad, and Mac. Not checked for Windows, ChromeOS, or Kindle. Their FAQ says filtering on Android and iOS can be turned on at no extra charge, and that blocking is off by default.'],
  ['Coaching', 'No coached plan. Peer accountability.', 'Optional. Essential is $99/month. Premium is $197/month.'],
  ['Source you can read', 'iOS app, backend, and website on GitHub.', 'Not published as open source on the plans or how-it-works pages.'],
]

export default function EverAccountableCompare() {
  return (
    <>
      <JsonLd data={comparisonSchema('Ever Accountable')} />
      <JsonLd data={pageFaqSchema} />

      <div className="pt-24 pb-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <Breadcrumbs items={[
            { name: 'Compare', url: 'https://www.remainfaithful.com/compare' },
            { name: 'Ever Accountable', url: 'https://www.remainfaithful.com/compare/ever-accountable' },
          ]} />

          <div className="mb-8">
            <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
              <p className="text-wax text-sm font-semibold uppercase tracking-widest">Direct Comparison</p>
              <UpdatedStamp />
            </div>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-6 leading-tight text-balance">
              Remain Faithful vs Ever Accountable
            </h1>
          </div>

          <section className="mb-8 p-5 rounded-sm border border-wax/25 bg-wax/5">
            <h2 className="font-serif text-xl font-bold text-ink mb-2">Short answer</h2>
            <p className="text-ink-soft leading-relaxed text-pretty">{faqs[0].a}</p>
          </section>

          <div className="mb-10 p-5 rounded-sm border border-hairline bg-paper-deep text-sm text-ink-soft leading-relaxed">
            <strong className="text-wax">Disclosure:</strong> We built Remain Faithful. The prices and features below are what each public site said on the date in the stamp.
          </div>

          <WhatRemainFaithfulDoes />

          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-ink mb-6">Side-by-side</h2>
            <p className="text-xs text-ink-faint mb-4">{COMPARE_PRICING_NOTE}</p>
            <div className="overflow-x-auto rounded-sm border border-hairline">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-hairline">
                    <th scope="col" className="text-left p-4 text-ink-soft font-semibold bg-paper-deep w-1/3">Feature</th>
                    <th scope="col" className="text-center p-4 text-wax font-semibold bg-paper-deep">Remain Faithful</th>
                    <th scope="col" className="text-center p-4 text-ink font-semibold bg-paper-deep">Ever Accountable</th>
                  </tr>
                </thead>
                <tbody>
                  {rows.map(([feature, rf, ea], i) => (
                    <tr key={feature} className={`border-b border-hairline ${i % 2 === 0 ? 'bg-paper-deep' : 'bg-paper'}`}>
                      <th scope="row" className="p-4 text-left text-ink-soft font-medium">{feature}</th>
                      <td className="p-4 text-center text-ink">{rf}</td>
                      <td className="p-4 text-center text-ink-soft">{ea}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-ink mb-5">What their public pages say</h2>
            <div className="space-y-4">
              <article className="p-5 rounded-sm border border-hairline bg-paper-deep">
                <h3 className="font-semibold text-ink mb-2">Price, trial, and coaching</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  The plans page lists $14.99 per month or $129 per year for one reporting user. The annual card also says $10.75 per month and a 28% discount. Both cards include unlimited monitored devices, unlimited accountability partners, and text and email alerts. The buttons say a 14-day free trial. A plans FAQ says they will refund payments for up to 90 days. Coaching is separate. Essential is $99 per month: one call, 30 minutes, one weekly check-in, and a weekly report review. Premium is $197 per month: two calls, 60 minutes each, twice-weekly check-ins, a weekly report review, and priority support after relapse. The Premium card strikes through 30 minutes and shows 2 hours. It is marked Best value.
                </p>
              </article>
              <article className="p-5 rounded-sm border border-hairline bg-paper-deep">
                <h3 className="font-semibold text-ink mb-2">One reporting user, and the device table</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  The plan card says one reporting user. The FAQ says the subscription covers all of your devices, with no device limit, so that you and your loved ones have no reason not to be protected. This page keeps the card’s “one reporting user” as the plan definition. The table columns are Android, iPhone/iPad, Windows, Mac, ChromeOS, and Kindle. The FAQ names Android, iOS, Windows, MacOS, Chromebook, and Kindle. A footer line says 5 platforms. This page uses the named list, not the footer count. Screenshot Accountability is checked in every column. “Shares content viewed inside apps” is checked for Android, Mac, and Kindle. The iPhone/iPad cell says Most Apps. The Windows cell says Common Apps. ChromeOS is not checked for in-app content. Blocking is checked for Android, iPhone/iPad, and Mac, and not for Windows, ChromeOS, or Kindle.
                </p>
              </article>
              <article className="p-5 rounded-sm border border-hairline bg-paper-deep">
                <h3 className="font-semibold text-ink mb-2">iPhone, in their words</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  The June 6, 2026 post says they rebuilt the iPhone app four times over six years. It says protection for a social app, or a browser other than Safari, starts with a button and takes about four seconds. It says the AI uses about 1% of the battery while it is on. The iOS explainer describes three layers you choose: accountability, on-screen warnings, and optional blocking. Protected apps stay locked until App Unlock. App Unlock uses an AI to watch the screen. Taking an app off the protected list is not an instant alert. Their explainer says time in unprotected apps that are known to contain sexualized content is included in the weekly report. The plans FAQ says filtering on Android and iOS can be turned on at no extra charge, that it is off by default, and that it uses safe search and filters known porn, malware, and spyware sites. The plans page footer says 700,000+ installs. That is their figure.
                </p>
              </article>
              <article className="p-5 rounded-sm border border-hairline bg-paper-deep">
                <h3 className="font-semibold text-ink mb-2">What we are not combining</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  How-it-works says partners have browsing history and will see blacklisted or questionable content. The plans table checks Screenshot Accountability, including on iPhone and iPad. The June 6 post and the iOS explainer say the on-phone AI is not a video anyone watches. Those pages are all public. This comparison does not decide which image, if any, a partner is shown.
                </p>
              </article>
            </div>
          </section>

          <section className="mb-14">
            <div className="grid md:grid-cols-2 gap-6">
              <div className="p-6 rounded-sm border border-hairline bg-paper-deep">
                <h2 className="font-serif text-lg font-bold text-ink mb-4">Ever Accountable fits when</h2>
                <ul className="space-y-2">
                  {[
                    'You need Android, Windows, Mac, Chromebook, or Kindle today',
                    'You want the install count they publish (700,000+)',
                    'You want the optional coaching plans on their plans page',
                    'You are ready to pay the public monthly or yearly price',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2 text-sm text-ink-soft">
                      <svg className="flex-shrink-0 mt-0.5 text-ink-soft" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><polyline points="20 6 9 17 4 12"/></svg>
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
              <div className="p-6 rounded-sm border border-wax/25 bg-paper-deep">
                <h2 className="font-serif text-lg font-bold text-ink mb-4">Remain Faithful fits when</h2>
                <ul className="space-y-2">
                  {[
                    'You are on iPhone and want protection with no subscription',
                    'You want screenshots and browsing history to stay off the partner report',
                    'You want to read the source',
                    'You want a group of up to 12 people included',
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

          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-ink mb-6">Frequently asked questions</h2>
            <div className="space-y-4">
              {faqs.map((faq) => (
                <div key={faq.q} className="rounded-sm border border-hairline bg-paper-deep p-6">
                  <h3 className="font-semibold text-ink mb-3">{faq.q}</h3>
                  <p className="text-sm text-ink-soft leading-relaxed">{faq.a}</p>
                </div>
              ))}
            </div>
          </section>

          <CheckedSources links={[
            { href: 'https://everaccountable.com/plans/', title: 'Ever Accountable plans' },
            { href: 'https://everaccountable.com/how-it-works/', title: 'Ever Accountable how it works' },
            { href: 'https://everaccountable.com/blog/how-ever-accountable-sees-inside-any-iphone-app/', title: 'June 6, 2026 post on iPhone protection' },
            { href: 'https://everaccountable.com/ios-nextgen-explainer/', title: 'Ever Accountable iOS explainer' },
          ]} />

          <section className="mb-14">
            <h2 className="font-serif text-xl font-bold text-ink mb-4">Other comparisons</h2>
            <div className="grid sm:grid-cols-2 gap-4">
              <Link href="/compare/covenant-eyes" className="p-4 rounded-sm border border-hairline bg-paper-deep hover:border-wax/40 transition-colors text-wax font-medium text-sm">
                Remain Faithful vs Covenant Eyes
              </Link>
              <Link href="/compare/accountable2you" className="p-4 rounded-sm border border-hairline bg-paper-deep hover:border-wax/40 transition-colors text-wax font-medium text-sm">
                Remain Faithful vs Accountable2You
              </Link>
            </div>
          </section>

          <div className="text-center p-10 rounded-sm border border-wax/25">
            <h2 className="font-serif text-2xl font-bold text-ink mb-4">Try Remain Faithful free</h2>
            <p className="text-ink-soft mb-6">No subscription. No credit card. Accountability on your iPhone.</p>
            <Link href="/#download" className="btn-wax">
              Get started
            </Link>
          </div>
        </div>
      </div>
    </>
  )
}
