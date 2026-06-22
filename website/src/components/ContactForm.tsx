'use client'

import { useState } from 'react'

export default function ContactForm() {
  const [submitted, setSubmitted] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [form, setForm] = useState({ name: '', email: '', subject: '', message: '' })

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
        body: JSON.stringify({ ...form, type: 'general' }),
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
      setError(msg + ' — or email us directly at jeff@hanokventures.co')
    } finally {
      setLoading(false)
    }
  }

  if (submitted) {
    return (
      <div className="text-center py-12 rounded-2xl border border-[#C9A84C]/20 bg-[#162235]">
        <div className="text-4xl mb-4">✅</div>
        <h3 className="font-serif text-xl font-semibold text-[#F0EDE8] mb-3">Message Sent</h3>
        <p className="text-[#8A9BB0] text-sm">We&apos;ll get back to you within 2–3 business days.</p>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      <div className="grid sm:grid-cols-2 gap-5">
        <div>
          <label className="block text-sm font-medium text-[#F0EDE8] mb-2">
            Name <span className="text-[#C9A84C]">*</span>
          </label>
          <input
            type="text"
            required
            value={form.name}
            onChange={(e) => update('name', e.target.value)}
            placeholder="Your name"
            className="input-field"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-[#F0EDE8] mb-2">
            Email <span className="text-[#C9A84C]">*</span>
          </label>
          <input
            type="email"
            required
            value={form.email}
            onChange={(e) => update('email', e.target.value)}
            placeholder="you@example.com"
            className="input-field"
          />
        </div>
      </div>
      <div>
        <label className="block text-sm font-medium text-[#F0EDE8] mb-2">
          Subject <span className="text-[#C9A84C]">*</span>
        </label>
        <select
          required
          value={form.subject}
          onChange={(e) => update('subject', e.target.value)}
          className="input-field"
        >
          <option value="">Select a subject</option>
          <option>General Question</option>
          <option>Bug Report</option>
          <option>Partnership Inquiry</option>
          <option>Media / Press</option>
          <option>Other</option>
        </select>
      </div>
      <div>
        <label className="block text-sm font-medium text-[#F0EDE8] mb-2">
          Message <span className="text-[#C9A84C]">*</span>
        </label>
        <textarea
          required
          rows={5}
          value={form.message}
          onChange={(e) => update('message', e.target.value)}
          placeholder="What's on your mind?"
          className="input-field resize-none"
        />
      </div>
      {error && <p className="text-red-400 text-sm">{error}</p>}
      <button
        type="submit"
        disabled={loading}
        className="w-full py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-[box-shadow,scale] duration-200 ease-out active:not-disabled:scale-[0.96] disabled:opacity-70 disabled:cursor-not-allowed"
      >
        {loading ? 'Sending...' : 'Send Message'}
      </button>
    </form>
  )
}
