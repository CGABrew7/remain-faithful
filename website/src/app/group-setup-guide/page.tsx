import type { Metadata } from 'next'
import PrintButton from '@/components/PrintButton'

export const metadata: Metadata = {
  title: 'Group Setup Guide',
  description: 'A complete guide for pastors and ministry leaders on setting up a Remain Faithful accountability group.',
  robots: { index: false },
}

export default function GroupSetupGuidePage() {
  return (
    <>
      {/* Print button — hidden when printing */}
      <div className="print:hidden fixed bottom-6 right-6 z-50">
        <PrintButton />
      </div>

      <div className="pt-28 pb-20 print:pt-0 print:pb-0">
        <div className="max-w-3xl mx-auto px-4 sm:px-6 print:px-0 print:max-w-full">

          {/* Cover */}
          <div
            className="rounded-sm p-12 mb-12 text-center print:rounded-none print:mb-8 print:page-break-after-avoid"
            style={{ boxShadow: 'var(--shadow-border)' }}
          >
            <div className="flex justify-center mb-6">
              <svg width="56" height="64" viewBox="0 0 32 36" fill="none">
                <path d="M16 0L2 6V18C2 26.284 8.268 33.916 16 36C23.732 33.916 30 26.284 30 18V6L16 0Z" fill="url(#guideShieldGrad)"/>
                <path d="M11 18L14.5 21.5L21 14" stroke="white" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"/>
                <defs>
                  <linearGradient id="guideShieldGrad" x1="0" y1="0" x2="32" y2="36" gradientUnits="userSpaceOnUse">
                    <stop stopColor="#7a2c28"/>
                    <stop offset="1" stopColor="#5c211e"/>
                  </linearGradient>
                </defs>
              </svg>
            </div>
            <p className="text-wax text-sm font-semibold uppercase tracking-widest mb-3">Ministry Leader Guide</p>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-4">
              Remain Faithful
            </h1>
            <p className="font-serif text-2xl text-wax mb-6">Group Setup Guide</p>
            <p className="text-ink-soft max-w-md mx-auto leading-relaxed">
              Everything you need to launch a structured accountability group in your men&apos;s ministry, discipleship cohort, or small group.
            </p>
            <div className="mt-8 flex justify-center gap-8 text-sm text-ink-soft">
              <span>Free Forever</span>
              <span className="text-wax">|</span>
              <span>Privacy-First</span>
              <span className="text-wax">|</span>
              <span>iOS App</span>
            </div>
          </div>

          {/* Section 1: What is Remain Faithful */}
          <GuideSection number="1" title="What Is Remain Faithful?">
            <p className="text-ink-soft leading-relaxed mb-4">
              Remain Faithful (RF) is a free iOS app that provides peer accountability for adults committed to purity. Always-on Family Controls filtering blocks the apps and categories you choose and notifies partners when a blocked category is attempted. Optional Deep Scan adds on-device AI for high-risk periods.
            </p>
            <p className="text-ink-soft leading-relaxed mb-6">
              Always-on filtering uses Apple Family Controls. It does not see screen content. Optional Deep Scan, if started, classifies non-DRM frames on-device. No screen content is ever transmitted to anyone. Partners receive only alert metadata: category, timestamp, and severity.
            </p>
            <div className="grid sm:grid-cols-3 gap-4">
              {[
                { title: 'Always Free', desc: 'No subscription, no paywalls. Sustained by voluntary donations.' },
                { title: 'On-Device Privacy', desc: 'Screen content stays on the user\'s device. Partners see metadata only.' },
                { title: 'Covenant-Based', desc: 'Every partner accepts a covenant before gaining any access.' },
              ].map((item) => (
                <div key={item.title} className="p-4 rounded-sm border border-hairline bg-paper-deep">
                  <h4 className="font-semibold text-wax text-sm mb-2">{item.title}</h4>
                  <p className="text-xs text-ink-soft leading-relaxed">{item.desc}</p>
                </div>
              ))}
            </div>
          </GuideSection>

          {/* Section 2: How Group Accountability Works */}
          <GuideSection number="2" title="How Group Accountability Works">
            <p className="text-ink-soft leading-relaxed mb-4">
              RF supports both one-to-one partnerships and groups of up to 20 members. In a group setting:
            </p>
            <ul className="space-y-3 mb-6">
              {[
                'Each member installs the app and enables always-on filtering on their own device.',
                'When something is flagged, all group members (or the leader only, depending on settings) receive a push notification with the alert metadata.',
                'Members can choose to share alerts with all group members or the leader only.',
                'No one sees which app was involved, screenshots, browsing history, or any content from another member\'s screen.',
              ].map((point, i) => (
                <li key={i} className="flex gap-3 text-sm text-ink-soft leading-relaxed">
                  <span className="text-wax font-bold flex-shrink-0">{i + 1}.</span>
                  {point}
                </li>
              ))}
            </ul>
            <div
              className="rounded-sm p-6 border border-wax/25"
              
            >
              <p className="text-sm text-ink font-semibold mb-2">What partners see when an alert fires:</p>
              <div className="grid grid-cols-3 gap-3 mt-3">
                {['Timestamp', 'Category (e.g., "Adult Content")', 'Severity Level'].map((item) => (
                  <div key={item} className="text-center p-3 rounded-lg bg-paper border border-hairline">
                    <p className="text-xs text-ink-soft">{item}</p>
                  </div>
                ))}
              </div>
              <p className="text-xs text-ink-soft mt-3 text-center">That&apos;s it. Nothing else is shared.</p>
            </div>
          </GuideSection>

          {/* Section 3: Setting Up Your Group */}
          <GuideSection number="3" title="Setting Up Your Group (Step by Step)">
            <p className="text-ink-soft leading-relaxed mb-6">
              The whole process takes about 15 minutes for you as the leader, and 5 minutes per member.
            </p>
            <div className="space-y-4">
              {[
                {
                  n: '1',
                  title: 'Download and Create Your Leader Account',
                  desc: 'Download Remain Faithful from the App Store (iOS 17 or later required). Create an account with your name and email. Your role as group leader is set at this step.',
                },
                {
                  n: '2',
                  title: 'Create Your Group',
                  desc: 'Tap the Group tab at the bottom of the screen, then tap New Group. Give your group a name (e.g., "Tuesday Men" or "Iron Sharpens Iron"), review the covenant text, and choose your leader visibility settings.',
                },
                {
                  n: '3',
                  title: 'Choose Visibility Settings',
                  desc: 'Members can share alerts with the whole group or with the leader only. Discuss with your group which setting fits your context before enabling. Partners receive category, timestamp, and severity — not which app.',
                },
                {
                  n: '4',
                  title: 'Share the Invite Code',
                  desc: 'Your group generates a 6-character invite code. Share it in your group chat, on a handout, or in person. Members download the app, create an account, and join using that code.',
                },
                {
                  n: '5',
                  title: 'Members Enable Filtering',
                  desc: 'Each member authorizes Family Controls so always-on filtering can block chosen apps and categories. Deep Scan is optional and, if used later, asks for screen broadcast permission separately.',
                },
              ].map((step) => (
                <div key={step.n} className="flex gap-5 p-5 rounded-sm border border-hairline bg-paper-deep">
                  <div className="w-8 h-8 rounded-sm bg-wax text-paper font-bold text-sm flex items-center justify-center flex-shrink-0">
                    {step.n}
                  </div>
                  <div>
                    <h4 className="font-semibold text-ink mb-1">{step.title}</h4>
                    <p className="text-sm text-ink-soft leading-relaxed">{step.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </GuideSection>

          {/* Section 4: Inviting Members */}
          <GuideSection number="4" title="Inviting Members and the Covenant">
            <p className="text-ink-soft leading-relaxed mb-4">
              Before you share the invite code, have the conversation. Don&apos;t send someone a link to an app without context. Cover these points:
            </p>
            <div className="space-y-3 mb-6">
              {[
                { label: 'What it does', text: 'Always-on filtering blocks chosen apps and categories and notifies the group when a blocked category is attempted. Optional Deep Scan adds on-device AI. Partners receive category, timestamp, and severity only.' },
                { label: 'What you see', text: 'Metadata only: timestamp, content category, and severity. No screenshots. No browsing history. No screen content.' },
                { label: 'The covenant', text: 'Everyone in the group, including you as leader, agrees to a covenant before gaining any access. This shapes how the group responds to alerts.' },
                { label: 'It\'s voluntary', text: 'Monitoring can be paused or disabled at any time. Joining the group is a choice, not a requirement.' },
              ].map((item) => (
                <div key={item.label} className="flex gap-4 p-4 rounded-sm border border-hairline bg-paper-deep">
                  <div className="w-2 rounded-sm bg-wax flex-shrink-0 mt-1" style={{ minHeight: 36 }}/>
                  <div>
                    <p className="font-semibold text-ink text-sm mb-1">{item.label}</p>
                    <p className="text-sm text-ink-soft">{item.text}</p>
                  </div>
                </div>
              ))}
            </div>

            <div
              className="rounded-sm p-7 border border-wax/25"
              
            >
              <h3 className="font-serif text-lg font-semibold text-ink mb-4">The Covenant</h3>
              <p className="text-sm text-ink-soft mb-4">Every member agrees to this before joining. Reviewing it together in your first meeting is worth the time.</p>
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
          </GuideSection>

          {/* Section 5: Running Your Group */}
          <GuideSection number="5" title="Running Your Group Well">
            <p className="text-ink-soft leading-relaxed mb-5">
              RF is infrastructure. The ministry is what happens around it. Here&apos;s how to use it well:
            </p>

            <div className="space-y-5">
              <GuideTip title="First Meeting (Within Two Weeks of Launch)">
                Normalize the awkwardness. Monitoring feels strange at first. Talk about it openly. Discuss the covenant together. Agree on how the group will respond when an alert comes in: a text, a call, a prayer, or a conversation. The goal isn&apos;t catching each other. It&apos;s removing the barrier to honest conversation.
              </GuideTip>

              <GuideTip title="When an Alert Comes In">
                Respond quickly and with grace. The covenant already set the expectation. A good first response is something like: &ldquo;Hey, I got an alert. Thinking of you. Want to talk?&rdquo; Keep it simple. Keep it relational. Alerts are conversation starters, not verdicts.
              </GuideTip>

              <GuideTip title="Regular Check-Ins (Leader)">
                Use your group&apos;s existing meeting rhythm to pray for one another and notice who may need a personal check-in. Alerts are conversation starters as they arrive.
              </GuideTip>

              <GuideTip title="Handling False Positives">
                The AI isn&apos;t perfect. False positives happen. If a member gets an alert for something innocent, that&apos;s a conversation, not a verdict. Normalize this in your first meeting so men aren&apos;t caught off guard.
              </GuideTip>

              <GuideTip title="Leaving a Group">
                Members can leave at any time via Settings → Groups → Leave Group. When a member leaves, all group members are notified. Their historical alert data is purged per their data retention settings.
              </GuideTip>
            </div>
          </GuideSection>

          {/* Section 6: FAQ */}
          <GuideSection number="6" title="Common Questions">
            <div className="space-y-5">
              {[
                {
                  q: 'Does this cost anything?',
                  a: 'No. Remain Faithful is free for everyone, forever. No subscription, no premium tier. The app is sustained by voluntary donations from users who find it valuable.',
                },
                {
                  q: 'What devices does it work on?',
                  a: 'Currently iOS only (iPhone with iOS 17 or later). Android support is planned for late 2026 or early 2027.',
                },
                {
                  q: 'Can I see what a member\'s screen looks like?',
                  a: 'No. No one sees screenshots, screen recordings, or any content from another member\'s device. You see alert metadata only: timestamp, content category, and severity level.',
                },
                {
                  q: 'What if a member doesn\'t want to share with the whole group?',
                  a: 'Members can choose to share alerts with the group leader only, rather than the entire group. This is configurable in their settings.',
                },
                {
                  q: 'How many people can be in a group?',
                  a: 'Up to 20 members per group.',
                },
                {
                  q: 'Is the app open source?',
                  a: 'Yes. The iOS app, Go backend, and this website are all publicly available on GitHub. Anyone can inspect the code to verify our privacy claims.',
                },
              ].map((item, i) => (
                <div key={i} className="p-5 rounded-sm border border-hairline bg-paper-deep">
                  <p className="font-semibold text-ink text-sm mb-2">{item.q}</p>
                  <p className="text-sm text-ink-soft leading-relaxed">{item.a}</p>
                </div>
              ))}
            </div>
          </GuideSection>

          {/* Footer */}
          <div className="mt-12 pt-8 border-t border-hairline text-center">
            <div className="flex justify-center mb-4">
              <svg width="28" height="32" viewBox="0 0 32 36" fill="none">
                <path d="M16 0L2 6V18C2 26.284 8.268 33.916 16 36C23.732 33.916 30 26.284 30 18V6L16 0Z" fill="url(#footerGuideGrad)"/>
                <path d="M11 18L14.5 21.5L21 14" stroke="white" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"/>
                <defs>
                  <linearGradient id="footerGuideGrad" x1="0" y1="0" x2="32" y2="36" gradientUnits="userSpaceOnUse">
                    <stop stopColor="#7a2c28"/>
                    <stop offset="1" stopColor="#5c211e"/>
                  </linearGradient>
                </defs>
              </svg>
            </div>
            <p className="font-serif text-lg text-ink mb-2">Remain Faithful</p>
            <p className="text-sm text-ink-soft mb-4">Free peer accountability for adults committed to purity.</p>
            <p className="text-xs text-ink-soft/60">
              remainfaithful.com &nbsp;|&nbsp; support@remainfaithful.com &nbsp;|&nbsp; Free forever
            </p>
          </div>

        </div>
      </div>

    </>
  )
}

function GuideSection({ number, title, children }: { number: string; title: string; children: React.ReactNode }) {
  return (
    <div className="mb-12 print:mb-8 print:page-break-inside-avoid">
      <div className="flex items-center gap-4 mb-6">
        <div className="w-10 h-10 rounded-sm bg-wax text-paper font-bold text-lg flex items-center justify-center flex-shrink-0">
          {number}
        </div>
        <h2 className="font-serif text-2xl font-bold text-ink">{title}</h2>
      </div>
      {children}
    </div>
  )
}

function GuideTip({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="flex gap-4 p-5 rounded-sm border border-hairline bg-paper-deep">
      <div className="w-2 rounded-sm bg-wax flex-shrink-0 mt-1" style={{ minHeight: 40 }}/>
      <div>
        <h4 className="font-semibold text-ink mb-2">{title}</h4>
        <p className="text-sm text-ink-soft leading-relaxed">{children}</p>
      </div>
    </div>
  )
}
