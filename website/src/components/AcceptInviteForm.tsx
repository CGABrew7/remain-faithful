'use client'

import { useState } from 'react'

export default function AcceptInviteForm({
  token,
  type,
}: {
  token: string
  type: 'partner' | 'group'
}) {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [company, setCompany] = useState('')
  const [loading, setLoading] = useState<'accept' | 'decline' | null>(null)
  const [error, setError] = useState('')
  const [done, setDone] = useState('')

  async function submit(decision: 'accept' | 'decline') {
    setError('')
    setLoading(decision)
    try {
      const res = await fetch('/api/accept-invite', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password, token, type, decision, company }),
      })
      const body = await res.json().catch(() => ({}))
      if (!res.ok) throw new Error(body.error || 'Could not update this invite.')
      setDone(decision === 'decline' ? 'Invite declined.' : 'Invite accepted. Open the app to continue.')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not update this invite.')
    } finally {
      setLoading(null)
    }
  }

  if (done) {
    return (
      <p className="mt-8 text-center text-white" role="status">
        {done}
      </p>
    )
  }

  return (
    <form
      className="relative mt-8 space-y-4 rounded-[14px] border border-white/10 bg-[#0C1D4C] p-5"
      onSubmit={(e) => {
        e.preventDefault()
        void submit('accept')
      }}
    >
      <h2 className="font-serif text-2xl text-white">Already have an account?</h2>
      <p className="text-sm leading-relaxed text-ink-soft">
        Sign in with the email this invite was sent to, then accept. Signing in on the iPhone by itself does not accept the invite.
      </p>
      <div className="absolute left-[-9999px] h-0 w-0 overflow-hidden" aria-hidden="true">
        <label htmlFor="invite-company">Company</label>
        <input id="invite-company" name="company" tabIndex={-1} autoComplete="off" value={company} onChange={(e) => setCompany(e.target.value)} />
      </div>
      <div>
        <label htmlFor="invite-email" className="mb-2 block text-sm text-white">
          Email
        </label>
        <input
          id="invite-email"
          type="email"
          required
          autoComplete="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="input-field w-full"
        />
      </div>
      <div>
        <label htmlFor="invite-password" className="mb-2 block text-sm text-white">
          Password
        </label>
        <input
          id="invite-password"
          type="password"
          required
          autoComplete="current-password"
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
      <div className="flex flex-col gap-3 sm:flex-row">
        <button type="submit" className="btn-wax" disabled={loading !== null}>
          {loading === 'accept' ? 'Accepting...' : 'Accept invite'}
        </button>
        <button
          type="button"
          className="btn-ghost"
          disabled={loading !== null}
          onClick={() => void submit('decline')}
        >
          {loading === 'decline' ? 'Declining...' : 'Decline'}
        </button>
      </div>
    </form>
  )
}
