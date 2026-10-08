import type { Metadata } from 'next'
import InviteLanding from '@/components/InviteLanding'

export const metadata: Metadata = {
  title: 'Join a group',
  description: 'How to join a Remain Faithful accountability group on iPhone.',
  robots: { index: false, follow: false },
  referrer: 'no-referrer',
}

// Shared from the app's Group tab invite sheet (www.remainfaithful.com/join/{groupId}).
// The id is not displayed or used.
export default function JoinGroupPage() {
  return <InviteLanding variant="link" />
}
