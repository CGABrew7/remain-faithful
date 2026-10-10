import type { Metadata } from 'next'
import { Suspense } from 'react'
import Link from 'next/link'
import DonateButton from '@/components/DonateButton'
import DonationSuccessBanner from '@/components/DonationSuccessBanner'
import FaqAccordion from '@/components/FaqAccordion'
import Ornament from '@/components/Ornament'
import PhoneFrame from '@/components/PhoneFrame'
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
    href: '/privacy-architecture',
    kicker: 'Privacy',
    title: 'Your screen stays on your device',
    body: 'Screenshots, OCR text, and raw screen content stay on the phone. Partners receive alert metadata and a short system-generated summary.',
    lead: true,
  },
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

      <section className="shell grid items-center gap-8 pb-16 pt-28 lg:grid-cols-[minmax(0,1.05fr)_minmax(260px,0.8fr)] lg:gap-x-16 lg:gap-y-6 lg:pb-24 lg:pt-32">
        <div className="lg:col-start-1 lg:row-start-1">
          <p className="kicker rise">Free · iPhone (iOS 17+)</p>
          <Ornament className="rise rise-2 mb-6 mt-5" />
          <h1 className="rise rise-2 font-serif text-[2.75rem] font-bold leading-[1.04] tracking-[-0.03em] text-white sm:text-6xl lg:text-[4.15rem]">
            Remain Faithful
          </h1>
          <p className="rise rise-3 mt-6 max-w-[22ch] font-serif text-2xl leading-snug text-white/90 sm:text-[1.7rem]">
            Free peer accountability for Christians committed to purity.
          </p>
          <div className="rise rise-4 mt-8 flex flex-col gap-3 sm:flex-row">
            <a href="#download" className="btn-wax">
              Get started
            </a>
            <Link href="/how-it-works" className="btn-ghost">
              How it works
              <Arrow />
            </Link>
          </div>
        </div>
        <div className="rise rise-3 lg:col-start-2 lg:row-span-2 lg:row-start-1 lg:justify-self-end">
          <PhoneFrame
            src="/screens/welcome.png"
            alt="Remain Faithful iPhone welcome screen: Remain Faithful, accountability between trusted friends, and Begin Your Journey."
            priority
          />
        </div>
        <p className="max-w-[46ch] text-[1.05rem] leading-relaxed text-ink-soft lg:col-start-1">
          Always-on filtering blocks the apps and categories you choose, and notifies your partners when a blocked category is attempted. Optional Deep Scan adds on-device AI for high-risk periods. Screenshots, OCR text, and raw screen content stay on your device. Partners receive a short system-generated summary.
        </p>
      </section>

      <section className="shell pb-6" aria-label="What you get">
        <ul className="proof-bento m-0 list-none p-0">
          {proof.map((item) => (
            <li key={item.href} className={item.lead ? 'proof-lead' : undefined}>
              <ProofLink href={item.href}>
                <span className="kicker">{item.kicker}</span>
                <span className="mt-3 font-serif text-[1.65rem] leading-snug text-white">{item.title}</span>
                <span className="mt-3 text-[0.98rem] leading-relaxed text-ink-soft">{item.body}</span>
              </ProofLink>
            </li>
          ))}
        </ul>
      </section>

      <section id="download" className="section scroll-mt-24">
        <div className="shell grid items-center gap-12 lg:grid-cols-[minmax(0,1fr)_280px]">
          <div className="surface p-8 sm:p-10">
            <p className="kicker mb-4">iPhone · iOS 17+</p>
            <h2 className="font-serif text-4xl font-bold text-white sm:text-5xl">Get started</h2>
            <p className="mt-4 max-w-[42ch] text-lg leading-relaxed text-ink-soft">
              Free forever. Always-on Family Controls filtering, optional Deep Scan, and partners you choose. Android is planned.
            </p>
            <dl className="mt-8 grid gap-5 sm:grid-cols-3">
              <Spec term="Platform" detail="iPhone, iOS 17+" />
              <Spec term="Price" detail="Free forever" />
              <Spec term="Privacy" detail="Screen stays on device" />
            </dl>
            <div className="mt-8 flex flex-col gap-2 sm:flex-row">
              <Link href="/how-it-works" className="btn-ghost">
                How it works
                <Arrow />
              </Link>
              <a href="#waitlist" className="btn-ghost">
                Android notify
                <Arrow />
              </a>
            </div>
          </div>
          <PhoneFrame
            src="/screens/create.png"
            alt="Remain Faithful create-account screen with name, email, and password fields."
          />
        </div>
      </section>

      <section className="section !pt-4">
        <div className="shell grid items-start gap-12 lg:grid-cols-[minmax(0,0.85fr)_minmax(0,1.15fr)]">
          <div>
            <p className="kicker mb-4">What it is</p>
            <h2 className="font-serif text-4xl font-bold text-white sm:text-5xl">
              Built on covenant, not a stage.
            </h2>
            <p className="mt-5 max-w-[42ch] text-lg leading-relaxed text-ink-soft">
              <Link href="/blog/why-accountability-fails" className="text-white underline decoration-white/25 underline-offset-4 hover:decoration-wax">Most accountability tools rely on shame or surveillance.</Link> Remain Faithful is built on <Link href="/blog/covenant-model" className="text-white underline decoration-white/25 underline-offset-4 hover:decoration-wax">covenant</Link>, trust, and genuine community.
            </p>
          </div>
          <ol className="feature-rail">
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
        </div>
      </section>

      <section className="band section" aria-labelledby="method-heading">
        <div className="shell">
          <div className="grid items-center gap-10 lg:grid-cols-[minmax(0,1fr)_260px]">
            <div>
              <p className="kicker mb-4">The method</p>
              <h2 id="method-heading" className="font-serif text-4xl font-bold text-white sm:text-5xl">
                Three steps to real accountability
              </h2>
              <p className="mt-4 max-w-[46ch] text-lg leading-relaxed text-ink-soft">
                Invite adult peers, then authorize Family Controls so always-on filtering can run. Each person joins and authorizes on their own iPhone.
              </p>
            </div>
            <PhoneFrame
              src="/screens/home.png"
              alt="Remain Faithful home screen with monitoring active, a 47-day clean streak, and no recent flags."
            />
          </div>

          <ol className="stepper mt-12">
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
          <Link href="/how-it-works" className="btn-ghost mt-8">
            Read the full breakdown
            <Arrow />
          </Link>

          <div className="mt-14 grid gap-4 lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)]">
            <div className="surface p-7 sm:p-8">
              <p className="kicker mb-4">Privacy</p>
              <h3 className="font-serif text-3xl text-white">Screenshots stay on your device</h3>
              <p className="mt-3 leading-relaxed text-ink-soft">
                Classification runs on your iPhone. You approve every partner, set thresholds, and can pause or remove access.
              </p>
              <Link href="/privacy-architecture" className="btn-ghost mt-6">
                Full Privacy Architecture
                <Arrow />
              </Link>
            </div>
            <div className="split-panel">
              <div className="p-7 sm:p-8">
                <h4 className="kicker mb-4">On your device</h4>
                <ul className="space-y-2.5 text-sm leading-relaxed text-ink-soft">
                  <li>Screenshots &amp; raw screen content</li>
                  <li>Browsing history &amp; page content</li>
                  <li>Passwords &amp; financial data</li>
                  <li>Message content</li>
                  <li>Photos &amp; videos</li>
                </ul>
              </div>
              <div className="border-t border-white/10 p-7 sm:border-l sm:border-t-0 sm:p-8">
                <h4 className="kicker mb-4">With partners</h4>
                <ul className="space-y-2.5 text-sm leading-relaxed text-white">
                  <li>Alert category</li>
                  <li>Severity (Low / Medium / High)</li>
                  <li>Timestamp</li>
                  <li>Short system-generated summary</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="waitlist" className="section scroll-mt-24">
        <div className="shell grid items-center gap-10 lg:grid-cols-2">
          <div>
            <p className="kicker mb-4">Updates</p>
            <h2 className="font-serif text-4xl font-bold text-white sm:text-5xl">Write your name down.</h2>
            <p className="mt-4 max-w-[40ch] text-lg leading-relaxed text-ink-soft">
              Remain Faithful is available for iPhone (iOS 17+). Leave your email for updates and Android notify.
            </p>
          </div>
          <div className="surface p-6 sm:p-8">
            <WaitlistForm variant="default" buttonText="Get updates" />
          </div>
        </div>
      </section>

      <section className="section !pt-0">
        <div className="shell max-w-3xl">
          <p className="kicker mb-4">Common questions</p>
          <h2 className="mb-8 font-serif text-4xl font-bold text-white sm:text-5xl">
            Frequently asked questions
          </h2>
          <FaqAccordion faqs={faqs} />
        </div>
      </section>

      <section id="donate" className="section scroll-mt-24">
        <div className="shell">
          <div className="surface relative overflow-hidden px-6 py-14 text-center sm:px-12">
            <div
              className="pointer-events-none absolute inset-0"
              aria-hidden="true"
              style={{ background: 'radial-gradient(50% 80% at 50% 0%, rgba(209,171,76,0.16), transparent 70%)' }}
            />
            <div className="relative">
              <p className="kicker mb-4">Sustain the work</p>
              <h2 className="font-serif text-4xl font-bold text-white sm:text-5xl">
                Keep Remain Faithful free
              </h2>
              <p className="mx-auto mt-4 max-w-[40ch] text-lg leading-relaxed text-ink-soft">
                We&apos;re committed to never charging for accountability. Your donation funds server costs, development, and ministry outreach.
              </p>
              <div className="mt-10">
                <DonateButton />
              </div>
            </div>
          </div>
        </div>
      </section>
    </>
  )
}

