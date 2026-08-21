'use client'

import { useState } from 'react'

export default function FaqAccordion({ faqs }: { faqs: { q: string; a: string }[] }) {
  const [open, setOpen] = useState<number | null>(0)
  return (
    <div>
      {faqs.map((faq, i) => (
        <div key={i} className="border-b border-hairline">
          <button
            className="w-full flex items-center justify-between py-5 text-left min-h-10"
            onClick={() => setOpen(open === i ? null : i)}
            aria-expanded={open === i}
          >
            <span className="font-serif text-lg text-ink pr-4">{faq.q}</span>
            <svg
              width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round"
              className={`shrink-0 text-wax transition-transform duration-200 ${open === i ? 'rotate-45' : ''}`}
            >
              <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
          </button>
          <div className={`overflow-hidden transition-[max-height] duration-300 ${open === i ? 'max-h-96' : 'max-h-0'}`}>
            <p className="pb-5 text-ink-soft leading-relaxed">{faq.a}</p>
          </div>
        </div>
      ))}
    </div>
  )
}
