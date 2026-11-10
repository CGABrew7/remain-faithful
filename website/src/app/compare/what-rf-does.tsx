import Link from 'next/link'
import { COMPARE_UPDATED_ON } from '@/lib/competitor-pricing'
import { ALERT_TRIGGER, FAQ_STREAMING } from '@/lib/product-copy'

const GITHUB_URL = 'https://github.com/CGABrew7/remain-faithful'

export function WhatRemainFaithfulDoes() {
  return (
    <section className="mb-14">
      <h2 className="font-serif text-2xl font-bold text-ink mb-5">What Remain Faithful does</h2>
      <p className="text-ink-soft leading-relaxed mb-5 text-pretty">
        Remain Faithful is free peer accountability for adults on iPhone. You invite a person you trust, or a group, and turn on protection. Filtering runs on the phone. Screenshots stay on the phone.
      </p>
      <div className="space-y-4">
        <article className="p-5 rounded-sm border border-wax/25 bg-paper-deep">
          <h3 className="font-semibold text-wax mb-2">Free on iPhone</h3>
          <p className="text-sm text-ink-soft leading-relaxed">
            There is no subscription. The current app is for iPhone on iOS 17 or later. You create an account, invite someone you choose, and turn protection on.
          </p>
        </article>
        <article className="p-5 rounded-sm border border-wax/25 bg-paper-deep">
          <h3 className="font-semibold text-wax mb-2">Alerts, blocking, and optional Deep Scan</h3>
          <p className="text-sm text-ink-soft leading-relaxed">
            {ALERT_TRIGGER} {FAQ_STREAMING}
          </p>
        </article>
        <article className="p-5 rounded-sm border border-wax/25 bg-paper-deep">
          <h3 className="font-semibold text-wax mb-2">What a partner can receive</h3>
          <p className="text-sm text-ink-soft leading-relaxed">
            A partner can see the name and email on your account, your streak, and alert metadata: time, category, severity, and a short summary. That summary does not name the app. It does not include screenshots, browsing history, or page titles. The{' '}
            <Link href="/privacy" className="text-wax hover:underline underline-offset-2">privacy page</Link>
            {' '}has the full list.
          </p>
        </article>
        <article className="p-5 rounded-sm border border-wax/25 bg-paper-deep">
          <h3 className="font-semibold text-wax mb-2">A group, and code you can read</h3>
          <p className="text-sm text-ink-soft leading-relaxed">
            You can invite up to 10 partners. A group is separate and holds up to 12 people. A group needs a name. A written covenant can be added later. The default wording asks members to respond promptly, with grace, when an alert arrives. The iOS app, backend, and website source are public on{' '}
            <a href={GITHUB_URL} className="text-wax hover:underline underline-offset-2" rel="noreferrer">GitHub</a>.
          </p>
        </article>
      </div>
    </section>
  )
}

export function CheckedSources({ links }: { links: { href: string; title: string }[] }) {
  return (
    <section className="mb-14">
      <h2 className="font-serif text-xl font-bold text-ink mb-3">Pages checked on {COMPARE_UPDATED_ON}</h2>
      <ul className="space-y-2 text-sm">
        {links.map((link) => (
          <li key={link.href}>
            <a href={link.href} className="text-wax hover:underline underline-offset-2" rel="noreferrer">
              {link.title}
            </a>
          </li>
        ))}
      </ul>
    </section>
  )
}
