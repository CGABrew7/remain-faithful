import { createHash, createHmac, randomBytes } from 'node:crypto'

const processSalt = randomBytes(32).toString('hex')

/** Master salt. A daily HMAC of this value rotates the visitor hash each UTC day. */
export function analyticsSalt(): string {
  return process.env.ANALYTICS_SALT || processSalt
}

export function utcDay(now = new Date()): string {
  return now.toISOString().slice(0, 10)
}

/** First language tag only. The rest of Accept-Language is discarded. */
export function primaryLanguage(header: string): string {
  const tag = (header.split(',')[0] || '').trim().toLowerCase()
  const match = /^[a-z]{2,8}(?:-[a-z0-9]{1,8})?/.exec(tag)
  return match ? match[0].slice(0, 16) : ''
}

/**
 * Daily anonymous visitor id. Inputs are a secret salt, the UTC date,
 * a truncated user agent, and the primary language. The IP address is not an input.
 */
export function dailyVisitorHash(input: {
  salt: string
  day: string
  userAgent: string
  language: string
}): string {
  const daily = createHmac('sha256', input.salt).update(input.day).digest('hex')
  const ua = input.userAgent.slice(0, 180)
  const lang = primaryLanguage(input.language)
  return createHash('sha256').update(`${daily}|${ua}|${lang}`).digest('hex')
}

export function deviceType(ua: string): 'mobile' | 'desktop' {
  if (/Mobile|Android|iPhone|iPod|IEMobile|BlackBerry/i.test(ua)) return 'mobile'
  return 'desktop'
}

/** Pathname only. Admin pages, query strings, and odd characters are dropped. */
export function trackablePath(path: string): string | null {
  const clean = path.split('?')[0].split('#')[0]
  if (!clean.startsWith('/') || clean.length > 120) return null
  if (clean === '/admin' || clean.startsWith('/admin/')) return null
  if (clean !== '/' && !/^\/(?:[A-Za-z0-9_-]+\/?)*$/.test(clean)) return null
  return clean
}

/** Hostname only. Same-site referrers count as direct. */
function stripCampaign(value: string, max: number) {
  return value.replace(/[\u0000-\u001f\u007f\u200e\u200f\u202a-\u202e\u2066-\u2069]/g, '').trim().slice(0, max)
}

export function referrerDomain(raw: string): string {
  const trimmed = stripCampaign(raw, 300)
  if (!trimmed) return ''
  try {
    const withScheme = trimmed.includes('://') ? trimmed : `https://${trimmed}`
    const host = new URL(withScheme).hostname.replace(/^www\./, '').toLowerCase()
    if (host === 'remainfaithful.com' || host === 'remainfaithful.app') return ''
    return host.slice(0, 253)
  } catch {
    return ''
  }
}

export function cleanUTM(value: unknown): string {
  if (typeof value !== 'string') return ''
  return stripCampaign(value, 80)
}
