import type { Metadata } from 'next'
import Link from 'next/link'
import { Breadcrumbs } from '@/components/Breadcrumbs'

export const metadata: Metadata = {
  title: 'Church Accountability Software | Free for Every Ministry',
  description: 'Bring structured accountability to your church, men\'s ministry, or small group. Remain Faithful is free for every church. Group mode, setup guides, and pastoral resources included.',
  alternates: { canonical: 'https://remainfaithful.com/church-accountability' },
}

export default function ChurchAccountability() {
  return (
    <div className="pt-24 pb-24">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <Breadcrumbs items={[{ name: 'Church Accountability', url: 'https://remainfaithful.com/church-accountability' }]} />

        {/* Hero */}
        <div className="mb-14">
          <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-4">For Pastors and Ministry Leaders</p>
          <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-6 leading-tight">
            Church Accountability Software, Free for Every Ministry
          </h1>
          <p className="text-[#8A9BB0] text-lg leading-relaxed mb-6">
            Remain Faithful brings structured, technology-assisted accountability to your church. It is free for every church, regardless of size. There are no licenses, no organizational tiers, and no cost to your congregation.
          </p>
          <div className="flex flex-wrap gap-4">
            <Link href="/partners" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96]">
              Register Your Church
            </Link>
            <Link href="/group-setup-guide" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#F0EDE8] border border-[#1E3050] hover:border-[#C9A84C]/50 hover:bg-[#162235] transition-colors">
              Download Group Setup Guide
            </Link>
          </div>
        </div>

        {/* Why church accountability programs fail */}
        <section className="mb-14">
          <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Why Accountability Programs Fail in Churches</h2>
          <p className="text-[#8A9BB0] leading-relaxed mb-4">
            Most church accountability efforts fail at a predictable point. A program launches with enthusiasm, groups meet for a few weeks, and then the check-ins trail off. The struggle continues; the accountability does not.
          </p>
          <div className="space-y-4">
            {[
              { title: 'The self-disclosure problem', body: 'Traditional accountability requires the struggling person to choose to disclose their failure. This is asking the hardest possible thing at the hardest possible moment. Shame makes silence easier. So people stay silent, and the accountability relationship becomes performative.' },
              { title: 'The meeting gap', body: 'Most church accountability happens once a week, at a meeting. But the moments of temptation are not scheduled for Sunday mornings. The gap between meetings is where the struggle happens, and where accountability most often fails to reach.' },
              { title: 'The response problem', body: 'When a partner does find out about a failure, they often do not know how to respond. Without a clear agreement about what the relationship requires, responses can range from inadequate to actively damaging.' },
            ].map((item) => (
              <div key={item.title} className="p-5 rounded-xl border border-[#1E3050] bg-[#162235]">
                <h3 className="font-semibold text-[#F0EDE8] mb-2">{item.title}</h3>
                <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.body}</p>
              </div>
            ))}
          </div>
        </section>

        {/* How RF addresses these */}
        <section className="mb-14">
          <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">How Remain Faithful Addresses Each Problem</h2>
          <div className="space-y-4">
            {[
              { problem: 'Self-disclosure problem', solution: 'Automatic alerts remove the choice to disclose. When something is flagged, partners know. Not because the struggling person chose to tell them. Because the system did. The shame barrier to disclosure is eliminated.' },
              { problem: 'Meeting gap', solution: 'Always-on Screen Time monitoring runs continuously between meetings. Alerts reach partners in real time. The gap between meetings is covered.' },
              { problem: 'Response problem', solution: 'The covenant is accepted before any partner gains access. It specifies: respond with grace. Partners know what is expected of them before the first alert ever arrives.' },
            ].map((item) => (
              <div key={item.problem} className="p-5 rounded-xl border border-[#C9A84C]/20 bg-[#162235]">
                <p className="text-xs text-[#C9A84C] font-semibold uppercase tracking-wider mb-1">Solves: {item.problem}</p>
                <p className="text-sm text-[#8A9BB0] leading-relaxed">{item.solution}</p>
              </div>
            ))}
          </div>
        </section>

        {/* How RF's covenant model works in ministry */}
        <section className="mb-14">
          <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">
            <Link href="/blog/covenant-model" className="hover:text-[#C9A84C] transition-colors">How the Covenant Model Works in Ministry</Link>
          </h2>
          <p className="text-[#8A9BB0] leading-relaxed mb-4">
            The covenant is the theological foundation of Remain Faithful. Before any partner gains access to accountability data, they must accept a covenant that frames the entire relationship.
          </p>
          <p className="text-[#8A9BB0] leading-relaxed mb-6">
            In a ministry context, this covenant language works naturally alongside discipleship culture. It is not a legal document; it is an explicit statement of what kind of community you are trying to build. <Link href="/blog/covenant-model" className="text-[#C9A84C] hover:underline underline-offset-2">Read more about the covenant model.</Link>
          </p>
        </section>

        {/* Implementation guide */}
        <section className="mb-14">
          <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Implementing RF in Your Ministry</h2>
          <div className="space-y-4">
            {[
              { n: '1', title: 'Start with leadership', body: 'The pastor or ministry leader installs Remain Faithful and uses it personally before asking anyone else to join. Leaders who can speak from personal experience introduce it very differently than those who are just recommending a tool.' },
              { n: '2', title: 'Read the covenant together', body: 'Before launching the app, discuss what the covenant means for your specific group. What does "respond with grace" look like in your community? This conversation shapes how alerts will be received.' },
              { n: '3', title: 'Start with one group', body: 'Choose one existing small group or discipleship cohort for a 60-day pilot. This group becomes your church\'s subject matter experts and can speak from inside experience when you expand.' },
              { n: '4', title: 'Use the Group Setup Guide', body: 'The free Group Setup Guide provides a complete walkthrough: covenant text, step-by-step setup, group norms, and a first-month FAQ designed for in-person launch conversations.' },
              { n: '5', title: 'Expand by cohort', body: 'After the pilot, expand one ministry at a time. College students, young adults, men\'s ministry. Each cohort gets its own launch conversation. One conversation at a time beats a church-wide rollout with no context.' },
            ].map((step) => (
              <div key={step.n} className="flex gap-4 p-5 rounded-2xl border border-[#1E3050] bg-[#162235]">
                <div className="w-8 h-8 rounded-full bg-gradient-to-br from-[#C9A84C] to-[#E8C87A] text-[#0F1B2D] font-bold text-sm flex items-center justify-center flex-shrink-0">
                  {step.n}
                </div>
                <div>
                  <h3 className="font-semibold text-[#F0EDE8] mb-1">{step.title}</h3>
                  <p className="text-sm text-[#8A9BB0] leading-relaxed">{step.body}</p>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Resources */}
        <section className="mb-14">
          <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">Ministry Resources</h2>
          <div className="grid sm:grid-cols-2 gap-4">
            {[
              { title: 'Group Setup Guide', desc: 'Free printable PDF for ministry leaders. Covenant text, setup steps, group norms, and first-month FAQ.', href: '/group-setup-guide', cta: 'Download PDF' },
              { title: 'Partners Program', desc: 'Register your church for early access and receive direct support from the founder during your initial launch.', href: '/partners', cta: 'Register Your Church' },
              { title: 'How RF Works', desc: 'A complete technical and theological breakdown of what Remain Faithful does, and how.', href: '/how-it-works', cta: 'Read the Breakdown' },
              { title: 'Ministry Accountability Guide', desc: 'An in-depth look at how churches are implementing structured accountability programs and what works.', href: '/blog/mens-ministry-accountability', cta: 'Read the Article' },
            ].map((item) => (
              <div key={item.title} className="p-6 rounded-2xl border border-[#1E3050] bg-[#162235]">
                <h3 className="font-serif text-lg font-semibold text-[#F0EDE8] mb-2">{item.title}</h3>
                <p className="text-sm text-[#8A9BB0] leading-relaxed mb-4">{item.desc}</p>
                <Link href={item.href} className="text-[#C9A84C] text-sm font-semibold hover:underline underline-offset-2">
                  {item.cta} →
                </Link>
              </div>
            ))}
          </div>
        </section>

        {/* CTA */}
        <div className="text-center p-10 rounded-3xl border border-[#C9A84C]/20" style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}>
          <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-4">Bring Remain Faithful to Your Church</h2>
          <p className="text-[#8A9BB0] mb-6">Free for your congregation. No license, no cost, no catch.</p>
          <div className="flex flex-wrap justify-center gap-4">
            <Link href="/partners" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96]">
              Register Your Ministry
            </Link>
            <Link href="/#waitlist" className="inline-flex items-center gap-2 px-7 py-3.5 rounded-full font-semibold text-[#F0EDE8] border border-[#1E3050] hover:border-[#C9A84C]/50 hover:bg-[#162235] transition-colors">
              Join the Waitlist
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}
