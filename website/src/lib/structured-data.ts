import { COMPARE_UPDATED_ON } from '@/lib/competitor-pricing'
import {
  ALERT_TRIGGER,
  FAQ_DIFFERENT,
  FAQ_PARTNERS_SEE,
  FAQ_SLIP,
  FAQ_STREAMING,
  PRIVACY_POINTER,
} from '@/lib/product-copy'

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
  description: 'Free peer accountability app for adults committed to purity. A partner you choose stands with you. Free on iPhone. Screenshots and screen content stay on the device.',
  offers: {
    '@type': 'Offer',
    price: '0',
    priceCurrency: 'USD',
  },
  featureList: [
    'Alerts: one partner alert that day after a chosen app or category is used for a minute',
    'Blocking: Screen Time shield when a chosen app opens, with no alert from the open itself',
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
        text: FAQ_PARTNERS_SEE,
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
        text: FAQ_SLIP,
      },
    },
    {
      '@type': 'Question',
      name: 'How is this different from other accountability apps?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: FAQ_DIFFERENT,
      },
    },
    {
      '@type': 'Question',
      name: 'Can it see streaming video?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: FAQ_STREAMING,
      },
    },
    {
      '@type': 'Question',
      name: 'Does this work on Android?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Remain Faithful is available for iPhone (iOS 17+). Android support is planned, with no launch date. Leave your email in Get Android updates on the home page.',
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
        text: PRIVACY_POINTER,
      },
    },
    {
      '@type': 'Question',
      name: 'What exactly does Remain Faithful monitor?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: `${ALERT_TRIGGER} ${FAQ_STREAMING}`,
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
    description: `A detailed comparison of Remain Faithful and ${competitor} for Christian accountability.`,
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
