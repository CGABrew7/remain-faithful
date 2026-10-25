import { backendJSON, clientAddress, tooMany } from '@/lib/backend'

export async function POST(request: Request) {
  if (tooMany(clientAddress(request))) {
    return Response.json({ error: 'Too many attempts. Wait a minute and try again.' }, { status: 429 })
  }
  const body = await request.json().catch(() => null)
  if (!body?.token) {
    return Response.json({ error: 'A confirmation link is required.' }, { status: 400 })
  }
  const result = await backendJSON('/auth/verify-email', { token: body.token }, undefined, request)
  if (!result.ok) {
    return Response.json(
      { error: result.data.error || 'That confirmation link is not valid.' },
      { status: result.status },
    )
  }
  return Response.json({ status: result.data.status || 'verified' })
}
