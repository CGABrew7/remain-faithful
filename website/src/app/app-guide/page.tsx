import type { Metadata } from 'next'
import Link from 'next/link'
import PhoneFrame from '@/components/PhoneFrame'
import { Breadcrumbs } from '@/components/Breadcrumbs'
import { guideParts, type GuideScreen } from './guide-content'

export const metadata: Metadata = {
  title: 'App Guide: Every Screen and Button Explained',
  description:
    'A plain walkthrough of the Remain Faithful iPhone app. What each screen shows, what each button does, and why it is there.',
  alternates: { canonical: 'https://remainfaithful.com/app-guide' },
}

export default function AppGuidePage() {
  return (
    <>
      <div className="mx-auto max-w-6xl px-5 pt-24 sm:px-6">
        <Breadcrumbs items={[{ name: 'App Guide', url: 'https://remainfaithful.com/app-guide' }]} />
      </div>

      {/* Hero */}
      <section className="border-b border-white/10 pb-14 pt-4">
        <div className="mx-auto max-w-6xl px-5 sm:px-6">
          <p className="kicker mb-4">App Guide</p>
          <h1 className="max-w-3xl font-serif text-4xl font-bold text-ink sm:text-5xl">
            Every screen and button, explained
          </h1>
          <p className="mt-6 max-w-2xl text-lg leading-relaxed text-ink-soft">
            This page walks through the Remain Faithful iPhone app one screen at a time. For each
            button you will see what it does and why it is there. Keep it open next to your phone the
            first time you set things up.
          </p>

          <div className="mt-10 grid gap-4 md:grid-cols-2">
            <div className="rounded-sm border border-wax/30 bg-paper-deep p-6">
              <h2 className="font-serif text-lg font-semibold text-ink">What leaves your phone</h2>
              <p className="mt-2 text-[0.98rem] leading-relaxed text-ink-soft">
                Only two things. Alert details (the type of alert, how serious it was, the time, and a
                short summary the app writes) and a regular check-in that says the app is still
                running. Screenshots and what is on your screen stay on your phone.
              </p>
            </div>
            <div className="rounded-sm border border-white/10 bg-paper-deep p-6">
              <h2 className="font-serif text-lg font-semibold text-ink">Who it is for</h2>
              <p className="mt-2 text-[0.98rem] leading-relaxed text-ink-soft">
                Adults (18+) who choose to be accountable on their own iPhone, with friends they
                trust. It is not a parental control app. It needs iOS 17 or later.
              </p>
            </div>
          </div>
        </div>
      </section>

      <div className="mx-auto max-w-6xl px-5 py-14 sm:px-6 lg:grid lg:grid-cols-[200px_minmax(0,1fr)] lg:gap-14">
        {/* Table of contents */}
        <nav aria-label="On this page" className="mb-12 lg:mb-0">
          <div className="lg:sticky lg:top-24">
            <p className="kicker mb-4">On this page</p>
            <ol className="flex flex-wrap gap-2 lg:flex-col lg:gap-0">
              {guideParts.map((part) => (
                <li key={part.id} className="lg:mb-3">
                  <a
                    href={`#${part.id}`}
                    className="inline-flex min-h-11 items-center rounded-sm border border-white/10 px-3 text-sm text-ink-soft transition-colors duration-200 hover:border-wax/40 hover:text-white lg:min-h-0 lg:border-0 lg:px-0 lg:text-[0.98rem] lg:font-medium lg:text-ink"
                  >
                    {part.label}
                  </a>
                  <ul className="mt-1.5 hidden space-y-1 border-l border-white/10 pl-3 lg:block">
                    {part.screens.map((screen) => (
                      <li key={screen.id}>
                        <a
                          href={`#${screen.id}`}
                          className="text-sm text-ink-faint transition-colors duration-200 hover:text-wax"
                        >
                          {screen.title}
                        </a>
                      </li>
                    ))}
                  </ul>
                </li>
              ))}
              <li className="lg:mt-3">
                <Link
                  href="/feedback"
                  className="inline-flex min-h-11 items-center rounded-sm border border-wax/30 px-3 text-sm text-wax transition-colors duration-200 hover:text-white lg:min-h-0 lg:border-0 lg:px-0"
                >
                  Suggest an idea
                </Link>
              </li>
            </ol>
          </div>
        </nav>

        <div>
          {guideParts.map((part, i) => (
            <section
              key={part.id}
              id={part.id}
              className={`scroll-mt-24 ${i > 0 ? 'mt-20 border-t border-white/10 pt-16' : ''}`}
            >
              <p className="kicker mb-3">Part {i + 1}</p>
              <h2 className="font-serif text-3xl font-bold text-ink sm:text-4xl">{part.title}</h2>
              <p className="mt-4 max-w-2xl text-lg leading-relaxed text-ink-soft">{part.summary}</p>

              <div className="mt-10 space-y-16">
                {part.screens.map((screen) => (
                  <ScreenBlock key={screen.id} screen={screen} />
                ))}
              </div>
            </section>
          ))}

          {/* Feedback CTA */}
          <section className="mt-20 rounded-sm border border-wax/30 bg-paper-deep p-8 sm:p-10">
            <h2 className="font-serif text-2xl font-bold text-ink">Something missing or confusing?</h2>
            <p className="mt-3 max-w-xl leading-relaxed text-ink-soft">
              Tell us. Ideas for new features, ways to make a screen clearer, or something that is not
              working. We read every message.
            </p>
            <div className="mt-6 flex flex-wrap gap-4">
              <Link href="/feedback" className="btn-wax">
                Suggest an idea or report a problem
              </Link>
              <Link href="/how-it-works" className="btn-ghost">
                How the protection works
              </Link>
            </div>
          </section>
        </div>
      </div>
    </>
  )
}

