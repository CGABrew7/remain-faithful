'use client'

import { useState } from 'react'

interface WaitlistFormProps {
  variant?: 'default' | 'inline' | 'footer'
  buttonText?: string
  heading?: string
  subheading?: string
}

export default function WaitlistForm({
  variant = 'default',
  buttonText = 'Notify Me',
  heading,
  subheading,
}: WaitlistFormProps) {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [submitted, setSubmitted] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      const res = await fetch('/api/waitlist', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, email }),
      })
      if (!res.ok) throw new Error('Failed')
      setSubmitted(true)
    } catch {
      setError('Something went wrong. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  if (submitted) {
    return (
      <div className="py-4">
        <p className="text-wax font-serif text-lg">You&apos;re on the list.</p>
        <p className="text-ink-soft mt-1">
          We&apos;ll write when Remain Faithful is ready.
        </p>
      </div>
    )
  }

  if (variant === 'footer') {
    return (
      <form onSubmit={handleSubmit} className="flex gap-2">
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="your@email.com"
          required
          className="input-field flex-1 min-w-0 !py-2.5"
        />
        <button
          type="submit"
          disabled={loading}
          className="btn-wax !py-2.5 whitespace-nowrap"
        >
          {loading ? '…' : 'Join'}
        </button>
        {error && <p className="text-xs text-wax mt-1">{error}</p>}
      </form>
    )
  }

  if (variant === 'inline') {
    return (
      <form onSubmit={handleSubmit} className="flex flex-col sm:flex-row gap-3">
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="your@email.com"
          required
          className="input-field flex-1"
        />
        <button
          type="submit"
          disabled={loading}
          className="btn-wax whitespace-nowrap"
        >
          {loading ? 'Joining…' : buttonText}
        </button>
        {error && <p className="text-xs text-wax mt-2">{error}</p>}
      </form>
    )
  }

  return (
    <div>
      {heading && (
        <h3 className="font-serif text-xl font-medium text-ink mb-2">{heading}</h3>
      )}
      {subheading && (
        <p className="text-ink-soft mb-5">{subheading}</p>
      )}
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <div className="flex flex-col sm:flex-row gap-3">
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Your name"
            required
            className="input-field flex-1"
          />
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="your@email.com"
            required
            className="input-field flex-1"
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="btn-wax w-full"
        >
          {loading ? 'Joining…' : buttonText}
        </button>
        {error && <p className="text-xs text-wax">{error}</p>}
      </form>
    </div>
  )
}
