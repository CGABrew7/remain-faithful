'use client'

import { useEffect, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import AppStoreLink from '@/components/AppStoreLink'

export default function VerifyEmailForm() {
  const params = useSearchParams()
  const token = params.get('token') || ''
  const purpose = params.get('purpose')
  const [status, setStatus] = useState<'working' | 'done' | 'error'>('working')
  const [message, setMessage] = useState('Confirming your email.')

  useEffect(() => {
    if (!token) {
      setStatus('error')
      setMessage('Open the confirmation link from your email.')
      return
    }
    let cancelled = false
    ;(async () => {
      try {
        const res = await fetch('/api/verify-email', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ token }),
        })
        const body = await res.json().catch(() => ({}))
        if (cancelled) return
        if (!res.ok) {
          setStatus('error')
          setMessage(body.error || 'That confirmation link is not valid.')
          return
        }
        setStatus('done')
        setMessage(
          purpose === 'change'
            ? 'Your email address is updated. Sign in again in the app.'
            : 'Email confirmed. Pending partner and group invites for this address are now accepted.',
        )
      } catch {
        if (!cancelled) {
          setStatus('error')
          setMessage('Could not reach the account service. Try the link again.')
        }
      }
    })()
    return () => {
      cancelled = true
    }
  }, [token, purpose])

  return (
    <section className="section pt-28 lg:pt-32">
      <div className="shell max-w-lg">
        <div className="surface px-6 py-12 sm:px-10">
          <p className="kicker">Account</p>
          <h1 className="mt-4 font-serif text-4xl font-bold text-white">Confirm your email</h1>
          <p className={`mt-4 ${status === 'error' ? 'form-error' : 'text-ink-soft'}`} role={status === 'error' ? 'alert' : 'status'}>
            {message}
          </p>
          {status === 'done' && (
            <div className="mt-8">
              <AppStoreLink />
            </div>
          )}
        </div>
      </div>
    </section>
  )
}
