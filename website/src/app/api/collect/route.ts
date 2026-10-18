import { backendBase } from '@/lib/backend'
import {
  analyticsSalt,
  cleanUTM,
  dailyVisitorHash,
  deviceType,
  referrerDomain,
  trackablePath,
  utcDay,
} from '@/lib/visitor'

export async function POST(request: Request) {
  const dnt = request.headers.get('dnt')
  const gpc = request.headers.get('sec-gpc')
  if (dnt === '1' || gpc === '1') {
    return new Response(null, { status: 204 })
  }

  let body: Record<string, unknown> = {}
  try {
    const parsed = await request.json()
    if (parsed && typeof parsed === 'object') body = parsed as Record<string, unknown>
  } catch {
    return new Response(null, { status: 204 })
  }

  const event = body.event
  if (event !== 'pageview' && event !== 'app_store_click' && event !== 'android_signup') {
    return new Response(null, { status: 204 })
  }
  const path = trackablePath(typeof body.path === 'string' ? body.path : '')
  if (!path) return new Response(null, { status: 204 })

  const countryHeader = (request.headers.get('x-vercel-ip-country') || '').toUpperCase()
  const regionHeader = (request.headers.get('x-vercel-ip-country-region') || '').toUpperCase()
  const country = /^[A-Z]{2}$/.test(countryHeader) ? countryHeader : ''
  const region = country === 'US' && /^[A-Z]{2}$/.test(regionHeader) ? regionHeader : ''
  const ua = request.headers.get('user-agent') || ''

  const payload = {
    event,
    path,
    referrer: referrerDomain(typeof body.referrer === 'string' ? body.referrer : ''),
    utm_source: cleanUTM(body.utm_source),
    utm_medium: cleanUTM(body.utm_medium),
    utm_campaign: cleanUTM(body.utm_campaign),
    country,
    region,
    device: deviceType(ua),
    visitor_hash: dailyVisitorHash({
      salt: analyticsSalt(),
      day: utcDay(),
      userAgent: ua,
      language: request.headers.get('accept-language') || '',
    }),
  }

  const base = backendBase()
  if (!base) return new Response(null, { status: 204 })

  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (dnt) headers.DNT = dnt
  if (gpc) headers['Sec-GPC'] = gpc

  await fetch(`${base}/analytics/collect`, {
    method: 'POST',
    headers,
    body: JSON.stringify(payload),
    cache: 'no-store',
  }).catch(() => undefined)

  return new Response(null, { status: 204 })
}
