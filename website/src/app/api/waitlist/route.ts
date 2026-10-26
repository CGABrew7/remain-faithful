import { NextResponse } from 'next/server'
import { backendJSON, clientAddress, tooMany } from '@/lib/backend'
import { plausibleEmail } from '@/lib/waitlist'

const SAVE_FAILED = 'We could not save your signup. Please try again.'

export async function POST(req: Request) {
  if (tooMany(clientAddress(req))) {
    return NextResponse.json({ error: 'Too many tries. Wait a minute and try again.' }, { status: 429 })
  }
  const body = await req.json().catch(() => null)
  if (!body || typeof body !== 'object') {
    return NextResponse.json({ error: 'Invalid request' }, { status: 400 })
  }
  const { email, name, company } = body as { email?: unknown; name?: unknown; company?: unknown }
  if (!plausibleEmail(email)) {
    return NextResponse.json({ error: 'Please enter a valid email address.' }, { status: 400 })
  }

  let result: Awaited<ReturnType<typeof backendJSON>>
  try {
    result = await backendJSON('/waitlist', {
      email: email.trim(),
      name: typeof name === 'string' ? name : '',
      platform: 'android',
      company: typeof company === 'string' ? company : '',
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
      return NextResponse.json({ error: 'Too many tries. Wait a minute and try again.' }, { status: 429 })
    }
    return NextResponse.json({ error: SAVE_FAILED }, { status: 502 })
  }
  return NextResponse.json({ success: true })
}
