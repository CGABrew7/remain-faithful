import type { Metadata } from 'next'
import { FAQ_STREAMING } from '@/lib/product-copy'

export const metadata: Metadata = {
  title: 'Privacy Policy',
  description: 'How Remain Faithful collects, uses, and protects your data. Screenshots and screen text stay on your device. Partners can receive a short summary with the category.',
  alternates: { canonical: 'https://www.remainfaithful.com/privacy' },
}

export default function PrivacyPage() {
  return (
    <div className="pt-32 pb-24">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-12">
          <p className="text-wax text-sm font-semibold uppercase tracking-widest mb-3">Legal</p>
          <h1 className="font-serif text-4xl font-bold text-ink mb-4">Privacy Policy</h1>
          <p className="text-ink-soft text-sm">Last updated: November 2026</p>
        </div>

        <div className="prose-content space-y-12">
          <Section title="Our Commitment">
            <p className="font-semibold text-ink leading-relaxed">
              We do not sell your name or email. People you accept as partners, and members of a group you join, can see the name and email on your account, plus alert metadata. Apple or Google may receive a device token so alerts can be delivered. Website analytics run only if you allow them.
            </p>
            <p>
              Remain Faithful is an accountability app for adults. This policy describes what the iPhone app, the Android app, the desktop app, and this website collect today, what stays on your device, and what a partner can see.
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
                <li>Your email address (used for sign-in and email)</li>
                <li>A hashed password, when you sign in with email (we never store plaintext passwords)</li>
              </ul>
              <p>
                You can also sign in with Apple or Google. We store the identifier those services give us, plus the name and email they share. We do not receive your Apple or Google password.
              </p>
            </Subsection>

            <Subsection title="Device Information">
              <p>
                The iPhone app registers an Apple Push Notification device token so alerts can be delivered. The Android app registers a push token for the same purpose. The desktop app sends a device id so a protection notice can name the computer, not so we can see the screen. These tokens do not reveal anything beyond the account they are attached to.
              </p>
            </Subsection>

            <Subsection title="Alert Metadata">
              <p>
                If Alerts is on, your partner is told when a chosen app or category is used for 1 minute, then again at 5, 15, and 30 minutes that day. If Blocking is on, Apple&apos;s adult web filter turns on. If you open an app you asked to block, your partner gets an alert that a blocked app was opened, and Apple&apos;s Screen Time shield appears. That alert does not include the app name. The web filter runs on the phone and does not send browsing history. Optional Deep Scan can also store an alert when it flags a frame. {FAQ_STREAMING} The same kind of alert metadata and heartbeat can leave the Android app and the desktop app. That alert metadata is:
              </p>
              <ul>
                <li>Timestamp of the alert</li>
                <li>Content category (for example, &quot;Adult Content&quot; or &quot;Explicit Text&quot;)</li>
                <li>Severity level (Low / Medium / High)</li>
                <li>A short system-generated summary string</li>
              </ul>
              <p>
                <strong className="text-ink">We do not collect, store, or transmit:</strong> screenshots, screen recordings, raw OCR text, browsing history, or raw screen content. Partner sharing does not include those items beyond the short summary above.
              </p>
              <p>When the app loads a streak, it sends the current UTC offset in minutes so the days match the phone. That number is used for the response and is not saved.</p>
            </Subsection>

            <Subsection title="Website Analytics">
              <p>
                Website analytics stay off until you tap Allow analytics. The choice is saved in this browser. If you allow them, the site may load Google Analytics and Microsoft Clarity. Microsoft Clarity sets first-party cookies (_clck and _clsk) only after you allow analytics. The iPhone, Android, and desktop apps do not include this website analytics.
              </p>
              <p>
                The site also keeps a cookieless count of pages, referrer domain, campaign tags, country, region, and device type for 13 months, with a daily anonymous visitor id and no IP address.
              </p>
            </Subsection>

            <Subsection title="Forms">
              <p>The website and the apps collect what you type into these forms:</p>
              <ul>
                <li>Waitlist: the name and email you submit so we can write when you ask us to</li>
                <li>Pilot sign-up: name, email, organization, role, estimated number of groups, and how you heard about the app</li>
                <li>Contact form: name, email, subject, and message</li>
                <li>Ideas and problem reports: the type you pick, the words you type, and the name and email on your account so we can reply. App version, device model, and operating system version are included only if you turn that on. The form does not send logs, alerts, screenshots, Screen Time choices, partner names, or group names.</li>
              </ul>
            </Subsection>

            <Subsection title="Payments">
              <p>
                We do not store payment card numbers. If a payment is processed on the website, Stripe handles the card. We may keep a confirmation of the amount on the account. That row is removed when the account is deleted. Stripe keeps its own payment records.
              </p>
            </Subsection>
          </Section>

          <Section title="What Stays On Your Device">
            <p>Always-on Family Controls filtering and DeviceActivity usage monitoring do not capture screen frames. If you start an optional Deep Scan session on iPhone, the following data is processed on your device and is not transmitted as raw screen content:</p>
            <ul>
              <li>Raw screen frames captured by the ReplayKit broadcast extension (Deep Scan only)</li>
              <li>OCR text extracted by Apple Vision</li>
              <li>SensitiveContentAnalysis classification results</li>
              <li>Local keyword classifier scores</li>
              <li>Full browsing URLs</li>
              <li>Any visual content from your screen</li>
            </ul>
            <p>
              When Deep Scan flags a frame, the same alert metadata listed above can be uploaded. Screenshots, OCR text, and raw screen content never leave the device. {FAQ_STREAMING} The iOS app, Android app, desktop app, backend, and website source are public on GitHub.
            </p>
          </Section>

          <Section title="What Your Partners See">
            <p>Accountability partners, and members of a group you belong to, can receive:</p>
            <ul>
              <li>Your display name and account email</li>
              <li>Alert metadata as described above: timestamp, category, severity level, and a short system-generated summary. A blocked-app alert does not include which app, and it does not include a bundle ID.</li>
              <li>Your streak count (aggregate, not individual events)</li>
            </ul>
            <p>They do not have access to your screen content, browsing history, or any data beyond alert metadata and the account details listed above.</p>
          </Section>

          <Section title="Data Retention">
            <p>We retain data for the following periods:</p>
            <ul>
              <li><strong className="text-ink">Account data:</strong> deleted when you delete the account</li>
              <li><strong className="text-ink">Alert history:</strong> kept while your account exists, and deleted with the account. Older events and alerts are also removed after 365 days. Leaving a group notifies the other members. It does not delete your alert history.</li>
              <li><strong className="text-ink">Partnership records:</strong> removed when either person ends the relationship</li>
              <li><strong className="text-ink">Website page counts:</strong> removed after 13 months</li>
              <li><strong className="text-ink">Payment confirmation rows on your account:</strong> removed with the account. Stripe keeps its own payment records.</li>
            </ul>
          </Section>

          <Section title="Who Helps Us Run This">
            <ul>
              <li>
                <strong className="text-ink">Apple:</strong> Sign in with Apple, push delivery on iPhone, and Screen Time on the device. Apple&apos;s privacy policy applies to their handling of device tokens and sign-in.
              </li>
              <li>
                <strong className="text-ink">Google:</strong> Sign in with Google, and push delivery for the Android app. Google Analytics runs on the website only after you allow analytics.
              </li>
              <li>
                <strong className="text-ink">Microsoft Clarity:</strong> Website session recordings and heatmaps, only after you allow analytics.
              </li>
              <li>
                <strong className="text-ink">Resend:</strong> Sends account email, including invites, address confirmation, password reset, and the messages from the contact, pilot, waitlist, and feedback forms.
              </li>
              <li>
                <strong className="text-ink">Stripe:</strong> Processes a website payment when one is made. We do not store card numbers.
              </li>
              <li>
                <strong className="text-ink">Fly:</strong> Hosts the API.
              </li>
              <li>
                <strong className="text-ink">Vercel:</strong> Hosts this website.
              </li>
            </ul>
          </Section>

          <Section title="Your Rights">
            <p>You may, at any time:</p>
            <ul>
              <li><strong className="text-ink">Access your data:</strong> Settings, then View My Activity Log, shows your flagged activity in the app. For anything else, email support@remainfaithful.com. There is no separate export button.</li>
              <li><strong className="text-ink">Correct your data:</strong> update your display name and email in Settings, then Edit Profile</li>
              <li><strong className="text-ink">Delete your data:</strong> delete your account in Settings, then Delete Account. The account, alerts, and payment confirmation rows tied to the account are removed then.</li>
              <li><strong className="text-ink">Change protection:</strong> you can change App Restrictions and turn off App Lockout in Settings. Pausing or turning off monitoring notifies your partners. Turning off App Lockout also sends a notice.</li>
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
              Remain Faithful is not intended for anyone under 18. We do not knowingly collect data from anyone under 18. If you believe we have collected data from a minor, contact us at support@remainfaithful.com and we will delete it.
            </p>
          </Section>

          <Section title="Security">
            <p>
              The app and the server talk over HTTPS. Passwords are stored as bcrypt hashes. The session token is kept in the iOS Keychain and in the app group so Screen Time extensions can send events. Server logs do not include reset links or the address an email was sent to.
            </p>
          </Section>

          <Section title="Changes to This Policy">
            <p>
              We will post any material changes to this policy on this page with a new &quot;Last updated&quot; date. For significant changes, we will notify active users by push notification or email. Your continued use of the app after a policy change constitutes acceptance of the updated policy.
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
      <div className="space-y-4 text-ink-soft leading-relaxed [&_ul]:space-y-2 [&_ul]:list-none [&_ul_li]:flex [&_ul_li]:gap-3 [&_ul_li]:before:content-['•'] [&_ul_li]:before:text-wax [&_ul_li]:before:flex-shrink-0">
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
