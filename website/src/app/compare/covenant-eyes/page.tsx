import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { comparisonSchema } from '@/lib/structured-data'
import { COMPARE_PRICING_NOTE } from '@/lib/competitor-pricing'
import { Breadcrumbs } from '@/components/Breadcrumbs'
import { UpdatedStamp } from '@/components/UpdatedStamp'
import { CheckedSources, WhatRemainFaithfulDoes } from '../what-rf-does'

export const metadata: Metadata = {
  title: 'Remain Faithful vs Covenant Eyes: Direct Comparison (2026)',
  description: 'Remain Faithful is free on iPhone and keeps screenshots on the device. Covenant Eyes Victory lists $18/month, $198/year, or $950 once, and sends screenshots to an ally. Rechecked 2026-11-10.',
  alternates: { canonical: 'https://www.remainfaithful.com/compare/covenant-eyes' },
}

const faqs = [
  {
    q: 'What is the main difference?',
    a: 'Remain Faithful is free on iPhone. You choose a partner or a group. It can block apps you pick, and optional Deep Scan classifies on the phone. Screenshots stay on the device. Covenant Eyes, on its Victory page, lists $18 a month, $198 a year, or $950 once, and Screen Accountability sends screenshots to an ally. On iPhone, those screenshots are listed for Safari only.',
  },
  {
    q: 'Can I switch from Covenant Eyes to Remain Faithful?',
    a: 'You can set up Remain Faithful on an iPhone running iOS 17 or later and invite your partners. Remain Faithful does not cover Android, Windows, or Mac. If those devices still need a report, Covenant Eyes lists them and this app does not. End the Covenant Eyes plan from their own account.',
  },
  {
    q: 'Does Remain Faithful work on the same devices as Covenant Eyes?',
    a: 'No. The Victory page lists Screen Accountability and blocking on iOS 15.1+, Mac 11+, Android 11+, and Windows 10+. It says Victory does not provide porn blocking and Screen Accountability on Chromebooks, Kindle devices and tablets, Smart TVs, VR headsets, or Linux. Remain Faithful’s current app is iPhone (iOS 17+).',
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
  ['Price', 'Free. No subscription.', '$18/month, $198/year, or $950 once, on the Victory page web checkout'],
  ['Who that covers', 'You, plus the partners you invite. A group holds up to 12.', 'Unlimited devices for up to 10 users in a family. Allies install the Victory app free.'],
  ['Devices', 'iPhone (iOS 17+)', 'iOS 15.1+, Mac 11+, Android 11+, Windows 10+. Not Chromebooks, Kindle devices and tablets, Smart TVs, VR headsets, or Linux.'],
  ['How it watches', 'Apple Family Controls for apps and categories you choose. Optional Deep Scan classifies on the phone.', 'Screen Accountability sends screenshots to an ally. On iOS, screenshots are listed for Safari only. The same iOS list also says it reviews, reports, and blocks explicit domains across apps.'],
  ['What a partner receives', 'Time, category, severity, and a short summary. Not a screenshot. The summary does not name the app.', 'Screenshots. Under the heading "Confident Security," the Victory page says, "Screenshots are blurred, distorted, and encrypted with AES 256-bit security." The page does not say how long those images are kept.'],
  ['Blocking', 'Chosen apps and categories, with Apple’s Screen Time shield. Opening a blocked app does not itself send an alert.', 'Explicit domains, plus custom block and allow lists. On iOS they say to add Screen Time limits and app restrictions.'],
  ['Source you can read', 'iOS app, backend, and website on GitHub.', 'Not published as open source on the Victory page.'],
]

export default function CovenantEyesCompare() {
  return (
    <>
      <JsonLd data={comparisonSchema('Covenant Eyes')} />
      <JsonLd data={pageFaqSchema} />

      <div className="pt-24 pb-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <Breadcrumbs items={[
            { name: 'Compare', url: 'https://www.remainfaithful.com/compare' },
            { name: 'Covenant Eyes', url: 'https://www.remainfaithful.com/compare/covenant-eyes' },
          ]} />

          <div className="mb-8">
            <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
              <p className="text-wax text-sm font-semibold uppercase tracking-widest">Direct Comparison</p>
              <UpdatedStamp />
            </div>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-6 leading-tight text-balance">
              Remain Faithful vs Covenant Eyes
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
                    <th scope="col" className="text-center p-4 text-ink font-semibold bg-paper-deep">Covenant Eyes</th>
                  </tr>
                </thead>
                <tbody>
                  {rows.map(([feature, rf, ce], i) => (
                    <tr key={feature} className={`border-b border-hairline ${i % 2 === 0 ? 'bg-paper-deep' : 'bg-paper'}`}>
                      <th scope="row" className="p-4 text-left text-ink-soft font-medium">{feature}</th>
                      <td className="p-4 text-center text-ink">{rf}</td>
                      <td className="p-4 text-center text-ink-soft">{ce}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-ink mb-5">What the Victory page says</h2>
            <div className="space-y-4">
              <article className="p-5 rounded-sm border border-hairline bg-paper-deep">
                <h3 className="font-semibold text-ink mb-2">Price, trial, and guarantee</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  The Victory page lists $18 per month, $198 per year, and a $950 one-time lifetime payment. Monthly and yearly show 7 days free. A 30-day money-back line sits under the prices. Its tooltip says to cancel within the first 30 days and request a refund, and that in-app purchases and reactivated accounts are not eligible. These are the web checkout amounts. This page does not add a separate in-app price.
                </p>
              </article>
              <article className="p-5 rounded-sm border border-hairline bg-paper-deep">
                <h3 className="font-semibold text-ink mb-2">Two lifetime savings sentences</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  The Victory page says the $950 lifetime price saves 50% over 10 years. The lifetime membership page says lifetime saves 21% after 5 years and over 60% after 10 years. Both sentences are theirs. This page does not pick a third figure.
                </p>
              </article>
              <article className="p-5 rounded-sm border border-hairline bg-paper-deep">
                <h3 className="font-semibold text-ink mb-2">Screenshots and devices</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  Screen Accountability takes screenshots and sends them to an ally on iOS, Mac, Android, and Windows. On iOS, the screenshot line says Safari only. The same iOS list says it reviews, reports, and blocks explicit domain activity across apps, with a custom block and allow list, and tells people to couple that with Screen Time limits and app restrictions. Mac, Android, and Windows list screenshots, explicit-domain blocking, and custom block and allow lists. Under the heading &quot;Confident Security,&quot; the Victory page says, &quot;Screenshots are blurred, distorted, and encrypted with AES 256-bit security.&quot; The Victory page does not say when those images are deleted. Allies install Victory for free. The page says blocking and Screen Accountability are not provided on Chromebooks, Kindle devices and tablets, Smart TVs, VR headsets, or Linux.
                </p>
              </article>
              <article className="p-5 rounded-sm border border-hairline bg-paper-deep">
                <h3 className="font-semibold text-ink mb-2">Figures they publish</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  The Victory page menu says the mission launched 25 years ago. The same page says to join over 1.7 million people who have used Covenant Eyes. Those are their figures, not a count we made.
                </p>
              </article>
            </div>
          </section>

          <section className="mb-14">
            <div className="grid md:grid-cols-2 gap-6">
              <div className="p-6 rounded-sm border border-hairline bg-paper-deep">
                <h2 className="font-serif text-lg font-bold text-ink mb-4">Covenant Eyes fits when</h2>
                <ul className="space-y-2">
                  {[
                    'You need Mac, Android, or Windows today, not only iPhone',
                    'You want one paid account for up to 10 users in a family',
                    'You want an ally to see screenshots',
                    'You want the history and the 1.7 million figure they state',
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
                    'You want screenshots to stay on the phone',
                    'You want a partner alert, not an image of the screen',
                    'You want to read the iOS app, backend, and website source',
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
            { href: 'https://www.covenanteyes.com/victory/', title: 'Victory by Covenant Eyes' },
            { href: 'https://www.covenanteyes.com/victory/lifetime-membership/', title: 'Victory lifetime membership' },
          ]} />

          <section className="mb-14">
            <h2 className="font-serif text-xl font-bold text-ink mb-4">Other comparisons</h2>
            <div className="grid sm:grid-cols-2 gap-4">
              <Link href="/compare/ever-accountable" className="p-4 rounded-sm border border-hairline bg-paper-deep hover:border-wax/40 transition-colors text-wax font-medium text-sm">
                Remain Faithful vs Ever Accountable
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
