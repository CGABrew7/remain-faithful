import { backendJSON } from '@/lib/backend'
import { guardPublicPost } from '@/lib/form-guard'

const RATE_LIMITED = 'Too many attempts. Wait a minute and try again.'

export async function POST(request: Request) {
  const guarded = await guardPublicPost(request, {
    actions: ['reset-password'],
    turnstile: 'verify',
    rateMessage: RATE_LIMITED,
    drop: () => Response.json({ status: 'ok' }),
  })
  if ('response' in guarded) return guarded.response

  const body = guarded.body
  if (!body.token || !body.password) {
    return Response.json({ error: 'A reset link and a new password are required.' }, { status: 400 })
  }
  const result = await backendJSON('/auth/reset-password', {
    token: body.token,
    password: body.password,
  }, undefined, request)
  if (!result.ok) {
    return Response.json(
      { error: result.data.error || 'That reset link is not valid.' },
      { status: result.status },
    )
  }
  return Response.json({ status: 'ok' })
}
