import { COMPARE_UPDATED_ON } from '@/lib/competitor-pricing'

export const organizationSchema = {
  '@context': 'https://schema.org',
  '@type': 'Organization',
  name: 'Remain Faithful',
  url: 'https://remainfaithful.com',
  description: 'Free peer accountability for Christians committed to purity. Always-on filtering, optional Deep Scan, privacy-first, open source.',
  founder: {
    '@type': 'Person',
    name: 'Jeff Brewer',
    url: 'https://remainfaithful.com/about',
  },
  sameAs: [
    'https://github.com/CGABrew7/remain-faithful',
  ],
  nonprofitStatus: 'Nonprofit501c3',
}

export const softwareApplicationSchema = {
  '@context': 'https://schema.org',
  '@type': 'SoftwareApplication',
  name: 'Remain Faithful',
  operatingSystem: 'iOS 17+',
  applicationCategory: 'LifestyleApplication',
  description: 'Free peer accountability app for Christians committed to purity. Always-on Family Controls filtering with partner notify. Optional on-device Deep Scan. No screenshots ever leave your device. Partners may receive a short system-generated summary in addition to category, timestamp, and severity.',
  offers: {
    '@type': 'Offer',
    price: '0',
    priceCurrency: 'USD',
  },
  featureList: [
    'Always-on Family Controls filtering with partner notify',
    'DeviceActivity usage monitoring as category events',
    'Chosen apps stay blocked all day',
    'Optional on-device Deep Scan',
    'One-to-one or group accountability',
    'Covenant-based partner system',
    'Open source codebase',
    'Free forever, donation funded',
  ],
  author: {
    '@type': 'Organization',
    name: 'Remain Faithful',
  },
}

export const websiteSchema = {
  '@context': 'https://schema.org',
  '@type': 'WebSite',
  name: 'Remain Faithful',
  url: 'https://remainfaithful.com',
  description: 'Free peer accountability for Christians committed to purity.',
}

export const homepageFaqSchema = {
  '@context': 'https://schema.org',
  '@type': 'FAQPage',
  mainEntity: [
    {
      '@type': 'Question',
      name: 'Can my accountability partners see what I was looking at?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'No. Remain Faithful never captures or shares screenshots. Partners receive a discreet alert with a category (like "Adult Content"), severity, timestamp, and a short system-generated summary — never screenshots, OCR text, or raw screen content. This protects partners from being exposed to harmful material.',
      },
    },
    {
      '@type': 'Question',
      name: 'Can my spouse be my accountability partner?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes. You can set up a one-to-one partnership with your spouse, a friend, a mentor, or a pastor. You choose who sees your alerts.',
      },
    },
    {
      '@type': 'Question',
      name: 'Can our church use this for small groups?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Absolutely. Remain Faithful was designed for small group accountability. Groups of 3 to 12 members can all monitor and encourage each other. A free Group Setup Guide is available for ministry leaders.',
      },
    },
    {
      '@type': 'Question',
      name: 'What happens if I slip up?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Your accountability partners receive a discreet alert. The goal is conversation, not condemnation. Every alert includes conversation starter prompts to help your partners respond with grace.',
      },
    },
    {
      '@type': 'Question',
      name: 'How is this different from other accountability apps?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Remain Faithful is different in three key ways: the credible core is always-on Family Controls filtering with partner notify (Deep Scan is optional), it is 100% free forever, and screenshots, OCR text, and raw screen content never leave your device. Partners may receive a short system-generated summary in addition to category, severity, and timestamp. The iOS app, backend, and website source are public on GitHub. Design notes and audit markdown in that repo are historical working papers, not product promises.',
      },
    },
    {
      '@type': 'Question',
      name: 'Does this work on Android?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Remain Faithful is available for iPhone (iOS 17+). Android support is planned, with no launch date. Leave your email for Android notify and updates.',
      },
    },
    {
      '@type': 'Question',
      name: 'Is Remain Faithful really free?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes, always. Remain Faithful has no subscription tier, no premium features, and no advertising. The app is sustained by voluntary donations from users who find it valuable.',
      },
    },
  ],
}

export const howItWorksFaqSchema = {
  '@context': 'https://schema.org',
  '@type': 'FAQPage',
  mainEntity: [
    {
      '@type': 'Question',
      name: 'Is Remain Faithful really free?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Yes, always. Remain Faithful has no subscription tier, no premium features, and no advertising. The app is sustained by voluntary donations from users who find it valuable. This model is committed to indefinitely.',
      },
    },
    {
      '@type': 'Question',
      name: 'Who sees my data when using Remain Faithful?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Your chosen accountability partners can see alert metadata: the timestamp, the category (e.g., "adult content"), the severity level, and a short system-generated summary. They do not see which app, a bundle ID, screenshots, browsing history, or raw OCR text. Screenshots, OCR text, and raw screen content are never transmitted off your device.',
      },
    },
    {
      '@type': 'Question',
      name: 'What exactly does Remain Faithful monitor?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'You choose the apps and categories to block. They stay blocked all day, through the lock screen, reboot, and app restarts. If a chosen app is used for a minute, your partner can get one alert that day with the category and a short summary. Optional Deep Scan classifies non-DRM frames on the phone. Screenshots and screen text stay on the phone. Deep Scan cannot see DRM streaming.',
      },
    },
    {
      '@type': 'Question',
      name: 'Can I be anonymous on Remain Faithful?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'You choose your display name when you create your account. However, accountability by design requires that your partners know who they are holding accountable. Anonymity defeats the purpose. Your partners see the name you provide, typically your real name.',
      },
    },
  ],
}

export function blogPostSchema(post: {
  title: string
  description: string
  slug: string
  datePublished: string
  readTime: string
}) {
  return {
    '@context': 'https://schema.org',
    '@type': 'Article',
    headline: post.title,
    description: post.description,
    url: `https://remainfaithful.com/blog/${post.slug}`,
    datePublished: post.datePublished,
    dateModified: post.datePublished,
    author: {
      '@type': 'Person',
      name: 'Jeff Brewer',
      url: 'https://remainfaithful.com/about',
    },
    publisher: {
      '@type': 'Organization',
      name: 'Remain Faithful',
      url: 'https://remainfaithful.com',
    },
    mainEntityOfPage: {
      '@type': 'WebPage',
      '@id': `https://remainfaithful.com/blog/${post.slug}`,
    },
    timeRequired: post.readTime,
  }
}

export function comparisonSchema(competitor: string) {
  return {
    '@context': 'https://schema.org',
    '@type': 'WebPage',
    name: `Remain Faithful vs ${competitor}: Detailed Comparison`,
    description: `An honest, detailed comparison of Remain Faithful and ${competitor} for Christian accountability.`,
    url: `https://remainfaithful.com/compare/${competitor.toLowerCase().replace(/\s+/g, '-')}`,
    dateModified: COMPARE_UPDATED_ON,
    mainEntity: {
      '@type': 'ItemList',
      itemListElement: [
        {
          '@type': 'SoftwareApplication',
          name: 'Remain Faithful',
          applicationCategory: 'LifestyleApplication',
          operatingSystem: 'iOS 17+',
          offers: { '@type': 'Offer', price: '0', priceCurrency: 'USD' },
        },
        {
          '@type': 'SoftwareApplication',
          name: competitor,
          applicationCategory: 'LifestyleApplication',
        },
      ],
    },
  }
}

export function breadcrumbSchema(items: { name: string; url: string }[]) {
  return {
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    itemListElement: items.map((item, index) => ({
      '@type': 'ListItem',
      position: index + 1,
      name: item.name,
      item: item.url,
    })),
  }
}
