'use client'

import { useEffect, useRef, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import AppStoreLink from '@/components/AppStoreLink'
import { useFormShield } from '@/components/FormShield'

const confirmedKey = (token: string) => `rf-email-confirmed:${token}`

export default function VerifyEmailForm() {
  const verifyShield = useFormShield('verify-email')
  const resendShield = useFormShield('resend-confirmation')
  const verifyFields = useRef(verifyShield.fields)
  const verifyReset = useRef(verifyShield.reset)
  verifyFields.current = verifyShield.fields
  verifyReset.current = verifyShield.reset
  const params = useSearchParams()
  const token = params.get('token') || ''
  const purpose = params.get('purpose')
  const [status, setStatus] = useState<'working' | 'done' | 'error'>('working')
  const [message, setMessage] = useState('Confirming your email.')
  const [email, setEmail] = useState('')
  const [resendNote, setResendNote] = useState('')
  const [sending, setSending] = useState(false)

  useEffect(() => {
    if (!token) {
      setStatus('error')
      setMessage('Open the confirmation link from your email.')
      return
    }
    if (typeof window !== 'undefined' && window.sessionStorage.getItem(confirmedKey(token)) === '1') {
      setStatus('done')
      setMessage('This email is already confirmed.')
      return
    }
    let cancelled = false
    const started = Date.now()
    const timer = window.setInterval(() => {
      if (cancelled) return
      const fields = verifyFields.current()
      const siteKey = process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY || ''
      const elapsed = Date.now() - fields.loadedAt
      if (fields.loadedAt === 0 || elapsed < 2000) return
      if (siteKey && !fields.turnstileToken && Date.now() - started < 20000) return
      window.clearInterval(timer)
      void (async () => {
        try {
          const res = await fetch('/api/verify-email', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ token, ...fields }),
          })
          const body = await res.json().catch(() => ({}))
          if (cancelled) return
          if (!res.ok) {
            setStatus('error')
            setMessage(body.error || 'That confirmation link is not valid. You can send a new one below.')
            return
          }
          window.sessionStorage.setItem(confirmedKey(token), '1')
          setStatus('done')
          setMessage(
            purpose === 'change'
              ? 'Your email address is updated. Sign in again in the app.'
              : 'Email confirmed. Pending partner and group invites for this address are now accepted.',
          )
        } catch {
          if (!cancelled) {
            setStatus('error')
            setMessage('We could not reach the server. Check your connection and try again.')
          }
        } finally {
          verifyReset.current()
        }
      })()
    }, 200)
    return () => {
      cancelled = true
      window.clearInterval(timer)
    }
  }, [token, purpose])

  async function resend(e: React.FormEvent) {
    e.preventDefault()
    setSending(true)
    setResendNote('')
    try {
      const res = await fetch('/api/resend-confirmation', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, ...resendShield.fields() }),
      })
      const body = await res.json().catch(() => ({}))
      if (!res.ok) {
        setResendNote(body.error || 'We could not send that email. Try again in a little while.')
        return
      }
      setResendNote('If that address still needs confirmation, a new email is on the way.')
    } catch {
      setResendNote('We could not reach the server. Check your connection and try again.')
    } finally {
      resendShield.reset()
      setSending(false)
    }
  }

  return (
    <section className="section pt-28 lg:pt-32">
      <div className="shell max-w-lg">
        <div className="surface relative px-6 py-12 sm:px-10" data-clarity-mask="true">
          {verifyShield.slot}
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
          {status === 'error' && (
            <form onSubmit={resend} className="relative mt-8">
              {resendShield.slot}
              <div className="space-y-4">
              <label htmlFor="resend-email" className="block text-sm text-white">
                Resend confirmation email
              </label>
              <input
                id="resend-email"
                type="email"
                required
                autoComplete="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input-field w-full"
              />
              {resendNote && <p className="text-sm text-ink-soft">{resendNote}</p>}
              <button type="submit" className="btn-wax w-full" disabled={sending}>
                {sending ? 'Sending...' : 'Resend confirmation email'}
              </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </section>
  )
}
