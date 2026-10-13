import { backendJSON, clientAddress, tooMany } from '@/lib/backend'

export async function POST(request: Request) {
  if (tooMany(clientAddress(request))) {
    return Response.json({ error: 'Too many attempts. Wait a minute and try again.' }, { status: 429 })
  }
  const body = await request.json().catch(() => null)
  if (typeof body?.company === 'string' && body.company.trim() !== '') {
    return Response.json({ status: 'accepted' })
  }
  if (!body?.email || !body?.password || !body?.token) {
    return Response.json({ error: 'Email, password, and the invite link are required.' }, { status: 400 })
  }
  const login = await backendJSON('/auth/login', {
    email: String(body.email),
    password: String(body.password),
  })
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
  const result = await backendJSON(path, { token: body.token }, login.data.token)
  if (!result.ok) {
    return Response.json(
      { error: result.data.error || 'Could not update this invite.' },
      { status: result.status },
    )
  }
  return Response.json({ status: result.data.status || (decline ? 'declined' : 'accepted') })
}
