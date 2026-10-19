export function backendBase() {
  return (process.env.BACKEND_URL || '').replace(/\/$/, '')
}

export async function backendJSON(
  path: string,
  body: unknown,
  token?: string,
): Promise<{ ok: boolean; status: number; data: { error?: string; token?: string; status?: string } }> {
  const base = backendBase()
  if (!base) {
    return { ok: false, status: 503, data: { error: 'Account service is not configured.' } }
  }
  const res = await fetch(`${base}${path}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify(body),
    cache: 'no-store',
  })
  let data: { error?: string; token?: string; status?: string } = {}
  try {
    data = await res.json()
  } catch {
    data = {}
  }
  return { ok: res.ok, status: res.status, data }
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
  // Vercel sets x-real-ip to the client address it observed. The docs describe
  // it as identical to x-forwarded-for after Vercel replaces that header, and
  // x-vercel-forwarded-for as the same address when another proxy rewrites
  // x-forwarded-for. The first forwarded hop is not used.
  // https://vercel.com/docs/headers/request-headers
  const real = request.headers.get('x-real-ip')?.split(',')[0]?.trim()
  if (real) return real
  const platform = request.headers.get('x-vercel-forwarded-for')?.split(',')[0]?.trim()
  if (platform) return platform
  return 'local'
}
