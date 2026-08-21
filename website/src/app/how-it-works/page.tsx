import type { Metadata } from 'next'
import Link from 'next/link'
import FaqAccordion from '@/components/FaqAccordion'
import { JsonLd } from '@/components/JsonLd'
import { howItWorksFaqSchema } from '@/lib/structured-data'
import { Breadcrumbs } from '@/components/Breadcrumbs'

export const metadata: Metadata = {
  title: 'How Remain Faithful Works: Always-On Filtering for iPhone',
  description: 'Always-on Family Controls filtering, DeviceActivity usage monitoring, time-window shielding, and optional Deep Scan. Partners receive category, timestamp, and severity — never screen content.',
  alternates: { canonical: 'https://remainfaithful.com/how-it-works' },
}

const faqs = [
  {
    q: 'Is it really free?',
    a: 'Yes, always. Remain Faithful has no subscription tier, no premium features, and no advertising. The app is sustained by voluntary donations from users who find it valuable. We have committed to this model indefinitely.',
  },
  {
    q: 'Who sees my data?',
    a: 'Your chosen accountability partners can see alert metadata: the timestamp, the category (e.g., "adult content"), and the severity level. They do not see which app, a bundle ID, screenshots, browsing history, app content, or raw OCR text. None of that data is ever transmitted off your device.',
  },
  {
    q: 'Can I be anonymous?',
    a: 'You choose your display name when you create your account. However, accountability by design requires that your partners know who they are holding accountable. Anonymity defeats the purpose. Your partners see the name you provide, typically your real name.',
  },
  {
    q: 'What exactly gets monitored?',
    a: 'Remain Faithful has four public layers. The credible core is always-on Family Controls filtering: you choose apps and categories to block, they stay blocked through lock and reboot, and partners are notified when a blocked category is attempted. DeviceActivity monitors usage as category events — not screen content and not which specific app. Time-window shielding can restrict chosen apps during hours you set. Opt-in Deep Scan uses ReplayKit and on-device AI (Vision OCR, SensitiveContentAnalysis) on non-DRM frames. Deep Scan cannot see DRM-protected video or banking apps — Netflix, Disney+, Hulu, Prime Video, Apple TV, HBO, and banking apps render as black frames. That is Apple FairPlay / platform DRM, not a Remain Faithful bug. Partners receive category, timestamp, and severity only.',
  },
  {
    q: 'How do I leave a group?',
    a: 'Navigate to Settings → Groups → select the group → Leave Group. When you leave a group, all group members are notified. Your data is not retained after leaving. Partners will no longer receive alerts from you. Your historical alerts within the group are purged per your data retention setting (default 30 days).',
  },
  {
    q: 'Does this work on Android?',
    a: 'Not yet. The current app requires iOS 17 or later due to its reliance on Apple-specific frameworks (ReplayKit, Vision, SensitiveContentAnalysis). Android support is planned for late 2026 or early 2027 but is not yet available.',
  },
  {
    q: 'What is the broadcast extension?',
    a: 'The broadcast extension is part of opt-in Deep Scan only. iOS ReplayKit allows a sandboxed extension to capture screen frames from non-DRM apps. The extension runs AI analysis entirely on-device and, when something is flagged, uploads only alert metadata (category, severity, and timestamp) — never the screen frame itself. Screen content, OCR text, and screenshots are never transmitted. Always-on filtering, usage monitoring, and time-window shielding use Apple\'s Family Controls and DeviceActivity frameworks — they do not use a broadcast extension and they do not see screen content.',
  },
]

