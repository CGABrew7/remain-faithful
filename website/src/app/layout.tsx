import type { Metadata } from 'next'
import { Playfair_Display, Inter } from 'next/font/google'
import './globals.css'
import Nav from '@/components/Nav'
import Footer from '@/components/Footer'
import Script from 'next/script'

const playfair = Playfair_Display({
  subsets: ['latin'],
  weight: ['400', '600', '700'],
  variable: '--font-playfair',
  display: 'swap',
})

const inter = Inter({
  subsets: ['latin'],
  weight: ['400', '500', '600'],
  variable: '--font-inter',
  display: 'swap',
})

export const metadata: Metadata = {
  title: {
    default: 'Remain Faithful: Accountability That Works',
    template: '%s | Remain Faithful',
  },
  description:
    'Free peer accountability for adults committed to purity. On-device AI, privacy-first, built on trust, not surveillance.',
  keywords: [
    'accountability app',
    'purity',
    'Christian men',
    'accountability partner',
    'purity app',
    'mens ministry',
    'screen monitoring',
    'privacy-first',
  ],
  authors: [{ name: 'Remain Faithful Team' }],
  creator: 'Remain Faithful',
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
      'Free peer accountability for adults committed to purity. On-device AI, privacy-first, built on trust, not surveillance.',
    images: [
      {
        url: '/og-image.png',
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
      'Free peer accountability for adults committed to purity. On-device AI, privacy-first.',
    images: ['/og-image.png'],
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
    <html lang="en" className={`scroll-smooth ${playfair.variable} ${inter.variable}`}>
      <head />
      <body className="bg-[#0F1B2D] text-[#F0EDE8] font-sans antialiased">
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
        <Nav />
        <main className="min-h-screen">{children}</main>
        <Footer />
      </body>
    </html>
  )
}
