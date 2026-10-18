import type { Metadata } from 'next'
import { Suspense } from 'react'
import Link from 'next/link'
import AppStoreLink from '@/components/AppStoreLink'
import DonateButton from '@/components/DonateButton'
import DonationSuccessBanner from '@/components/DonationSuccessBanner'
import FaqAccordion from '@/components/FaqAccordion'
import Ornament from '@/components/Ornament'
import PhoneFrame from '@/components/PhoneFrame'
import WaitlistForm from '@/components/WaitlistForm'
import { JsonLd } from '@/components/JsonLd'
import { softwareApplicationSchema, homepageFaqSchema } from '@/lib/structured-data'
import {
  ALERT_TRIGGER,
  FAQ_SLIP,
  FAQ_STREAMING,
  PRIVACY_PROMISE,
} from '@/lib/product-copy'

export const metadata: Metadata = {
  title: 'Remain Faithful | Free Christian Accountability App for iPhone',
  description: 'A partner you trust stands with you. Free peer accountability for adults committed to purity, on iPhone. Screenshots and screen content stay on your phone.',
  keywords: ['free accountability app', 'Christian accountability app', 'purity app', 'accountability partner app', 'church accountability', 'iPhone accountability app', 'on-device AI', 'privacy-first accountability', 'open source accountability'],
  openGraph: {
    title: 'Remain Faithful | Free Christian Accountability App',
    description: 'A partner you trust stands with you. Free on iPhone. Screenshots stay on your phone.',
    url: 'https://remainfaithful.com',
    siteName: 'Remain Faithful',
    type: 'website',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Remain Faithful | Free Christian Accountability App',
    description: 'A partner you trust stands with you. Free on iPhone. Screenshots stay on your phone.',
  },
  alternates: {
    canonical: 'https://remainfaithful.com',
  },
}

const proof = [
  {
    href: '#privacy',
    kicker: 'Privacy',
    title: 'Your screen stays on your phone',
    body: 'Screenshots and screen content stay on the device.',
    lead: true,
  },
  {
    href: '#download',
    kicker: 'iPhone',
    title: 'Free on the App Store',
    body: 'For iPhone on iOS 17 or later. Android is not available yet.',
  },
  {
    href: '/app-guide',
    kicker: 'How it works',
    title: 'Three short steps',
    body: 'Create an account, invite a trusted partner, and turn on protection.',
  },
]

const faqs = [
  {
    q: 'Can my accountability partners see what I was looking at?',
    a: (
      <>
        No. Screenshots and screen content stay on your phone.{' '}
        <a href="#privacy" className="text-white underline decoration-white/25 underline-offset-4 hover:decoration-wax">
          The privacy section
        </a>{' '}
        says what a partner can receive.
      </>
    ),
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
    a: FAQ_SLIP,
  },
  {
    q: 'How is this different from other accountability apps?',
    a: (
      <>
        It is free, it is for adults who choose their own partners, and screenshots stay on the phone.{' '}
        <a href="#privacy" className="text-white underline decoration-white/25 underline-offset-4 hover:decoration-wax">
          The privacy section
        </a>{' '}
        says what a partner can receive. The iOS app, backend, and website source are public on GitHub.
      </>
    ),
  },
  {
    q: 'Can it see streaming video?',
    a: FAQ_STREAMING,
  },
  {
    q: 'Does this work on Android?',
    a: 'Remain Faithful is available for iPhone (iOS 17+). Android support is planned, with no launch date. Leave your email in Get Android updates on this page.',
  },
]

const steps = [
  { n: '1', title: 'Create an account', body: 'Free, on your iPhone.' },
  { n: '2', title: 'Invite a trusted partner', body: 'A friend, spouse, or mentor you choose.' },
  { n: '3', title: 'Turn on protection', body: 'Choose the apps, then turn on Alerts or Blocking.' },
]

