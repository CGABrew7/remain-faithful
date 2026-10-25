import { backendJSON, clientAddress, tooMany } from '@/lib/backend'

export async function POST(request: Request) {
  if (tooMany(clientAddress(request))) {
    return Response.json({ error: 'Too many attempts. Wait a minute and try again.' }, { status: 429 })
  }
  const body = await request.json().catch(() => null)
  const email = typeof body?.email === 'string' ? body.email : ''
  const result = await backendJSON('/auth/resend-confirmation', { email }, undefined, request)
  if (!result.ok && result.status !== 200) {
    return Response.json(
      { error: result.data.error || 'We could not send that email. Try again in a little while.' },
      { status: result.status || 502 },
    )
  }
  return Response.json({ status: 'ok' })
}
