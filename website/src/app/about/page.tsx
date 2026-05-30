'use client'

import { useState } from 'react'
import Link from 'next/link'

export default function AboutPage() {
  return (
    <>
      {/* Hero */}
      <section className="pt-32 pb-20 border-b border-[#1E3050]">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-4">Our Story</p>
              <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-6">
                Built for the Struggle
              </h1>
              <p className="text-[#8A9BB0] text-lg leading-relaxed">
                We built Remain Faithful because the tools for accountability were either expensive, invasive, or ineffective. And because the men we know deserved something better.
              </p>
            </div>
            <div
              className="rounded-2xl p-8 border border-[#C9A84C]/20"
              style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}
            >
              <p className="font-serif text-lg text-[#F0EDE8] italic leading-relaxed">
                &ldquo;Accountability without relationship is just surveillance. We built RF to be the former, never the latter.&rdquo;
              </p>
              <p className="text-sm text-[#8A9BB0] mt-4">— The RF Team</p>
            </div>
          </div>
        </div>
      </section>

      {/* Mission */}
      <section className="py-20">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 gap-12">
            <div>
              <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Why We Built This</h2>
              <p className="text-[#8A9BB0] leading-relaxed mb-4">
                The existing accountability software market has a few products, but they tend toward one of two failure modes: either they&apos;re heavy-handed surveillance tools that damage trust, or they&apos;re check-in apps that depend entirely on self-reporting honesty — defeating the point of accountability.
              </p>
              <p className="text-[#8A9BB0] leading-relaxed mb-4">
                We wanted something that creates automatic, consistent signals — so that disclosure isn&apos;t a choice that a struggling person has to summon the courage to make in their worst moment — while simultaneously preserving dignity and privacy.
              </p>
              <p className="text-[#8A9BB0] leading-relaxed">
                The on-device AI model was the key insight. Your screen content never leaves your device. Partners see metadata, not surveillance. That changes everything about what&apos;s possible in this space.
              </p>
            </div>
            <div>
              <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Mission & Vision</h2>
              <div className="space-y-5">
                <div className="p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
                  <h3 className="font-semibold text-[#C9A84C] mb-2 text-sm uppercase tracking-wide">Mission</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">
                    To provide every person committed to sexual purity with a free, dignified, and effective accountability tool — regardless of their economic situation.
                  </p>
                </div>
                <div className="p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
                  <h3 className="font-semibold text-[#C9A84C] mb-2 text-sm uppercase tracking-wide">Vision</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">
                    A world where no man or woman faces the struggle alone — where accountability is normalized, technology-assisted, and built on genuine community rather than shame.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Values */}
      <section className="py-20 bg-[#0A1420] border-y border-[#1E3050]">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">What We Stand For</h2>
          </div>
          <div className="grid sm:grid-cols-2 gap-5">
            {[
              {
                icon: '🔒',
                title: 'Privacy-First',
                desc: 'We designed the privacy model before the first line of code was written. On-device AI is not a feature — it is the architecture. Your content stays on your device. This is non-negotiable.',
              },
              {
                icon: '❤️',
                title: 'Radically Free',
                desc: 'No one should have to choose between their financial situation and access to accountability. We will never charge for this app. Donations sustain the project; they do not gate it.',
              },
              {
                icon: '🤝',
                title: 'Peer-Centered',
                desc: 'Accountability must be relational. RF is a tool for people who already have — or are building — relationships of trust. The app facilitates; the people do the real work.',
              },
              {
                icon: '✝️',
                title: 'Faith-Grounded',
                desc: 'RF is built from a Christian perspective and reflects a biblical view of sexuality, covenant, and community. That said, the app itself is open to anyone who finds its accountability model valuable.',
              },
            ].map((v) => (
              <div key={v.title} className="flex gap-5 p-6 rounded-2xl border border-[#1E3050] bg-[#162235]">
                <div className="text-3xl flex-shrink-0">{v.icon}</div>
                <div>
                  <h3 className="font-serif text-lg font-semibold text-[#F0EDE8] mb-2">{v.title}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{v.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Open Source */}
      <section className="py-20">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Open Source Commitment</h2>
              <p className="text-[#8A9BB0] leading-relaxed mb-4">
                Remain Faithful is open source. Every line of code — iOS app, Go backend, and this website — is publicly available on GitHub.
              </p>
              <p className="text-[#8A9BB0] leading-relaxed mb-4">
                For an app that handles intimate personal behavioral data, open source isn&apos;t optional. You should be able to verify our privacy claims by reading the code. Anyone who tells you to &ldquo;just trust us&rdquo; with this kind of data is asking too much.
              </p>
              <p className="text-[#8A9BB0] leading-relaxed mb-6">
                Security researchers, privacy advocates, and curious developers are all welcome to inspect, fork, and contribute.
              </p>
              <a
                href="https://github.com/remainfaithful"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 px-6 py-3 rounded-xl border border-[#1E3050] bg-[#162235] text-[#F0EDE8] text-sm font-semibold hover:border-[#C9A84C]/50 transition-colors"
              >
                <svg width="18" height="18" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/>
                </svg>
                View on GitHub
              </a>
            </div>
            <div
              className="rounded-2xl p-6 border border-[#1E3050] font-mono text-xs"
              style={{ background: '#0A1015' }}
            >
              <div className="flex items-center gap-1.5 mb-4">
                <div className="w-2.5 h-2.5 rounded-full bg-red-500"/>
                <div className="w-2.5 h-2.5 rounded-full bg-yellow-500"/>
                <div className="w-2.5 h-2.5 rounded-full bg-green-500"/>
                <span className="ml-2 text-[#8A9BB0]">SampleHandler.swift</span>
              </div>
              <pre className="text-[#8A9BB0] overflow-x-auto leading-relaxed">
                <span className="text-[#6BBF6B]">// On-device only.</span>{'\n'}
                <span className="text-[#6BBF6B]">// This data never leaves the device.</span>{'\n'}
                {'\n'}
                <span className="text-[#E8C87A]">private func</span>{' '}
                <span className="text-[#F0EDE8]">classifyFrame</span>
                {'(_ frame: CMSampleBuffer) {'}{'\n'}
                {'  '}<span className="text-[#E8C87A]">let</span> tier1 = URLBlocklist
                {'\n'}
                {'    '}.check(frame){'\n'}
                {'  '}<span className="text-[#E8C87A]">guard</span> !tier1.flagged{' '}
                <span className="text-[#E8C87A]">else</span> {'{'}
                {'\n'}
                {'    '}sendAlert(tier1.result){'\n'}
                {'    '}<span className="text-[#E8C87A]">return</span>{'\n'}
                {'  }'}{'\n'}
                {'\n'}
                {'  '}<span className="text-[#E8C87A]">let</span> tier2 = OnDeviceClassifier
                {'\n'}
                {'    '}.classify(frame){'\n'}
                {'  '}<span className="text-[#E8C87A]">if</span> tier2.confidence {'>'} 0.85 {'{'}
                {'\n'}
                {'    '}sendAlert(tier2.result){'\n'}
                {'  }'}{'\n'}
                {'}'}
              </pre>
            </div>
          </div>
        </div>
      </section>

      {/* Donation Model */}
      <section className="py-20 bg-[#0A1420] border-y border-[#1E3050]">
        <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-5">Why Donations, Not Subscriptions</h2>
          <p className="text-[#8A9BB0] leading-relaxed mb-8">
            A subscription model creates a perverse incentive: we benefit when you stay subscribed, not necessarily when you succeed. Donations reverse this. We depend on users who find real value in the app — users who are succeeding.
          </p>
          <div className="grid sm:grid-cols-3 gap-4 text-left mb-10">
            {[
              { label: 'Server costs', amount: '~$80/mo', desc: 'Go backend, PostgreSQL, APNs relay' },
              { label: 'Stripe fees', amount: '2.9% + 30¢', desc: 'Per donation transaction' },
              { label: 'Development', amount: 'Volunteer', desc: 'Core team donates time' },
            ].map((c) => (
              <div key={c.label} className="p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
                <p className="text-xs text-[#8A9BB0] uppercase tracking-wide mb-1">{c.label}</p>
                <p className="font-serif font-bold text-[#C9A84C] text-xl mb-1">{c.amount}</p>
                <p className="text-xs text-[#8A9BB0]">{c.desc}</p>
              </div>
            ))}
          </div>
          <Link
            href="/#donate"
            className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all duration-200"
          >
            Support the Project
          </Link>
        </div>
      </section>

      {/* Contact */}
      <section id="contact" className="py-20">
        <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">Get in Touch</h2>
            <p className="text-[#8A9BB0]">
              Questions, feedback, bug reports, or partnership inquiries — we read everything.
            </p>
          </div>
          <ContactForm />
          <div className="mt-10 flex flex-col sm:flex-row gap-4 text-center sm:text-left">
            <div className="flex-1 p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
              <p className="text-xs text-[#8A9BB0] uppercase tracking-wide mb-1">Email</p>
              <a href="mailto:support@remainfaithful.com" className="text-[#C9A84C] text-sm hover:underline">
                support@remainfaithful.com
              </a>
            </div>
            <div className="flex-1 p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
              <p className="text-xs text-[#8A9BB0] uppercase tracking-wide mb-1">Bug Reports</p>
              <a
                href="https://github.com/remainfaithful"
                target="_blank"
                rel="noopener noreferrer"
                className="text-[#C9A84C] text-sm hover:underline"
              >
                GitHub Issues
              </a>
            </div>
          </div>
        </div>
      </section>
    </>
  )
}

function ContactForm() {
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
      if (!res.ok) throw new Error('failed')
      setSubmitted(true)
    } catch {
      setError('Something went wrong. Please try again or email us directly.')
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
          <label className="block text-sm font-medium text-[#F0EDE8] mb-2">Name <span className="text-[#C9A84C]">*</span></label>
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
          <label className="block text-sm font-medium text-[#F0EDE8] mb-2">Email <span className="text-[#C9A84C]">*</span></label>
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
        <label className="block text-sm font-medium text-[#F0EDE8] mb-2">Subject <span className="text-[#C9A84C]">*</span></label>
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
        <label className="block text-sm font-medium text-[#F0EDE8] mb-2">Message <span className="text-[#C9A84C]">*</span></label>
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
        className="w-full py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all duration-200 disabled:opacity-70 disabled:cursor-not-allowed"
      >
        {loading ? 'Sending...' : 'Send Message'}
      </button>
    </form>
  )
}
