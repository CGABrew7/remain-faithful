import { backendJSON } from '@/lib/backend'
import { guardPublicPost } from '@/lib/form-guard'

const RATE_LIMITED = 'Too many attempts. Wait a minute and try again.'

export async function POST(request: Request) {
  const guarded = await guardPublicPost(request, {
    actions: ['verify-email'],
    turnstile: 'verify',
    rateMessage: RATE_LIMITED,
    drop: () => Response.json({ status: 'verified' }),
  })
  if ('response' in guarded) return guarded.response

  const body = guarded.body
  if (!body.token) {
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
