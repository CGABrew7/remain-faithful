'use client'

import { useState } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'

export default function ResetPasswordForm() {
  const params = useSearchParams()
  const token = params.get('token') || ''
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [done, setDone] = useState(false)
  const [email, setEmail] = useState('')
  const [resetNote, setResetNote] = useState('')
  const [sending, setSending] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    if (!token) {
      setError('Open the reset link from your email. This page needs that link.')
      return
    }
    setLoading(true)
    try {
      const res = await fetch('/api/reset-password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token, password }),
      })
      const body = await res.json().catch(() => ({}))
      if (!res.ok) {
        const raw = body.error || 'Could not reset the password.'
        if (String(raw).toLowerCase().includes('invalid or expired')) {
          throw new Error('This reset link is no longer active. Send a new one below.')
        }
        if (String(raw).toLowerCase().includes('failed to fetch')) {
          throw new Error('We could not reach the server. Check your connection and try again.')
        }
        throw new Error(raw)
      }
      setDone(true)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Could not reset the password.'
      setError(message.includes('Failed to fetch') ? 'We could not reach the server. Check your connection and try again.' : message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="section pt-28 lg:pt-32">
      <div className="shell max-w-lg">
        <div className="surface px-6 py-12 sm:px-10">
          <p className="kicker">Account</p>
          <h1 className="mt-4 font-serif text-4xl font-bold text-white">Choose a new password</h1>
          <p className="mt-4 text-ink-soft">
            The link in your email works for one hour. After you save a new password, older sign-ins on other devices stop working.
          </p>
          {done ? (
            <p className="mt-8 text-white" role="status">
              Password saved. Open Remain Faithful on your iPhone and sign in with the new password.
            </p>
          ) : (
            <form onSubmit={handleSubmit} className="mt-8 space-y-4">
              <div>
                <label htmlFor="new-password" className="mb-2 block text-sm text-white">
                  New password
                </label>
                <input
                  id="new-password"
                  type="password"
                  required
                  minLength={8}
                  autoComplete="new-password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="input-field w-full"
                />
              </div>
              {error && (
                <p className="form-error text-sm" role="alert">
                  {error}
                </p>
              )}
              <button type="submit" className="btn-wax w-full" disabled={loading}>
                {loading ? 'Saving...' : 'Save password'}
              </button>
            </form>
          )}
          {!done && (
            <form
              className="mt-8 space-y-3 border-t border-white/10 pt-6"
              onSubmit={async (e) => {
                e.preventDefault()
                setSending(true)
                setResetNote('')
                try {
                  const res = await fetch('/api/forgot-password', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email }),
                  })
                  const body = await res.json().catch(() => ({}))
                  if (!res.ok) {
                    setResetNote(body.error || 'We could not send that email. Try again in a little while.')
                    return
                  }
                  setResetNote('If that address has an account, a new reset link is on the way.')
                } catch {
                  setResetNote('We could not reach the server. Check your connection and try again.')
                } finally {
                  setSending(false)
                }
              }}
            >
              <label htmlFor="reset-email" className="block text-sm text-white">
                Send a new reset link
              </label>
              <input
                id="reset-email"
                type="email"
                required
                autoComplete="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input-field w-full"
              />
              {resetNote && <p className="text-sm text-ink-soft">{resetNote}</p>}
              <button type="submit" className="btn-ghost w-full" disabled={sending}>
                {sending ? 'Sending...' : 'Send a new reset link'}
              </button>
            </form>
          )}
          <p className="mt-6 text-sm text-ink-faint">
            <Link href="/app-guide" className="underline underline-offset-2">
              App guide
            </Link>
          </p>
        </div>
      </div>
    </section>
  )
}
