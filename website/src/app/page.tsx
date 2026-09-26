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
  description: 'Free peer accountability for Christians committed to purity. Always-on Family Controls filtering with partner notify. Optional Deep Scan. Screenshots and raw screen content never leave your device. Partners may receive a short system-generated summary plus category, severity, and timestamp. 100% free, forever.',
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

const proof = [
  {
    href: '#download',
    kicker: 'iPhone',
    title: 'Get started on iOS 17+',
    body: 'Free forever. Always-on Family Controls filtering with partner notify. Android is planned.',
  },
  {
    href: '/how-it-works',
    kicker: 'How it works',
    title: 'Partners, then filtering',
    body: 'Invite adult peers, authorize Family Controls, and receive alerts with category, severity, timestamp, and a short system-generated summary.',
  },
  {
    href: '/privacy-architecture',
    kicker: 'Privacy',
    title: 'Your screen stays on your device',
    body: 'Screenshots, OCR text, and raw screen content stay on the phone. Partners receive alert metadata and a short system-generated summary.',
  },
]

const faqs = [
  {
    q: 'Can my accountability partners see what I was looking at?',
    a: 'No. Remain Faithful never captures or shares screenshots. Partners receive a discreet alert with a category (like "Adult Content"), severity, timestamp, and a short system-generated summary — never screenshots, OCR text, or raw screen content. This protects partners from being exposed to harmful material.',
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
    a: 'Remain Faithful is different in three key ways: the credible core is always-on Family Controls filtering with partner notify (Deep Scan is optional), it is 100% free forever, and screenshots, OCR text, and raw screen content never leave your device. Partners may receive a short system-generated summary in addition to category, severity, and timestamp. The iOS app, backend, and website source are public on GitHub. Design notes and audit markdown in that repo are historical working papers, not product promises.',
  },
  {
    q: 'Does this work on Android?',
    a: 'Remain Faithful is available for iPhone (iOS 17+). Android support is planned, with no launch date. Leave your email for Android notify and updates.',
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
          <p className="kicker">Fort Wayne · iPhone (iOS 17+)</p>
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

        <p className="text-[1.05rem] leading-[1.75] text-ink-soft mb-8 max-w-[42ch]">
          Always-on filtering blocks the apps and categories you choose, and notifies your partners when a blocked category is attempted. Optional Deep Scan adds on-device AI for high-risk periods. Screenshots, OCR text, and raw screen content stay on your device. Partners receive a short system-generated summary.
        </p>

        <div className="flex flex-col sm:flex-row sm:items-stretch gap-3 mb-10 max-w-xl">
          <a href="#download" className="btn-wax w-full sm:w-auto">
            <span className="text-left">
              <span className="block font-mono text-[10px] font-semibold tracking-[0.14em] uppercase opacity-80">iPhone · iOS 17+</span>
              <span className="block text-[0.95rem] leading-tight">Get started</span>
            </span>
          </a>
          <Link href="/how-it-works" className="btn-ghost w-full sm:w-auto pl-6 pr-[1.375rem]">
            How it works
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" aria-hidden="true">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </Link>
        </div>

        <ul className="grid sm:grid-cols-3 gap-3 m-0 p-0 list-none">
          {proof.map((item) => (
            <li key={item.href}>
              <ProofLink href={item.href}>
                <span className="kicker">{item.kicker}</span>
                <span className="font-serif text-[1.35rem] leading-snug text-ink mt-3 text-balance">{item.title}</span>
                <span className="text-sm leading-relaxed text-ink-soft mt-2 text-pretty">{item.body}</span>
              </ProofLink>
            </li>
          ))}
        </ul>
      </LetterAtmosphere>

      <hr className="rule" />

      <section id="download" className="page-section !pt-0 scroll-mt-24">
        <div className="download-card">
          <div className="download-pane">
            <p className="kicker mb-4">iPhone · iOS 17+</p>
            <h2 className="font-serif font-normal text-3xl sm:text-4xl text-ink mb-4">
              Get started
            </h2>
            <p className="text-ink-soft text-lg leading-relaxed max-w-[40ch]">
              Free forever. Always-on Family Controls filtering, optional Deep Scan, and partners you choose. Android is planned.
            </p>

            <div className="grid sm:grid-cols-3 gap-4 mt-8 mb-2">
              <Spec term="Platform" detail="iPhone, iOS 17+" />
              <Spec term="Price" detail="Free forever" />
              <Spec term="Privacy" detail="Screen stays on device" />
            </div>

            <div className="flex flex-col sm:flex-row sm:items-stretch gap-1 mt-4">
              <Link href="/how-it-works" className="tap-row sm:flex-1">
                How it works
                <Arrow />
              </Link>
              <a href="#waitlist" className="tap-row sm:flex-1">
                Android notify
                <Arrow />
              </a>
            </div>
          </div>
        </div>
      </section>

      <section id="waitlist" className="page-section">
        <p className="kicker mb-4">Updates</p>
        <h2 className="font-serif font-normal text-3xl sm:text-4xl text-ink mb-4">
          Write your name down.
        </h2>
        <p className="text-ink-soft text-lg leading-relaxed mb-8 max-w-[40ch]">
          Remain Faithful is available for iPhone (iOS 17+). Leave your email for updates and Android notify.
        </p>
        <WaitlistForm variant="default" buttonText="Get updates" />
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
            body="Optional Deep Scan classification runs locally using Apple's Vision and SensitiveContentAnalysis frameworks. Screenshots, OCR text, and raw screen content are never transmitted. Partners may receive a short system-generated summary in addition to category, timestamp, and severity — not which app, and not your screen."
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

      <section className="page-section !pt-0" aria-labelledby="method-heading">
        <p className="kicker mb-4">The method</p>
        <h2 id="method-heading" className="font-serif font-normal text-3xl sm:text-4xl text-ink mb-4">
          Three steps to real accountability
        </h2>
        <p className="text-ink-soft text-lg leading-relaxed mb-8 max-w-[46ch]">
          Invite adult peers, then authorize Family Controls so always-on filtering can run. Each person joins and authorizes on their own iPhone.
        </p>

        <div className="method-card">
          <div className="grid md:grid-cols-2 gap-2">
            <div className="method-pane">
              <p className="kicker mb-4">How it works</p>
              <ol className="m-0 p-0 list-none">
                <Step
                  n="I"
                  title="Choose your partners"
                  body="Invite trusted friends, a mentor, a spouse, or a small group. They accept a covenant before gaining any access."
                />
                <Step
                  n="II"
                  title="Enable filtering"
                  body="Always-on Family Controls blocks the apps and categories you choose and notifies partners when a blocked category is attempted. Optional Deep Scan is for high-risk periods."
                />
                <Step
                  n="III"
                  title="Stay accountable"
                  body="A discreet alert carries category, severity, timestamp, and a short system-generated summary. The aim is an honest conversation."
                />
              </ol>
              <Link href="/how-it-works" className="tap-row">
                Read the full breakdown
                <Arrow />
              </Link>
            </div>

            <div className="method-pane">
              <p className="kicker mb-4">Privacy</p>
              <h3 className="font-serif text-2xl text-ink mb-3">
                Screenshots stay on your device
              </h3>
              <p className="text-ink-soft leading-relaxed mb-6">
                Classification runs on your iPhone. You approve every partner, set thresholds, and can pause or remove access.
              </p>
              <div className="grid grid-cols-1 min-[420px]:grid-cols-2 gap-5">
                <div>
                  <h4 className="kicker mb-3">On your device</h4>
                  <ul className="space-y-2 text-sm text-ink-soft leading-relaxed">
                    <li>Screenshots &amp; raw screen content</li>
                    <li>Browsing history &amp; page content</li>
                    <li>Passwords &amp; financial data</li>
                    <li>Message content</li>
                    <li>Photos &amp; videos</li>
                  </ul>
                </div>
                <div>
                  <h4 className="kicker mb-3">With partners</h4>
                  <ul className="space-y-2 text-sm text-ink leading-relaxed">
                    <li>Alert category</li>
                    <li>Severity (Low / Medium / High)</li>
                    <li>Timestamp</li>
                    <li>Short system-generated summary</li>
                  </ul>
                </div>
              </div>
              <Link href="/privacy-architecture" className="tap-row">
                Full Privacy Architecture
                <Arrow />
              </Link>
            </div>
          </div>
        </div>
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

function ProofLink({ href, children }: { href: string; children: React.ReactNode }) {
  const className = 'proof-card h-full'
  if (href.startsWith('#')) {
    return <a href={href} className={className}>{children}</a>
  }
  return <Link href={href} className={className}>{children}</Link>
}

function Spec({ term, detail }: { term: string; detail: string }) {
  return (
    <div className="min-h-11">
      <p className="kicker mb-1">{term}</p>
      <p className="text-ink leading-snug m-0">{detail}</p>
    </div>
  )
}

function Arrow() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" aria-hidden="true">
      <path d="M5 12h14M12 5l7 7-7 7"/>
    </svg>
  )
}

function Feature({
  n, title, body, href,
}: {
  n: string; title: string; body: string; href?: string
}) {
  const heading = href ? (
    <Link href={href} className="inline-flex items-center min-h-11 font-serif text-2xl text-ink hover:text-wax transition-colors duration-200">
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
    <li className="py-4 border-b border-hairline last:border-b-0">
      <p className="font-mono text-[11px] tracking-[0.18em] uppercase text-wax mb-1.5">{n}</p>
      <h3 className="font-serif text-xl text-ink mb-1.5">{title}</h3>
      <p className="text-sm text-ink-soft leading-relaxed text-pretty">{body}</p>
    </li>
  )
}