function ScreenBlock({ screen }: { screen: GuideScreen }) {
  const hasShots = !!screen.shots && screen.shots.length > 0
  return (
    <article id={screen.id} className="scroll-mt-24">
      <div
        className={
          hasShots
            ? 'grid items-start gap-10 md:grid-cols-[minmax(0,1fr)_240px] xl:grid-cols-[minmax(0,1fr)_260px]'
            : ''
        }
      >
        <div>
          <p className="font-mono text-[11px] uppercase tracking-[0.1em] text-ink-faint">{screen.where}</p>
          <h3 className="mt-2 font-serif text-2xl font-semibold text-ink">{screen.title}</h3>
          <p className="mt-3 leading-relaxed text-ink-soft">{screen.intro}</p>

          <dl className="mt-6 divide-y divide-white/10 rounded-sm border border-white/10 bg-paper-deep">
            {screen.buttons.map((b) => (
              <div key={b.name} className="px-5 py-4">
                <dt className="font-semibold text-wax">{b.name}</dt>
                <dd className="mt-1 leading-relaxed text-ink-soft">
                  {b.does}
                  {b.why && <span className="mt-1.5 block text-ink">{b.why}</span>}
                </dd>
              </div>
            ))}
          </dl>

          {screen.note && (
            <p className="mt-4 border-l-2 border-wax/60 pl-4 text-[0.95rem] leading-relaxed text-ink-soft">
              {screen.note}
            </p>
          )}
        </div>

        {hasShots && (
          <div
            className={`mx-auto grid w-full gap-8 [&_.phone]:!w-full ${
              screen.shots!.length > 1
                ? 'max-w-[440px] grid-cols-2 gap-5 md:max-w-[260px] md:grid-cols-1 md:gap-8'
                : 'max-w-[240px] md:max-w-[260px]'
            }`}
          >
            {screen.shots!.map((shot) => (
              <div key={shot.src}>
                <PhoneFrame src={shot.src} alt={shot.alt} />
                <p className="mt-3 text-center text-sm text-ink-faint">{shot.caption}</p>
              </div>
            ))}
          </div>
        )}
      </div>
    </article>
  )
}
