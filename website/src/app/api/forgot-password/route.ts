import { backendJSON } from '@/lib/backend'
import { guardPublicPost } from '@/lib/form-guard'

const RATE_LIMITED = 'Too many attempts. Wait a minute and try again.'

export async function POST(request: Request) {
  const guarded = await guardPublicPost(request, {
    actions: ['forgot-password'],
    turnstile: 'verify',
    rateMessage: RATE_LIMITED,
    email: (body) => (typeof body.email === 'string' ? body.email : ''),
    drop: () => Response.json({ status: 'ok' }),
  })
  if ('response' in guarded) return guarded.response

  const email = typeof guarded.body.email === 'string' ? guarded.body.email : ''
  const result = await backendJSON('/auth/forgot-password', { email }, undefined, request)
  if (!result.ok && result.status !== 200) {
    return Response.json(
      { error: result.data.error || 'We could not send that email. Try again in a little while.' },
      { status: result.status || 502 },
    )
  }
  return Response.json({ status: 'ok' })
}
