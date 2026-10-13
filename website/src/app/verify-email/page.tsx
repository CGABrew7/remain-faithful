import type { Metadata } from 'next'
import { Suspense } from 'react'
import VerifyEmailForm from '@/components/VerifyEmailForm'

export const metadata: Metadata = {
  title: 'Confirm your email',
  description: 'Confirm the email address on your Remain Faithful account.',
  robots: { index: false, follow: false },
  alternates: { canonical: 'https://remainfaithful.com/verify-email' },
}

export default function VerifyEmailPage() {
  return (
    <Suspense fallback={null}>
      <VerifyEmailForm />
    </Suspense>
  )
}
