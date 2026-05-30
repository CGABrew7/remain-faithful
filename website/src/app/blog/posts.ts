export interface Post {
  slug: string
  title: string
  excerpt: string
  category: string
  date: string
  readTime: string
  body: string
}

export const posts: Post[] = [
  {
    slug: 'why-accountability-fails',
    title: 'Why Accountability Fails (And How RF Fixes It)',
    excerpt:
      "Most accountability systems fail within months. Not because the men aren't serious — but because the structure places the burden of disclosure on the person most likely to avoid it.",
    category: 'Accountability',
    date: 'May 15, 2025',
    readTime: '5 min read',
    body: `
Most accountability programs follow a predictable pattern. A men's group meets, they talk about the struggle, they make commitments to each other, someone volunteers to be an accountability partner, phone numbers are exchanged. For a few weeks, there are check-in texts. Then the texts get shorter. Then less frequent. Then they stop.

This isn't a character failure. It's a structural one.

## The Disclosure Problem

Traditional accountability requires the person who is struggling to be the one who discloses their struggle. This is backwards. The moment when someone needs accountability most — in the aftermath of a failure — is precisely the moment when the shame barrier is highest.

Asking someone to text their accountability partner and say "I fell again" is asking them to do the hardest possible thing at the hardest possible moment. Many won't. And over time, the knowledge that they won't creates a loop where the accountability relationship becomes performative — a structure that exists to make people feel safer without actually creating accountability.

## The Observer Effect

There's a well-documented phenomenon in behavioral psychology: behavior changes when we know we're being observed. This is why cameras at intersections reduce red-light violations even when the cameras aren't actively monitored. The possibility of observation changes behavior in the moment.

Traditional accountability systems create observation-by-report. You observe yourself, then report what you observed. This eliminates the observer effect entirely — because you're both the observed and the observer, and you can choose what to report.

Remain Faithful introduces actual observation. Not by another person watching your screen — but by an automatic system that removes the choice to disclose. When something is flagged, your partners know. Not because you told them. Because the system did.

This sounds intrusive. But consider: this is exactly what covenant accountability looks like in practice. "I'll know, and you'll know" is the implicit covenant in every accountability relationship. RF makes it explicit and automatic.

## The Shame-Spiral Problem

Here's the other failure mode: some accountability systems work too well. A partner finds out, responds poorly, and the accountability seeker feels worse than they did before. They disengage from the relationship. They disengage from the community. The accountability infrastructure that was supposed to help them becomes the source of additional trauma.

We built Remain Faithful with this failure mode in mind. The covenant model requires both parties to agree upfront to respond with grace. Partners see metadata — not content, not screenshots, not a record of exactly what happened. "An alert was triggered" is enough information to open a conversation. It's not enough to enable judgment.

The alert says: "Something happened." The conversation that follows determines what happens next.

## What RF Does Differently

Remain Faithful flips three things:

**Disclosure moves from opt-in to opt-out.** When you install RF and enable monitoring, disclosure is automatic. You can disable monitoring at any time — but the default is transparency.

**The shame barrier is removed from the disclosure moment.** You don't have to choose to tell your partner. The system tells them. The conversation that follows starts from a different place than "I have to confess something."

**Response is governed by covenant, not by reaction.** Partners agree before they gain access. They know what they're signing up for and what's expected of them. This shapes how they respond when an alert comes in.

None of this replaces the relational work. The technology facilitates; the humans do the real thing. But it changes what's possible — and for men who are serious about change, that change in what's possible makes all the difference.
    `.trim(),
  },
  {
    slug: 'setting-up-your-first-group',
    title: 'Setting Up Your First Accountability Group',
    excerpt:
      "A step-by-step guide to launching a Remain Faithful group in your men's ministry, friendship circle, or discipleship cohort.",
    category: 'Guide',
    date: 'May 10, 2025',
    readTime: '4 min read',
    body: `
Getting a group up and running in Remain Faithful takes about fifteen minutes for the leader and five minutes per member. Here's the complete walkthrough, including what to say when you're inviting people.

## Before You Start: The Conversation

Don't send someone a link to an app without context. The conversation matters. Here's what to cover:

- **What the app does:** It monitors for concerning content using on-device AI and notifies your accountability partners when something is flagged.
- **What partners see:** Metadata only — timestamp, category, severity. Not your screen, not your browsing history.
- **The covenant:** Both parties commit to honesty, grace, and confidentiality. You'll see this in the app.
- **It's voluntary:** Anyone can disable monitoring at any time. This is a choice, not a surveillance program.

If someone isn't ready for that conversation, they probably aren't ready for the app. That's okay.

## Step 1: The Leader Creates a Group

After creating your account:

1. Tap the **Group** tab at the bottom of the screen
2. Tap **New Group**
3. Enter a group name (e.g., "Tuesday Men," "Iron Sharpens Iron," or whatever your group calls itself)
4. Review the covenant text — this is what every member will agree to before joining
5. Choose your leader visibility settings:
   - **Alert summaries:** receive aggregate weekly data (recommended)
   - **Individual alerts:** see each alert as it occurs (use with care — discuss with your group first)
6. Tap **Create Group**

You'll see a 6-character invite code. Screenshot it or write it down.

## Step 2: Invite Your Members

Share the invite code however makes sense for your group — group chat, in person, email. Include these instructions:

> *Download Remain Faithful from the App Store. Create an account with your name and email. Tap Group → Join Group → enter code: [YOUR CODE]. Read and accept the covenant. Then enable monitoring in the Settings tab.*

That's it. The app walks them through the rest.

## Step 3: Enable Monitoring (Each Member)

Each member goes to **Settings → Monitoring → Enable**. The app will:

1. Ask for Screen Recording permission (a standard iOS dialog)
2. Walk through a brief explanation of what gets monitored
3. Confirm that monitoring is active with a green indicator on the dashboard

The first few days, encourage members to keep an eye on the dashboard to make sure everything is working. The green "Monitoring active" status in the dashboard is the confirmation.

## Step 4: First Group Meeting

Plan a meeting within the first two weeks of launch. Use it to:

- **Normalize the awkwardness.** Monitoring can feel weird at first. That's normal. Talk about it.
- **Discuss the covenant.** What does "respond with grace" mean in practice for your group?
- **Set response expectations.** When an alert comes in, what should the response look like? A text? A call? A prayer? Agree on norms.
- **Talk about the goal.** This isn't about catching each other. It's about removing the barrier to honest conversation.

## Common Questions from Members

**"What if I get an alert for something innocent?"**
False positives happen. The AI isn't perfect. If your partner gets an alert for something that wasn't what it looks like, that's a conversation — not a verdict. Talk about it.

**"Can I see what my partner sees?"**
No. The app doesn't share screen content. You see the same metadata your partners see: timestamp, category, severity.

**"What if I need to step away for a bit?"**
Monitoring can be paused anytime from the Settings tab. Partners will see that monitoring is paused — which is its own form of transparency.

## The First Month

The first month is when habits form. Encourage your group to stay engaged — respond to alerts quickly, check in with each other weekly, and have the honest conversations the app makes available. The technology opens doors; the group has to walk through them.
    `.trim(),
  },
  {
    slug: 'science-of-peer-accountability',
    title: 'The Science Behind Peer Accountability',
    excerpt:
      "Research consistently shows that peer accountability outperforms self-regulation for behavior change. Here's what the literature says — and what it means for how RF works.",
    category: 'Research',
    date: 'May 5, 2025',
    readTime: '6 min read',
    body: `
The research on peer accountability is remarkably consistent: people who make commitments to others — and who know those commitments are being observed — change their behavior more durably than those relying on willpower alone.

## Commitment Devices and Self-Control

Behavioral economists have studied "commitment devices" extensively. A commitment device is any mechanism that makes future misbehavior more costly — by creating social accountability, financial stakes, or observable commitments. The research shows these work significantly better than intentions alone.

Remain Faithful is, in technical terms, a commitment device. By installing the app and enabling monitoring, users are creating a structure that makes future misbehavior observable to people they care about. This changes the moment-of-temptation calculus in their favor.

## Social Norming Effects

Research from the field of social norms suggests that people's behavior is strongly influenced by what they believe others around them are doing. Accountability groups that talk honestly about struggle tend to normalize the struggle — which reduces shame and increases the likelihood of honest disclosure.

RF creates a structure for this conversation to happen naturally, starting from a concrete alert event rather than requiring someone to volunteer difficult information.
    `.trim(),
  },
  {
    slug: 'on-device-privacy-explained',
    title: 'On-Device AI: Why Your Content Never Leaves Your Phone',
    excerpt:
      'A technical explanation of how Apple\'s Vision framework, SensitiveContentAnalysis, and our local classifier work together to keep your screen private.',
    category: 'Technology',
    date: 'April 28, 2025',
    readTime: '7 min read',
    body: `
The centerpiece of Remain Faithful's privacy model is on-device AI classification. Here's what that means in practice, and why it matters.

## The ReplayKit Sandbox

iOS's ReplayKit framework creates a separate extension process for screen recording. This process is sandboxed — it cannot make network requests, cannot access your files, and cannot communicate with the outside world. Its only channel is a shared app group container that connects it to the main RF app.

This architectural constraint means your screen content is physically incapable of being transmitted over the network from the broadcast extension.

## Three Layers of On-Device Classification

**Layer 1: URL Blocklist + Regex**
When a browser is detected, visible URLs are checked against a local blocklist of known adult domains. Visible text is pattern-matched against regex rules for explicit content categories. This is fast, deterministic, and 100% local.

**Layer 2: Apple Vision + SensitiveContentAnalysis**
Apple provides two relevant frameworks: Vision OCR (which extracts text from screen frames) and SensitiveContentAnalysis (which detects nudity and explicit images). Both run on the device's Neural Engine — the dedicated AI chip in modern iPhones. No server involved.

**Layer 3: Local Keyword Classifier**
Our open-source keyword classifier assigns weighted scores across content categories based on the OCR output. The weights are tuned for common explicit content patterns without requiring image analysis.

Only if all three layers are uncertain does an anonymized category query reach our cloud classifier — and that query contains no screen content.
    `.trim(),
  },
  {
    slug: 'covenant-model',
    title: 'The Covenant Model: More Than an Agreement',
    excerpt:
      "The covenant that partners accept before gaining access to accountability data isn't legal boilerplate. It's the theological foundation of how RF works.",
    category: 'Faith',
    date: 'April 20, 2025',
    readTime: '5 min read',
    body: `
Every partner in Remain Faithful — before gaining access to any accountability data — accepts a covenant. This is not a terms-of-service agreement. It's a different kind of commitment.

## Why Covenant, Not Contract

A contract specifies what happens when things go wrong. A covenant specifies what kind of relationship we're committing to. The biblical covenant language is deliberate: we believe accountability is fundamentally relational, and relationship is fundamentally covenantal.

The RF covenant asks partners to commit to: honesty, grace in response, confidentiality, and the genuine flourishing of the person they're accountable for. It doesn't describe remedies for breach. It describes the character of the relationship.

## What It Changes

Partners who have explicitly agreed to respond with grace before they see a single alert are more likely to actually respond with grace when the alert comes. The covenant sets a posture before the moment of test arrives.

This is why we require covenant acceptance before any data access — not just as a consent mechanism, but as a formative moment in the partnership.
    `.trim(),
  },
  {
    slug: 'mens-ministry-accountability',
    title: 'Modernizing Men\'s Ministry Accountability',
    excerpt:
      'How churches are moving from informal check-ins to structured accountability programs — and what the next generation of men\'s ministry looks like.',
    category: 'Ministry',
    date: 'April 12, 2025',
    readTime: '8 min read',
    body: `
Men's ministry has a long history of accountability groups. What's changed is the landscape of temptation — and the tools available to meet it.

## The Smartphone as the Arena

A generation ago, accountability conversations were about what men were watching on cable TV or renting from video stores. Today, the internet puts near-unlimited explicit content in every pocket, available at 2am, with no external accountability whatsoever.

The informal accountability group hasn't kept pace. Monthly meetings and occasional texts are not matched to the moment-by-moment nature of digital temptation.

## What Structured Technology Accountability Looks Like

Churches that have implemented RF-style accountability report a consistent pattern: the conversations get better, not worse. Men who initially feel exposed by the monitoring typically report — after a few months — that the relief of not having to decide whether to disclose outweighs the discomfort of the monitoring itself.

The goal isn't surveillance. It's removing the burden of disclosure from the person who is most likely to avoid it, in the moment when they most need support.

## Building It Into Your Ministry Structure

The most successful implementations integrate RF into existing discipleship structures rather than making it a standalone program. Men's groups that were already meeting use RF as an accountability layer for what happens between meetings. One-on-one discipleship relationships use the partner accountability feature.

The technology is a backbone. The ministry is what happens around it.
    `.trim(),
  },
]
