import type { Metadata } from 'next'
import Link from 'next/link'
import PilotForm from '@/components/PilotForm'

export const metadata: Metadata = {
  title: 'Partners',
  description:
    "Equip your church's accountability ministry with Remain Faithful. Free structured accountability for small groups, discipleship cohorts, recovery programs, and ministry-wide rollouts.",
  openGraph: {
    title: 'Remain Faithful for Churches & Ministries',
    description:
      'Free structured accountability for small groups, discipleship cohorts, and ministry-wide programs. No cost to your church.',
  },
}

export default function PartnersPage() {
  return (
    <>
      {/* Hero */}
      <section className="pt-32 pb-20 border-b border-[#1E3050]">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-4">For Pastors &amp; Ministry Leaders</p>
          <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-6">
            Equip Your Church&apos;s Accountability Ministry
          </h1>
          <p className="text-[#8A9BB0] text-lg max-w-2xl mx-auto">
            Remain Faithful brings real accountability technology to small groups, discipleship cohorts, and ministry programs — at no cost to your church.
          </p>
          <div className="flex flex-wrap justify-center gap-4 mt-8">
            <a
              href="#pilot"
              className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all duration-200"
            >
              Request a Church Pilot
            </a>
            <a
              href="/group-setup-guide"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#F0EDE8] border border-[#1E3050] hover:border-[#C9A84C]/50 hover:bg-[#162235] transition-all duration-200"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3"/>
              </svg>
              Download Group Setup Guide
            </a>
          </div>
        </div>
      </section>

      {/* Benefits */}
      <section className="py-20">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">Why Ministries Choose RF</h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
              From small accountability triads to church-wide programs, RF scales to your structure.
            </p>
          </div>
          <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-5">
            {[
              {
                title: 'Structured Accountability',
                desc: 'Replace vague monthly check-ins with consistent, automatic alerts that keep conversations grounded in reality.',
              },
              {
                title: 'Easy Group Setup',
                desc: 'Create a group, generate an invite code, share it with your members. They\'re monitoring within minutes. No IT required.',
              },
              {
                title: 'Pastoral Oversight',
                desc: 'As group leader, you stay informed about group activity. Members control whether you see individual alerts.',
              },
              {
                title: 'Privacy by Design',
                desc: 'Screen content stays on member devices. You see alert metadata, not surveillance footage. Dignity is preserved.',
              },
            ].map((b) => (
              <div key={b.title} className="rounded-2xl p-6 border border-[#1E3050] bg-[#162235]">
                <h3 className="font-serif text-lg font-semibold text-[#F0EDE8] mb-2">{b.title}</h3>
                <p className="text-sm text-[#8A9BB0] leading-relaxed">{b.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Group Setup Guide CTA (prominent) */}
      <section className="py-8">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div
            className="rounded-3xl p-8 sm:p-10 flex flex-col sm:flex-row items-center gap-8 border border-[#C9A84C]/30"
            style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}
          >
            <div className="w-16 h-16 rounded-2xl bg-[#C9A84C]/15 border border-[#C9A84C]/30 flex items-center justify-center flex-shrink-0">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#C9A84C" strokeWidth="2" strokeLinecap="round">
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
                <polyline points="10 9 9 9 8 9"/>
              </svg>
            </div>
            <div className="flex-1 text-center sm:text-left">
              <h3 className="font-serif text-xl font-bold text-[#F0EDE8] mb-2">Free Group Setup Guide</h3>
              <p className="text-[#8A9BB0] text-sm leading-relaxed">
                A complete printable guide for ministry leaders: covenant text, step-by-step setup, group norms, and FAQ. Designed for in-person group launches.
              </p>
            </div>
            <a
              href="/group-setup-guide"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-6 py-3 rounded-xl font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all duration-200 whitespace-nowrap text-sm"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3"/>
              </svg>
              Download PDF Guide
            </a>
          </div>
        </div>
      </section>

      {/* Men's Ministry Implementation */}
      <section className="py-20 bg-[#0A1420] border-y border-[#1E3050]">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">Implementation Guide</p>
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">Men&apos;s Ministry Implementation</h2>
            <p className="text-[#8A9BB0] max-w-2xl mx-auto">
              How to successfully roll out Remain Faithful in your men&apos;s small group, from first conversation to ongoing culture.
            </p>
          </div>

          <div className="grid md:grid-cols-2 gap-8">
            <div>
              <h3 className="font-serif text-xl font-semibold text-[#F0EDE8] mb-4">Phase 1: Preparing Leadership</h3>
              <ul className="space-y-3">
                {[
                  'The leader installs and uses RF personally before asking anyone else to join.',
                  'Review the Group Setup Guide and covenant text with co-leaders.',
                  'Decide on leader visibility settings before launch (alert summaries vs. individual alerts).',
                  'Brief your pastor on the approach — his awareness helps normalize the conversation.',
                ].map((item, i) => (
                  <li key={i} className="flex items-start gap-3 text-sm text-[#8A9BB0]">
                    <span className="flex-shrink-0 w-5 h-5 rounded-full bg-[#C9A84C]/20 flex items-center justify-center text-[#C9A84C] text-xs font-bold">{i + 1}</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>
            <div>
              <h3 className="font-serif text-xl font-semibold text-[#F0EDE8] mb-4">Phase 2: Group Launch</h3>
              <ul className="space-y-3">
                {[
                  'Introduce RF at an existing group meeting — not as an app, but as an accountability covenant.',
                  'Read the covenant text together. Discuss what "responding with grace" means for your group.',
                  'Allow members to ask questions. Invite honest hesitation — this is a trust-building moment.',
                  'Set a 30-day trial. No pressure to continue; just try it and see.',
                ].map((item, i) => (
                  <li key={i} className="flex items-start gap-3 text-sm text-[#8A9BB0]">
                    <span className="flex-shrink-0 w-5 h-5 rounded-full bg-[#C9A84C]/20 flex items-center justify-center text-[#C9A84C] text-xs font-bold">{i + 1}</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>
            <div>
              <h3 className="font-serif text-xl font-semibold text-[#F0EDE8] mb-4">Phase 3: Sustaining the Culture</h3>
              <ul className="space-y-3">
                {[
                  'Begin every meeting with a brief accountability check: "Did anything come up this week?"',
                  'When an alert fires, respond within 24 hours — a short text asking "How are you doing?" is enough.',
                  'After 90 days, revisit the covenant. Adjust based on what you\'ve learned.',
                ].map((item, i) => (
                  <li key={i} className="flex items-start gap-3 text-sm text-[#8A9BB0]">
                    <span className="flex-shrink-0 w-5 h-5 rounded-full bg-[#C9A84C]/20 flex items-center justify-center text-[#C9A84C] text-xs font-bold">{i + 1}</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>
            <div>
              <h3 className="font-serif text-xl font-semibold text-[#F0EDE8] mb-4">What Not To Do</h3>
              <ul className="space-y-3">
                {[
                  'Don\'t mandate RF without conversation — it must be chosen, not required.',
                  'Don\'t use alerts as discipline tools. They are conversation starters, not verdicts.',
                  'Don\'t respond to an alert with shame. The covenant requires grace.',
                  'Don\'t expect the app to do the relational work. RF is infrastructure, not relationship.',
                ].map((item, i) => (
                  <li key={i} className="flex items-start gap-3 text-sm text-[#8A9BB0]">
                    <svg className="flex-shrink-0 mt-0.5 text-red-400" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                      <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
                    </svg>
                    {item}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      </section>

      {/* Recovery Ministry Integration */}
      <section className="py-20 border-b border-[#1E3050]">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">Recovery Ministry</p>
              <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-5">
                Recovery Ministry Integration
              </h2>
              <p className="text-[#8A9BB0] leading-relaxed mb-4">
                Remain Faithful integrates naturally with Celebrate Recovery and similar programs. For participants working through sexual purity or addiction to pornography, RF provides an automatic accountability layer that works between weekly meetings.
              </p>
              <p className="text-[#8A9BB0] leading-relaxed mb-4">
                The key alignment with 12-step and recovery models: RF removes the decision to self-disclose at the moment of greatest shame. It makes honesty automatic, which is exactly what recovery programs teach but struggle to enforce between meetings.
              </p>
              <p className="text-[#8A9BB0] leading-relaxed mb-6">
                RF can pair with any existing sponsor or accountability partner structure — the sponsor relationship already exists; RF simply adds a consistent signal layer to it.
              </p>
              <div className="space-y-3">
                {[
                  'Compatible with Celebrate Recovery accountability structures',
                  'Works alongside existing sponsor relationships',
                  'Provides automatic accountability between weekly meetings',
                  'Shame-resistant — partners see metadata, not content',
                  'Supports one-to-one sponsor relationships or CR small groups',
                ].map((item) => (
                  <div key={item} className="flex items-center gap-2 text-sm text-[#8A9BB0]">
                    <svg className="flex-shrink-0 text-[#C9A84C]" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                      <polyline points="20 6 9 17 4 12"/>
                    </svg>
                    {item}
                  </div>
                ))}
              </div>
            </div>
            <div className="space-y-4">
              <div
                className="rounded-2xl p-7 border border-[#C9A84C]/20"
                style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}
              >
                <h3 className="font-semibold text-[#F0EDE8] mb-3">Suggested CR Integration</h3>
                <ol className="space-y-2">
                  {[
                    'Introduce RF at Step 5 or Step 8 groups where accountability already exists.',
                    'Pair RF with existing accountability partner assignments.',
                    'Use the Group feature for CR small groups of 4–8 participants.',
                    'CR leader receives the weekly digest as pastoral oversight.',
                    'Use RF alerts as a bridge topic at weekly CR meetings.',
                  ].map((step, i) => (
                    <li key={i} className="flex items-start gap-2 text-sm text-[#8A9BB0]">
                      <span className="text-[#C9A84C] font-semibold flex-shrink-0">{i + 1}.</span>
                      {step}
                    </li>
                  ))}
                </ol>
              </div>
              <div className="rounded-2xl p-5 border border-[#1E3050] bg-[#162235]">
                <p className="text-xs text-[#8A9BB0]/70 uppercase tracking-wide mb-2">Important Note</p>
                <p className="text-sm text-[#8A9BB0] leading-relaxed">
                  RF is not a clinical tool and is not a replacement for counseling, therapy, or professional addiction treatment. It is a peer accountability tool designed for voluntary use within healthy community structures.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Small Group Deployment Checklist */}
      <section className="py-20 bg-[#0A1420] border-b border-[#1E3050]">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">Small Group</p>
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">Small Group Deployment Checklist</h2>
            <p className="text-[#8A9BB0] max-w-xl mx-auto">
              Everything a group leader needs to confirm before, during, and after launch.
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            {[
              {
                phase: 'Before Launch',
                color: '#C9A84C',
                items: [
                  'Leader has installed RF and enabled monitoring personally',
                  'Leader has read the Group Setup Guide',
                  'Covenant text has been reviewed and is ready to discuss',
                  'Leader has briefed their pastor or supervisor',
                  'A group meeting date is scheduled for the launch conversation',
                ],
              },
              {
                phase: 'During Launch Week',
                color: '#C9A84C',
                items: [
                  'Covenant discussed and agreed on by all members',
                  'All members have installed RF and created accounts',
                  'Invite code shared and all members have joined the group',
                  'All members have enabled monitoring successfully',
                  'Group has agreed on alert response norms',
                ],
              },
              {
                phase: 'First 30 Days',
                color: '#C9A84C',
                items: [
                  'At least one alert has been responded to with grace',
                  'Group has had at least one conversation that started with an alert',
                  '30-day check-in meeting has occurred',
                  'Any false-positive alerts have been discussed and normalized',
                ],
              },
            ].map((col) => (
              <div key={col.phase} className="rounded-2xl border border-[#1E3050] bg-[#162235] overflow-hidden">
                <div className="h-1.5 bg-gradient-to-r from-[#C9A84C] to-[#E8C87A]" />
                <div className="p-6">
                  <h3 className="font-semibold text-[#F0EDE8] mb-4">{col.phase}</h3>
                  <ul className="space-y-3">
                    {col.items.map((item, i) => (
                      <li key={i} className="flex items-start gap-2 text-sm text-[#8A9BB0]">
                        <div className="w-4 h-4 mt-0.5 rounded border border-[#C9A84C]/40 flex-shrink-0" />
                        {item}
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Church-Wide Accountability Program */}
      <section className="py-20 border-b border-[#1E3050]">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14">
            <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">Church-Wide</p>
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">Church-Wide Accountability Program</h2>
            <p className="text-[#8A9BB0] max-w-2xl mx-auto">
              How to adopt Remain Faithful across multiple groups, ministries, or age cohorts within a single church body.
            </p>
          </div>

          <div className="space-y-6">
            {[
              {
                step: '1',
                title: 'Start with a Pilot Group',
                desc: 'Select one existing small group or discipleship cohort for a 60-day pilot. This group becomes your church\'s RF champions and can speak to the experience from the inside.',
              },
              {
                step: '2',
                title: 'Document the Pilot Experience',
                desc: 'After 60 days, gather feedback from the pilot group. What worked? What felt awkward? What conversation happened because of RF that wouldn\'t have happened otherwise? This documentation guides the broader rollout.',
              },
              {
                step: '3',
                title: 'Train Group Leaders',
                desc: 'Host a 30-minute leader training session — in person or via video — walking through the app, the covenant, and how to respond to alerts. Leaders who understand the "why" launch groups that stick.',
              },
              {
                step: '4',
                title: 'Roll Out by Ministry Cohort',
                desc: 'Expand one ministry at a time — college students, young adults, adult small groups — rather than church-wide all at once. Each cohort gets its own launch conversation and covenant discussion.',
              },
              {
                step: '5',
                title: 'Maintain Pastoral Oversight',
                desc: 'Assign one staff member or elder as the RF point person. They receive aggregate data from group leaders, address questions, and can escalate to professional care when alerts suggest deeper need.',
              },
            ].map((step) => (
              <div key={step.step} className="flex gap-5 p-6 rounded-2xl border border-[#1E3050] bg-[#162235]">
                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-[#C9A84C] to-[#E8C87A] text-[#0F1B2D] font-bold text-sm flex items-center justify-center flex-shrink-0">
                  {step.step}
                </div>
                <div>
                  <h3 className="font-semibold text-[#F0EDE8] mb-1">{step.title}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{step.desc}</p>
                </div>
              </div>
            ))}
          </div>

          <div
            className="mt-10 p-7 rounded-2xl border border-[#C9A84C]/20"
            style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}
          >
            <h3 className="font-semibold text-[#F0EDE8] mb-3">A Note on Scale</h3>
            <p className="text-sm text-[#8A9BB0] leading-relaxed">
              Remain Faithful is free regardless of how many groups your church runs. Whether you have 2 groups or 20, there is no cost. Our model is donor-funded so that no church faces a barrier to adoption. If your church has found RF valuable, consider{' '}
              <a href="/#donate" className="text-[#C9A84C] hover:underline">making a donation</a> to support the infrastructure for others.
            </p>
          </div>
        </div>
      </section>

      {/* Group Setup Steps */}
      <section id="group-setup" className="py-20 bg-[#0A1420] border-b border-[#1E3050]">
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
              { n: '2', title: 'Create a Group', desc: 'Tap Group tab → New Group. Give it a name (e.g., "Tuesday Accountability Group"), set your covenant expectations, and choose your leader visibility settings.' },
              { n: '3', title: 'Share the Invite Code', desc: 'Your group generates a 6-character invite code. Share it in your group chat, bulletin, or Sunday handout. Members join with that code.' },
              { n: '4', title: 'Members Enable Monitoring', desc: 'Each member enables monitoring in the app and grants screen recording permission. The setup walkthrough guides them through it step by step.' },
              { n: '5', title: 'Stay Connected', desc: 'Members can alert you directly. Use RF as the backbone of your regular accountability conversations.' },
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
            <a
              href="/group-setup-guide"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-6 py-3 rounded-xl border border-[#C9A84C]/40 bg-[#C9A84C]/10 text-[#C9A84C] text-sm font-medium hover:bg-[#C9A84C]/20 transition-colors"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3"/>
              </svg>
              Download Group Setup Guide (PDF)
            </a>
          </div>
        </div>
      </section>

      {/* Church Pilot Form */}
      <section id="pilot" className="py-20 bg-[#0A1420] border-t border-[#1E3050]">
        <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">Get Started</p>
            <h2 className="font-serif text-3xl font-bold text-[#F0EDE8] mb-4">
              Request a Church Pilot
            </h2>
            <p className="text-[#8A9BB0]">
              Want hands-on support setting up Remain Faithful for your ministry? We offer guided pilot programs with direct leader support, customized rollout plans, and follow-up check-ins.
            </p>
          </div>
          <PilotForm />
        </div>
      </section>
    </>
  )
}
