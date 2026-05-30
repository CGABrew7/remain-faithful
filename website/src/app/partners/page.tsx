import type { Metadata } from 'next'
import PilotForm from '@/components/PilotForm'

export const metadata: Metadata = {
  title: 'Partners',
  description:
    'Equip your men\'s ministry with Remain Faithful. Free structured accountability for churches, small groups, and discipleship cohorts.',
}

export default function PartnersPage() {
  return (
    <>
      {/* Hero */}
      <section className="pt-32 pb-20 border-b border-[#1E3050]">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-4">For Pastors & Ministry Leaders</p>
          <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-6">
            Equip Your Men&apos;s Ministry
          </h1>
          <p className="text-[#8A9BB0] text-lg max-w-2xl mx-auto">
            Remain Faithful brings real accountability technology to small groups, discipleship cohorts, and men&apos;s ministry programs — at no cost to your church.
          </p>
        </div>
      </section>

      {/* Benefits */}
      <section className="py-20">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">Why Ministries Choose RF</h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
              From small accountability triads to 20-person men&apos;s groups, RF scales to your structure.
            </p>
          </div>
          <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-5">
            {[
              {
                icon: '🤝',
                title: 'Structured Accountability',
                desc: 'Replace vague monthly check-ins with consistent, automatic alerts that keep conversations grounded in reality.',
              },
              {
                icon: '⚙️',
                title: 'Easy Group Setup',
                desc: 'Create a group, generate an invite code, share it with your men. They\'re monitoring within minutes — no IT required.',
              },
              {
                icon: '👁',
                title: 'Pastoral Oversight',
                desc: 'As group leader, you receive a weekly aggregate digest. Members control whether you see individual alerts.',
              },
              {
                icon: '🔒',
                title: 'Privacy by Design',
                desc: 'Screen content stays on member devices. You see alert metadata, not surveillance footage. Dignity is preserved.',
              },
            ].map((b) => (
              <div key={b.title} className="rounded-2xl p-6 border border-[#1E3050] bg-[#162235]">
                <div className="text-3xl mb-4">{b.icon}</div>
                <h3 className="font-serif text-lg font-semibold text-[#F0EDE8] mb-2">{b.title}</h3>
                <p className="text-sm text-[#8A9BB0] leading-relaxed">{b.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Group Setup Steps */}
      <section id="group-setup" className="py-20 bg-[#0A1420] border-y border-[#1E3050]">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">
              How to Set Up a Group
            </h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
              Five steps from download to your first accountable group.
            </p>
          </div>
          <div className="space-y-4">
            {[
              { n: '1', title: 'Create Your Leader Account', desc: 'Download Remain Faithful from the App Store and create an account. Your role as group leader is established at this step.' },
              { n: '2', title: 'Create a Group', desc: 'Tap Group tab → New Group. Give it a name (e.g., "Tuesday Mens Group"), set your covenant expectations, and choose your leader visibility settings.' },
              { n: '3', title: 'Share the Invite Code', desc: 'Your group generates a 6-character invite code. Share it in your group chat, bulletin, or Sunday handout. Members join with that code.' },
              { n: '4', title: 'Members Enable Monitoring', desc: 'Each member enables monitoring in the app and grants screen recording permission. The setup walkthrough guides them through it step by step.' },
              { n: '5', title: 'Stay Connected', desc: 'You receive a weekly digest every Monday morning. Members can alert you directly. Use RF as the backbone of your regular accountability conversations.' },
            ].map((step) => (
              <div key={step.n} className="flex gap-5 p-6 rounded-2xl border border-[#1E3050] bg-[#162235]">
                <div className="w-8 h-8 rounded-full bg-gradient-to-br from-[#C9A84C] to-[#E8C87A] text-[#0F1B2D] font-bold text-sm flex items-center justify-center flex-shrink-0">
                  {step.n}
                </div>
                <div>
                  <h3 className="font-semibold text-[#F0EDE8] mb-1">{step.title}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{step.desc}</p>
                </div>
              </div>
            ))}
          </div>

          <div className="mt-8 text-center">
            <span className="inline-flex items-center gap-2 px-6 py-3 rounded-xl border border-[#1E3050] text-[#8A9BB0] text-sm">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3"/>
              </svg>
              Group Setup Guide — Coming Soon
            </span>
          </div>
        </div>
      </section>

      {/* Testimonials */}
      <section className="py-20">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">From Ministry Leaders</h2>
          </div>
          <div className="grid md:grid-cols-2 gap-6">
            {[
              {
                quote: 'Our men\'s group had been meeting for two years without any real accountability mechanism. RF gave us that mechanism without making it weird. The alerts happen automatically — conversations happen naturally.',
                author: 'Pastor James W.',
                role: 'Bethel Community Church, Colorado',
              },
              {
                quote: 'I was skeptical of any tech solution to a spiritual problem. But RF doesn\'t replace relationship — it undergirds it. My men say the app makes them feel like someone actually notices when they struggle.',
                author: 'Rev. Marcus T.',
                role: 'Men\'s Ministry Director, Georgia',
              },
            ].map((t, i) => (
              <div key={i} className="rounded-2xl p-7 border border-[#1E3050] bg-[#162235]">
                <p className="text-[#F0EDE8] leading-relaxed mb-5 text-sm italic">
                  &ldquo;{t.quote}&rdquo;
                </p>
                <div>
                  <p className="font-semibold text-[#F0EDE8] text-sm">{t.author}</p>
                  <p className="text-xs text-[#8A9BB0]">{t.role}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Pilot Program Form */}
      <section className="py-20 bg-[#0A1420] border-t border-[#1E3050]">
        <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">
              Request a Pilot Program
            </h2>
            <p className="text-[#8A9BB0]">
              Want hands-on support setting up RF for your ministry? We offer guided pilot programs for churches and organizations.
            </p>
          </div>
          <PilotForm />
        </div>
      </section>
    </>
  )
}
