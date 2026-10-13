import type { Metadata } from 'next'
import { Suspense } from 'react'
import ResetPasswordForm from '@/components/ResetPasswordForm'

export const metadata: Metadata = {
  title: 'Choose a new password',
  description: 'Set a new Remain Faithful password from the link in your email.',
  robots: { index: false, follow: false },
  alternates: { canonical: 'https://remainfaithful.com/reset-password' },
}

export default function ResetPasswordPage() {
  return (
    <Suspense fallback={null}>
      <ResetPasswordForm />
    </Suspense>
  )
}
