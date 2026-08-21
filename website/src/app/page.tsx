import type { Metadata } from 'next'
import { Suspense } from 'react'
import Link from 'next/link'
import DonateButton from '@/components/DonateButton'
import DonationSuccessBanner from '@/components/DonationSuccessBanner'
import LetterAtmosphere from '@/components/LetterAtmosphere'
import WaitlistForm from '@/components/WaitlistForm'
import { JsonLd } from '@/components/JsonLd'
import { softwareApplicationSchema, homepageFaqSchema } from '@/lib/structured-data'

export const metadata: Metadata = {
  title: 'Remain Faithful | Free Christian Accountability App for iPhone',
  description: 'Free peer accountability for Christians committed to purity. Always-on Family Controls filtering with partner notify. Optional Deep Scan. No screen content ever leaves your device. Partners get category, timestamp, and severity only. 100% free, forever.',
  keywords: ['free accountability app', 'Christian accountability app', 'purity app', 'accountability partner app', 'church accountability', 'iPhone accountability app', 'on-device AI', 'privacy-first accountability', 'open source accountability'],
  openGraph: {
    title: 'Remain Faithful | Free Christian Accountability App',
    description: 'Free peer accountability for Christians committed to purity. Always-on filtering, optional Deep Scan, privacy-first, open source.',
    url: 'https://remainfaithful.com',
    siteName: 'Remain Faithful',
    type: 'website',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Remain Faithful | Free Christian Accountability App',
    description: 'Free peer accountability for Christians committed to purity. Always-on filtering, optional Deep Scan, privacy-first.',
  },
  alternates: {
    canonical: 'https://remainfaithful.com',
  },
}

const faqs = [
  {
    q: 'Can my accountability partners see what I was looking at?',
    a: 'No. Remain Faithful never captures or shares screenshots. Partners receive a discreet alert with a category (like "Adult Content") and severity level, never the actual content. This protects partners from being exposed to harmful material.',
  },
  {
    q: 'Can my spouse be my accountability partner?',
    a: 'Yes. You can set up a one-to-one partnership with your spouse, a friend, a mentor, or a pastor. You choose who sees your alerts.',
  },
  {
    q: 'Can our church use this for small groups?',
    a: 'Absolutely. Remain Faithful was designed for small group accountability. Groups of 3 to 12 members can all monitor and encourage each other. We provide a free Group Setup Guide for ministry leaders.',
  },
  {
    q: 'What happens if I slip up?',
    a: 'Your accountability partners receive a discreet alert. The goal is conversation, not condemnation. Every alert includes conversation starter prompts to help your partners respond with grace.',
  },
  {
    q: 'How is this different from other accountability apps?',
    a: 'Remain Faithful is different in three key ways: the credible core is always-on Family Controls filtering with partner notify (Deep Scan is optional), it is 100% free forever, and no screen content ever leaves your device. The entire codebase is open source so anyone can verify exactly what is and is not transmitted.',
  },
  {
    q: 'Does this work on Android?',
    a: 'Android support is planned for Fall 2026. Currently Remain Faithful is available for iPhone (iOS 17+). Join the waitlist to be notified when Android launches.',
  },
]

