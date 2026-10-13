import { backendJSON, clientAddress, tooMany } from '@/lib/backend'

export async function POST(request: Request) {
  if (tooMany(clientAddress(request))) {
    return Response.json({ error: 'Too many attempts. Wait a minute and try again.' }, { status: 429 })
  }
  const body = await request.json().catch(() => null)
  if (!body?.token || !body?.password) {
    return Response.json({ error: 'A reset link and a new password are required.' }, { status: 400 })
  }
  const result = await backendJSON('/auth/reset-password', {
    token: body.token,
    password: body.password,
  })
  if (!result.ok) {
    return Response.json(
      { error: result.data.error || 'That reset link is not valid.' },
      { status: result.status },
    )
  }
  return Response.json({ status: 'ok' })
}
