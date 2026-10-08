import type { Metadata } from 'next'
import InviteLanding from '@/components/InviteLanding'

export const metadata: Metadata = {
  title: 'Accept your invite',
  description: 'How to accept a Remain Faithful group or partner invite on iPhone.',
  robots: { index: false, follow: false },
  referrer: 'no-referrer',
}

// Emailed invite links land here: /invite?token=...&type=group|partner.
// The token is never read or sent anywhere from the website; the app redeems
// invites automatically when the invited email creates an account.
export default function InvitePage({ searchParams }: { searchParams: { type?: string } }) {
  const variant = searchParams?.type === 'partner' ? 'partner' : 'group'
  return <InviteLanding variant={variant} />
}
