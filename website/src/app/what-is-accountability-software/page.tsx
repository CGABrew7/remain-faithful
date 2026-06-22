import type { Metadata } from 'next'
import Link from 'next/link'
import { Breadcrumbs } from '@/components/Breadcrumbs'
import { JsonLd } from '@/components/JsonLd'

export const metadata: Metadata = {
  title: 'What Is Accountability Software? A Complete Guide (2026)',
  description: 'Accountability software monitors device activity and shares reports with a trusted partner. Learn how it works, who uses it, and how to choose the right tool for your situation.',
  alternates: { canonical: 'https://remainfaithful.com/what-is-accountability-software' },
}

const faqSchema = {
  '@context': 'https://schema.org',
  '@type': 'FAQPage',
  mainEntity: [
    {
      '@type': 'Question',
      name: 'What is accountability software?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Accountability software monitors device activity and shares a report with a trusted person, called an accountability partner. The goal is to create consistent visibility into device use so that struggling behavior is not hidden. Different tools use different approaches: some capture screenshots, some route traffic through a VPN, some use on-device AI to classify screen content without transmitting it.',
      },
    },
    {
      '@type': 'Question',
      name: 'How does accountability software work?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Accountability software works by monitoring device activity and generating reports or alerts for a designated partner. The monitoring approach varies by product: screenshot-based tools capture images of the screen at intervals; DNS-based tools block content at the network level; VPN-based tools route traffic through a local VPN to inspect URLs; on-device AI tools classify screen content using machine learning running entirely on the device. Partners receive reports via email, app notifications, or in-app dashboards depending on the product.',
      },
    },
    {
      '@type': 'Question',
      name: 'Who uses accountability software?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'The primary users are individuals and families seeking to address pornography or explicit content use, often within a religious context. Christians seeking purity accountability represent a large portion of the market. Parents use accountability software to monitor children\'s device activity. Recovery programs for sexual addiction integrate accountability software as part of treatment. Church small groups and men\'s ministries use it as a structured accountability layer.',
      },
    },
    {
      '@type': 'Question',
      name: 'Is accountability software legal?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes, on your own devices or on devices you own and provide to minors. Using accountability software on another adult\'s device without their knowledge or consent would raise serious legal and ethical concerns. All reputable accountability software requires the monitored person to install and enable the software themselves. Monitoring should always be consensual for adults.',
      },
    },
    {
      '@type': 'Question',
      name: 'What is the best free accountability app?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Remain Faithful is the only fully free accountability app for iPhone with no subscription tier, no premium features, and no advertising. It uses on-device AI so screen content never leaves your device. Other options like Covenant Eyes ($16.99/month), Ever Accountable ($99/year), and Accountable2You (~$80/year) all charge subscription fees.',
      },
    },
  ],
}

