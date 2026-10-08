import type { Metadata } from 'next'
import InviteLanding from '@/components/InviteLanding'

export const metadata: Metadata = {
  title: 'Accept your partner invite',
  description: 'How to accept a Remain Faithful partner invite on iPhone.',
  robots: { index: false, follow: false },
  referrer: 'no-referrer',
}

// Partner invite emails to people who already have an account link here.
export default function AcceptInvitePage({ searchParams }: { searchParams: { type?: string } }) {
  const variant = searchParams?.type === 'group' ? 'group' : 'partner'
  return <InviteLanding variant={variant} />
}
