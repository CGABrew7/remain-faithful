import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { softwareApplicationSchema } from '@/lib/structured-data'
import { Breadcrumbs } from '@/components/Breadcrumbs'

export const metadata: Metadata = {
  title: 'Free Accountability App for iPhone | No Subscription, No Ads',
  description: 'Remain Faithful is a free accountability app for iPhone. No subscription tiers. No premium features. No ads. On-device AI, privacy-first, built for Christians serious about purity.',
  alternates: { canonical: 'https://remainfaithful.com/free-accountability-app' },
}

const faqSchema = {
  '@context': 'https://schema.org',
  '@type': 'FAQPage',
  mainEntity: [
    {
      '@type': 'Question',
      name: 'Why is Remain Faithful free?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Remain Faithful is funded by voluntary donations through the Woodfield Foundation, a registered 501(c)(3) nonprofit. The founder believes no one should face a financial barrier to accountability. Donors who find the app valuable choose to support it; those who cannot afford to donate use it free of charge.',
      },
    },
    {
      '@type': 'Question',
      name: 'Will Remain Faithful stay free?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes. The commitment is indefinite. There will never be a subscription tier, a premium version, or paywalled features. The donation model is designed specifically to avoid creating financial incentives that would conflict with keeping the app free.',
      },
    },
    {
      '@type': 'Question',
      name: 'How does Remain Faithful make money if it is free?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'It does not make money. Remain Faithful is a nonprofit project sustained by voluntary donations. Server costs run approximately $80 per month. Donors who find the app valuable contribute to cover those costs. The founder personally funds any shortfall.',
      },
    },
    {
      '@type': 'Question',
      name: 'Is there a catch?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'No. No ads, no data selling, no premium upsells. The codebase is open source so anyone can verify this. The only ask is an optional voluntary donation if you find the app valuable.',
      },
    },
  ],
}

