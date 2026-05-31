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
      <div className="text-center py-4">
        <p className="text-[#C9A84C] font-semibold">You&apos;re on the list!</p>
        <p className="text-sm text-[#8A9BB0] mt-1">
          We&apos;ll notify you the moment Remain Faithful launches.
        </p>
      </div>
    )
  }

  if (variant === 'footer') {
    return (
      <form onSubmit={handleSubmit} className="flex gap-2 mt-3">
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="your@email.com"
          required
          className="flex-1 px-3 py-2 rounded-lg bg-[#162235] border border-[#1E3050] text-sm text-[#F0EDE8] placeholder-[#8A9BB0]/60 focus:outline-none focus:border-[#C9A84C]/50 min-w-0"
        />
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 rounded-lg text-sm font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] disabled:opacity-70 whitespace-nowrap"
        >
          {loading ? '…' : 'Join'}
        </button>
        {error && <p className="text-xs text-red-400 mt-1">{error}</p>}
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
          className="flex-1 px-4 py-3 rounded-xl bg-[#0F1B2D] border border-[#1E3050] text-[#F0EDE8] placeholder-[#8A9BB0]/60 focus:outline-none focus:border-[#C9A84C]/50"
        />
        <button
          type="submit"
          disabled={loading}
          className="px-6 py-3 rounded-xl font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all duration-200 disabled:opacity-70 whitespace-nowrap"
        >
          {loading ? 'Joining…' : buttonText}
        </button>
        {error && <p className="text-xs text-red-400 mt-2">{error}</p>}
      </form>
    )
  }

  return (
    <div>
      {heading && (
        <h3 className="font-serif text-xl font-bold text-[#F0EDE8] mb-2">{heading}</h3>
      )}
      {subheading && (
        <p className="text-[#8A9BB0] text-sm mb-5">{subheading}</p>
      )}
      <form onSubmit={handleSubmit} className="flex flex-col gap-4 max-w-md mx-auto">
        <div className="flex flex-col sm:flex-row gap-3">
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Your name"
            required
            className="flex-1 px-4 py-3 rounded-xl bg-[#162235] border border-[#1E3050] text-[#F0EDE8] placeholder-[#8A9BB0]/60 focus:outline-none focus:border-[#C9A84C]/50"
          />
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="your@email.com"
            required
            className="flex-1 px-4 py-3 rounded-xl bg-[#162235] border border-[#1E3050] text-[#F0EDE8] placeholder-[#8A9BB0]/60 focus:outline-none focus:border-[#C9A84C]/50"
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="w-full py-3.5 rounded-xl font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all duration-200 shadow-lg shadow-[#C9A84C]/20 disabled:opacity-70"
        >
          {loading ? 'Joining…' : buttonText}
        </button>
        {error && <p className="text-xs text-red-400">{error}</p>}
      </form>
    </div>
  )
}
