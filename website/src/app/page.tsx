import type { Metadata } from 'next'
import { Suspense } from 'react'
import Link from 'next/link'
import AppMockup from '@/components/AppMockup'
import DonateButton from '@/components/DonateButton'
import DonationSuccessBanner from '@/components/DonationSuccessBanner'

export const metadata: Metadata = {
  title: 'Remain Faithful: Accountability That Works',
  description:
    'Free peer accountability for adults committed to purity. On-device AI, privacy-first, built on trust, not surveillance.',
}

export default function HomePage() {
  return (
    <>
      <Suspense fallback={null}>
        <DonationSuccessBanner />
      </Suspense>

      {/* ── Hero ── */}
      <section className="relative min-h-screen flex items-center overflow-hidden">
        {/* Background radial */}
        <div
          className="absolute inset-0 pointer-events-none"
          aria-hidden="true"
          style={{
            background:
              'radial-gradient(ellipse 80% 60% at 70% 50%, rgba(201,168,76,0.08) 0%, transparent 70%)',
          }}
        />
        <div
          className="absolute top-0 left-0 right-0 h-px"
          style={{ background: 'linear-gradient(90deg, transparent, #1E3050, transparent)' }}
        />

        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 pt-32 w-full">
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            {/* Left: Copy */}
            <div className="max-w-xl">
              <h1 className="font-serif text-5xl sm:text-6xl font-bold leading-[1.1] text-[#F0EDE8] mb-6">
                Accountability<br />
                <span
                  style={{
                    background: 'linear-gradient(135deg, #C9A84C, #E8C87A)',
                    WebkitBackgroundClip: 'text',
                    WebkitTextFillColor: 'transparent',
                    backgroundClip: 'text',
                  }}
                >
                  That Works
                </span>
              </h1>

              <p className="text-lg text-[#8A9BB0] leading-relaxed mb-10">
                Free peer accountability for adults committed to purity.
                Built on trust, not surveillance. Your content never leaves your device.
              </p>

              <div className="flex flex-wrap gap-4">
                <a
                  href="https://apps.apple.com"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all duration-200 shadow-lg shadow-[#C9A84C]/20"
                >
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.8-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
                  </svg>
                  Download for iPhone
                </a>
                <Link
                  href="/how-it-works"
                  className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#F0EDE8] border border-[#1E3050] hover:border-[#C9A84C]/50 hover:bg-[#162235] transition-all duration-200"
                >
                  Learn How It Works
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                    <path d="M5 12h14M12 5l7 7-7 7"/>
                  </svg>
                </Link>
              </div>

              {/* Trust indicators */}
              <div className="grid grid-cols-2 gap-x-6 gap-y-3 mt-10 pt-8 border-t border-[#1E3050]">
                <TrustItem>100% Free, Forever</TrustItem>
                <TrustItem>Donation Funded</TrustItem>
                <TrustItem>On-Device AI &amp; Security</TrustItem>
                <TrustItem>Zero Data Shared Off Your Device</TrustItem>
              </div>
            </div>

            {/* Right: Phone mockup */}
            <div className="flex justify-center lg:justify-end">
              <AppMockup />
            </div>
          </div>
        </div>
      </section>

      {/* ── Feature Cards ── */}
      <section className="py-24 border-t border-[#1E3050]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl sm:text-4xl font-bold text-[#F0EDE8] mb-4">
              Built Different
            </h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
              Most accountability tools rely on shame or surveillance. Remain Faithful is built on covenant, trust, and genuine community.
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-6">
            <FeatureCard
              icon={
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2" strokeLinecap="round">
                  <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
                  <circle cx="9" cy="7" r="4"/>
                  <path d="M23 21v-2a4 4 0 00-3-3.87M16 3.13a4 4 0 010 7.75"/>
                </svg>
              }
              title="One-to-One or Group"
              body="Choose a single trusted partner or set up a small group. RF works for close friendships, mentorship relationships, and men's ministry groups alike."
            />
            <FeatureCard
              icon={
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2" strokeLinecap="round">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0110 0v4"/>
                </svg>
              }
              title="On-Device Privacy"
              body="All AI classification runs locally using Apple's Vision and SensitiveContentAnalysis frameworks. Your screen content is never transmitted. Partners see metadata, not your screen."
            />
            <FeatureCard
              icon={
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2" strokeLinecap="round">
                  <path d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z"/>
                </svg>
              }
              title="Always Free"
              body="Remain Faithful is free today and will remain free forever. No subscription tiers, no paywalls, no premium features. Sustained entirely by voluntary donations."
            />
          </div>
        </div>
      </section>

      {/* ── How It Works ── */}
      <section className="py-24 bg-[#0A1420]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">Simple Process</p>
            <h2 className="font-serif text-3xl sm:text-4xl font-bold text-[#F0EDE8] mb-4">
              Three Steps to Real Accountability
            </h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
              No complicated setup. No long onboarding. Start holding each other accountable in minutes.
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8 relative">
            {/* Connector line (desktop) */}
            <div className="hidden md:block absolute top-10 left-[calc(16.667%+20px)] right-[calc(16.667%+20px)] h-px bg-gradient-to-r from-[#C9A84C]/40 via-[#C9A84C]/20 to-[#C9A84C]/40" />

            <Step
              number="01"
              icon={
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                  <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
                  <circle cx="9" cy="7" r="4"/>
                  <line x1="19" y1="8" x2="19" y2="14"/>
                  <line x1="22" y1="11" x2="16" y2="11"/>
                </svg>
              }
              title="Choose Your Partners"
              body="Invite one or more trusted friends, mentors, or group members to be your accountability partners. They accept a covenant before gaining any access."
            />
            <Step
              number="02"
              icon={
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                  <circle cx="12" cy="12" r="10"/>
                  <path d="M12 8v4M12 16h.01"/>
                </svg>
              }
              title="Enable Monitoring"
              body="Grant screen broadcast permission. The app runs quietly in the background, using on-device AI to classify content according to your agreed standards."
            />
            <Step
              number="03"
              icon={
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                  <path d="M22 16.92v3a2 2 0 01-2.18 2 19.79 19.79 0 01-8.63-3.07A19.5 19.5 0 013.86 9.5a19.79 19.79 0 01-3.07-8.67A2 2 0 012.81 0h3a2 2 0 012 1.72c.127.96.361 1.903.7 2.81a2 2 0 01-.45 2.11L6.91 7.91a16 16 0 006.29 6.29l1.28-1.28a2 2 0 012.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0122 16.92z"/>
                </svg>
              }
              title="Stay Accountable"
              body="When something is flagged, you and your partners receive a discreet alert. No surprises, no shame spirals. Just honest accountability."
            />
          </div>

          <div className="text-center mt-12">
            <Link
              href="/how-it-works"
              className="inline-flex items-center gap-2 text-[#C9A84C] font-semibold hover:underline underline-offset-4"
            >
              Read the full breakdown
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                <path d="M5 12h14M12 5l7 7-7 7"/>
              </svg>
            </Link>
          </div>
        </div>
      </section>

      {/* ── Privacy Callout ── */}
      <section className="py-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div
            className="rounded-3xl p-10 sm:p-14 text-center"
            style={{
              background: 'linear-gradient(135deg, #162235 0%, #1A2A40 100%)',
              border: '1px solid rgba(201,168,76,0.25)',
              boxShadow: '0 0 80px rgba(201,168,76,0.04)',
            }}
          >
            <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-[#C9A84C]/10 border border-[#C9A84C]/20 mb-6">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2" strokeLinecap="round">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
              </svg>
            </div>
            <h2 className="font-serif text-3xl sm:text-4xl font-bold text-[#F0EDE8] mb-4">
              Your content never leaves your device
            </h2>
            <p className="text-[#8A9BB0] text-lg mb-10 max-w-xl mx-auto">
              We designed the privacy model first, then built the app around it. Zero compromise.
            </p>
            <div className="grid sm:grid-cols-3 gap-6 text-left">
              <PrivacyPoint
                title="On-Device AI"
                body="Apple Vision OCR and SensitiveContentAnalysis run entirely on your hardware. No server-side processing of your screen."
              />
              <PrivacyPoint
                title="Alert Metadata Only"
                body="Partners receive: timestamp, app category, and severity level. Never a screenshot. Never raw content. Never your browsing history."
              />
              <PrivacyPoint
                title="You Control Access"
                body="You approve every partner. You set alert thresholds. You can pause monitoring or remove partners at any time, instantly."
              />
            </div>
          </div>
        </div>
      </section>

      {/* ── Testimonials ── */}
      <section className="py-24 border-t border-[#1E3050]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">Beta Users</p>
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8]">
              What People Are Saying
            </h2>
          </div>
          <div className="grid md:grid-cols-3 gap-6">
            <Testimonial
              quote="RF changed our men's group. We went from vague commitments to real accountability. The alerts make the invisible visible, without the shame spiral."
              author="Beta User"
              role="Kansas City, Missouri"
            />
            <Testimonial
              quote="I finally have partners who actually know the struggle. The app makes it real, not just talk. And knowing they'll get an alert... that changes the moment."
              author="Beta User"
              role="Austin, Texas"
            />
            <Testimonial
              quote="Set it up in five minutes. My accountability partner got the hang of it immediately. Simple, honest, and it actually works."
              author="Beta User"
              role="Columbus, Ohio"
            />
          </div>
        </div>
      </section>

      {/* ── Download ── */}
      <section id="download" className="py-24 bg-[#0A1420] border-t border-[#1E3050]">
        <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">Get Started</p>
          <h2 className="font-serif text-3xl sm:text-4xl font-bold text-[#F0EDE8] mb-4">
            Start Your Accountability Journey
          </h2>
          <p className="text-[#8A9BB0] mb-12">
            Free forever. No subscription. No ads. Just accountability.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            {/* App Store */}
            <a
              href="https://apps.apple.com"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-3 px-6 py-3.5 rounded-2xl bg-black border border-[#2A3F5F] hover:border-[#C9A84C]/50 transition-colors group"
            >
              <svg width="28" height="28" viewBox="0 0 24 24" fill="#F0EDE8" className="group-hover:scale-105 transition-transform">
                <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.8-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
              </svg>
              <div className="text-left">
                <p className="text-[10px] text-[#8A9BB0] leading-none">Download on the</p>
                <p className="text-[15px] font-semibold text-[#F0EDE8] leading-tight">App Store</p>
              </div>
            </a>

            {/* Android Coming Soon */}
            <div className="relative flex items-center gap-3 px-6 py-3.5 rounded-2xl bg-[#162235] border border-[#1E3050] opacity-60 cursor-not-allowed">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="#8A9BB0">
                <path d="M6 18c0 .55.45 1 1 1h1v3.5c0 .83.67 1.5 1.5 1.5s1.5-.67 1.5-1.5V19h2v3.5c0 .83.67 1.5 1.5 1.5s1.5-.67 1.5-1.5V19h1c.55 0 1-.45 1-1V8H6v10zM3.5 8C2.67 8 2 8.67 2 9.5v7c0 .83.67 1.5 1.5 1.5S5 17.33 5 16.5v-7C5 8.67 4.33 8 3.5 8zm17 0c-.83 0-1.5.67-1.5 1.5v7c0 .83.67 1.5 1.5 1.5s1.5-.67 1.5-1.5v-7c0-.83-.67-1.5-1.5-1.5zm-4.97-5.84l1.3-1.3c.2-.2.2-.51 0-.71-.2-.2-.51-.2-.71 0l-1.48 1.48C14.15 1.23 13.1 1 12 1c-1.1 0-2.15.23-3.12.63L7.4.15c-.2-.2-.51-.2-.71 0-.2.2-.2.51 0 .71l1.31 1.31C6.1 3.26 5 5.01 5 7h14c0-1.99-1.1-3.74-2.47-4.84zM10 5H9V4h1v1zm5 0h-1V4h1v1z"/>
              </svg>
              <div className="text-left">
                <p className="text-[10px] text-[#8A9BB0] leading-none">Coming Soon to</p>
                <p className="text-[15px] font-semibold text-[#8A9BB0] leading-tight">Google Play</p>
              </div>
              <div className="absolute -top-2 -right-2 px-2 py-0.5 rounded-full bg-[#C9A84C] text-[#0F1B2D] text-[9px] font-bold uppercase tracking-wide">
                Soon
              </div>
            </div>
          </div>

          <p className="text-xs text-[#8A9BB0]/60 mt-6">
            Requires iOS 17.0 or later. Free. No in-app purchases.
          </p>
        </div>
      </section>

      {/* ── Donation CTA ── */}
      <section id="donate" className="py-24" style={{ background: 'linear-gradient(180deg, #0F1B2D 0%, #111D2E 100%)' }}>
        <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <div className="inline-flex items-center justify-center w-12 h-12 rounded-2xl bg-[#C9A84C]/10 border border-[#C9A84C]/20 mb-6">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="#C9A84C">
              <path d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z"/>
            </svg>
          </div>
          <h2 className="font-serif text-3xl sm:text-4xl font-bold text-[#F0EDE8] mb-4">
            Keep Remain Faithful Free
          </h2>
          <p className="text-[#8A9BB0] text-lg mb-10">
            We&apos;re committed to never charging for accountability. Your donation funds server costs, development, and ministry outreach.
          </p>
          <DonateButton />
          <p className="text-xs text-[#8A9BB0]/50 mt-6">
            Donations are processed securely via Stripe.
          </p>
        </div>
      </section>
    </>
  )
}

function TrustItem({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex items-center gap-2">
      <svg className="flex-shrink-0" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2.5" strokeLinecap="round">
        <polyline points="20 6 9 17 4 12"/>
      </svg>
      <span className="text-xs font-medium text-[#8A9BB0]">{children}</span>
    </div>
  )
}

function FeatureCard({ icon, title, body }: { icon: React.ReactNode; title: string; body: string }) {
  return (
    <div
      className="rounded-2xl p-8 border border-[#1E3050] hover:border-[#C9A84C]/30 transition-colors duration-300 group"
      style={{ background: 'linear-gradient(135deg, #162235 0%, #131F30 100%)' }}
    >
      <div className="w-12 h-12 rounded-xl bg-[#C9A84C]/10 border border-[#C9A84C]/20 flex items-center justify-center mb-5 group-hover:bg-[#C9A84C]/15 transition-colors">
        {icon}
      </div>
      <h3 className="font-serif text-xl font-semibold text-[#F0EDE8] mb-3">{title}</h3>
      <p className="text-[#8A9BB0] leading-relaxed text-sm">{body}</p>
    </div>
  )
}

function Step({
  number, icon, title, body,
}: {
  number: string; icon: React.ReactNode; title: string; body: string
}) {
  return (
    <div className="flex flex-col items-center text-center">
      <div
        className="relative w-20 h-20 rounded-full flex items-center justify-center mb-6 border border-[#C9A84C]/30 text-[#C9A84C]"
        style={{ background: 'radial-gradient(circle, #162235 60%, #0F1B2D)' }}
      >
        {icon}
        <span className="absolute -top-2 -right-2 w-6 h-6 rounded-full bg-[#C9A84C] text-[#0F1B2D] text-xs font-bold flex items-center justify-center">
          {number.slice(1)}
        </span>
      </div>
      <h3 className="font-serif text-xl font-semibold text-[#F0EDE8] mb-3">{title}</h3>
      <p className="text-[#8A9BB0] text-sm leading-relaxed max-w-xs">{body}</p>
    </div>
  )
}

function PrivacyPoint({ title, body }: { title: string; body: string }) {
  return (
    <div>
      <div className="flex items-center gap-2 mb-2">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2.5" strokeLinecap="round">
          <polyline points="20 6 9 17 4 12"/>
        </svg>
        <h4 className="font-semibold text-[#F0EDE8] text-sm">{title}</h4>
      </div>
      <p className="text-[#8A9BB0] text-sm leading-relaxed">{body}</p>
    </div>
  )
}

function Testimonial({ quote, author, role }: { quote: string; author: string; role: string }) {
  return (
    <div
      className="rounded-2xl p-7 border border-[#1E3050]"
      style={{ background: '#162235' }}
    >
      <div className="flex gap-0.5 mb-4">
        {[...Array(5)].map((_, i) => (
          <svg key={i} width="14" height="14" viewBox="0 0 24 24" fill="#C9A84C">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
          </svg>
        ))}
      </div>
      <p className="text-[#F0EDE8] leading-relaxed mb-5 text-sm">
        &ldquo;{quote}&rdquo;
      </p>
      <div>
        <p className="font-semibold text-[#F0EDE8] text-sm">{author}</p>
        <p className="text-xs text-[#8A9BB0]">{role}</p>
      </div>
    </div>
  )
}