export default function HomePage() {
  return (
    <>
      <JsonLd data={softwareApplicationSchema} />
      <JsonLd data={homepageFaqSchema} />

      <Suspense fallback={null}>
        <DonationSuccessBanner />
      </Suspense>

      <section className="shell grid items-center gap-8 pb-12 pt-28 lg:grid-cols-[minmax(0,1.05fr)_minmax(260px,0.8fr)] lg:gap-x-16 lg:gap-y-6 lg:pb-16 lg:pt-32">
        <div className="lg:col-start-1 lg:row-start-1">
          <p className="kicker rise">Free · iPhone (iOS 17+)</p>
          <Ornament className="rise rise-2 mb-6 mt-5" />
          <h1 className="rise rise-2 font-serif text-[2.75rem] font-bold leading-[1.04] tracking-[-0.03em] text-white sm:text-6xl lg:text-[4.15rem]">
            Remain Faithful
          </h1>
          <p className="rise rise-3 mt-6 max-w-[22ch] font-serif text-2xl leading-snug text-white/90 sm:text-[1.7rem]">
            A partner you trust stands with you.
          </p>
          <p className="mt-4 max-w-[42ch] text-[1.05rem] leading-relaxed text-ink-soft">
            Free peer accountability for adults committed to purity. You do not have to carry the hard moments alone.
          </p>
          <div className="rise rise-4 mt-8 flex flex-col gap-3 sm:flex-row sm:flex-wrap">
            <AppStoreLink />
            <Link href="/app-guide" className="btn-ghost">
              See how it works
              <Arrow />
            </Link>
            <a href="#android-updates" className="btn-ghost">
              Get Android updates
              <Arrow />
            </a>
          </div>
        </div>
        <div className="rise rise-3 lg:col-start-2 lg:row-start-1 lg:justify-self-end">
          <PhoneFrame
            src="/screens/welcome.png"
            alt="Remain Faithful iPhone welcome screen: Remain Faithful, accountability between trusted friends, and Begin Your Journey."
            priority
          />
        </div>
      </section>

      <section id="how-it-works" className="shell scroll-mt-24 pb-8" aria-labelledby="how-strip-heading">
        <h2 id="how-strip-heading" className="font-serif text-3xl font-bold text-white sm:text-4xl">
          How it works
        </h2>
        <ol className="stepper mt-6">
          {steps.map((step) => (
            <li key={step.n} className="surface p-6">
              <p className="font-serif text-sm tabular-nums tracking-[0.14em] text-wax">{step.n}</p>
              <h3 className="mt-3 font-serif text-2xl text-white">{step.title}</h3>
              <p className="mt-2 text-sm leading-relaxed text-ink-soft">{step.body}</p>
            </li>
          ))}
        </ol>
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
            <h2 className="font-serif text-4xl font-bold text-white sm:text-5xl">Get the app</h2>
            <p className="mt-4 max-w-[42ch] text-lg leading-relaxed text-ink-soft">
              Free on iPhone. Invite a partner you trust, turn on protection, and keep going together.
            </p>
            <dl className="mt-8 grid gap-5 sm:grid-cols-3">
              <Spec term="Platform" detail="iPhone, iOS 17+" />
              <Spec term="Price" detail="Free forever" />
              <Spec term="Who it is for" detail="Adults, 18+" />
            </dl>
            <div className="mt-8 flex flex-col gap-2 sm:flex-row sm:flex-wrap">
              <AppStoreLink />
              <Link href="/app-guide" className="btn-ghost">
                See how it works
                <Arrow />
              </Link>
              <a href="#android-updates" className="btn-ghost">
                Get Android updates
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
              body="Screenshots and screen content stay on your phone."
              href="#privacy"
            />
            <Feature
              n="03"
              title="Alerts and Blocking"
              body="Alerts and Blocking are separate switches. The next section says when a partner is told."
            />
            <Feature
              n="04"
              title="Always free"
              body="Remain Faithful is free today and will remain free forever. No subscription tiers, no paywalls, no premium features. Sustained entirely by voluntary donations."
            />
          </ol>
        </div>
      </section>

      <section className="band section" aria-labelledby="alert-heading">
        <div className="shell">
          <div className="grid items-center gap-10 lg:grid-cols-[minmax(0,1fr)_260px]">
            <div>
              <p className="kicker mb-4">Partner alerts</p>
              <h2 id="alert-heading" className="font-serif text-4xl font-bold text-white sm:text-5xl">
                What your partner hears
              </h2>
              <p className="mt-4 max-w-[52ch] text-lg leading-relaxed text-ink-soft">
                {ALERT_TRIGGER}
              </p>
              <Link href="/app-guide" className="btn-ghost mt-8">
                See how it works
                <Arrow />
              </Link>
            </div>
            <PhoneFrame
              src="/screens/home.png"
              alt="Remain Faithful home screen with monitoring active, a 47-day clean streak, and no recent flags."
            />
          </div>

          <div id="privacy" className="mt-14 scroll-mt-24 grid gap-4 lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)]">
            <div className="surface p-7 sm:p-8">
              <p className="kicker mb-4">Privacy</p>
              <h3 className="font-serif text-3xl text-white">What leaves your phone</h3>
              <p className="mt-3 leading-relaxed text-ink-soft">
                {PRIVACY_PROMISE}
              </p>
              <div className="mt-6 flex flex-col gap-2 sm:flex-row sm:flex-wrap">
                <Link href="/privacy" className="btn-ghost">
                  Privacy policy
                  <Arrow />
                </Link>
                <Link href="/privacy-architecture" className="btn-ghost">
                  Full Privacy Architecture
                  <Arrow />
                </Link>
              </div>
            </div>
            <div className="split-panel">
              <div className="p-7 sm:p-8">
                <h4 className="kicker mb-4">Stays on your phone</h4>
                <ul className="space-y-2.5 text-sm leading-relaxed text-ink-soft">
                  <li>Screenshots and screen content</li>
                  <li>Browsing history and page content</li>
                  <li>Passwords and financial data</li>
                  <li>Message content</li>
                  <li>Photos and videos</li>
                </ul>
              </div>
              <div className="border-t border-white/10 p-7 sm:border-l sm:border-t-0 sm:p-8">
                <h4 className="kicker mb-4">Can leave the phone</h4>
                <ul className="space-y-2.5 text-sm leading-relaxed text-white">
                  <li>Category</li>
                  <li>Severity</li>
                  <li>Timestamp</li>
                  <li>A short system summary</li>
                  <li>A heartbeat while protection is on</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="android-updates" className="section scroll-mt-24">
        <div id="waitlist" className="shell grid scroll-mt-24 items-center gap-10 lg:grid-cols-2">
          <div>
            <p className="kicker mb-4">Android</p>
            <h2 className="font-serif text-4xl font-bold text-white sm:text-5xl">Get Android updates</h2>
            <p className="mt-4 max-w-[40ch] text-lg leading-relaxed text-ink-soft">
              Remain Faithful is on the App Store for iPhone. Leave your email and we will write when Android is ready. There is no launch date yet.
            </p>
          </div>
          <div className="surface p-6 sm:p-8">
            <WaitlistForm variant="default" buttonText="Get Android updates" />
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
                We are committed to never charging for accountability. Your donation funds server costs, development, and ministry outreach.
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
