import type { Metadata } from 'next'
import Link from 'next/link'
import { JsonLd } from '@/components/JsonLd'
import { Breadcrumbs } from '@/components/Breadcrumbs'
import { UpdatedStamp } from '@/components/UpdatedStamp'
import { COMPARE_UPDATED_ON } from '@/lib/competitor-pricing'

export const metadata: Metadata = {
  title: 'Compare Remain Faithful with Other Accountability Apps',
  description: 'Side-by-side comparisons of Remain Faithful with Covenant Eyes, Ever Accountable, and Accountable2You. Prices and features rechecked on 2026-11-10.',
  alternates: { canonical: 'https://www.remainfaithful.com/compare' },
}

const indexSchema = {
  '@context': 'https://schema.org',
  '@type': 'WebPage',
  name: 'Compare Remain Faithful',
  description: 'Comparisons of Remain Faithful with Covenant Eyes, Ever Accountable, and Accountable2You.',
  url: 'https://www.remainfaithful.com/compare',
  dateModified: COMPARE_UPDATED_ON,
}

const pages = [
  {
    href: '/compare/covenant-eyes',
    title: 'Remain Faithful vs Covenant Eyes',
    body: 'Free on-device accountability on iPhone, compared with Victory by Covenant Eyes. Their Victory page lists $18 a month, $198 a year, or $950 once, and Screen Accountability sends screenshots to an ally.',
  },
  {
    href: '/compare/ever-accountable',
    title: 'Remain Faithful vs Ever Accountable',
    body: 'Their plans page lists $14.99 a month or $129 a year for one reporting user. Partners are described as receiving browsing history, and the device table includes Screenshot Accountability.',
  },
  {
    href: '/compare/accountable2you',
    title: 'Remain Faithful vs Accountable2You',
    body: 'Personal is $11 a month or $121 a year. Family is $16 a month or $176 a year. They report browsing history and page titles, and say they do not capture screenshots.',
  },
]

export default function CompareIndexPage() {
  return (
    <>
      <JsonLd data={indexSchema} />
      <div className="pt-24 pb-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <Breadcrumbs items={[{ name: 'Compare', url: 'https://www.remainfaithful.com/compare' }]} />
          <div className="mb-10">
            <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
              <p className="text-wax text-sm font-semibold uppercase tracking-widest">Direct Comparison</p>
              <UpdatedStamp />
            </div>
            <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-6 leading-tight text-balance">
              Compare Remain Faithful
            </h1>
            <p className="text-ink-soft text-lg leading-relaxed text-pretty">
              Remain Faithful is free on iPhone. Optional Deep Scan classifies on the phone, and screenshots stay on the device. The pages below compare that with three paid accountability products. Prices and features were rechecked on their public sites on {COMPARE_UPDATED_ON}.
            </p>
          </div>

          <div className="grid gap-4">
            {pages.map((page) => (
              <Link
                key={page.href}
                href={page.href}
                className="block p-5 rounded-sm border border-hairline bg-paper-deep hover:border-wax/40 transition-colors"
              >
                <h2 className="font-serif text-xl font-bold text-ink mb-2">{page.title}</h2>
                <p className="text-sm text-ink-soft leading-relaxed mb-3">{page.body}</p>
                <span className="text-wax font-medium text-sm">Read the comparison</span>
              </Link>
            ))}
          </div>
        </div>
      </div>
    </>
  )
}