export default function FreeAccountabilityApp() {
  return (
    <>
      <JsonLd data={softwareApplicationSchema} />
      <JsonLd data={faqSchema} />

      <div className="pt-24 pb-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <Breadcrumbs items={[{ name: 'Free Accountability App', url: 'https://remainfaithful.com/free-accountability-app' }]} />

          {/* Hero */}
          <div className="mb-14">
            <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-4">100% Free, Forever</p>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-6 leading-tight">
              The Only Free Accountability App That Actually Works
            </h1>
            <p className="text-[#8A9BB0] text-lg leading-relaxed mb-6">
              Remain Faithful is a free accountability app for iPhone. No subscription. No premium features. No advertising. On-device AI keeps your screen content private. Built for Christians serious about purity.
            </p>
            <Link href="/#waitlist" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96] shadow-lg shadow-[#C9A84C]/20">
              Join the Waitlist (Free)
            </Link>
          </div>

          {/* Why "free" usually isn't */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Why Most "Free" Accountability Apps Aren&apos;t Really Free</h2>
            <p className="text-[#8A9BB0] leading-relaxed mb-4">
              Most apps that advertise as free use one of three models, and none of them are actually free in any meaningful sense.
            </p>
            <div className="space-y-4">
              {[
                { label: 'Freemium', body: 'The app is free but the features that actually work are locked behind a subscription. You download it, discover it is limited, and pay. This is the most common model in the accountability app space.' },
                { label: 'Free trial', body: 'The app is free for 7 or 14 days, then automatically charges. You have to remember to cancel. Many people do not.' },
                { label: 'Ad-supported', body: 'The app is free because your attention is the product. Ads are served to you. For a purity accountability app, this creates an obvious problem: the ads themselves may be the temptation.' },
              ].map((item) => (
                <div key={item.label} className="p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
                  <h3 className="font-semibold text-[#F0EDE8] mb-2">{item.label}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.body}</p>
                </div>
              ))}
            </div>
            <p className="text-[#8A9BB0] leading-relaxed mt-6">
              Remain Faithful is none of these. It is free in the same way that a library is free: because someone decided the service should be accessible to everyone, and funded it accordingly.
            </p>
          </section>

          {/* How it is funded */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">How Remain Faithful Is Funded</h2>
            <p className="text-[#8A9BB0] leading-relaxed mb-4">
              Remain Faithful is a project of the Woodfield Foundation, a registered 501(c)(3) nonprofit organization. Donations are tax-deductible. The model is simple: people who find the app valuable donate voluntarily. People who cannot afford to donate use it free.
            </p>
            <div className="grid sm:grid-cols-3 gap-4 mb-6">
              {[
                { label: 'Monthly server costs', value: '~$80', note: 'Go backend, database, push notifications' },
                { label: 'Development', value: 'Volunteer', note: 'Founder donates time' },
                { label: 'Required donation', value: '$0', note: 'Never required, always optional' },
              ].map((item) => (
                <div key={item.label} className="p-5 rounded-xl border border-[#1E3050] bg-[#162235] text-center">
                  <p className="text-xs text-[#8A9BB0] uppercase tracking-wide mb-1">{item.label}</p>
                  <p className="font-serif font-bold text-[#C9A84C] text-2xl mb-1">{item.value}</p>
                  <p className="text-xs text-[#8A9BB0]">{item.note}</p>
                </div>
              ))}
            </div>
          </section>

          {/* Full feature list */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Everything Free. No Asterisk.</h2>
            <div className="grid sm:grid-cols-2 gap-3">
              {[
                'Always-on Screen Time monitoring',
                'On-device AI screen analysis (Deep Scan)',
                'One-to-one accountability partnerships',
                'Group accountability mode (up to 12 members)',
                'Covenant-based partner system',
                'Discreet alert notifications to partners',
                'App and category blocking',
                'Open source codebase on GitHub',
                'No ads, ever',
                'No data selling, ever',
                'No premium tier, ever',
                'Privacy-first architecture',
              ].map((feature) => (
                <div key={feature} className="flex items-center gap-3 p-3 rounded-lg border border-[#1E3050] bg-[#162235]">
                  <svg className="flex-shrink-0 text-[#C9A84C]" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><polyline points="20 6 9 17 4 12"/></svg>
                  <span className="text-sm text-[#8A9BB0]">{feature}</span>
                </div>
              ))}
            </div>
          </section>

          {/* Cost comparison */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">What You&apos;d Pay for Equivalent Features Elsewhere</h2>
            <div className="overflow-x-auto rounded-2xl border border-[#1E3050]">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-[#1E3050] bg-[#0A1420]">
                    <th className="text-left p-4 text-[#8A9BB0] font-semibold">App</th>
                    <th className="text-center p-4 text-[#8A9BB0] font-semibold">Annual Cost (1 person)</th>
                    <th className="text-center p-4 text-[#8A9BB0] font-semibold">Annual Cost (family of 4)</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    ['Remain Faithful', '$0', '$0'],
                    ['Covenant Eyes', '~$204', '~$816'],
                    ['Ever Accountable', '~$99', '~$396'],
                    ['Accountable2You', '~$80', '~$320'],
                  ].map(([app, ind, fam], i) => (
                    <tr key={i} className={`border-b border-[#1E3050] ${i === 0 ? 'bg-[#C9A84C]/5' : i % 2 === 0 ? 'bg-[#162235]' : 'bg-[#0F1B2D]'}`}>
                      <td className={`p-4 font-medium ${i === 0 ? 'text-[#C9A84C]' : 'text-[#8A9BB0]'}`}>{app}</td>
                      <td className={`p-4 text-center ${i === 0 ? 'text-[#C9A84C] font-bold' : 'text-[#8A9BB0]'}`}>{ind}</td>
                      <td className={`p-4 text-center ${i === 0 ? 'text-[#C9A84C] font-bold' : 'text-[#8A9BB0]'}`}>{fam}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            <p className="text-xs text-[#8A9BB0]/60 mt-2">Pricing as of June 2026. Check each provider&apos;s website for current pricing.</p>
          </section>

          {/* FAQ */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-6">Common Questions About the Free Model</h2>
            <div className="space-y-4">
              {[
                { q: 'Why is Remain Faithful free?', a: 'Remain Faithful is funded by voluntary donations through the Woodfield Foundation, a registered 501(c)(3) nonprofit. The founder believes no one should face a financial barrier to accountability. Donors who find the app valuable choose to support it; those who cannot afford to donate use it free of charge.' },
                { q: 'Will Remain Faithful stay free?', a: 'Yes. The commitment is indefinite. There will never be a subscription tier, a premium version, or paywalled features. The donation model is designed specifically to avoid creating financial incentives that would conflict with keeping the app free.' },
                { q: 'How does Remain Faithful make money if it is free?', a: 'It does not make money. Remain Faithful is a nonprofit project sustained by voluntary donations. Server costs run approximately $80 per month. Donors who find the app valuable contribute to cover those costs. The founder personally funds any shortfall.' },
                { q: 'Is there a catch?', a: 'No. No ads, no data selling, no premium upsells. The codebase is open source so anyone can verify this. The only ask is an optional voluntary donation if you find the app valuable.' },
              ].map((faq) => (
                <div key={faq.q} className="rounded-2xl border border-[#1E3050] bg-[#162235] p-6">
                  <h3 className="font-semibold text-[#F0EDE8] mb-3">{faq.q}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{faq.a}</p>
                </div>
              ))}
            </div>
          </section>

          {/* CTA */}
          <div className="text-center p-10 rounded-3xl border border-[#C9A84C]/20" style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}>
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-4">Join the Waitlist. It&apos;s Free.</h2>
            <p className="text-[#8A9BB0] mb-6">No credit card. No subscription. No catch.</p>
            <Link href="/#waitlist" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96]">
              Get Early Access
            </Link>
            <p className="text-xs text-[#8A9BB0]/60 mt-4">
              Compare alternatives: <Link href="/compare/covenant-eyes" className="text-[#C9A84C] hover:underline">vs Covenant Eyes</Link> &middot; <Link href="/compare/ever-accountable" className="text-[#C9A84C] hover:underline">vs Ever Accountable</Link> &middot; <Link href="/compare/accountable2you" className="text-[#C9A84C] hover:underline">vs Accountable2You</Link>
            </p>
          </div>
        </div>
      </div>
    </>
  )
}
