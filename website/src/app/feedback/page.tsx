import type { Metadata } from 'next'
import Link from 'next/link'
import FeedbackForm from '@/components/FeedbackForm'
import { Breadcrumbs } from '@/components/Breadcrumbs'

export const metadata: Metadata = {
  title: 'Suggest an Idea or Report a Problem',
  description:
    'Send the Remain Faithful team an idea for a future update, a way to improve the app, or a problem you ran into.',
  alternates: { canonical: 'https://remainfaithful.com/feedback' },
}

export default function FeedbackPage() {
  return (
    <>
      <div className="mx-auto max-w-5xl px-5 pt-24 sm:px-6">
        <Breadcrumbs items={[{ name: 'Feedback', url: 'https://remainfaithful.com/feedback' }]} />
      </div>
      <section className="pb-20 pt-4">
        <div className="mx-auto grid max-w-5xl gap-12 px-5 sm:px-6 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.6fr)]">
          <div>
            <p className="kicker mb-4">Feedback</p>
            <h1 className="font-serif text-4xl font-bold text-ink sm:text-5xl">
              Suggest an idea or report a problem
            </h1>
            <p className="mt-6 text-lg leading-relaxed text-ink-soft">
              Remain Faithful gets better because people who use it tell us what they need. Send an
              idea for a future update, a way to make something clearer, or a problem you ran into.
            </p>
            <ul className="mt-8 space-y-4 text-ink-soft">
              <li className="border-l-2 border-wax/60 pl-4">A real person reads every message.</li>
              <li className="border-l-2 border-wax/60 pl-4">
                We use your email only to reply. We do not add you to any list.
              </li>
              <li className="border-l-2 border-wax/60 pl-4">
                Not sure what a button does? The{' '}
                <Link href="/app-guide" className="text-wax underline-offset-2 hover:underline">
                  App Guide
                </Link>{' '}
                walks through every screen.
              </li>
            </ul>
            <p className="mt-8 text-sm text-ink-faint">
              Rather email? Write to{' '}
              <a href="mailto:support@remainfaithful.com" className="text-ink-soft hover:text-white">
                support@remainfaithful.com
              </a>
              .
            </p>
          </div>
          <div className="rounded-sm border border-white/10 bg-paper-raised/40 p-6 sm:p-8">
            <FeedbackForm />
          </div>
        </div>
      </section>
    </>
  )
}
