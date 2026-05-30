export default function AppMockup() {
  return (
    <div className="relative mx-auto" style={{ width: 260, height: 540 }}>
      {/* Glow behind phone */}
      <div
        className="absolute inset-0 rounded-[48px] blur-3xl opacity-30"
        style={{ background: 'radial-gradient(ellipse at center, #C9A84C 0%, transparent 70%)' }}
      />

      {/* Phone frame */}
      <div
        className="relative rounded-[44px] border-2 border-[#2A3F5F] overflow-hidden shadow-2xl"
        style={{
          width: 260,
          height: 540,
          background: '#0A1628',
          boxShadow: '0 0 0 1px #1E3050, 0 40px 80px rgba(0,0,0,0.6)',
        }}
      >
        {/* Dynamic Island */}
        <div
          className="absolute top-3 left-1/2 -translate-x-1/2 rounded-full bg-black z-10"
          style={{ width: 88, height: 24 }}
        />

        {/* Status bar */}
        <div className="flex justify-between items-center px-6 pt-10 pb-1">
          <span className="text-[10px] text-[#8A9BB0] font-medium">9:41</span>
          <div className="flex items-center gap-1">
            <SignalIcon />
            <WifiIcon />
            <BatteryIcon />
          </div>
        </div>

        {/* Screen content */}
        <div className="px-4 pt-2 pb-4 space-y-3">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div>
              <p className="text-[10px] text-[#8A9BB0]">Good morning,</p>
              <p className="text-[15px] font-serif font-semibold text-[#F0EDE8]">Jeff</p>
            </div>
            <div className="w-8 h-8 rounded-full bg-gradient-to-br from-[#C9A84C] to-[#E8C87A] flex items-center justify-center">
              <svg width="14" height="16" viewBox="0 0 32 36" fill="none">
                <path d="M16 0L2 6V18C2 26.284 8.268 33.916 16 36C23.732 33.916 30 26.284 30 18V6L16 0Z" fill="white" fillOpacity="0.9"/>
                <path d="M11 18L14.5 21.5L21 14" stroke="#C9A84C" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
            </div>
          </div>

          {/* Streak card */}
          <div
            className="rounded-2xl p-3"
            style={{ background: 'linear-gradient(135deg, #1A2F1A, #1E3A1E)', border: '1px solid #2A5A2A' }}
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-[9px] text-[#6BBF6B] uppercase tracking-wider font-semibold">Current Streak</p>
                <div className="flex items-baseline gap-1 mt-0.5">
                  <span className="text-[28px] font-serif font-bold text-[#F0EDE8]">14</span>
                  <span className="text-[11px] text-[#8A9BB0]">days</span>
                </div>
              </div>
              <div className="text-right">
                <div className="text-2xl">🔥</div>
                <p className="text-[9px] text-[#6BBF6B] mt-0.5">Best: 21</p>
              </div>
            </div>
          </div>

          {/* Status card */}
          <div
            className="rounded-2xl p-3 flex items-center gap-2.5"
            style={{ background: '#162235', border: '1px solid #1E3050' }}
          >
            <div className="w-7 h-7 rounded-full bg-[#1A3A1A] border border-[#2A5A2A] flex items-center justify-center flex-shrink-0">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#6BBF6B" strokeWidth="2.5" strokeLinecap="round">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
            </div>
            <div>
              <p className="text-[11px] font-medium text-[#F0EDE8]">No alerts today</p>
              <p className="text-[9px] text-[#8A9BB0]">Monitoring active</p>
            </div>
            <div className="ml-auto flex items-center gap-1">
              <div className="w-1.5 h-1.5 rounded-full bg-[#6BBF6B] animate-pulse"/>
            </div>
          </div>

          {/* Partners section */}
          <div>
            <p className="text-[9px] font-semibold uppercase tracking-wider text-[#8A9BB0] mb-2">Partners</p>
            <div className="space-y-2">
              <PartnerCard name="Mike R." streak="8 days" status="active" initial="M" color="#2A4A8A" />
              <PartnerCard name="David C." streak="22 days" status="active" initial="D" color="#4A2A6A" />
            </div>
          </div>

          {/* Weekly digest */}
          <div
            className="rounded-xl p-2.5 flex items-center justify-between"
            style={{ background: '#1A2535', border: '1px solid #1E3050' }}
          >
            <div className="flex items-center gap-2">
              <span className="text-base">📊</span>
              <div>
                <p className="text-[10px] font-medium text-[#F0EDE8]">Weekly Digest</p>
                <p className="text-[9px] text-[#8A9BB0]">Sent to partners</p>
              </div>
            </div>
            <div className="w-8 h-4 rounded-full bg-[#C9A84C] flex items-center justify-end pr-0.5">
              <div className="w-3 h-3 rounded-full bg-white"/>
            </div>
          </div>
        </div>

        {/* Bottom nav */}
        <div
          className="absolute bottom-0 left-0 right-0 flex items-center justify-around px-4 py-3 border-t border-[#1E3050]"
          style={{ background: '#0A1628' }}
        >
          <NavItem icon="home" active />
          <NavItem icon="group" />
          <NavItem icon="settings" />
        </div>
      </div>
    </div>
  )
}

function PartnerCard({
  name, streak, initial, color
}: {
  name: string; streak: string; status: string; initial: string; color: string
}) {
  return (
    <div
      className="rounded-xl px-3 py-2 flex items-center gap-2.5"
      style={{ background: '#162235', border: '1px solid #1E3050' }}
    >
      <div
        className="w-6 h-6 rounded-full flex items-center justify-center flex-shrink-0 text-[9px] font-bold text-white"
        style={{ background: color }}
      >
        {initial}
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-[11px] font-medium text-[#F0EDE8] truncate">{name}</p>
        <p className="text-[9px] text-[#8A9BB0]">{streak} clean</p>
      </div>
      <div className="w-1.5 h-1.5 rounded-full bg-[#6BBF6B]"/>
    </div>
  )
}

function NavItem({ icon, active }: { icon: string; active?: boolean }) {
  const color = active ? '#C9A84C' : '#8A9BB0'
  if (icon === 'home') return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill={active ? color : 'none'} stroke={color} strokeWidth="2" strokeLinecap="round">
      <path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z"/>
      <polyline points="9 22 9 12 15 12 15 22"/>
    </svg>
  )
  if (icon === 'group') return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke={color} strokeWidth="2" strokeLinecap="round">
      <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
      <circle cx="9" cy="7" r="4"/>
      <path d="M23 21v-2a4 4 0 00-3-3.87M16 3.13a4 4 0 010 7.75"/>
    </svg>
  )
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke={color} strokeWidth="2" strokeLinecap="round">
      <circle cx="12" cy="12" r="3"/>
      <path d="M19.07 4.93l-1.41 1.41M5.34 18.66l-1.41 1.41M2 12h2m16 0h2M4.93 4.93l1.41 1.41M18.66 18.66l1.41 1.41"/>
    </svg>
  )
}