export default function HomePage() {
  return (
    <>
      <JsonLd data={softwareApplicationSchema} />
      <JsonLd data={homepageFaqSchema} />

      <Suspense fallback={null}>
        <DonationSuccessBanner />
      </Suspense>

      <LetterAtmosphere>
        <header className="relative flex items-start justify-between gap-4 mb-8">
          <p className="kicker">Fort Wayne · Pre-launch · iPhone</p>
          <span className="wax-seal inline-flex items-center justify-center shrink-0" aria-hidden="true">
            <svg width="14" height="16" viewBox="0 0 32 36" fill="none">
              <path d="M16 0L2 6V18C2 26.284 8.268 33.916 16 36C23.732 33.916 30 26.284 30 18V6L16 0Z" fill="currentColor" />
            </svg>
          </span>
        </header>

        <h1
          id="letter-heading"
          className="font-serif font-normal text-[2.6rem] sm:text-6xl leading-[1.08] tracking-[-0.02em] text-ink mb-6"
        >
          Remain Faithful
        </h1>

        <p className="text-lg sm:text-xl leading-[1.7] text-ink mb-6 max-w-[36ch]">
          Free peer accountability for Christians committed to purity.
        </p>

        <p className="text-[1.05rem] leading-[1.75] text-ink-soft mb-10 max-w-[42ch]">
          Always-on filtering blocks the apps and categories you choose, and notifies your partners when a blocked category is attempted. Optional Deep Scan adds on-device AI for high-risk periods. No screen content ever leaves your device.
        </p>

        <div className="flex flex-wrap gap-3 mb-10">
          <a href="#waitlist" className="btn-wax">
            Join the Waitlist
          </a>
          <Link href="/how-it-works" className="btn-ghost">
            How it works
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" aria-hidden="true">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </Link>
        </div>

        <ul className="grid sm:grid-cols-2 gap-x-6 gap-y-2 pt-6 border-t border-hairline">
          <TrustItem>Always-on filtering</TrustItem>
          <TrustItem>100% free, forever</TrustItem>
          <TrustItem>No screen content ever leaves your device</TrustItem>
          <TrustItem>Optional Deep Scan</TrustItem>
        </ul>
      </LetterAtmosphere>

      <section id="waitlist" className="page-section">
        <p className="kicker mb-4">Coming soon</p>
        <h2 className="font-serif font-normal text-3xl sm:text-4xl text-ink mb-4">
          Write your name down.
        </h2>
        <p className="text-ink-soft text-lg leading-relaxed mb-8 max-w-[40ch]">
          Enter your email and we&apos;ll notify you the moment Remain Faithful is available.
        </p>
        <WaitlistForm variant="default" buttonText="Join the Waitlist" />
      </section>

      <hr className="rule" />

      <section className="page-section !pt-0">
        <p className="kicker mb-4">What it is</p>
        <h2 className="font-serif font-normal text-3xl sm:text-4xl text-ink mb-5">
          Built on covenant, not a stage.
        </h2>
        <p className="text-ink-soft text-lg leading-relaxed mb-10 max-w-[46ch]">
          <Link href="/blog/why-accountability-fails" className="text-ink underline decoration-hairline underline-offset-4 hover:decoration-wax">Most accountability tools rely on shame or surveillance.</Link> Remain Faithful is built on <Link href="/blog/covenant-model" className="text-ink underline decoration-hairline underline-offset-4 hover:decoration-wax">covenant</Link>, trust, and genuine community.
        </p>

        <ol className="m-0 p-0 list-none">
          <Feature
            n="01"
            title="One-to-one or group"
            body="Choose a single trusted partner or set up a small group. RF works for close friendships, mentorship relationships, and accountability groups alike."
            href="/blog/setting-up-your-first-group"
          />
          <Feature
            n="02"
            title="On-device privacy"
            body="Optional Deep Scan classification runs locally using Apple's Vision and SensitiveContentAnalysis frameworks. Your screen content is never transmitted. Partners see category, timestamp, and severity — not which app, and not your screen."
            href="/blog/on-device-privacy-explained"
          />
          <Feature
            n="03"
            title="Always-on filtering"
            body="Select the apps and categories you want blocked. They stay shielded continuously — through lock screen, reboot, and app restarts. Your partners are notified when a blocked category is attempted."
          />
          <Feature
            n="04"
            title="Always free"
            body="Remain Faithful is free today and will remain free forever. No subscription tiers, no paywalls, no premium features. Sustained entirely by voluntary donations."
          />
        </ol>
      </section>

      <hr className="rule" />

      <section className="page-section !pt-0">
        <p className="kicker mb-4">The method</p>
        <h2 className="font-serif font-normal text-3xl sm:text-4xl text-ink mb-5">
          Three steps to real accountability
        </h2>
        <p className="text-ink-soft text-lg leading-relaxed mb-10 max-w-[42ch]">
          No complicated setup. No long onboarding. Start holding each other accountable in minutes.
        </p>

        <ol className="m-0 p-0 list-none space-y-8">
          <Step
            n="I"
            title="Choose your partners"
            body="Invite one or more trusted friends, mentors, or group members to be your accountability partners. They accept a covenant before gaining any access."
          />
          <Step
            n="II"
            title="Enable filtering"
            body="Always-on Family Controls filtering blocks the apps and categories you choose and notifies your partners when a blocked category is attempted. DeviceActivity monitors usage as category events — not screen content. For high-risk periods you can start an optional Deep Scan session."
          />
          <Step
            n="III"
            title="Stay accountable"
            body="When something is flagged, you and your partners receive a discreet alert. No surprises, no shame spirals. Just honest accountability."
          />
        </ol>

        <p className="mt-10">
          <Link
            href="/how-it-works"
            className="font-mono text-[12px] tracking-[0.08em] uppercase text-wax hover:text-wax-deep"
          >
            Read the full breakdown →
          </Link>
        </p>
      </section>

      <hr className="rule" />

      <section className="page-section !pt-0">
        <p className="kicker mb-4">Privacy</p>
        <h2 className="font-serif font-normal text-3xl sm:text-4xl text-ink mb-5">
          Your content never leaves your device
        </h2>
        <p className="text-ink-soft text-lg leading-relaxed mb-10 max-w-[42ch]">
          We designed the privacy model first, then built the app around it.
        </p>

        <div className="space-y-6 mb-12">
          <PrivacyPoint
            title="On-device AI"
            body="Apple Vision OCR and SensitiveContentAnalysis run entirely on your hardware. No server-side processing of your screen."
          />
          <PrivacyPoint
            title="Alert metadata only"
            body="Partners receive: timestamp, category, and severity level. Never which app. Never a screenshot. Never raw content. Never your browsing history."
          />
          <PrivacyPoint
            title="You control access"
            body="You approve every partner. You set alert thresholds. You can pause monitoring or remove partners at any time, instantly."
          />
        </div>

        <div className="grid sm:grid-cols-2 gap-8 sm:gap-12">
          <div>
            <h3 className="kicker mb-4">Stays on your device</h3>
            <ul className="space-y-2 text-ink-soft leading-relaxed">
              <li>Screenshots &amp; raw screen content</li>
              <li>Browsing history &amp; page content</li>
              <li>Passwords &amp; financial data</li>
              <li>Message content</li>
              <li>Photos &amp; videos</li>
            </ul>
          </div>
          <div>
            <h3 className="kicker mb-4">Shared with partners</h3>
            <ul className="space-y-2 text-ink leading-relaxed">
              <li>Alert category (e.g. &ldquo;Adult Content&rdquo;)</li>
              <li>Severity level (Low / Medium / High)</li>
              <li>Timestamp</li>
            </ul>
          </div>
        </div>

        <p className="mt-8 text-sm text-ink-soft leading-relaxed">
          Classification is fully on-device. Only alert metadata (category, severity, and timestamp) is ever uploaded — no screen content, no OCR text, no screenshots, no system-generated description.{' '}
          <Link href="/privacy-architecture" className="text-wax hover:text-wax-deep">
            Full Privacy Architecture
          </Link>
        </p>
      </section>

      <hr className="rule" />

      <section id="download" className="page-section !pt-0">
        <p className="kicker mb-4">Get started</p>
        <h2 className="font-serif font-normal text-3xl sm:text-4xl text-ink mb-5">
          Start your accountability journey
        </h2>
        <p className="text-ink-soft text-lg leading-relaxed mb-8 max-w-[36ch]">
          Free forever. No subscription. No ads. Just accountability.
        </p>

        <div className="flex flex-col sm:flex-row gap-3">
          <a href="#waitlist" className="btn-wax">
            <span className="text-left">
              <span className="block font-mono text-[10px] tracking-[0.12em] uppercase opacity-80">iPhone (iOS 17+)</span>
              <span>Join the Waitlist</span>
            </span>
          </a>
          <div className="btn-ghost opacity-70 pointer-events-none" aria-disabled="true">
            <span className="text-left">
              <span className="block font-mono text-[10px] tracking-[0.12em] uppercase">Android launching</span>
              <span>Fall 2026</span>
            </span>
          </div>
        </div>
        <p className="font-mono text-[11px] tracking-[0.1em] uppercase text-ink-faint mt-6">
          Android launching Fall 2026. Join the waitlist to be notified.
        </p>
      </section>

      <hr className="rule" />

      <section className="page-section !pt-0">
        <p className="kicker mb-4">Common questions</p>
        <h2 className="font-serif font-normal text-3xl sm:text-4xl text-ink mb-10">
          Frequently asked questions
        </h2>

        <dl className="space-y-0">
          {faqs.map((faq) => (
            <div key={faq.q} className="plate-row !grid-cols-1 gap-2 py-6">
              <dt className="font-serif text-xl text-ink">{faq.q}</dt>
              <dd className="text-ink-soft leading-relaxed m-0">{faq.a}</dd>
            </div>
          ))}
        </dl>
      </section>

      <hr className="rule" />

      <section id="donate" className="page-section !pt-0 pb-24">
        <p className="kicker mb-4">Sustain the work</p>
        <h2 className="font-serif font-normal text-3xl sm:text-4xl text-ink mb-5">
          Keep Remain Faithful free
        </h2>
        <p className="text-ink-soft text-lg leading-relaxed mb-10 max-w-[40ch]">
          We&apos;re committed to never charging for accountability. Your donation funds server costs, development, and ministry outreach.
        </p>
        <DonateButton />
      </section>
    </>
  )
}

