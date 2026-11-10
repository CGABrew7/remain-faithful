import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { comparisonSchema } from '@/lib/structured-data'
import { COMPARE_PRICING_NOTE } from '@/lib/competitor-pricing'
import { Breadcrumbs } from '@/components/Breadcrumbs'
import { UpdatedStamp } from '@/components/UpdatedStamp'
import { CheckedSources, WhatRemainFaithfulDoes } from '../what-rf-does'

export const metadata: Metadata = {
  title: 'Remain Faithful vs Accountable2You: Direct Comparison (2026)',
  description: 'Remain Faithful is free on iPhone and does not send browsing history or screenshots to a partner. Accountable2You lists Personal at $11/month and Family at $16/month. Rechecked 2026-11-10.',
  alternates: { canonical: 'https://www.remainfaithful.com/compare/accountable2you' },
}

const faqs = [
  {
    q: 'What is the main difference?',
    a: 'Remain Faithful is free on iPhone. It blocks apps you choose, and optional Deep Scan classifies on the phone. A partner gets a short summary, not a screenshot, a page title, or a browsing history. Accountable2You’s pricing page lists Personal at $11 a month or $121 a year, and Family at $16 a month or $176 a year. They say they do not capture screenshots. They do report browsing history, including Incognito, and full web page titles.',
  },
  {
    q: 'Does Accountable2You block websites?',
    a: 'Their pricing FAQ says no. They describe Accountable2You as an accountability and reporting tool that does not block websites or filter content. Remain Faithful can block apps and categories you choose, using Apple Family Controls. Opening a blocked app shows the Screen Time shield and does not itself send an alert.',
  },
  {
    q: 'Can I switch from Accountable2You to Remain Faithful?',
    a: 'You can set up Remain Faithful on an iPhone running iOS 17 or later and invite your partners. Remain Faithful’s current app is iPhone only. Accountable2You’s iOS page covers iPhone and iPad, and the site also lists Windows, Mac, Android, Chromebook, and Linux. Their public pages checked on 2026-11-10 do not describe a required VPN-profile removal step. End the Accountable2You plan from their own account.',
  },
  {
    q: 'Does Remain Faithful send web page titles the way Accountable2You does?',
    a: 'No. Remain Faithful does not log or transmit web page titles or browsing history. A partner can receive time, category, severity, and a short summary. The privacy page says that summary does not name the app. Accountable2You’s features page says reports include full web page titles, and that Incognito and private browsing are visible to partners.',
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
  ['Price', 'Free. No subscription.', 'Personal: $11/month or $121/year ($10.08/month, billed annually). Family: $16/month or $176/year ($14.67/month, billed annually).'],
  ['Who that covers', 'You, plus the partners you invite. A group holds up to 12.', 'Personal: up to 6 active installs and unlimited partners. Family: up to 20 active installs, parent and child accounts, and time-monitoring alerts.'],
  ['Devices', 'iPhone (iOS 17+)', 'iOS page: iPhone and iPad. Site list: Windows, Mac, Android, iOS (iPhone), Chromebook, and Linux.'],
  ['How it watches', 'Apple Family Controls for apps and categories you choose. Optional Deep Scan classifies on the phone. No VPN.', 'iOS page: VPN technology records browser activity in Safari, Chrome, or any browser. They say iOS cannot reliably monitor activity inside non-browser apps.'],
  ['Browsing history and page titles', 'Not sent. Not collected.', 'Yes. Browsing history, including Incognito, and full web page titles. Features page: stored 15 days, with export after that.'],
  ['Screenshots', 'Never leave the phone.', 'Their pricing FAQ says they do not capture screenshots.'],
  ['What a partner receives', 'Time, category, severity, and a short summary. Not a page title.', 'Activity reports and alerts. You can set full reports, alerts only, or specific devices.'],
  ['Blocking', 'Chosen apps and categories, with Apple’s Screen Time shield.', 'Pricing FAQ: they do not block websites or filter content.'],
  ['Source you can read', 'iOS app, backend, and website on GitHub.', 'Not published as open source on the pricing, features, or iOS pages.'],
]

export default function Accountable2YouCompare() {
  return (
    <>
      <JsonLd data={comparisonSchema('Accountable2You')} />
      <JsonLd data={pageFaqSchema} />

      <div className="pt-24 pb-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <Breadcrumbs items={[
            { name: 'Compare', url: 'https://www.remainfaithful.com/compare' },
            { name: 'Accountable2You', url: 'https://www.remainfaithful.com/compare/accountable2you' },
          ]} />

          <div className="mb-8">
            <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
              <p className="text-wax text-sm font-semibold uppercase tracking-widest">Direct Comparison</p>
              <UpdatedStamp />
            </div>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-6 leading-tight text-balance">
              Remain Faithful vs Accountable2You
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
                    <th scope="col" className="text-center p-4 text-ink font-semibold bg-paper-deep">Accountable2You</th>
                  </tr>
                </thead>
                <tbody>
                  {rows.map(([feature, rf, a2y], i) => (
                    <tr key={feature} className={`border-b border-hairline ${i % 2 === 0 ? 'bg-paper-deep' : 'bg-paper'}`}>
                      <th scope="row" className="p-4 text-left text-ink-soft font-medium">{feature}</th>
                      <td className="p-4 text-center text-ink">{rf}</td>
                      <td className="p-4 text-center text-ink-soft">{a2y}</td>
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
                <h3 className="font-semibold text-ink mb-2">Personal, Family, and the other plans</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  The individual pricing page offers a 10-day free trial. Annual cards say GET 1 MONTH FREE. Personal is $11 per month or $121 per year, shown as $10.08 per month billed annually, with up to 6 active installs and unlimited partners. Family is $16 per month or $176 per year, shown as $14.67 per month billed annually, with up to 20 active installs, parent and child accounts, and time-monitoring alerts. The same site lists more plans. The business page lists an organization plan at $50 per month plus $2 per active installation per month, and says the $50 base fee does not include installs. The ministry page lists Foundation at 25 credits for $155 per month, Launch at 40 credits for $240, Growth at 60 credits for $345, and Impact at 100 credits for $550. A personal account uses 1 credit and up to 6 devices. A family account uses 2 credits and up to 20 devices. The school page does not list a price.
                </p>
              </article>
              <article className="p-5 rounded-sm border border-hairline bg-paper-deep">
                <h3 className="font-semibold text-ink mb-2">Reports, screenshots, and blocking</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  The pricing page says reports include internet browsing history, even in Incognito, plus app and device usage. The features page says reports include full web page titles, with color ratings, and that every action is stored for 15 days, with export after that. The pricing FAQ says they do not capture screenshots, and that they do not block websites or filter content. A partner can receive full activity reports, alerts only, or reports from specific devices. The pricing page says you can opt out of recording sensitive apps on Windows, Mac, and Android. The features page says an enhanced app from support can show texts and calls on Android only.
                </p>
              </article>
              <article className="p-5 rounded-sm border border-hairline bg-paper-deep">
                <h3 className="font-semibold text-ink mb-2">iPhone, in their words</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  The iOS page says VPN technology records internet activity in Safari, Chrome, or any browser, and shows the title of each web page. It says that, because of iOS limits, they cannot reliably monitor activity inside non-browser apps, including Facebook, Instagram, and TikTok. Alerts on that page are Reports app notifications and email. The monitoring table marks iOS app usage as limited, YouTube as limited, and social media as in the browser only. Calls and texts are not marked for iOS on that table.
                </p>
              </article>
            </div>
          </section>

          <section className="mb-14">
            <div className="grid md:grid-cols-2 gap-6">
              <div className="p-6 rounded-sm border border-hairline bg-paper-deep">
                <h2 className="font-serif text-lg font-bold text-ink mb-4">Accountable2You fits when</h2>
                <ul className="space-y-2">
                  {[
                    'You need iPad, Android, Windows, Mac, Linux, or Chromebook coverage',
                    'A partner should see browsing history and page titles',
                    'The public Personal, Family, ministry, or organization price fits',
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
                    'You are on iPhone and want on-device classification at no cost',
                    'You do not want browsing history or page titles sent to a partner',
                    'You do not want screenshots sent to a partner',
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
            { href: 'https://accountable2you.com/pricing/', title: 'Accountable2You pricing' },
            { href: 'https://accountable2you.com/features/', title: 'Accountable2You features' },
            { href: 'https://accountable2you.com/ios', title: 'Accountable2You for iPhone and iPad' },
            { href: 'https://accountable2you.com/monitoring/', title: 'What Accountable2You monitors' },
            { href: 'https://accountable2you.com/business/', title: 'Accountable2You organization plan' },
            { href: 'https://accountable2you.com/ministry/', title: 'Accountable2You ministry plans' },
            { href: 'https://accountable2you.com/school/', title: 'Accountable2You for schools' },
          ]} />

          <section className="mb-14">
            <h2 className="font-serif text-xl font-bold text-ink mb-4">Other comparisons</h2>
            <div className="grid sm:grid-cols-2 gap-4">
              <Link href="/compare/covenant-eyes" className="p-4 rounded-sm border border-hairline bg-paper-deep hover:border-wax/40 transition-colors text-wax font-medium text-sm">
                Remain Faithful vs Covenant Eyes
              </Link>
              <Link href="/compare/ever-accountable" className="p-4 rounded-sm border border-hairline bg-paper-deep hover:border-wax/40 transition-colors text-wax font-medium text-sm">
                Remain Faithful vs Ever Accountable
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