function SignalIcon() {
  return (
    <svg width="12" height="10" viewBox="0 0 12 10" fill="#8A9BB0">
      <rect x="0" y="6" width="2" height="4" rx="0.5"/>
      <rect x="3.5" y="4" width="2" height="6" rx="0.5"/>
      <rect x="7" y="2" width="2" height="8" rx="0.5"/>
      <rect x="10.5" y="0" width="1.5" height="10" rx="0.5"/>
    </svg>
  )
}

function WifiIcon() {
  return (
    <svg width="12" height="10" viewBox="0 0 24 24" fill="none" stroke="#8A9BB0" strokeWidth="2.5" strokeLinecap="round">
      <path d="M5 12.55a11 11 0 0114.08 0M1.42 9a16 16 0 0121.16 0M8.53 16.11a6 6 0 016.95 0M12 20h.01"/>
    </svg>
  )
}

function BatteryIcon() {
  return (
    <svg width="18" height="10" viewBox="0 0 18 10" fill="none">
      <rect x="0.5" y="0.5" width="14" height="9" rx="2" stroke="#8A9BB0"/>
      <rect x="2" y="2" width="9" height="6" rx="1" fill="#8A9BB0"/>
      <path d="M16 3.5v3" stroke="#8A9BB0" strokeWidth="1.5" strokeLinecap="round"/>
    </svg>
  )
}
