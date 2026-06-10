import type { Metadata } from 'next'
import Link from 'next/link'
import FaqAccordion from '@/components/FaqAccordion'

export const metadata: Metadata = {
  title: 'How It Works',
  description:
    'A complete breakdown of Remain Faithful\'s accountability model, on-device monitoring pipeline, covenant system, and frequently asked questions.',
}

const faqs = [
  {
    q: 'Is it really free?',
    a: 'Yes, always. Remain Faithful has no subscription tier, no premium features, and no advertising. The app is sustained by voluntary donations from users who find it valuable. We have committed to this model indefinitely.',
  },
  {
    q: 'Who sees my data?',
    a: 'Your chosen accountability partners can see alert metadata: the timestamp, the category (e.g., "adult content"), and the severity level. They do not see screenshots, browsing history, app content, or raw OCR text. None of that data is ever transmitted off your device.',
  },
  {
    q: 'Can I be anonymous?',
    a: 'You choose your display name when you create your account. However, accountability by design requires that your partners know who they are holding accountable. Anonymity defeats the purpose. Your partners see the name you provide, typically your real name.',
  },
  {
    q: 'What exactly gets monitored?',
    a: 'When you enable monitoring and grant screen broadcast permission, the app analyzes frames from your screen using Apple\'s built-in Vision OCR and SensitiveContentAnalysis frameworks. It looks for concerning content across all apps. The analysis is performed entirely on-device. Nothing is sent to a server for classification (unless the on-device classifiers are uncertain, in which case a category-only query (no content) may reach our cloud classifier).',
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
    a: 'iOS\'s ReplayKit framework allows apps to record the screen via a separate extension process. This extension runs in isolation from the main app. It cannot access the internet or your personal data directly. It captures frames, runs AI analysis locally, and only sends a structured alert event (not the frame) to the main app if something is flagged.',
  },
]

