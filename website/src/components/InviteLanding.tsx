import Link from 'next/link'
import Ornament from '@/components/Ornament'
import AppStoreLink from '@/components/AppStoreLink'
import AcceptInviteForm from '@/components/AcceptInviteForm'

export type InviteVariant = 'group' | 'partner' | 'link'

const copy: Record<InviteVariant, { kicker: string; title: string; lede: string; steps: { title: string; body: string }[] }> = {
  group: {
    kicker: 'Group invite',
    title: 'Join your accountability group',
    lede: 'Someone you know invited you to their Remain Faithful group. Accept with the same email, then confirm that email.',
    steps: [
      {
        title: 'Get the app',
        body: 'Remain Faithful is on the App Store for iPhone (iOS 17+).',
      },
      {
        title: 'Create your account with this email',
        body: 'Sign up with the same email this invite was sent to, then open the confirmation email. The group invite is accepted after you confirm. Signing in alone does not add you.',
      },
      {
        title: 'Already have an account',
        body: 'Use the form on this page to sign in and accept. Opening the app does not accept the invite by itself.',
      },
    ],
  },
  partner: {
    kicker: 'Partner invite',
    title: 'Accept your partner invite',
    lede: 'Someone you know asked you to be their accountability partner. Accept with the same email, then confirm that email.',
    steps: [
      {
        title: 'Get the app',
        body: 'Remain Faithful is on the App Store for iPhone (iOS 17+).',
      },
      {
        title: 'Create your account with this email',
        body: 'Sign up with the same email this invite was sent to, then open the confirmation email. The partner invite is accepted after you confirm. Signing in alone does not connect you.',
      },
      {
        title: 'Already have an account',
        body: 'Use the form on this page to sign in and accept. Opening the app does not accept the invite by itself.',
      },
    ],
  },
  link: {
    kicker: 'Group invite',
    title: 'You were invited to a group',
    lede: 'Someone shared their Remain Faithful group with you. This link does not add you to the group.',
    steps: [
      {
        title: 'Get the app',
        body: 'Remain Faithful is on the App Store for iPhone (iOS 17+).',
      },
      {
        title: 'Ask to be invited by email',
        body: 'Send the person who shared this link the email address you will use. They invite you from the Group tab.',
      },
      {
        title: 'Confirm that email',
        body: 'Create the account with that email and confirm it. The invite stays pending until you accept. This page does not add you.',
      },
    ],
  },
}

export default function InviteLanding({
  variant,
  token,
}: {
  variant: InviteVariant
  token?: string
}) {
  const c = copy[variant]
  return (
    <section className="section pt-28 lg:pt-32">
      <div className="shell max-w-3xl">
        <div className="surface relative overflow-hidden px-6 py-12 sm:px-12 sm:py-14">
          <div
            className="pointer-events-none absolute inset-0"
            aria-hidden="true"
            style={{ background: 'radial-gradient(60% 70% at 50% 0%, rgba(209,171,76,0.14), transparent 70%)' }}
          />
          <div className="relative text-center">
            <p className="kicker">{c.kicker}</p>
            <Ornament className="mx-auto mb-6 mt-5 justify-center" />
            <h1 className="font-serif text-4xl font-bold leading-tight text-white sm:text-5xl">{c.title}</h1>
            <p className="mx-auto mt-4 max-w-[46ch] text-lg leading-relaxed text-ink-soft">{c.lede}</p>
          </div>

          <ol className="relative mt-10 space-y-4">
            {c.steps.map((step, i) => (
              <li key={step.title} className="flex gap-5 rounded-sm border border-hairline bg-paper-deep/70 p-5">
                <span className="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-sm bg-wax text-sm font-bold text-navy">
                  {i + 1}
                </span>
                <div>
                  <h2 className="mb-1 font-semibold text-white">{step.title}</h2>
                  <p className="text-sm leading-relaxed text-ink-soft">{step.body}</p>
                </div>
              </li>
            ))}
          </ol>

          <p className="relative mt-8 text-center text-sm leading-relaxed text-ink-faint">
            Remain Faithful is for adults (18+) who choose accountability for themselves. Each adult authorizes Family
            Controls on their own iPhone. It is not parental controls.
          </p>

          {token && variant !== 'link' && (
            <AcceptInviteForm token={token} type={variant === 'group' ? 'group' : 'partner'} />
          )}

          <div className="relative mt-8 flex flex-col justify-center gap-3 sm:flex-row">
            <AppStoreLink />
            <Link href="/how-it-works" className="btn-ghost">
              How it works
            </Link>
          </div>
        </div>
      </div>
    </section>
  )
}
