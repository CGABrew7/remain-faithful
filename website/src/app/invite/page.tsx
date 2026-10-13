import type { Metadata } from 'next'
import InviteLanding from '@/components/InviteLanding'

export const metadata: Metadata = {
  title: 'Accept your invite',
  description: 'Accept a Remain Faithful group or partner invite. Signing in alone does not accept it.',
  robots: { index: false, follow: false },
  referrer: 'no-referrer',
  alternates: { canonical: 'https://remainfaithful.com/invite' },
}

export default async function InvitePage({
  searchParams,
}: {
  searchParams: Promise<{ type?: string; token?: string }>
}) {
  const query = await searchParams
  const variant = query?.type === 'partner' ? 'partner' : 'group'
  return <InviteLanding variant={variant} token={query?.token} />
}
