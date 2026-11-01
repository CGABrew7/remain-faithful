import { NextResponse } from 'next/server'
import { backendJSON } from '@/lib/backend'
import { guardPublicPost } from '@/lib/form-guard'
import { plausibleEmail } from '@/lib/waitlist'

const SAVE_FAILED = 'We could not save your signup. Please try again.'
const RATE_LIMITED = 'Too many tries. Wait a minute and try again.'

export async function POST(req: Request) {
  const guarded = await guardPublicPost(req, {
    actions: ['waitlist'],
    turnstile: 'forward',
    rateMessage: RATE_LIMITED,
    email: (body) => (typeof body.email === 'string' ? body.email : ''),
    drop: () => NextResponse.json({ success: true }),
  })
  if ('response' in guarded) return guarded.response

  const { email, name } = guarded.body
  if (!plausibleEmail(email)) {
    return NextResponse.json({ error: 'Please enter a valid email address.' }, { status: 400 })
  }

  let result: Awaited<ReturnType<typeof backendJSON>>
  try {
    result = await backendJSON('/waitlist', {
      email: (email as string).trim(),
      name: typeof name === 'string' ? name : '',
      platform: 'android',
      turnstileToken: guarded.turnstileToken,
      formAction: guarded.formAction,
      loadedAt: guarded.loadedAt,
    }, undefined, req)
  } catch (err) {
    console.error('[waitlist] backend unreachable:', err)
    return NextResponse.json({ error: SAVE_FAILED }, { status: 502 })
  }

  if (!result.ok) {
    console.error('[waitlist] backend error:', result.status)
    if (result.status === 400) {
      return NextResponse.json({ error: 'Please enter a valid email address.' }, { status: 400 })
    }
    if (result.status === 429) {
      return NextResponse.json({ error: RATE_LIMITED }, { status: 429 })
    }
    return NextResponse.json({ error: SAVE_FAILED }, { status: 502 })
  }
  return NextResponse.json({ success: true })
}
