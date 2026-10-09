'use client'

import { useState } from 'react'

const TYPES = [
  { value: 'Idea', label: 'Idea', hint: 'Something new you would like the app to do' },
  { value: 'Improvement', label: 'Improvement', hint: 'A way to make something better or clearer' },
  { value: 'Problem', label: 'Problem', hint: 'Something is broken or not working right' },
] as const

const AREAS = [
  'Setup / sign in',
  'Home tab',
  'Group tab',
  'Settings tab',
  'App Restrictions / blocking',
  'Deep Scan',
  'Alerts and notifications',
  'Website',
  'Other',
]

type FeedbackType = (typeof TYPES)[number]['value']

export default function FeedbackForm() {
  const [type, setType] = useState<FeedbackType>('Idea')
  const [form, setForm] = useState({
    name: '',
    email: '',
    area: '',
    summary: '',
    details: '',
    device: '',
    company: '', // honeypot, left empty by people
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [submitted, setSubmitted] = useState(false)

  function update(key: keyof typeof form, val: string) {
    setForm((f) => ({ ...f, [key]: val }))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError('')

    const lines = [
      `Type: ${type}`,
      `Area: ${form.area || 'Not given'}`,
      type === 'Problem' ? `iPhone / iOS: ${form.device || 'Not given'}` : null,
      '',
      form.details,
    ].filter((l) => l !== null)

    try {
      const res = await fetch('/api/contact', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: form.name,
          email: form.email,
          subject: `[${type}] ${form.summary}`,
          message: lines.join('\n'),
          type: 'feedback',
          company: form.company,
        }),
      })
      if (!res.ok) {
        let msg = `Error ${res.status}`
        try {
          const body = await res.json()
          if (body?.error) msg = body.error
        } catch {}
        throw new Error(msg)
      }
      setSubmitted(true)
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Unknown error'
      setError(`${msg}. You can also email support@remainfaithful.com.`)
    } finally {
      setLoading(false)
    }
  }

  if (submitted) {
    return (
      <div className="rounded-sm border border-wax/30 bg-paper-deep p-8" role="status">
        <h2 className="font-serif text-2xl font-semibold text-ink">Thank you. We got it.</h2>
        <p className="mt-3 leading-relaxed text-ink-soft">
          Every message is read by a real person. If we need more detail, we will email you at{' '}
          <span className="text-ink">{form.email}</span>.
        </p>
        <button
          type="button"
          className="btn-ghost mt-6"
          onClick={() => {
            setSubmitted(false)
            setForm((f) => ({ ...f, summary: '', details: '', area: '', device: '' }))
          }}
        >
          Send another
        </button>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6" noValidate={false}>
      <fieldset>
        <legend className="mb-3 block font-mono text-[12px] tracking-[0.06em] text-ink">
          What kind of feedback? <span className="text-wax">*</span>
        </legend>
        <div className="grid gap-3 sm:grid-cols-3">
          {TYPES.map((t) => {
            const active = type === t.value
            return (
              <label
                key={t.value}
                className={`flex min-h-11 cursor-pointer flex-col rounded-sm border p-4 transition-colors duration-200 ${
                  active
                    ? 'border-wax bg-wax/10'
                    : 'border-white/15 bg-paper-deep hover:border-white/30'
                }`}
              >
                <input
                  type="radio"
                  name="feedback-type"
                  value={t.value}
                  checked={active}
                  onChange={() => setType(t.value)}
                  className="sr-only"
                />
                <span className={`font-semibold ${active ? 'text-wax' : 'text-ink'}`}>{t.label}</span>
                <span className="mt-1 text-sm leading-snug text-ink-soft">{t.hint}</span>
              </label>
            )
          })}
        </div>
      </fieldset>

      <div>
        <label htmlFor="fb-area" className="mb-2 block font-mono text-[12px] tracking-[0.06em] text-ink">
          Which part of the app?
        </label>
        <select
          id="fb-area"
          value={form.area}
          onChange={(e) => update('area', e.target.value)}
          className="input-field"
        >
          <option value="">Choose one (optional)</option>
          {AREAS.map((a) => (
            <option key={a}>{a}</option>
          ))}
        </select>
      </div>

      <div>
        <label htmlFor="fb-summary" className="mb-2 block font-mono text-[12px] tracking-[0.06em] text-ink">
          {type === 'Problem' ? 'What went wrong, in a few words' : 'Your idea, in a few words'}{' '}
          <span className="text-wax">*</span>
        </label>
        <input
          id="fb-summary"
          type="text"
          required
          maxLength={120}
          value={form.summary}
          onChange={(e) => update('summary', e.target.value)}
          placeholder={type === 'Problem' ? 'e.g. Invite button shows an error' : 'e.g. A weekly check-in reminder'}
          className="input-field"
        />
      </div>

      <div>
        <label htmlFor="fb-details" className="mb-2 block font-mono text-[12px] tracking-[0.06em] text-ink">
          Tell us more <span className="text-wax">*</span>
        </label>
        <textarea
          id="fb-details"
          required
          rows={6}
          maxLength={4000}
          value={form.details}
          onChange={(e) => update('details', e.target.value)}
          placeholder={
            type === 'Problem'
              ? 'What were you doing, what did you expect, and what happened instead?'
              : 'How would it work, and how would it help you?'
          }
          className="input-field resize-y"
        />
        <p className="mt-2 text-sm text-ink-faint">
          Please do not include anything private about someone else.
        </p>
      </div>

      {type === 'Problem' && (
        <div>
          <label htmlFor="fb-device" className="mb-2 block font-mono text-[12px] tracking-[0.06em] text-ink">
            iPhone model and iOS version
          </label>
          <input
            id="fb-device"
            type="text"
            maxLength={80}
            value={form.device}
            onChange={(e) => update('device', e.target.value)}
            placeholder="Optional, e.g. iPhone 15, iOS 18.1"
            className="input-field"
          />
        </div>
      )}

      <div className="grid gap-5 sm:grid-cols-2">
        <div>
          <label htmlFor="fb-name" className="mb-2 block font-mono text-[12px] tracking-[0.06em] text-ink">
            Name <span className="text-wax">*</span>
          </label>
          <input
            id="fb-name"
            type="text"
            required
            maxLength={80}
            autoComplete="name"
            value={form.name}
            onChange={(e) => update('name', e.target.value)}
            placeholder="Your name"
            className="input-field"
          />
        </div>
        <div>
          <label htmlFor="fb-email" className="mb-2 block font-mono text-[12px] tracking-[0.06em] text-ink">
            Email <span className="text-wax">*</span>
          </label>
          <input
            id="fb-email"
            type="email"
            required
            maxLength={120}
            autoComplete="email"
            value={form.email}
            onChange={(e) => update('email', e.target.value)}
            placeholder="you@example.com"
            className="input-field"
          />
          <p className="mt-2 text-sm text-ink-faint">Only used to reply to you.</p>
        </div>
      </div>

      {/* Honeypot: hidden from people, bots tend to fill it in */}
      <div className="hidden" aria-hidden="true">
        <label htmlFor="fb-company">Company</label>
        <input
          id="fb-company"
          type="text"
          tabIndex={-1}
          autoComplete="off"
          value={form.company}
          onChange={(e) => update('company', e.target.value)}
        />
      </div>

      {error && (
        <p className="text-sm text-wax" role="alert">
          {error}
        </p>
      )}

      <button type="submit" disabled={loading} className="btn-wax w-full sm:w-auto">
        {loading ? 'Sending...' : `Send ${type.toLowerCase()}`}
      </button>
    </form>
  )
}
