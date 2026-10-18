import { MetadataRoute } from 'next'

export default function robots(): MetadataRoute.Robots {
  return {
    rules: [
      {
        userAgent: '*',
        allow: '/',
        disallow: ['/api/', '/_next/', '/admin', '/admin/'],
      },
    ],
    sitemap: 'https://remainfaithful.com/sitemap.xml',
  }
}
