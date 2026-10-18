'use client'

import { useState } from 'react'
import { trackEvent } from '@/lib/track'

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
      trackEvent('android_signup')
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
          We&apos;ll write when there is Android news.
        </p>
      </div>
    )
  }

  if (variant === 'footer') {
    return (
      <form onSubmit={handleSubmit} className="flex flex-wrap gap-2">
        <label htmlFor="waitlist-footer-email" className="sr-only">
          Email
        </label>
        <input
          id="waitlist-footer-email"
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
          {loading ? '...' : 'Join'}
        </button>
        {error && (
          <p className="form-error text-xs mt-1 w-full" role="alert">
            {error}
          </p>
        )}
      </form>
    )
  }

  if (variant === 'inline') {
    return (
      <form onSubmit={handleSubmit} className="flex flex-col sm:flex-row gap-3">
        <label htmlFor="waitlist-inline-email" className="sr-only">
          Email
        </label>
        <input
          id="waitlist-inline-email"
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
          {loading ? 'Joining...' : buttonText}
        </button>
        {error && (
          <p className="form-error text-xs mt-2" role="alert">
            {error}
          </p>
        )}
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
          <div className="flex-1">
            <label htmlFor="waitlist-name" className="mb-2 block text-sm text-white">
              Name
            </label>
            <input
              id="waitlist-name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Your name"
              required
              className="input-field w-full"
            />
          </div>
          <div className="flex-1">
            <label htmlFor="waitlist-email" className="mb-2 block text-sm text-white">
              Email
            </label>
            <input
              id="waitlist-email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="your@email.com"
              required
              className="input-field w-full"
            />
          </div>
        </div>
        <button
          type="submit"
          disabled={loading}
          className="btn-wax w-full"
        >
          {loading ? 'Joining...' : buttonText}
        </button>
        {error && (
          <p className="form-error text-sm" role="alert">
            {error}
          </p>
        )}
      </form>
    </div>
  )
}
