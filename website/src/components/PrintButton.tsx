'use client'

export default function PrintButton() {
  return (
    <button
      onClick={() => window.print()}
      className="flex items-center gap-2 px-5 py-3 rounded-xl bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] text-[#0F1B2D] font-semibold text-sm shadow-lg hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all"
    >
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
        <polyline points="6 9 6 2 18 2 18 9"/>
        <path d="M6 18H4a2 2 0 01-2-2v-5a2 2 0 012-2h16a2 2 0 012 2v5a2 2 0 01-2 2h-2"/>
        <rect x="6" y="14" width="12" height="8"/>
      </svg>
      Print / Save as PDF
    </button>
  )
}
