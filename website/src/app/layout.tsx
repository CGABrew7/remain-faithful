import type { Metadata } from 'next'
import { Newsreader, IBM_Plex_Mono } from 'next/font/google'
import './globals.css'
import Nav from '@/components/Nav'
import Footer from '@/components/Footer'
import Script from 'next/script'
import { JsonLd } from '@/components/JsonLd'
import { organizationSchema, websiteSchema } from '@/lib/structured-data'

const newsreader = Newsreader({
  subsets: ['latin'],
  weight: ['300', '400', '500', '600'],
  style: ['normal', 'italic'],
  variable: '--font-newsreader',
  display: 'swap',
  adjustFontFallback: false,
})

const plexMono = IBM_Plex_Mono({
  subsets: ['latin'],
  weight: ['400', '500'],
  variable: '--font-plex-mono',
  display: 'swap',
})

export const metadata: Metadata = {
  title: {
    default: 'Remain Faithful: Accountability That Works',
    template: '%s | Remain Faithful',
  },
  description:
    'Free peer accountability for Christians committed to purity. On-device AI, privacy-first, built on trust, not surveillance. Join the waitlist.',
  keywords: [
    'accountability app',
    'purity',
    'Christian accountability',
    'accountability partner',
    'purity app',
    'church accountability',
    'small group accountability',
    'screen monitoring',
    'privacy-first',
    'open source accountability',
    'free accountability app',
  ],
  authors: [{ name: 'Jeff Brewer' }],
  creator: 'Jeff Brewer',
  metadataBase: new URL(
    process.env.NEXT_PUBLIC_SITE_URL || 'https://remainfaithful.com'
  ),
  openGraph: {
    type: 'website',
    locale: 'en_US',
    url: '/',
    siteName: 'Remain Faithful',
    title: 'Remain Faithful: Accountability That Works',
    description:
      'Free peer accountability for Christians committed to purity. On-device AI, privacy-first, built on trust, not surveillance.',
    images: [
      {
        url: '/opengraph-image',
        width: 1200,
        height: 630,
        alt: 'Remain Faithful: Accountability That Works',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Remain Faithful: Accountability That Works',
    description:
      'Free peer accountability for Christians committed to purity. On-device AI, privacy-first.',
    images: ['/opengraph-image'],
  },
  robots: {
    index: true,
    follow: true,
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const gaId = process.env.NEXT_PUBLIC_GA_ID

  return (
    <html lang="en" className={`scroll-smooth ${newsreader.variable} ${plexMono.variable}`}>
      <head />
      <body className="bg-paper text-ink font-serif">
        {gaId && (
          <>
            <Script
              src={`https://www.googletagmanager.com/gtag/js?id=${gaId}`}
              strategy="afterInteractive"
            />
            <Script id="google-analytics" strategy="afterInteractive">
              {`
                window.dataLayer = window.dataLayer || [];
                function gtag(){dataLayer.push(arguments);}
                gtag('js', new Date());
                gtag('config', '${gaId}');
              `}
            </Script>
          </>
        )}
        <JsonLd data={organizationSchema} />
        <JsonLd data={websiteSchema} />
        <Nav />
        <main className="min-h-screen">{children}</main>
        <Footer />
      </body>
    </html>
  )
}
