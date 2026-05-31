import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Privacy Policy',
  description: 'How Remain Faithful collects, uses, and protects your data. Your screen content stays on your device.',
}

export default function PrivacyPage() {
  return (
    <div className="pt-32 pb-24">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-12">
          <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">Legal</p>
          <h1 className="font-serif text-4xl font-bold text-[#F0EDE8] mb-4">Privacy Policy</h1>
          <p className="text-[#8A9BB0] text-sm">Last updated: May 2025</p>
        </div>

        <div className="prose-content space-y-12">
          <Section title="Our Commitment">
            <p className="font-semibold text-[#F0EDE8] leading-relaxed">
              Your name, email address, and personal information will never be shared, sold, or disclosed to any third party, at any time, for any reason. This is a permanent commitment.
            </p>
            <p>
              Remain Faithful is an accountability app for adults committed to purity. Because of the deeply personal nature of what this app does, we have designed our privacy model first and built the technology to match it. This policy describes exactly what data we collect, what stays on your device, and what your partners can see.
            </p>
            <p>
              If you have any questions about this policy, contact us at{' '}
              <a href="mailto:jeff@hanokventures.co" className="text-[#C9A84C] hover:underline">
                jeff@hanokventures.co
              </a>.
            </p>
          </Section>

          <Section title="What We Collect">
            <Subsection title="Account Information">
              <p>When you create an account, we collect:</p>
              <ul>
                <li>Your display name (the name your partners will see)</li>
                <li>Your email address (used for authentication and communications)</li>
                <li>A hashed password (we never store plaintext passwords)</li>
              </ul>
            </Subsection>

            <Subsection title="Device Information">
              <p>We collect your APNs (Apple Push Notification service) device token to deliver accountability alerts to your partners. This token identifies your device to Apple&apos;s notification system but does not reveal your identity to us beyond what your account already provides.</p>
            </Subsection>

            <Subsection title="Alert Metadata">
              <p>When your device&apos;s monitoring detects concerning content, an alert is generated. We store and transmit the following alert metadata:</p>
              <ul>
                <li>Timestamp of the alert</li>
                <li>Content category (e.g., &quot;Adult Content&quot;, &quot;Explicit Text&quot;)</li>
                <li>Severity level (Low / Medium / High)</li>
                <li>Which app or browser was active (app bundle ID or domain category, not full URL)</li>
              </ul>
              <p>
                <strong className="text-[#F0EDE8]">We do not collect, store, or transmit:</strong> screenshots, screen recordings, raw OCR text, browsing history, app content, or any representation of what was on your screen.
              </p>
            </Subsection>

            <Subsection title="Usage Analytics">
              <p>
                We use Google Analytics (GA4) to understand how users navigate the website. This is optional. You can opt out via your browser&apos;s Do Not Track setting or a GA opt-out extension. The app itself does not include analytics SDKs.
              </p>
            </Subsection>

            <Subsection title="Donation Information">
              <p>
                Donations are processed by Stripe. We do not store your payment card information. Stripe stores payment data per their own privacy policy. We receive only a confirmation of donation amount for internal records.
              </p>
            </Subsection>
          </Section>

          <Section title="What Stays On Your Device">
            <p>The following data is processed entirely on your device and is never transmitted to our servers or to your partners:</p>
            <ul>
              <li>Raw screen frames captured by the broadcast extension</li>
              <li>OCR text extracted by Apple Vision</li>
              <li>SensitiveContentAnalysis classification results</li>
              <li>Local keyword classifier scores</li>
              <li>Full browsing URLs</li>
              <li>Any visual content from your screen</li>
            </ul>
            <p>
              The broadcast extension runs in a sandboxed process that cannot access the internet. It can only communicate with the main Remain Faithful app process via a shared app group container, ensuring that no screen content can be exfiltrated.
            </p>
          </Section>

          <Section title="What Your Partners See">
            <p>Your accountability partners have access to:</p>
            <ul>
              <li>Your display name and account email (as provided by you when you initiated the partnership)</li>
              <li>Alert metadata as described above: timestamp, category, severity level</li>
              <li>Your streak count and weekly digest summary (aggregate counts, not individual events)</li>
            </ul>
            <p>Partners do not have access to your screen content, browsing history, or any data beyond alert metadata and account-level statistics you have explicitly made visible.</p>
          </Section>

          <Section title="Data Retention">
            <p>We retain data for the following periods:</p>
            <ul>
              <li><strong className="text-[#F0EDE8]">Account data:</strong> retained while your account is active, deleted within 30 days of account deletion</li>
              <li><strong className="text-[#F0EDE8]">Alert history:</strong> configurable by you in Settings → Data Retention. Options: 7, 14, 30, or 90 days. Alerts older than your selected window are automatically purged.</li>
              <li><strong className="text-[#F0EDE8]">Partnership records:</strong> removed immediately when either partner ends the relationship</li>
              <li><strong className="text-[#F0EDE8]">Donation records:</strong> retained for 7 years per financial record-keeping obligations</li>
            </ul>
          </Section>

          <Section title="Third-Party Services">
            <ul>
              <li>
                <strong className="text-[#F0EDE8]">Apple (APNs):</strong> We use Apple&apos;s Push Notification service to deliver alerts. Apple&apos;s privacy policy applies to their handling of device tokens.
              </li>
              <li>
                <strong className="text-[#F0EDE8]">Stripe:</strong> Donation payment processing. Stripe is PCI-compliant. We do not store card data. See Stripe&apos;s privacy policy for details.
              </li>
              <li>
                <strong className="text-[#F0EDE8]">Anthropic:</strong> In limited cases where on-device classification is uncertain, a text category query (never image content) may be sent to Anthropic&apos;s API for classification. No personally identifiable information is included in these requests.
              </li>
              <li>
                <strong className="text-[#F0EDE8]">Google Analytics:</strong> Website analytics only. Not present in the iOS app.
              </li>
            </ul>
          </Section>

          <Section title="Your Rights">
            <p>You may, at any time:</p>
            <ul>
              <li><strong className="text-[#F0EDE8]">Access your data:</strong> request a full export of your account data via Settings → Export My Data, or by emailing jeff@hanokventures.co</li>
              <li><strong className="text-[#F0EDE8]">Correct your data:</strong> update your display name and email in Settings → Edit Profile</li>
              <li><strong className="text-[#F0EDE8]">Delete your data:</strong> delete your account in Settings → Delete Account. This triggers immediate deletion of your account data and alerts, with confirmation within 30 days</li>
              <li><strong className="text-[#F0EDE8]">Withdraw consent:</strong> disable monitoring at any time from within the app. Partners will no longer receive alerts immediately upon disabling.</li>
            </ul>
            <p>
              For data requests or concerns, contact{' '}
              <a href="mailto:jeff@hanokventures.co" className="text-[#C9A84C] hover:underline">
                jeff@hanokventures.co
              </a>. We respond within 30 days.
            </p>
          </Section>

          <Section title="Children">
            <p>
              Remain Faithful is not intended for users under the age of 18. We do not knowingly collect data from anyone under 18. If you believe we have inadvertently collected data from a minor, contact us immediately at jeff@hanokventures.co and we will delete it promptly.
            </p>
          </Section>

          <Section title="Security">
            <p>
              All data in transit is encrypted using TLS 1.3. Account passwords are hashed using bcrypt. JWT authentication tokens are stored in the iOS Keychain, not in UserDefaults or plain storage. We do not log sensitive data in our server logs.
            </p>
          </Section>

          <Section title="Changes to This Policy">
            <p>
              We will post any material changes to this policy on this page with a new &quot;Last updated&quot; date. For significant changes, we will notify active users via push notification or email. Your continued use of the app after a policy change constitutes acceptance of the updated policy.
            </p>
          </Section>

          <Section title="Contact">
            <p>
              Remain Faithful<br />
              Email: <a href="mailto:jeff@hanokventures.co" className="text-[#C9A84C] hover:underline">jeff@hanokventures.co</a><br />
              GitHub: <a href="https://github.com/remainfaithful" className="text-[#C9A84C] hover:underline" target="_blank" rel="noopener noreferrer">github.com/remainfaithful</a>
            </p>
          </Section>
        </div>
      </div>
    </div>
  )
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="border-t border-[#1E3050] pt-10">
      <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-5">{title}</h2>
      <div className="space-y-4 text-[#8A9BB0] leading-relaxed [&_ul]:space-y-2 [&_ul]:list-none [&_ul_li]:flex [&_ul_li]:gap-3 [&_ul_li]:before:content-['–'] [&_ul_li]:before:text-[#C9A84C] [&_ul_li]:before:flex-shrink-0">
        {children}
      </div>
    </div>
  )
}

function Subsection({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div>
      <h3 className="font-semibold text-[#F0EDE8] mb-2">{title}</h3>
      <div className="space-y-3">{children}</div>
    </div>
  )
}
