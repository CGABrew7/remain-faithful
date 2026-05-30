'use client'

import { useState } from 'react'

export default function FaqAccordion({ faqs }: { faqs: { q: string; a: string }[] }) {
  const [open, setOpen] = useState<number | null>(0)
  return (
    <div className="space-y-3">
      {faqs.map((faq, i) => (
        <div key={i} className="rounded-xl border border-[#1E3050] bg-[#162235] overflow-hidden">
          <button
            className="w-full flex items-center justify-between px-6 py-5 text-left"
            onClick={() => setOpen(open === i ? null : i)}
            aria-expanded={open === i}
          >
            <span className="font-semibold text-[#F0EDE8] text-sm sm:text-base pr-4">{faq.q}</span>
            <svg
              width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2.5" strokeLinecap="round"
              className={`flex-shrink-0 transition-transform duration-200 ${open === i ? 'rotate-45' : ''}`}
            >
              <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
          </button>
          <div className={`overflow-hidden transition-all duration-300 ${open === i ? 'max-h-96' : 'max-h-0'}`}>
            <p className="px-6 pb-5 text-sm text-[#8A9BB0] leading-relaxed">{faq.a}</p>
          </div>
        </div>
      ))}
    </div>
  )
}
