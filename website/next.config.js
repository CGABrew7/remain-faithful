/** @type {import('next').NextConfig} */

function contentSecurityPolicy() {
  const gaEnabled = Boolean(process.env.NEXT_PUBLIC_GA_ID)
  const scriptSrc = ["'self'", "'unsafe-inline'"]
  const imgSrc = ["'self'", 'data:']
  const connectSrc = ["'self'"]

  if (gaEnabled) {
    scriptSrc.push(
      'https://www.googletagmanager.com',
      'https://www.google-analytics.com',
    )
    imgSrc.push(
      'https://www.google-analytics.com',
      'https://www.googletagmanager.com',
    )
    connectSrc.push(
      'https://www.google-analytics.com',
      'https://analytics.google.com',
      'https://*.google-analytics.com',
      'https://*.analytics.google.com',
      'https://www.googletagmanager.com',
    )
  }

  return [
    "default-src 'self'",
    `script-src ${scriptSrc.join(' ')}`,
    "style-src 'self' 'unsafe-inline'",
    `img-src ${imgSrc.join(' ')}`,
    "font-src 'self'",
    `connect-src ${connectSrc.join(' ')}`,
    "object-src 'none'",
    "base-uri 'self'",
    "form-action 'self'",
    "frame-ancestors 'none'",
    "frame-src 'none'",
    'upgrade-insecure-requests',
  ].join('; ')
}

const securityHeaders = [
  { key: 'X-Frame-Options', value: 'DENY' },
  { key: 'X-Content-Type-Options', value: 'nosniff' },
  { key: 'Content-Security-Policy', value: contentSecurityPolicy() },
]

const nextConfig = {
  reactStrictMode: true,
  images: {
    remotePatterns: [],
  },
  async headers() {
    return [
      { source: '/', headers: securityHeaders },
      { source: '/:path*', headers: securityHeaders },
    ]
  },
}

module.exports = nextConfig
