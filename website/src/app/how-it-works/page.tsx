import type { Metadata } from 'next'
import Link from 'next/link'
import FaqAccordion from '@/components/FaqAccordion'
import FeaturePlate from '@/components/FeaturePlate'
import { JsonLd } from '@/components/JsonLd'
import { howItWorksFaqSchema } from '@/lib/structured-data'
import { Breadcrumbs } from '@/components/Breadcrumbs'
import { ALERT_TRIGGER, FAQ_STREAMING, PRIVACY_POINTER } from '@/lib/product-copy'
import AppStoreLink from '@/components/AppStoreLink'

export const metadata: Metadata = {
  title: 'How Remain Faithful Works: Always-On Filtering for iPhone',
  description: 'Create an account, invite a trusted partner, and turn on protection. Free on iPhone for adults.',
  alternates: { canonical: 'https://remainfaithful.com/how-it-works' },
}

const faqs = [
  {
    q: 'Is it really free?',
    a: 'Yes, always. Remain Faithful has no subscription tier, no premium features, and no advertising. The app is sustained by voluntary donations from users who find it valuable. We have committed to this model indefinitely.',
  },
  {
    q: 'Who sees my data?',
    a: (
      <>
        Screenshots and screen content stay on the device. What a partner can receive is on the{' '}
        <Link href="/privacy" className="text-white underline decoration-white/25 underline-offset-4 hover:decoration-wax">
          privacy page
        </Link>
        .
      </>
    ),
  },
  {
    q: 'Can I be anonymous?',
    a: 'You choose your display name when you create your account. However, accountability by design requires that your partners know who they are holding accountable. Anonymity defeats the purpose. Your partners see the name you provide, typically your real name.',
  },
  {
    q: 'What exactly gets monitored?',
    a: `${ALERT_TRIGGER} ${FAQ_STREAMING}`,
  },
  {
    q: 'How do I leave a group?',
    a: 'Open Settings → Manage Groups → Leave Group. The other members are notified. Leaving does not delete your alert history. Alerts stay while your account exists and are deleted when you delete the account.',
  },
  {
    q: 'Does this work on Android?',
    a: 'Not yet. The current app requires iOS 17 or later due to its reliance on Apple-specific frameworks (ReplayKit, Vision, SensitiveContentAnalysis). Android support is planned for late 2026 or early 2027 but is not yet available.',
  },
  {
    q: 'What is the broadcast extension?',
    a: `${FAQ_STREAMING} ${PRIVACY_POINTER}`,
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
      <section className="border-b border-white/10 pb-20 pt-8">
        <div className="mx-auto grid max-w-6xl items-center gap-12 px-4 sm:px-6 lg:grid-cols-[minmax(0,1fr)_20rem] lg:px-8">
          <div>
            <p className="kicker mb-4">The Method</p>
            <h1 className="mb-6 font-serif text-4xl font-bold text-ink sm:text-5xl lg:text-6xl">
              How Remain Faithful Works
            </h1>
            <p className="max-w-xl text-lg text-ink-soft">
              Always-on Family Controls filtering is the credible core. Deep Scan is optional.
              Built on peer trust, on-device privacy, and the covenant model.
            </p>
          </div>
          <FeaturePlate
            kicker="Setup"
            title="Three short steps"
            lines={[
              { label: 'Create an account', detail: 'Free, on your iPhone.' },
              { label: 'Invite a trusted partner', detail: 'A friend, spouse, or mentor you choose.' },
              { label: 'Turn on protection', detail: 'Choose the apps, then turn on Alerts or Blocking.' },
            ]}
          />
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
                { label: 'Vulnerability', desc: 'Inviting partners in requires a clear agreement to the covenant first.' },
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
              <span className="inline-flex items-center gap-2 px-3 py-1 rounded-sm bg-wax/10 border border-wax/25 text-wax text-xs font-bold uppercase tracking-widest">
                <span className="w-1.5 h-1.5 rounded-sm bg-wax" />
                Layer 1: Always-on filtering
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
                title: 'One alert after a minute of use',
                body: ALERT_TRIGGER,
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
              <span className="inline-flex items-center gap-2 px-3 py-1 rounded-sm bg-wax/10 border border-wax/25 text-wax text-xs font-bold uppercase tracking-widest">
                <span className="w-1.5 h-1.5 rounded-sm bg-wax" />
                Layer 2: DeviceActivity usage monitoring
              </span>
            </div>

            <div className="flex gap-6 p-5 rounded-sm border border-hairline bg-paper-deep ml-4">
              <div className="w-8 h-8 rounded-sm bg-wax flex items-center justify-center text-paper font-bold text-sm flex-shrink-0 -ml-8 border-2 border-paper">
                3
              </div>
              <div>
                <h3 className="font-semibold text-ink mb-1">Category Events, Not Screen Content</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  DeviceActivity watches usage as category-level events in the background. It does not capture screen frames or page content. {ALERT_TRIGGER}
                </p>
              </div>
            </div>

            <div className="ml-4 pb-1 pt-4">
              <span className="inline-flex items-center gap-2 px-3 py-1 rounded-sm bg-wax/10 border border-wax/25 text-wax text-xs font-bold uppercase tracking-widest">
                <span className="w-1.5 h-1.5 rounded-sm bg-wax" />
                Layer 3: blocked all day
              </span>
            </div>

            <div className="flex gap-6 p-5 rounded-sm border border-hairline bg-paper-deep ml-4">
              <div className="w-8 h-8 rounded-sm bg-wax flex items-center justify-center text-paper font-bold text-sm flex-shrink-0 -ml-8 border-2 border-paper">
                4
              </div>
              <div>
                <h3 className="font-semibold text-ink mb-1">Chosen apps stay blocked all day</h3>
                <p className="text-sm text-ink-soft leading-relaxed">
                  The apps and categories you choose stay blocked from the start of the day through the end. The block holds through the lock screen, reboot, and app restarts.
                </p>
              </div>
            </div>

            <div className="ml-4 pb-1 pt-4">
              <span className="inline-flex items-center gap-2 px-3 py-1 rounded-sm bg-wax/10 border border-wax/25 text-wax text-xs font-bold uppercase tracking-widest">
                Layer 4: Deep Scan (optional)
              </span>
            </div>

            {[
              {
                step: '5',
                title: 'You Start the Session',
                body: 'Deep Scan is never automatic. You start it intentionally when you want stronger scrutiny, such as a high-risk period, a travel trip, or a season of struggle. iOS asks for explicit screen broadcast permission each time. You are always in control.',
                note: undefined as string | undefined,
              },
              {
                step: '6',
                title: 'On-Device AI Analyzes Non-DRM Frames',
                body: "A sandboxed ReplayKit broadcast extension handles frames from apps that allow it. Each frame is analyzed by Apple Vision (OCR), SensitiveContentAnalysis (Apple's nudity detector), and a local keyword classifier, all running on your device's Neural Engine. Classification is entirely on-device. Frames are never stored or transmitted.",
                note: 'DRM video apps such as Netflix, Disney+, Hulu, Prime Video, Apple TV, and HBO do not share frames, so there is nothing from those apps to classify. Banking apps work the same way. That is Apple FairPlay and platform DRM, not a Remain Faithful bug. Browsers, photos, social media, and most other apps can be classified on the phone. Screenshots and screen content stay on the phone.',
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
                  When a layer flags something, you receive a notification first. Then a push notification can go to each partner. {PRIVACY_POINTER}
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
                  'I will be open with my partner, even when it is difficult.',
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
          <p className="text-ink-soft mb-8">Free on iPhone (iOS 17+).</p>
          <AppStoreLink />
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
            <svg className="flex-shrink-0 mt-0.5" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#D1AB4C" strokeWidth="2.5" strokeLinecap="round">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
            {p}
          </li>
        ))}
      </ul>
    </div>
  )
}
