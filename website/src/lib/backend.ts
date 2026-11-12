export function backendBase() {
  return (process.env.BACKEND_URL || '').replace(/\/$/, '')
}

export const SERVER_UNREACHABLE = 'We could not reach the server. Check your connection and try again.'

export function proxyHeaders(request?: Request): Record<string, string> {
  const secret = process.env.SITE_PROXY_SECRET || ''
  if (!secret || !request) return {}
  return {
    'X-RF-Proxy-Secret': secret,
    'X-RF-Client-IP': clientAddress(request),
  }
}

export async function backendJSON(
  path: string,
  body: unknown,
  token?: string,
  request?: Request,
): Promise<{ ok: boolean; status: number; data: { error?: string; token?: string; status?: string } }> {
  const base = backendBase()
  if (!base) {
    return { ok: false, status: 503, data: { error: 'Account service is not configured.' } }
  }
  try {
    const res = await fetch(`${base}${path}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...proxyHeaders(request),
      },
      body: JSON.stringify(body),
      cache: 'no-store',
      signal: AbortSignal.timeout(10_000),
    })
    let data: { error?: string; token?: string; status?: string } = {}
    try {
      data = await res.json()
    } catch {
      data = {}
    }
    return { ok: res.ok, status: res.status, data }
  } catch {
    return { ok: false, status: 504, data: { error: SERVER_UNREACHABLE } }
  }
}

const hits = new Map<string, { n: number; start: number }>()

export function tooMany(ip: string, limit = 10) {
  const now = Date.now()
  const window = hits.get(ip)
  if (!window || now - window.start > 60_000) {
    hits.set(ip, { n: 1, start: now })
    return false
  }
  window.n += 1
  return window.n > limit
}

function platformAddress(raw: string | null): string {
  const hop = raw?.split(',')[0]?.trim() ?? ''
  if (!hop || hop.length > 64 || /\s/.test(hop)) return ''
  return hop
}

export function clientAddress(request: Request) {
  // x-vercel-forwarded-for is the visitor address from the platform.
  // x-real-ip is used when that header is absent. The API trusts this value
  // only when the site sends the shared proxy secret.
  const platform = platformAddress(request.headers.get('x-vercel-forwarded-for'))
  if (platform) return platform
  const real = platformAddress(request.headers.get('x-real-ip'))
  if (real) return real
  return 'local'
}
