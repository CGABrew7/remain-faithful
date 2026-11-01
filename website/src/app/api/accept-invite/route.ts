import { backendJSON } from '@/lib/backend'
import { guardPublicPost } from '@/lib/form-guard'

const RATE_LIMITED = 'Too many attempts. Wait a minute and try again.'

export async function POST(request: Request) {
  const guarded = await guardPublicPost(request, {
    actions: ['accept-invite'],
    turnstile: 'verify',
    rateMessage: RATE_LIMITED,
    email: (body) => (typeof body.email === 'string' ? body.email : ''),
    drop: () => Response.json({ status: 'accepted' }),
  })
  if ('response' in guarded) return guarded.response

  const body = guarded.body
  if (!body.email || !body.password || !body.token) {
    return Response.json({ error: 'Email, password, and the invite link are required.' }, { status: 400 })
  }
  const login = await backendJSON('/auth/login', {
    email: String(body.email),
    password: String(body.password),
  }, undefined, request)
  if (!login.ok || !login.data.token) {
    return Response.json(
      { error: login.data.error || 'That email or password did not match.' },
      { status: login.status || 401 },
    )
  }
  const group = body.type === 'group'
  const decline = body.decision === 'decline'
  const path = group
    ? decline
      ? '/groups/decline-invite'
      : '/groups/accept-invite'
    : decline
      ? '/relationships/decline-invite'
      : '/relationships/accept-invite'
  const result = await backendJSON(path, { token: body.token }, login.data.token, request)
  if (!result.ok) {
    return Response.json(
      { error: result.data.error || 'Could not update this invite.' },
      { status: result.status },
    )
  }
  return Response.json({ status: result.data.status || (decline ? 'declined' : 'accepted') })
}
