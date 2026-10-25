import { cookies } from 'next/headers'
import { backendBase, proxyHeaders, SERVER_UNREACHABLE } from '@/lib/backend'

const COOKIE = 'rf_admin'
const EIGHT_HOURS = 60 * 60 * 8

function cookieOptions(maxAge: number) {
  return {
    httpOnly: true,
    sameSite: 'strict' as const,
    secure: process.env.NODE_ENV === 'production',
    path: '/',
    maxAge,
  }
}

/** Forwards an owner-dashboard call to the API. The session cookie stays on this site. */
export async function proxyAdmin(request: Request, path: string) {
  const base = backendBase()
  if (!base) {
    return Response.json({ error: 'The account service is not configured.' }, { status: 503 })
  }

  const jar = await cookies()
  const token = jar.get(COOKIE)?.value
  const headers = new Headers()
  headers.set('Content-Type', 'application/json')
  if (token) headers.set('Authorization', `Bearer ${token}`)
  for (const [key, value] of Object.entries(proxyHeaders(request))) {
    headers.set(key, value)
  }

  const url = new URL(request.url)
  const init: RequestInit = {
    method: request.method,
    headers,
    cache: 'no-store',
  }
  if (request.method !== 'GET' && request.method !== 'HEAD') {
    init.body = await request.text()
  }

  let res: Response
  try {
    res = await fetch(`${base}/admin/${path}${url.search}`, { ...init, signal: AbortSignal.timeout(10_000) })
  } catch {
    return Response.json({ error: SERVER_UNREACHABLE }, { status: 504 })
  }
  const text = await res.text()

  if (path === 'login') {
    if (!res.ok) {
      return new Response(text || JSON.stringify({ error: 'Sign-in failed' }), {
        status: res.status,
        headers: { 'Content-Type': 'application/json', 'Cache-Control': 'no-store' },
      })
    }
    let data: { token?: string } = {}
    try {
      data = JSON.parse(text) as { token?: string }
    } catch {
      data = {}
    }
    if (!data.token) {
      return Response.json({ error: 'Sign-in failed' }, { status: 502 })
    }
    jar.set(COOKIE, data.token, cookieOptions(EIGHT_HOURS))
    return Response.json({ ok: true })
  }

  if (path === 'logout') {
    jar.set(COOKIE, '', cookieOptions(0))
  }

  return new Response(text, {
    status: res.status,
    headers: {
      'Content-Type': res.headers.get('content-type') || 'application/json',
      'Cache-Control': 'no-store',
    },
  })
}