function ProofLink({ href, children }: { href: string; children: React.ReactNode }) {
  const className = 'surface surface-link h-full'
  if (href.startsWith('#')) {
    return <a href={href} className={className}>{children}</a>
  }
  return <Link href={href} className={className}>{children}</Link>
}

function Spec({ term, detail }: { term: string; detail: string }) {
  return (
    <div>
      <dt className="kicker mb-1">{term}</dt>
      <dd className="m-0 text-white">{detail}</dd>
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
    <Link href={href} className="inline-flex min-h-11 items-center font-serif text-2xl text-white transition-colors duration-200 hover:text-wax">
      {title}
    </Link>
  ) : (
    <h3 className="font-serif text-2xl text-white">{title}</h3>
  )

  return (
    <li>
      <span className="pt-1 font-serif text-sm tabular-nums text-wax">{n}</span>
      <div>
        {heading}
        <p className="mt-2 leading-relaxed text-ink-soft">{body}</p>
      </div>
    </li>
  )
}

function Step({ n, title, body }: { n: string; title: string; body: string }) {
  return (
    <li className="surface p-6">
      <p className="font-serif text-sm tracking-[0.14em] text-wax">{n}</p>
      <h3 className="mt-3 font-serif text-2xl text-white">{title}</h3>
      <p className="mt-2 text-sm leading-relaxed text-ink-soft">{body}</p>
    </li>
  )
}
