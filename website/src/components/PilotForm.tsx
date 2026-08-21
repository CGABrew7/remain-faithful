'use client'

import { useState } from 'react'

export default function PilotForm() {
  const [submitted, setSubmitted] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [form, setForm] = useState({
    name: '',
    church: '',
    email: '',
    role: '',
    numGroups: '',
    referral: '',
  })

  function update(key: string, val: string) {
    setForm((f) => ({ ...f, [key]: val }))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      const res = await fetch('/api/contact', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ...form, type: 'church_pilot' }),
      })
      if (!res.ok) throw new Error('Submit failed')
      setSubmitted(true)
    } catch {
      setError('Something went wrong. Please try again or email us directly at support@remainfaithful.com')
    } finally {
      setLoading(false)
    }
  }

  if (submitted) {
    return (
      <div className="py-10 border-t border-b border-hairline">
        <h3 className="font-serif text-xl font-medium text-ink mb-3">You&apos;re registered</h3>
        <p className="text-ink-soft max-w-sm">
          We&apos;ll reach out personally when Remain Faithful is ready to onboard your church.
        </p>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      <div className="grid sm:grid-cols-2 gap-5">
        <Field label="Your Name" required>
          <input
            type="text"
            required
            value={form.name}
            onChange={(e) => update('name', e.target.value)}
            placeholder="Your name"
            className="input-field"
          />
        </Field>
        <Field label="Church / Organization" required>
          <input
            type="text"
            required
            value={form.church}
            onChange={(e) => update('church', e.target.value)}
            placeholder="Grace Community Church"
            className="input-field"
          />
        </Field>
      </div>
      <Field label="Email Address" required>
        <input
          type="email"
          required
          value={form.email}
          onChange={(e) => update('email', e.target.value)}
          placeholder="pastor@yourchurch.com"
          className="input-field"
        />
      </Field>
      <div className="grid sm:grid-cols-2 gap-5">
        <Field label="Your Role" required>
          <select
            required
            value={form.role}
            onChange={(e) => update('role', e.target.value)}
            className="input-field"
          >
            <option value="">Select role</option>
            <option>Pastor</option>
            <option>Ministry Leader</option>
            <option>Small Group Leader</option>
            <option>Counselor</option>
            <option>Other</option>
          </select>
        </Field>
        <Field label="Estimated Number of Groups" required>
          <select
            required
            value={form.numGroups}
            onChange={(e) => update('numGroups', e.target.value)}
            className="input-field"
          >
            <option value="">Select</option>
            <option>1 group</option>
            <option>2–3 groups</option>
            <option>4–10 groups</option>
            <option>10+ groups</option>
          </select>
        </Field>
      </div>
      <Field label="How did you hear about Remain Faithful?">
        <input
          type="text"
          value={form.referral}
          onChange={(e) => update('referral', e.target.value)}
          placeholder="Word of mouth, social media, conference, etc."
          className="input-field"
        />
      </Field>

      {error && <p className="text-wax text-sm">{error}</p>}

      <button
        type="submit"
        disabled={loading}
        className="btn-wax w-full"
      >
        {loading ? 'Submitting...' : 'Register Your Church'}
      </button>
      <p className="text-center font-mono text-[11px] tracking-[0.08em] uppercase text-ink-faint">
        We&apos;ll be in touch before launch. No spam, ever.
      </p>
    </form>
  )
}

function Field({ label, required, children }: { label: string; required?: boolean; children: React.ReactNode }) {
  return (
    <div>
      <label className="block font-mono text-[12px] tracking-[0.06em] text-ink mb-2">
        {label}
        {required && <span className="text-wax ml-1">*</span>}
      </label>
      {children}
    </div>
  )
}