function TrustItem({ children }: { children: React.ReactNode }) {
  return (
    <li className="flex items-baseline gap-2 font-mono text-[11px] tracking-[0.08em] uppercase text-ink-soft">
      <span className="text-wax" aria-hidden="true">—</span>
      <span>{children}</span>
    </li>
  )
}

function Feature({
  n, title, body, href,
}: {
  n: string; title: string; body: string; href?: string
}) {
  const heading = href ? (
    <Link href={href} className="font-serif text-2xl text-ink hover:text-wax transition-colors duration-200">
      {title}
    </Link>
  ) : (
    <h3 className="font-serif text-2xl text-ink">{title}</h3>
  )

  return (
    <li className="plate-row">
      <span className="font-mono text-[11px] tabular-nums tracking-[0.12em] text-ink-faint">{n}</span>
      <div>
        {heading}
        <p className="text-ink-soft leading-relaxed mt-2">{body}</p>
      </div>
    </li>
  )
}

function Step({ n, title, body }: { n: string; title: string; body: string }) {
  return (
    <li>
      <p className="font-mono text-[11px] tracking-[0.18em] uppercase text-wax mb-2">{n}</p>
      <h3 className="font-serif text-2xl text-ink mb-2">{title}</h3>
      <p className="text-ink-soft leading-relaxed max-w-[46ch]">{body}</p>
    </li>
  )
}

function PrivacyPoint({ title, body }: { title: string; body: string }) {
  return (
    <div>
      <h3 className="font-serif text-xl text-ink mb-1">{title}</h3>
      <p className="text-ink-soft leading-relaxed">{body}</p>
    </div>
  )
}
