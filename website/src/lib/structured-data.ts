export const organizationSchema = {
  '@context': 'https://schema.org',
  '@type': 'Organization',
  name: 'Remain Faithful',
  url: 'https://remainfaithful.com',
  logo: 'https://remainfaithful.com/logo.png',
  description: 'Free peer accountability for Christians committed to purity. On-device AI, privacy-first, open source.',
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
  description: 'Free peer accountability app for Christians committed to purity. On-device AI monitors screen content privately. No screenshots ever leave your device. Partners receive discreet alerts only.',
  offers: {
    '@type': 'Offer',
    price: '0',
    priceCurrency: 'USD',
  },
  featureList: [
    'On-device AI content classification',
    'Always-on Screen Time monitoring',
    'One-to-one or group accountability',
    'Covenant-based partner system',
    'Privacy-first architecture',
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
        text: 'No. Remain Faithful never captures or shares screenshots. Partners receive a discreet alert with a category (like "Adult Content") and severity level, never the actual content. This protects partners from being exposed to harmful material.',
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
        text: 'Remain Faithful is different in three key ways: it is 100% free forever, all content classification happens entirely on your device so no screen content ever leaves your phone, and the entire codebase is open source so anyone can verify exactly what is and is not transmitted.',
      },
    },
    {
      '@type': 'Question',
      name: 'Does this work on Android?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Android support is planned for Fall 2026. Currently Remain Faithful is available for iPhone (iOS 17+). Join the waitlist to be notified when Android launches.',
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
        text: 'Your chosen accountability partners can see alert metadata: the timestamp, the category (e.g., "adult content"), and the severity level. They do not see screenshots, browsing history, app content, or raw OCR text. None of that data is ever transmitted off your device.',
      },
    },
    {
      '@type': 'Question',
      name: 'What exactly does Remain Faithful monitor?',
      acceptedAnswer: {
        '@type': 'Answer',
        text: 'Remain Faithful uses two monitoring layers. Layer A is always-on Screen Time monitoring that watches which apps you open and which web categories you visit. It runs persistently in the background, requires no screen broadcast permission, and survives device restarts. Layer B is Deep Scan, started intentionally for high-risk periods, using Apple ReplayKit to run on-device AI (Vision OCR, SensitiveContentAnalysis) on screen frames. All classification is on-device.',
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
