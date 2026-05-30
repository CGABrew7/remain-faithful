import type { Metadata } from 'next'
import Link from 'next/link'

export const metadata: Metadata = {
  title: 'Terms of Service',
  description: 'Terms of Service for Remain Faithful, the free, privacy-first accountability app.',
}

export default function TermsPage() {
  return (
    <div className="pt-32 pb-24">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-12">
          <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">Legal</p>
          <h1 className="font-serif text-4xl font-bold text-[#F0EDE8] mb-4">Terms of Service</h1>
          <p className="text-[#8A9BB0] text-sm">Last updated: May 2025</p>
        </div>

        <div className="space-y-10 text-[#8A9BB0] leading-relaxed">
          <Section title="Acceptance of Terms">
            <p>
              By downloading, installing, or using Remain Faithful (the &quot;App&quot;) or this website, you agree to be bound by these Terms of Service. If you do not agree to these terms, do not use the App or website.
            </p>
          </Section>

          <Section title="Use of the App">
            <p>
              Remain Faithful is designed for adults (18+) who voluntarily enter into accountability relationships with trusted partners. You are responsible for your use of the App and for ensuring that your use complies with applicable laws.
            </p>
            <p>
              You may not use the App to monitor another person without their explicit consent, or in any way that violates applicable privacy laws.
            </p>
          </Section>

          <Section title="Accounts">
            <p>
              You are responsible for maintaining the security of your account credentials. Remain Faithful is not liable for any loss resulting from unauthorized access to your account.
            </p>
          </Section>

          <Section title="Donations">
            <p>
              Donations are voluntary and non-refundable. Remain Faithful is not a registered nonprofit organization; donations are not tax-deductible. Payment processing is handled by Stripe, Inc. We do not store payment card information.
            </p>
          </Section>

          <Section title="Disclaimer of Warranties">
            <p>
              The App is provided &quot;as is&quot; without warranty of any kind. We do not warrant that the App will be error-free, uninterrupted, or free of harmful components. Use the App at your own risk.
            </p>
          </Section>

          <Section title="Limitation of Liability">
            <p>
              To the maximum extent permitted by law, Remain Faithful and its contributors shall not be liable for any indirect, incidental, or consequential damages arising from your use of the App.
            </p>
          </Section>

          <Section title="Changes to These Terms">
            <p>
              We may update these Terms at any time. Material changes will be communicated via push notification or email. Continued use after changes constitutes acceptance.
            </p>
          </Section>

          <Section title="Contact">
            <p>
              Questions about these Terms may be directed to{' '}
              <a href="mailto:support@remainfaithful.com" className="text-[#C9A84C] hover:underline">
                support@remainfaithful.com
              </a>.
            </p>
          </Section>
        </div>

        <div className="mt-12 pt-8 border-t border-[#1E3050]">
          <Link href="/privacy" className="text-[#C9A84C] text-sm hover:underline">
            ← View Privacy Policy
          </Link>
        </div>
      </div>
    </div>
  )
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="border-t border-[#1E3050] pt-8">
      <h2 className="font-serif text-2xl font-bold text-[#F0EDE8] mb-4">{title}</h2>
      <div className="space-y-4">{children}</div>
    </div>
  )
}
