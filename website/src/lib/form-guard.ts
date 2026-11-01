import { clientAddress, tooMany } from './backend.ts'
import { emailLimitKey, honeypotFilled, submittedTooFast } from './public-form.ts'

export const FORM_RATE_MESSAGE = 'Too many tries. Wait a minute and try again.'
export const FORM_WAIT_MESSAGE = 'Please wait a moment and try again.'
export const FORM_RETRY_MESSAGE = 'Please try again.'

const SITEVERIFY_URL = 'https://challenges.cloudflare.com/turnstile/v0/siteverify'

let turnstileWarningLogged = false

export function resetTurnstileWarningForTests() {
  turnstileWarningLogged = false
}

function warnTurnstileSkipped() {
  if (turnstileWarningLogged) return
  turnstileWarningLogged = true
  console.warn('[turnstile] TURNSTILE_SECRET_KEY is not set; skipping verification')
}

export type TurnstileVerdict = 'skipped' | 'verified' | 'rejected'

export async function verifyTurnstile(input: {
  token: string
  remoteip: string
  action: string
  endpoint?: string
  fetchImpl?: typeof fetch
}): Promise<TurnstileVerdict> {
  const secret = (process.env.TURNSTILE_SECRET_KEY || '').trim()
  if (!secret) {
    warnTurnstileSkipped()
    return 'skipped'
  }
  if (input.token.length === 0 || input.token.length > 2048 || input.action.length === 0) {
    return 'rejected'
  }
  const body = new URLSearchParams({
    secret,
    response: input.token,
  })
  const remoteip = input.remoteip.trim()
  if (remoteip && remoteip !== 'local') body.set('remoteip', remoteip)
  const doFetch = input.fetchImpl ?? fetch
  try {
    const response = await doFetch(input.endpoint || SITEVERIFY_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body,
      signal: AbortSignal.timeout(10_000),
    })
    if (!response.ok) return 'rejected'
    const result = (await response.json()) as { success?: unknown; action?: unknown }
    if (result.success !== true || result.action !== input.action) return 'rejected'
    return 'verified'
  } catch {
    console.error('[turnstile] siteverify failed')
    return 'rejected'
  }
}

export type PublicFormPass = {
  body: Record<string, unknown>
  turnstileToken: string
  formAction: string
  loadedAt: number
}

export type PublicFormStop = { response: Response }

type GuardOptions = {
  actions: readonly string[]
  turnstile: 'verify' | 'forward'
  email?: (body: Record<string, unknown>) => string
  rateMessage?: string
  drop?: () => Response
  fetchImpl?: typeof fetch
  now?: number
  endpoint?: string
}

function rateLimited(message: string) {
  return Response.json({ error: message }, { status: 429, headers: { 'Retry-After': '60' } })
}

// guardPublicPost applies the honeypot, rate limits, submit-time, and Turnstile
// checks for a public form route. Contact and waitlist forward the token so the
// API can verify it. Other routes verify here.
export async function guardPublicPost(request: Request, options: GuardOptions): Promise<PublicFormPass | PublicFormStop> {
  const overLimit = tooMany(clientAddress(request))
  let body: Record<string, unknown> | null = null
  try {
    const parsed = await request.json()
    if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
      body = parsed as Record<string, unknown>
    }
  } catch {
    body = null
  }

  if (body && honeypotFilled(body.company)) {
    return { response: options.drop ? options.drop() : Response.json({ success: true }) }
  }
  const rateMessage = options.rateMessage || FORM_RATE_MESSAGE
  if (overLimit) return { response: rateLimited(rateMessage) }
  if (!body) return { response: Response.json({ error: 'Invalid request' }, { status: 400 }) }

  const email = options.email ? options.email(body).trim() : ''
  if (email && tooMany(emailLimitKey(email))) {
    return { response: rateLimited(rateMessage) }
  }
  if (submittedTooFast(body.loadedAt, options.now)) {
    return { response: Response.json({ error: FORM_WAIT_MESSAGE }, { status: 400 }) }
  }

  const siteKey = (process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY || '').trim()
  const token = typeof body.turnstileToken === 'string' ? body.turnstileToken.trim() : ''
  const action = typeof body.formAction === 'string' ? body.formAction.trim() : ''
  if (siteKey && (token.length === 0 || token.length > 2048 || !options.actions.includes(action))) {
    return { response: Response.json({ error: FORM_RETRY_MESSAGE }, { status: 400 }) }
  }

  if (options.turnstile === 'verify') {
    const verdict = await verifyTurnstile({
      token,
      remoteip: clientAddress(request),
      action,
      endpoint: options.endpoint,
      fetchImpl: options.fetchImpl,
    })
    if (verdict === 'rejected') {
      return { response: Response.json({ error: FORM_RETRY_MESSAGE }, { status: 400 }) }
    }
  }

  return {
    body,
    turnstileToken: token,
    formAction: action,
    loadedAt: body.loadedAt as number,
  }
}