export default function HowItWorksPage() {
  return (
    <>
      {/* Hero */}
      <section className="pt-32 pb-20 border-b border-[#1E3050]">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-4">The Method</p>
          <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-6">
            How Remain Faithful Works
          </h1>
          <p className="text-[#8A9BB0] text-lg max-w-2xl mx-auto">
            A complete accountability system built on peer trust, on-device privacy, and the covenant model.
            Here is exactly what happens under the hood.
          </p>
        </div>
      </section>

      {/* Accountability Model */}
      <section className="py-20">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-5">
                The Accountability Model
              </h2>
              <p className="text-[#8A9BB0] leading-relaxed mb-4">
                Lasting change in purity requires three things: vulnerability, consistent visibility, and a community of trust. Most accountability fails because it operates on the honor system: you report what you choose to report, when you choose to report it.
              </p>
              <p className="text-[#8A9BB0] leading-relaxed mb-4">
                Remain Faithful creates a consistent, automatic signal that removes the decision to disclose. When your device flags something, your partners know, regardless of whether you would have told them. This isn&apos;t surveillance. It&apos;s the covenant made real.
              </p>
              <p className="text-[#8A9BB0] leading-relaxed">
                The app is a tool, not a replacement for relationship. Alerts are conversation starters, not verdicts.
              </p>
            </div>
            <div className="space-y-4">
              {[
                { label: 'Vulnerability', desc: 'Inviting partners to see you clearly requires honest agreement to the covenant upfront.' },
                { label: 'Consistent Visibility', desc: 'Automatic alerts remove the shame barrier of self-disclosure without removing personal responsibility.' },
                { label: 'Community of Trust', desc: 'Partners accept terms before gaining any access. Relationships must exist before accountability can work.' },
              ].map((item) => (
                <div key={item.label} className="flex gap-4 p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
                  <div className="w-2 rounded-full bg-gradient-to-b from-[#C9A84C] to-[#E8C87A] flex-shrink-0 mt-1" style={{ minHeight: 40 }} />
                  <div>
                    <h3 className="font-semibold text-[#F0EDE8] mb-1">{item.label}</h3>
                    <p className="text-sm text-[#8A9BB0]">{item.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* One-to-One vs Group */}
      <section className="py-20 bg-[#0A1420] border-y border-[#1E3050]">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">
              One-to-One or Group Accountability
            </h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
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
                'Private message thread for following up on alerts',
              ]}
            />
            <ModeCard
              title="Small Group"
              points={[
                'Up to 12 members in a single group',
                'Group leader (pastor or facilitator) can view aggregate dashboards',
                'Members can choose to share with all group members or leader only',
                'Ideal for men\'s ministry, discipleship cohorts, or Bible study groups',
                'Weekly group digest sent to all members',
              ]}
            />
          </div>
        </div>
      </section>

      {/* How Monitoring Works */}
      <section className="py-20">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">
              How Monitoring Works
            </h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
              A three-tier pipeline that prioritizes on-device privacy above all else.
            </p>
          </div>

          <div className="space-y-4 relative">
            <div className="absolute left-7 top-8 bottom-8 w-px bg-gradient-to-b from-[#C9A84C]/60 via-[#C9A84C]/20 to-transparent" />
            {[
              {
                step: '1',
                title: 'Broadcast Extension Captures Frames',
                body: 'When you enable monitoring, iOS activates a ReplayKit broadcast extension, a sandboxed process separate from the main app. It captures screen frames at intervals without access to the network or your personal data.',
              },
              {
                step: '2',
                title: 'Tier 1: Blocklist + Pattern Matching',
                body: 'The frame URL (if a browser) and visible text are checked against a local blocklist and regex patterns. Known adult domains and explicit keywords trigger an immediate flag. This is fast and 100% local.',
              },
              {
                step: '3',
                title: 'Tier 2: On-Device AI Classification',
                body: 'If Tier 1 doesn\'t catch it, the frame is analyzed by Apple Vision (OCR), SensitiveContentAnalysis (Apple\'s nudity/explicit content detector), and a local keyword-weighted text classifier. All of this runs on your device\'s Neural Engine.',
              },
              {
                step: '4',
                title: 'Tier 3: Cloud Classifier (Category Only)',
                body: 'If the on-device classifiers are uncertain and the content appears borderline, a category query (containing only the text category, never the image) is sent to our Claude-based classifier at /classify. No screenshot. No personal content.',
              },
              {
                step: '5',
                title: 'Alert Delivered to Partners',
                body: 'If any tier flags the content: you receive a notification first. Then a push notification goes to each partner containing only: timestamp, app category, and severity level. Partners see nothing beyond those three data points.',
              },
            ].map((item) => (
              <div key={item.step} className="flex gap-6 p-5 rounded-xl border border-[#1E3050] bg-[#162235] ml-4">
                <div className="w-8 h-8 rounded-full bg-gradient-to-br from-[#C9A84C] to-[#E8C87A] flex items-center justify-center text-[#0F1B2D] font-bold text-sm flex-shrink-0 -ml-8 border-2 border-[#0F1B2D]">
                  {item.step}
                </div>
                <div>
                  <h3 className="font-semibold text-[#F0EDE8] mb-1">{item.title}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.body}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Covenant Model */}
      <section className="py-20 bg-[#0A1420] border-y border-[#1E3050]">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">The Foundation</p>
              <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-5">
                The Covenant Model
              </h2>
              <p className="text-[#8A9BB0] leading-relaxed mb-4">
                Before any partner gains access to your account, they must agree to a covenant. This isn&apos;t a terms-of-service checkbox. It&apos;s a statement of intent that frames the entire relationship.
              </p>
              <p className="text-[#8A9BB0] leading-relaxed">
                You also agree to it on your end. Accountability is bilateral. The covenant frames both the monitoring and the response to it.
              </p>
            </div>
            <div
              className="rounded-2xl p-8 border border-[#C9A84C]/20"
              style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}
            >
              <h3 className="font-serif text-xl font-semibold text-[#F0EDE8] mb-5">The Covenant</h3>
              <ul className="space-y-3">
                {[
                  'I will be honest with my partner, even when it is difficult.',
                  'I will not use this app to condemn my partner.',
                  'I will respond to alerts with grace and genuine care.',
                  'I will not share my partner\'s alerts with others.',
                  'I will pursue my partner\'s flourishing above my own curiosity.',
                ].map((line, i) => (
                  <li key={i} className="flex gap-3 text-sm text-[#8A9BB0] leading-relaxed">
                    <span className="text-[#C9A84C] mt-0.5 flex-shrink-0">✦</span>
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
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">
              Frequently Asked Questions
            </h2>
          </div>
          <FaqAccordion faqs={faqs} />
        </div>
      </section>

      {/* CTA */}
      <section className="py-16 border-t border-[#1E3050]">
        <div className="max-w-xl mx-auto px-4 text-center">
          <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-4">Ready to start?</h2>
          <p className="text-[#8A9BB0] mb-8">Join the waitlist and be among the first to use Remain Faithful when it launches.</p>
          <a
            href="/#waitlist"
            className="inline-flex items-center gap-2 px-8 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all duration-200"
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
    <div className="rounded-2xl p-8 border border-[#1E3050] bg-[#162235]">
      <h3 className="font-serif text-xl font-semibold text-[#F0EDE8] mb-5">{title}</h3>
      <ul className="space-y-2.5">
        {points.map((p, i) => (
          <li key={i} className="flex gap-3 text-sm text-[#8A9BB0] leading-relaxed">
            <svg className="flex-shrink-0 mt-0.5" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2.5" strokeLinecap="round">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
            {p}
          </li>
        ))}
      </ul>
    </div>
  )
}
