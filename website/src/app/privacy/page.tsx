import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Privacy Policy',
  description: 'How Remain Faithful collects, uses, and protects your data. Screenshots and screen text stay on your phone. Partners can receive a short summary with the category.',
  alternates: { canonical: 'https://remainfaithful.com/privacy' },
}

export default function PrivacyPage() {
  return (
    <div className="pt-32 pb-24">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-12">
          <p className="text-wax text-sm font-semibold uppercase tracking-widest mb-3">Legal</p>
          <h1 className="font-serif text-4xl font-bold text-ink mb-4">Privacy Policy</h1>
          <p className="text-ink-soft text-sm">Last updated: October 2026</p>
        </div>

        <div className="prose-content space-y-12">
          <Section title="Our Commitment">
            <p className="font-semibold text-ink leading-relaxed">
              We do not sell your name or email. Partners you accept can see the name and email on your account, plus alert metadata. Apple receives a device token so alerts can be delivered. Stripe processes donations on the website. Website analytics run only if you allow them.
            </p>
            <p>
              Remain Faithful is an accountability app for adults committed to purity. Because of the deeply personal nature of what this app does, we have designed our privacy model first and built the technology to match it. This policy describes exactly what data we collect, what stays on your device, and what your partners can see.
            </p>
            <p>
              If you have any questions about this policy, contact us at{' '}
              <a href="mailto:support@remainfaithful.com" className="text-wax hover:underline">
                support@remainfaithful.com
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
              <p>If a chosen app is used for a minute, or if optional Deep Scan flags a non-DRM frame, an alert can be stored. The phone also sends a heartbeat while monitoring is on. That alert metadata is:</p>
              <ul>
                <li>Timestamp of the alert</li>
                <li>Content category (e.g., &quot;Adult Content&quot;, &quot;Explicit Text&quot;)</li>
                <li>Severity level (Low / Medium / High)</li>
                <li>A short system-generated summary string</li>
              </ul>
              <p>
                <strong className="text-ink">We do not collect, store, or transmit:</strong> screenshots, screen recordings, raw OCR text, browsing history, or raw screen content. Partner sharing does not include those items beyond the short summary above.
              </p>
            </Subsection>

            <Subsection title="Usage Analytics">
              <p>
                Website analytics stay off until you tap Allow analytics. The choice is saved in this browser. The iPhone app does not include analytics.
              </p>
            </Subsection>

            <Subsection title="Donation Information">
              <p>
                Donations are processed by Stripe. We do not store your payment card information. Stripe stores payment data per their own privacy policy. We receive only a confirmation of donation amount for internal records.
              </p>
            </Subsection>

            <Subsection title="Support Messages">
              <p>
                If you use Send ideas or report a problem, we receive the type you pick, the words you type, and the name and email on your account so we can reply. App version, iPhone model, and iOS version are included only if you turn that on. The form does not send logs, alerts, screenshots, Screen Time choices, partner names, or group names.
              </p>
            </Subsection>
          </Section>

          <Section title="What Stays On Your Device">
            <p>Always-on Family Controls filtering and DeviceActivity usage monitoring do not capture screen frames. If you start an optional Deep Scan session, the following data is processed entirely on your device and is never transmitted to our servers or to your partners:</p>
            <ul>
              <li>Raw screen frames captured by the ReplayKit broadcast extension (Deep Scan only)</li>
              <li>OCR text extracted by Apple Vision</li>
              <li>SensitiveContentAnalysis classification results</li>
              <li>Local keyword classifier scores</li>
              <li>Full browsing URLs</li>
              <li>Any visual content from your screen</li>
            </ul>
            <p>
              When Deep Scan flags a non-DRM frame, alert metadata is uploaded: category, severity, timestamp, and a short system-generated summary string. Screenshots, OCR text, and raw screen content never leave the device. The iOS app, backend, and website source are public on GitHub. Design notes and audit markdown in that repo are historical working papers, not product promises.
            </p>
          </Section>

          <Section title="What Your Partners See">
            <p>Your accountability partners have access to:</p>
            <ul>
              <li>Your display name and account email (as provided by you when you initiated the partnership)</li>
              <li>Alert metadata as described above: timestamp, category, severity level, and a short system-generated summary — not which app, and not a bundle ID</li>
              <li>Your streak count (aggregate, not individual events)</li>
            </ul>
            <p>Partners do not have access to your screen content, browsing history, or any data beyond alert metadata and account-level statistics you have explicitly made visible.</p>
          </Section>

          <Section title="Data Retention">
            <p>We retain data for the following periods:</p>
            <ul>
              <li><strong className="text-ink">Account data:</strong> deleted when you delete the account</li>
              <li><strong className="text-ink">Alert history:</strong> kept while your account exists, and deleted with the account. Older events and alerts are also removed after 365 days. Leaving a group notifies the other members. It does not delete your alert history.</li>
              <li><strong className="text-ink">Partnership records:</strong> removed when either person ends the relationship</li>
              <li><strong className="text-ink">Donation rows on your account:</strong> removed with the account. Stripe keeps its own payment records.</li>
            </ul>
          </Section>

          <Section title="Third-Party Services">
            <ul>
              <li>
                <strong className="text-ink">Apple (APNs):</strong> We use Apple&apos;s Push Notification service to deliver alerts. Apple&apos;s privacy policy applies to their handling of device tokens.
              </li>
              <li>
                <strong className="text-ink">Stripe:</strong> Donation payment processing. Stripe is PCI-compliant. We do not store card data. See Stripe&apos;s privacy policy for details.
              </li>
              <li>
                <strong className="text-ink">Google Analytics:</strong> Website analytics only. Not present in the iOS app.
              </li>
            </ul>
          </Section>

          <Section title="Your Rights">
            <p>You may, at any time:</p>
            <ul>
              <li><strong className="text-ink">Access your data:</strong> Settings → View My Activity Log shows your flagged activity in the app. For anything else, email support@remainfaithful.com. There is no separate export button.</li>
              <li><strong className="text-ink">Correct your data:</strong> update your display name and email in Settings → Edit Profile</li>
              <li><strong className="text-ink">Delete your data:</strong> delete your account in Settings, then Delete Account. The account, alerts, and donation rows tied to the account are removed then.</li>
              <li><strong className="text-ink">Change protection:</strong> you can change App Restrictions and turn off App Lockout in Settings. Turning off App Lockout sends your partners a notice. Stopping Deep Scan notifies partners only when App Lockout is on.</li>
            </ul>
            <p>
              For data requests or concerns, contact{' '}
              <a href="mailto:support@remainfaithful.com" className="text-wax hover:underline">
                support@remainfaithful.com
              </a>. We respond within 30 days.
            </p>
          </Section>

          <Section title="Children">
            <p>
              Remain Faithful is not intended for users under the age of 18. We do not knowingly collect data from anyone under 18. If you believe we have inadvertently collected data from a minor, contact us immediately at support@remainfaithful.com and we will delete it promptly.
            </p>
          </Section>

          <Section title="Security">
            <p>
              The app and the server talk over HTTPS. Passwords are stored as bcrypt hashes. The session token is kept in the iOS Keychain and in the app group so Screen Time extensions can send events. Server logs do not include reset links or the address an email was sent to.
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
              Email: <a href="mailto:support@remainfaithful.com" className="text-wax hover:underline">support@remainfaithful.com</a><br />
              GitHub: <a href="https://github.com/CGABrew7/remain-faithful" className="text-wax hover:underline" target="_blank" rel="noopener noreferrer">github.com/CGABrew7/remain-faithful</a>
            </p>
          </Section>
        </div>
      </div>
    </div>
  )
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="border-t border-hairline pt-10">
      <h2 className="font-serif text-2xl font-bold text-ink mb-5">{title}</h2>
      <div className="space-y-4 text-ink-soft leading-relaxed [&_ul]:space-y-2 [&_ul]:list-none [&_ul_li]:flex [&_ul_li]:gap-3 [&_ul_li]:before:content-['–'] [&_ul_li]:before:text-wax [&_ul_li]:before:flex-shrink-0">
        {children}
      </div>
    </div>
  )
}

function Subsection({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div>
      <h3 className="font-semibold text-ink mb-2">{title}</h3>
      <div className="space-y-3">{children}</div>
    </div>
  )
}
