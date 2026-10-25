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

export function clientAddress(request: Request) {
  // Vercel sets x-real-ip to the visitor it observed. When that header is
  // missing, the first hop of x-forwarded-for is the same address on Vercel.
  // x-vercel-forwarded-for is the fallback when another proxy rewrites the
  // forwarded list. The API only trusts this value when the site sends the
  // shared proxy secret.
  const real = request.headers.get('x-real-ip')?.split(',')[0]?.trim()
  if (real) return real
  const forwarded = request.headers.get('x-forwarded-for')?.split(',')[0]?.trim()
  if (forwarded) return forwarded
  const platform = request.headers.get('x-vercel-forwarded-for')?.split(',')[0]?.trim()
  if (platform) return platform
  return 'local'
}