export default function WhatIsAccountabilitySoftware() {
  return (
    <>
      <JsonLd data={faqSchema} />

      <div className="pt-24 pb-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <Breadcrumbs items={[{ name: 'What Is Accountability Software?', url: 'https://remainfaithful.com/what-is-accountability-software' }]} />

          {/* Hero */}
          <div className="mb-10">
            <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-4">Educational Guide</p>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-6 leading-tight">
              What Is Accountability Software? A Complete Guide (2026)
            </h1>
            <p className="text-[#8A9BB0] text-lg leading-relaxed">
              Accountability software monitors device activity and shares a report with a trusted person, called an accountability partner. This guide explains how it works technically, who uses it, and how to evaluate the different options available.
            </p>
          </div>

          {/* Table of contents */}
          <nav className="mb-14 p-6 rounded-2xl border border-[#1E3050] bg-[#162235]">
            <h2 className="font-semibold text-[#F0EDE8] mb-3 text-sm uppercase tracking-wide">Contents</h2>
            <ol className="space-y-1.5 text-sm text-[#C9A84C]">
              {[
                ['#definition', 'Definition: what accountability software does'],
                ['#how-it-works', 'How accountability software works technically'],
                ['#approaches', 'The main monitoring approaches'],
                ['#use-cases', 'Who uses it and why'],
                ['#choosing', 'How to choose the right tool'],
                ['#major-tools', 'Major tools in the market'],
                ['#faq', 'Frequently asked questions'],
              ].map(([href, label]) => (
                <li key={href}>
                  <a href={href} className="hover:underline underline-offset-2">{label}</a>
                </li>
              ))}
            </ol>
          </nav>

          {/* Definition */}
          <section className="mb-14" id="definition">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Definition: What Accountability Software Does</h2>
            <p className="text-[#8A9BB0] leading-relaxed mb-4">
              Accountability software is any application designed to monitor a person&apos;s device activity and share that information with a designated partner. The core function is visibility: the monitored person consents to sharing behavioral data with someone they trust, creating a structure where continued use of problematic content or apps cannot be easily hidden.
            </p>
            <p className="text-[#8A9BB0] leading-relaxed mb-4">
              The category covers a range of products with different technical approaches, different levels of privacy, and different use cases. The common thread is the accountability relationship: one person consents to be monitored, and another person receives information about their activity.
            </p>
            <p className="text-[#8A9BB0] leading-relaxed">
              Accountability software should be distinguished from parental controls, which are typically non-consensual and applied by a parent to a child&apos;s device. Adult accountability software requires the monitored person to install and enable it voluntarily.
            </p>
          </section>

          {/* How it works technically */}
          <section className="mb-14" id="how-it-works">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">How Accountability Software Works Technically</h2>
            <p className="text-[#8A9BB0] leading-relaxed mb-4">
              The technical implementation varies significantly by product. The core workflow is similar across all of them:
            </p>
            <div className="space-y-3 mb-6">
              {[
                { step: '1', label: 'Installation and setup', desc: 'The user installs the app and grants required permissions. Depending on the monitoring approach, this may include screen recording permission, VPN configuration, or Screen Time integration.' },
                { step: '2', label: 'Partner invitation', desc: 'The user invites one or more accountability partners. Partners create accounts or receive email reports. In most products, partners must accept some form of agreement before gaining access.' },
                { step: '3', label: 'Monitoring runs continuously', desc: 'The app monitors device activity in the background. The monitoring approach (see below) determines what data is collected and where the analysis happens.' },
                { step: '4', label: 'Reports or alerts delivered', desc: 'When something is flagged, the partner receives a notification. Depending on the product, this may be an immediate alert, a daily email report, or a weekly summary.' },
                { step: '5', label: 'Accountability conversation', desc: 'The partner and monitored person discuss the alert. The quality of this conversation is what determines whether accountability actually works.' },
              ].map((item) => (
                <div key={item.step} className="flex gap-4 p-4 rounded-xl border border-[#1E3050] bg-[#162235]">
                  <div className="w-7 h-7 rounded-full bg-gradient-to-br from-[#C9A84C] to-[#E8C87A] text-[#0F1B2D] font-bold text-xs flex items-center justify-center flex-shrink-0">{item.step}</div>
                  <div>
                    <p className="font-semibold text-[#F0EDE8] text-sm mb-0.5">{item.label}</p>
                    <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </section>

          {/* Main approaches */}
          <section className="mb-14" id="approaches">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">The Main Monitoring Approaches</h2>
            <p className="text-[#8A9BB0] leading-relaxed mb-6">
              There are four primary technical approaches used by accountability software today. Each has different privacy implications, battery impacts, and coverage capabilities.
            </p>
            <div className="space-y-5">
              {[
                {
                  title: '1. Screenshot-based monitoring',
                  pros: 'Easy for partners to understand what was seen. Works across most content types.',
                  cons: 'Screenshots of private or harmful content are transmitted to partner devices and potentially company servers. Partners are exposed to the same harmful content the monitored person was viewing. Large amounts of sensitive data leave the device.',
                  examples: 'Covenant Eyes',
                },
                {
                  title: '2. DNS/network-based monitoring',
                  pros: 'Blocks content at the network level before it loads. Works across all apps.',
                  cons: 'Cannot see inside apps or encrypted HTTPS traffic. Usually focused on blocking rather than reporting. Does not monitor app activity within apps.',
                  examples: 'Various parental control and filtering products',
                },
                {
                  title: '3. VPN-based monitoring',
                  pros: 'Can inspect URLs and domain traffic. Works across multiple platforms. Provides detailed URL and page title logs.',
                  cons: 'Increases battery drain. Can conflict with corporate or school VPNs. Cannot see inside end-to-end encrypted apps. VPN must remain active continuously.',
                  examples: 'Accountable2You',
                },
                {
                  title: '4. On-device AI classification',
                  pros: 'Screen content never leaves the device. Classification happens entirely on-device using local AI models. Partners receive only metadata (category, severity) rather than screenshots. No battery penalty from VPN. Can monitor screen content across all apps.',
                  cons: 'Requires sufficient device processing power. Cannot monitor DRM-protected streaming video (iOS limitation). Newer approach with less real-world track record than VPN or screenshot methods.',
                  examples: 'Remain Faithful, Ever Accountable',
                },
              ].map((item) => (
                <div key={item.title} className="p-6 rounded-2xl border border-[#1E3050] bg-[#162235]">
                  <h3 className="font-semibold text-[#F0EDE8] mb-3">{item.title}</h3>
                  <div className="grid sm:grid-cols-2 gap-3 mb-3">
                    <div>
                      <p className="text-xs text-green-400 font-semibold uppercase tracking-wide mb-1">Advantages</p>
                      <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.pros}</p>
                    </div>
                    <div>
                      <p className="text-xs text-red-400 font-semibold uppercase tracking-wide mb-1">Limitations</p>
                      <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.cons}</p>
                    </div>
                  </div>
                  <p className="text-xs text-[#8A9BB0]/60">Examples: {item.examples}</p>
                </div>
              ))}
            </div>
          </section>

          {/* Use cases */}
          <section className="mb-14" id="use-cases">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Who Uses Accountability Software</h2>
            <div className="grid sm:grid-cols-2 gap-4">
              {[
                { title: 'Individuals seeking purity', body: 'The largest segment of users are adults seeking to address pornography use, typically within a religious context. Christians seeking purity accountability represent a substantial portion of the accountability software market.' },
                { title: 'Church small groups', body: 'Churches and men\'s ministries use accountability software as infrastructure for existing accountability relationships. Groups that already meet use the software to create accountability between meetings.' },
                { title: 'Married couples', body: 'Spouses who want mutual transparency use accountability software as one partner, providing consistent visibility without requiring self-disclosure.' },
                { title: 'Recovery programs', body: 'Programs addressing sexual addiction or pornography use, including Celebrate Recovery and similar ministry-based approaches, integrate accountability software as part of their structured support framework.' },
                { title: 'Parents monitoring older children', body: 'Some parents use accountability software rather than parental controls for teenagers who have earned more autonomy. The consensual nature of accountability software is an advantage in building trust.' },
                { title: 'Mentorship relationships', body: 'Pastors, discipleship leaders, and mentors use accountability software as a structured layer for one-on-one accountability relationships.' },
              ].map((item) => (
                <div key={item.title} className="p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
                  <h3 className="font-semibold text-[#F0EDE8] mb-2">{item.title}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.body}</p>
                </div>
              ))}
            </div>
          </section>

          {/* How to choose */}
          <section className="mb-14" id="choosing">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">How to Choose the Right Accountability Software</h2>
            <p className="text-[#8A9BB0] leading-relaxed mb-6">Five questions determine which tool fits your situation best.</p>
            <div className="space-y-4">
              {[
                { q: '1. Which platforms do you and your partner use?', a: 'If you need cross-platform coverage (Windows, Mac, Android), your options are Covenant Eyes, Ever Accountable, and Accountable2You. If you use iPhone exclusively, Remain Faithful is also an option.' },
                { q: '2. What should your partner actually receive?', a: 'Screenshot tools send actual screenshots of flagged content, exposing partners to harmful material. Category-only tools (like Remain Faithful) send only a category label and severity level, protecting partners from exposure while still creating accountability.' },
                { q: '3. What does privacy mean to you?', a: 'Screenshot tools transmit screen content to cloud servers and partner devices. VPN tools route traffic externally. On-device AI tools (Remain Faithful, Ever Accountable) keep screen content on your device. If privacy is important, on-device AI is the strongest option.' },
                { q: '4. What is your budget?', a: 'Costs range from $0 (Remain Faithful) to $204/year (Covenant Eyes for one person). If cost is a constraint, Remain Faithful is the only option with no cost at all. Others range from $80 to $204 per year per person.' },
                { q: '5. What kind of relationship framework do you want?', a: 'Most tools are monitoring tools. Remain Faithful adds a covenant layer in which partners explicitly agree to respond with grace before gaining access. If the relational framework matters to you, this distinction is significant.' },
              ].map((item) => (
                <div key={item.q} className="p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
                  <h3 className="font-semibold text-[#F0EDE8] mb-2">{item.q}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.a}</p>
                </div>
              ))}
            </div>
          </section>

          {/* Major tools */}
          <section className="mb-14" id="major-tools">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Major Tools in the Market (2026)</h2>
            <div className="overflow-x-auto rounded-2xl border border-[#1E3050]">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-[#1E3050] bg-[#0A1420]">
                    <th className="text-left p-4 text-[#8A9BB0] font-semibold">Tool</th>
                    <th className="text-left p-4 text-[#8A9BB0] font-semibold">Price</th>
                    <th className="text-left p-4 text-[#8A9BB0] font-semibold">Approach</th>
                    <th className="text-left p-4 text-[#8A9BB0] font-semibold">Platforms</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    ['Remain Faithful', 'Free', 'On-device AI', 'iOS 17+'],
                    ['Covenant Eyes', '~$17/mo', 'Screenshots + cloud AI', 'iOS, Android, Windows, Mac, Chromebook'],
                    ['Ever Accountable', '~$99/yr', 'On-device AI', 'iOS, Android, Windows, Mac'],
                    ['Accountable2You', '~$80/yr', 'VPN-based', 'iOS, Android, Windows, Mac, Kindle'],
                  ].map(([tool, price, approach, platforms], i) => (
                    <tr key={i} className={`border-b border-[#1E3050] ${i % 2 === 0 ? 'bg-[#162235]' : 'bg-[#0F1B2D]'}`}>
                      <td className="p-4 text-[#F0EDE8] font-medium">{tool}</td>
                      <td className="p-4 text-[#8A9BB0]">{price}</td>
                      <td className="p-4 text-[#8A9BB0]">{approach}</td>
                      <td className="p-4 text-[#8A9BB0]">{platforms}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            <p className="text-xs text-[#8A9BB0]/60 mt-2">Pricing as of June 2026. Check each provider&apos;s website for current pricing.</p>
            <div className="mt-4 flex flex-wrap gap-3 text-sm">
              <Link href="/compare/covenant-eyes" className="text-[#C9A84C] hover:underline underline-offset-2">RF vs Covenant Eyes →</Link>
              <Link href="/compare/ever-accountable" className="text-[#C9A84C] hover:underline underline-offset-2">RF vs Ever Accountable →</Link>
              <Link href="/compare/accountable2you" className="text-[#C9A84C] hover:underline underline-offset-2">RF vs Accountable2You →</Link>
            </div>
          </section>

          {/* FAQ */}
          <section className="mb-14" id="faq">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-6">Frequently Asked Questions</h2>
            <div className="space-y-4">
              {[
                { q: 'What is accountability software?', a: 'Accountability software monitors device activity and shares a report with a trusted person, called an accountability partner. The goal is to create consistent visibility into device use so that struggling behavior is not hidden. Different tools use different approaches: some capture screenshots, some route traffic through a VPN, some use on-device AI to classify screen content without transmitting it.' },
                { q: 'How does accountability software work?', a: 'Accountability software works by monitoring device activity and generating reports or alerts for a designated partner. The monitoring approach varies by product: screenshot-based tools capture images of the screen at intervals; DNS-based tools block content at the network level; VPN-based tools route traffic through a local VPN to inspect URLs; on-device AI tools classify screen content using machine learning running entirely on the device. Partners receive reports via email, app notifications, or in-app dashboards depending on the product.' },
                { q: 'Who uses accountability software?', a: 'The primary users are individuals and families seeking to address pornography or explicit content use, often within a religious context. Christians seeking purity accountability represent a large portion of the market. Parents use accountability software to monitor children\'s device activity. Recovery programs for sexual addiction integrate accountability software as part of treatment. Church small groups and men\'s ministries use it as a structured accountability layer.' },
                { q: 'Is accountability software legal?', a: 'Yes, on your own devices or on devices you own and provide to minors. Using accountability software on another adult\'s device without their knowledge or consent would raise serious legal and ethical concerns. All reputable accountability software requires the monitored person to install and enable the software themselves. Monitoring should always be consensual for adults.' },
                { q: 'What is the best free accountability app?', a: 'Remain Faithful is the only fully free accountability app for iPhone with no subscription tier, no premium features, and no advertising. It uses on-device AI so screen content never leaves your device. Other options like Covenant Eyes ($16.99/month), Ever Accountable ($99/year), and Accountable2You (~$80/year) all charge subscription fees.' },
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
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-4">Try Remain Faithful</h2>
            <p className="text-[#8A9BB0] mb-2">The only fully free accountability app for iPhone. On-device AI. Open source. No subscription.</p>
            <p className="text-[#8A9BB0] mb-6">
              <Link href="/how-it-works" className="text-[#C9A84C] hover:underline underline-offset-2">Learn exactly how it works</Link> before you decide.
            </p>
            <Link href="/#waitlist" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96]">
              Join the Waitlist (Free)
            </Link>
          </div>
        </div>
      </div>
    </>
  )
}