export default function HowItWorksPage() {
  return (
    <>
      <JsonLd data={howItWorksFaqSchema} />
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 pt-24">
        <Breadcrumbs items={[{ name: 'How It Works', url: 'https://remainfaithful.com/how-it-works' }]} />
      </div>
      {/* Hero */}
      <section className="pt-8 pb-20 border-b border-hairline">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p className="text-wax text-sm font-semibold uppercase tracking-widest mb-4">The Method</p>
          <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-6">
            How Remain Faithful Works
          </h1>
          <p className="text-ink-soft text-lg max-w-2xl mx-auto">
            Always-on Family Controls filtering is the credible core. Deep Scan is optional.
            Built on peer trust, on-device privacy, and the covenant model.
          </p>
        </div>
      </section>

      {/* Accountability Model */}
      <section className="py-20">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <h2 className="font-serif text-3xl font-bold text-ink mb-5">
                The Accountability Model
              </h2>
              <p className="text-ink-soft leading-relaxed mb-4">
                Lasting change in purity requires three things: vulnerability, consistent visibility, and a community of trust. Most accountability fails because it operates on the honor system: you report what you choose to report, when you choose to report it.
              </p>
              <p className="text-ink-soft leading-relaxed mb-4">
                Remain Faithful creates a consistent, automatic signal that removes the decision to disclose. When your device flags something, your partners know, regardless of whether you would have told them. This isn&apos;t surveillance. It&apos;s the covenant made real.
              </p>
              <p className="text-ink-soft leading-relaxed">
                The app is a tool, not a replacement for relationship. Alerts are conversation starters, not verdicts.
              </p>
            </div>
            <div className="space-y-4">
              {[
                { label: 'Vulnerability', desc: 'Inviting partners to see you clearly requires honest agreement to the covenant upfront.' },
                { label: 'Consistent Visibility', desc: 'Automatic alerts remove the shame barrier of self-disclosure without removing personal responsibility.' },
                { label: 'Community of Trust', desc: 'Partners accept terms before gaining any access. Relationships must exist before accountability can work.' },
              ].map((item) => (
                <div key={item.label} className="flex gap-4 p-5 rounded-sm border border-hairline bg-paper-deep">
                  <div className="w-2 rounded-sm bg-wax flex-shrink-0 mt-1" style={{ minHeight: 40 }} />
                  <div>
                    <h3 className="font-semibold text-ink mb-1">{item.label}</h3>
                    <p className="text-sm text-ink-soft">{item.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* One-to-One vs Group */}
      <section className="py-20 bg-paper-deep border-y border-hairline">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-ink mb-4">
              One-to-One or Group Accountability
            </h2>
            <p className="text-ink-soft max-w-xl mx-auto">
              RF supports both models. Choose what fits your relationships and community structure.
            </p>
          </div>
          <div className="grid md:grid-cols-2 gap-6">
            <ModeCard
              title="One-to-One Partnership"
              points={[
                'Maximum privacy between two people',
                'Deep, focused relationship built on mutual trust',
                'Ideal for close friends, mentors, or spouses',
                'Each partner can monitor the other (reciprocal) or one direction',
              ]}
            />
            <ModeCard
              title="Small Group"
              points={[
                'Up to 12 members in a single group',
                'Ideal for men\'s ministry, discipleship cohorts, or Bible study groups',
                'All members receive alerts when any member is flagged',
                'Group admin manages membership and invite codes',
              ]}
            />
          </div>
        </div>
      </section>

      {/* How Monitoring Works */}
      <section className="py-20">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-ink mb-4">
              How Monitoring Works
            </h2>
            <p className="text-ink-soft max-w-xl mx-auto">
              Four public layers. Always-on Family Controls filtering is the credible core. Deep Scan is optional.
            </p>
          </div>

          <div className="space-y-3 relative">
            <div className="absolute left-7 top-10 bottom-10 w-px bg-wax/30" />

            <div className="ml-4 pb-1 pt-2">
              <span className="inline-flex items-center gap-2 px-3 py-1 rounded-sm bg-green-500/10 border border-green-500/25 text-green-400 text-xs font-bold uppercase tracking-widest">
                <span className="w-1.5 h-1.5 rounded-sm bg-green-400" />
                Layer 1 — Always-On Filtering
              </span>
            </div>

            {[
              {
                step: '1',
                title: 'Family Controls Blocks What You Choose',
                body: 'You select the apps and categories to block. Apple Family Controls keeps them shielded through lock screen, reboot, and app restarts. No screen broadcast permission. Force-quitting the app does not lift the block.',
              },
              {
                step: '2',
                title: 'Partners Are Notified on a Blocked Attempt',
                body: 'When a blocked category is attempted, an alert is created. Partners receive category, timestamp, and severity — not which app, not a bundle ID, and never screen content.',
              },
            ].map((item) => (
              <div key={item.step} className="flex gap-6 p-5 rounded-sm border border-hairline bg-paper-deep ml-4">
                <div className="w-8 h-8 rounded-sm bg-wax flex items-center justify-center text-paper font-bold text-sm flex-shrink-0 -ml-8 border-2 border-paper">
                  {item.step}
                </div>
                <div>
                  <h3 className="font-semibold text-ink mb-1">{item.title}</h3>
                  <p className="text-sm text-ink-soft leading-relaxed">{item.body}</p>
                </div>
              </div>
            ))}

            <div className="ml-4 pb-1 pt-4">
              <span className="inline-flex items-center gap-2 px-3 py-1 rounded-sm bg-green-500/10 border border-green-500/25 text-green-400 text-xs font-bold uppercase tracking-widest">
                <span className="w-1.5 h-1.5 rounded-sm bg-green-400" />
                Layer 2 — DeviceActivity Usage Monitoring
              </span>
            </div>

            <div className="flex gap-6 p-5 rounded-sm border border-hairline bg-paper-deep ml-4">
              <div className="w-8 h-8 rounded-sm bg-wax flex items-center justify-center text-paper font-bold text-sm flex-shrink-0 -ml-8 border-2 border-paper">
                3
              </div>
              <div>
                <h3 className="font-semibold text-ink mb-1">Category Events, Not Screen Content</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  DeviceActivity watches usage as category-level events in the background. It does not capture screen frames, page content, or which specific app was opened. Partners still receive only category, timestamp, and severity.
                </p>
              </div>
            </div>

            <div className="ml-4 pb-1 pt-4">
              <span className="inline-flex items-center gap-2 px-3 py-1 rounded-sm bg-green-500/10 border border-green-500/25 text-green-400 text-xs font-bold uppercase tracking-widest">
                <span className="w-1.5 h-1.5 rounded-sm bg-green-400" />
                Layer 3 — Time-Window Shielding
              </span>
            </div>

            <div className="flex gap-6 p-5 rounded-sm border border-hairline bg-paper-deep ml-4">
              <div className="w-8 h-8 rounded-sm bg-wax flex items-center justify-center text-paper font-bold text-sm flex-shrink-0 -ml-8 border-2 border-paper">
                4
              </div>
              <div>
                <h3 className="font-semibold text-ink mb-1">Shield Chosen Apps During Hours You Set</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  You can restrict chosen apps to time windows — evenings, travel, or a season of struggle. Shielding uses the same Family Controls stack as always-on filtering. It does not see screen content.
                </p>
              </div>
            </div>

            <div className="ml-4 pb-1 pt-4">
              <span className="inline-flex items-center gap-2 px-3 py-1 rounded-sm bg-wax/10 border border-wax/25 text-wax text-xs font-bold uppercase tracking-widest">
                Layer 4 — Deep Scan (Opt-In)
              </span>
            </div>

            {[
              {
                step: '5',
                title: 'You Start the Session',
                body: 'Deep Scan is never automatic. You start it intentionally when you want stronger scrutiny — a high-risk period, a travel trip, or a season of struggle. iOS asks for explicit screen broadcast permission each time. You are always in control.',
                note: undefined as string | undefined,
              },
              {
                step: '6',
                title: 'On-Device AI Analyzes Non-DRM Frames',
                body: "A sandboxed ReplayKit broadcast extension captures screen frames from apps that allow it. Each frame is analyzed by Apple Vision (OCR), SensitiveContentAnalysis (Apple's nudity detector), and a local keyword classifier — all running on your device's Neural Engine. Classification is entirely on-device. Frames are never stored or transmitted.",
                note: 'Deep Scan cannot see DRM-protected video or banking apps. Netflix, Disney+, Hulu, Prime Video, Apple TV, HBO, and banking apps render as black frames. That is Apple FairPlay / platform DRM — unbypassable, not a Remain Faithful bug. It can analyze browsers, photos, social media, and most non-DRM apps.',
              },
            ].map((item) => (
              <div key={item.step} className="flex gap-6 p-5 rounded-sm border border-hairline bg-paper-deep ml-4">
                <div className="w-8 h-8 rounded-sm bg-wax flex items-center justify-center text-paper font-bold text-sm flex-shrink-0 -ml-8 border-2 border-paper">
                  {item.step}
                </div>
                <div>
                  <h3 className="font-semibold text-ink mb-1">{item.title}</h3>
                  <p className="text-sm text-ink-soft leading-relaxed">{item.body}</p>
                  {item.note && (
                    <p className="text-sm text-wax/80 leading-relaxed mt-2 flex items-start gap-1.5">
                      <svg className="flex-shrink-0 mt-0.5" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                        <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
                      </svg>
                      {item.note}
                    </p>
                  )}
                </div>
              </div>
            ))}

            <div className="flex gap-6 p-5 rounded-sm border border-wax/25 bg-paper-deep ml-4 mt-3">
              <div className="w-8 h-8 rounded-sm bg-wax flex items-center justify-center text-paper font-bold text-sm flex-shrink-0 -ml-8 border-2 border-paper">
                7
              </div>
              <div>
                <h3 className="font-semibold text-ink mb-1">Alert Delivered to Partners</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  When a layer flags something, you receive a notification first. Then a push notification goes to each partner containing only: the alert category, severity level, and timestamp. Partners do not receive which app, a bundle ID, or a system-generated description. No screenshots. No raw content. Ever.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Covenant Model */}
      <section className="py-20 bg-paper-deep border-y border-hairline">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <p className="text-wax text-sm font-semibold uppercase tracking-widest mb-3">The Foundation</p>
              <h2 className="font-serif text-3xl font-bold text-ink mb-5">
                The Covenant Model
              </h2>
              <p className="text-ink-soft leading-relaxed mb-4">
                Before any partner gains access to your account, they must agree to a covenant. This isn&apos;t a terms-of-service checkbox. It&apos;s a statement of intent that frames the entire relationship.{' '}
                <Link href="/blog/covenant-model" className="text-wax hover:underline underline-offset-2">
                  Read more about the covenant model.
                </Link>
              </p>
              <p className="text-ink-soft leading-relaxed">
                You also agree to it on your end. Accountability is bilateral. The covenant frames both the monitoring and the response to it.
              </p>
            </div>
            <div className="rounded-sm p-8" style={{ boxShadow: 'var(--shadow-border)' }}>
              <h3 className="font-serif text-xl font-semibold text-ink mb-5">The Covenant</h3>
              <ul className="space-y-3">
                {[
                  'I will be honest with my partner, even when it is difficult.',
                  'I will not use this app to condemn my partner.',
                  'I will respond to alerts with grace and genuine care.',
                  'I will not share my partner\'s alerts with others.',
                  'I will pursue my partner\'s flourishing above my own curiosity.',
                ].map((line, i) => (
                  <li key={i} className="flex gap-3 text-sm text-ink-soft leading-relaxed">
                    <span className="text-wax mt-0.5 flex-shrink-0">✦</span>
                    {line}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      </section>

      {/* FAQ */}
      <section className="py-20">
        <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-ink mb-4">
              Frequently Asked Questions
            </h2>
          </div>
          <FaqAccordion faqs={faqs} />
        </div>
      </section>

      {/* CTA */}
      <section className="py-16 border-t border-hairline">
        <div className="max-w-xl mx-auto px-4 text-center">
          <h2 className="font-serif text-2xl font-bold text-ink mb-4">Ready to start?</h2>
          <p className="text-ink-soft mb-8">Join the waitlist and be among the first to use Remain Faithful when it launches.</p>
          <a
            href="/#waitlist"
            className="inline-flex items-center gap-2 px-8 py-3.5 rounded-sm font-semibold text-paper bg-wax hover:bg-wax-deep transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96]"
          >
            Join the Waitlist
          </a>
        </div>
      </section>
    </>
  )
}

function ModeCard({ title, points }: { title: string; points: string[] }) {
  return (
    <div className="rounded-sm p-8 border border-hairline bg-paper-deep">
      <h3 className="font-serif text-xl font-semibold text-ink mb-5">{title}</h3>
      <ul className="space-y-2.5">
        {points.map((p, i) => (
          <li key={i} className="flex gap-3 text-sm text-ink-soft leading-relaxed">
            <svg className="flex-shrink-0 mt-0.5" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#7a2c28" strokeWidth="2.5" strokeLinecap="round">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
            {p}
          </li>
        ))}
      </ul>
    </div>
  )
}
