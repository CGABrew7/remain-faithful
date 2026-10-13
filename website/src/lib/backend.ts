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
  const forwarded = request.headers.get('x-forwarded-for')
  if (forwarded) return forwarded.split(',')[0].trim()
  return 'local'
}
