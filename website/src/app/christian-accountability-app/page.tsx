import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { softwareApplicationSchema } from '@/lib/structured-data'
import { Breadcrumbs } from '@/components/Breadcrumbs'

export const metadata: Metadata = {
  title: 'Christian Accountability App | Built on Covenant, Not Surveillance',
  description: 'Remain Faithful is an accountability app built for Christians. Covenant-based partner system, on-device AI, small group support, and pastoral resources. Free forever.',
  alternates: { canonical: 'https://remainfaithful.com/christian-accountability-app' },
}

export default function ChristianAccountabilityApp() {
  return (
    <>
      <JsonLd data={softwareApplicationSchema} />

      <div className="pt-24 pb-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <Breadcrumbs items={[{ name: 'Christian Accountability App', url: 'https://remainfaithful.com/christian-accountability-app' }]} />

          {/* Hero */}
          <div className="mb-14">
            <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-4">Built for the Church</p>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-6 leading-tight">
              Christian Accountability App Built on Covenant, Not Surveillance
            </h1>
            <p className="text-[#8A9BB0] text-lg leading-relaxed mb-6">
              Most accountability tools treat the problem as a monitoring problem. Remain Faithful treats it as a relational one. The technology facilitates; the covenant does the real work.
            </p>
            <div className="flex flex-wrap gap-4">
              <Link href="/#waitlist" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96]">
                Get Early Access (Free)
              </Link>
              <Link href="/how-it-works" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#F0EDE8] border border-[#1E3050] hover:border-[#C9A84C]/50 hover:bg-[#162235] transition-colors">
                How It Works
              </Link>
            </div>
          </div>

          {/* Biblical basis */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">The Biblical Case for Accountability</h2>
            <p className="text-[#8A9BB0] leading-relaxed mb-6">
              Accountability is not a modern invention. Scripture has always described it as central to how believers pursue holiness together.
            </p>
            <div className="space-y-4">
              {[
                {
                  verse: 'Proverbs 27:17',
                  text: '"As iron sharpens iron, so one person sharpens another."',
                  commentary: 'Accountability requires friction. It requires someone close enough to push back, to ask hard questions, to say what a struggling person needs to hear. Remain Faithful builds the infrastructure for that kind of relationship.',
                },
                {
                  verse: 'James 5:16',
                  text: '"Therefore confess your sins to each other and pray for each other so that you may be healed."',
                  commentary: 'Confession to one another is a spiritual practice, not just a therapeutic tool. Remain Faithful\'s automatic alerts change the structure of confession: instead of requiring a struggling person to choose to disclose at their most vulnerable moment, the system creates the disclosure automatically, so the conversation that follows can focus on healing.',
                },
                {
                  verse: 'Ecclesiastes 4:9-12',
                  text: '"Two are better than one... If either of them falls down, one can help the other up."',
                  commentary: 'This is the simplest case for peer accountability. One person, trying alone, is easier to defeat. Two people, committed to each other, create a different structure for pursuing faithfulness.',
                },
              ].map((item) => (
                <div key={item.verse} className="p-6 rounded-2xl border border-[#C9A84C]/20 bg-[#162235]">
                  <p className="text-[#C9A84C] text-xs font-semibold uppercase tracking-wider mb-2">{item.verse}</p>
                  <p className="font-serif text-[#F0EDE8] italic mb-3 leading-relaxed">{item.text}</p>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.commentary}</p>
                </div>
              ))}
            </div>
          </section>

          {/* The covenant model */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">
              <Link href="/blog/covenant-model" className="hover:text-[#C9A84C] transition-colors">The Covenant Model</Link>
            </h2>
            <p className="text-[#8A9BB0] leading-relaxed mb-4">
              Before any partner gains access to your accountability data, they agree to a covenant. Not a terms-of-service checkbox. A covenant is a statement about the kind of relationship they intend to have with you.
            </p>
            <div className="p-8 rounded-2xl border border-[#C9A84C]/20 mb-6" style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}>
              <h3 className="font-serif text-lg font-semibold text-[#F0EDE8] mb-4">The Covenant Partners Accept</h3>
              <ul className="space-y-3">
                {[
                  'I will be honest with my partner, even when it is difficult.',
                  'I will not use this app to condemn my partner.',
                  'I will respond to alerts with grace and genuine care.',
                  'I will not share my partner\'s alerts with others.',
                  'I will pursue my partner\'s flourishing above my own curiosity.',
                ].map((line, i) => (
                  <li key={i} className="flex gap-3 text-sm text-[#8A9BB0] leading-relaxed">
                    <span className="text-[#C9A84C] mt-0.5 flex-shrink-0">✦</span>
                    {line}
                  </li>
                ))}
              </ul>
            </div>
            <p className="text-[#8A9BB0] leading-relaxed">
              This covenant is why the alerts work. A partner who has already committed to grace before seeing a single alert is a different kind of partner than one who agreed to a terms-of-service page. <Link href="/blog/covenant-model" className="text-[#C9A84C] hover:underline underline-offset-2">Read more about the covenant model.</Link>
            </p>
          </section>

          {/* What makes RF different */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">How RF Differs from Secular Monitoring Tools</h2>
            <div className="grid md:grid-cols-2 gap-4">
              {[
                { label: 'Framework', secular: 'Behavioral compliance', rf: 'Covenantal relationship' },
                { label: 'Goal', secular: 'Reduce problematic behavior', rf: 'Pursue genuine holiness together' },
                { label: 'Alert response', secular: 'Predetermined rules', rf: 'Agreed covenant to respond with grace' },
                { label: 'Partner access', secular: 'Usually automatic upon install', rf: 'Requires explicit covenant acceptance' },
                { label: 'Failure handling', secular: 'Report generated, consequences possible', rf: 'Alert opens a conversation, not a verdict' },
                { label: 'Cost model', secular: 'Subscription (product = accountability)', rf: 'Free, donation-funded (mission = accountability)' },
              ].map((row) => (
                <div key={row.label} className="rounded-xl border border-[#1E3050] bg-[#162235] overflow-hidden">
                  <div className="px-4 py-2 border-b border-[#1E3050] bg-[#0A1420]">
                    <p className="text-xs text-[#8A9BB0] font-semibold uppercase tracking-wide">{row.label}</p>
                  </div>
                  <div className="p-4 grid grid-cols-2 gap-3 text-xs">
                    <div>
                      <p className="text-[#8A9BB0]/60 mb-1">Secular tools</p>
                      <p className="text-[#8A9BB0]">{row.secular}</p>
                    </div>
                    <div>
                      <p className="text-[#C9A84C] mb-1">Remain Faithful</p>
                      <p className="text-[#F0EDE8]">{row.rf}</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </section>

          {/* Church and small group */}
          <section className="mb-14">
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">For Churches and Small Groups</h2>
            <p className="text-[#8A9BB0] leading-relaxed mb-6">
              Remain Faithful was built with the church in mind. The group mode supports up to 12 members, making it ideal for men&apos;s ministry triads, discipleship cohorts, and accountability circles within a congregation.
            </p>
            <div className="grid sm:grid-cols-2 gap-4 mb-6">
              {[
                { title: 'Built-in group mode', body: 'Create a group, share an invite code, and have your entire small group monitoring and supporting each other within minutes.' },
                { title: 'Pastoral oversight tools', body: 'Group leaders can receive weekly summaries or individual alerts, with member-controlled visibility settings.' },
                { title: 'No cost to your church', body: 'Remain Faithful is free for churches, small groups, and ministries of any size. No organizational license required.' },
                { title: 'Group Setup Guide', body: 'A free printable guide for ministry leaders covers covenant discussion, setup walkthrough, and first-month expectations.' },
              ].map((item) => (
                <div key={item.title} className="p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
                  <h3 className="font-semibold text-[#F0EDE8] mb-2">{item.title}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.body}</p>
                </div>
              ))}
            </div>
            <div className="flex flex-wrap gap-4">
              <Link href="/partners" className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl border border-[#C9A84C]/30 bg-[#C9A84C]/10 text-[#C9A84C] text-sm font-semibold hover:bg-[#C9A84C]/20 transition-colors">
                Church and Ministry Resources →
              </Link>
              <Link href="/blog/setting-up-your-first-group" className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl border border-[#1E3050] bg-[#162235] text-[#8A9BB0] text-sm font-semibold hover:border-[#C9A84C]/40 transition-colors">
                Group Setup Guide →
              </Link>
            </div>
          </section>

          {/* CTA */}
          <div className="text-center p-10 rounded-3xl border border-[#C9A84C]/20" style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}>
            <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-4">Start Your Accountability Journey</h2>
            <p className="text-[#8A9BB0] mb-6">Free for individuals, couples, friends, and churches. Always.</p>
            <Link href="/#waitlist" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96]">
              Join the Waitlist
            </Link>
          </div>
        </div>
      </div>
    </>
  )
}
