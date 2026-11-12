import { MetadataRoute } from 'next'
import { SITE_ORIGIN } from '@/lib/site'

export default function robots(): MetadataRoute.Robots {
  return {
    rules: [
      {
        userAgent: '*',
        allow: '/',
        disallow: ['/api/', '/_next/', '/admin', '/admin/'],
      },
    ],
    sitemap: `${SITE_ORIGIN}/sitemap.xml`,
  }
}
