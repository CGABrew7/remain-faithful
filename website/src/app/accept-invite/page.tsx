import type { Metadata } from 'next'
import InviteLanding from '@/components/InviteLanding'

export const metadata: Metadata = {
  title: 'Accept your partner invite',
  description: 'Accept a Remain Faithful partner or group invite. Signing in alone does not accept it.',
  robots: { index: false, follow: false },
  referrer: 'no-referrer',
  alternates: { canonical: 'https://remainfaithful.com/accept-invite' },
}

export default async function AcceptInvitePage({
  searchParams,
}: {
  searchParams: Promise<{ type?: string; token?: string }>
}) {
  const query = await searchParams
  const variant = query?.type === 'group' ? 'group' : 'partner'
  return <InviteLanding variant={variant} token={query?.token} />
}
