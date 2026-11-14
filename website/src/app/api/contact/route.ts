import { backendJSON } from '@/lib/backend'
import { guardPublicPost } from '@/lib/form-guard'
import { pilotContact } from '@/lib/pilot'

const NOTE = 'Please add your name, email, and a short note.'

function filled(value: unknown) {
  return typeof value === 'string' && value.trim() !== ''
}

export async function POST(request: Request) {
  const guarded = await guardPublicPost(request, {
    actions: ['contact', 'church-pilot', 'feedback'],
    turnstile: 'forward',
    email: (body) => (typeof body.email === 'string' ? body.email : ''),
    drop: () => Response.json({ success: true }),
  })
  if ('response' in guarded) return guarded.response

  const body = guarded.body
  delete body.company

  let payload: Record<string, unknown>
  if (body.type === 'church_pilot') {
    const pilot = pilotContact(body)
    if (!pilot.ok) {
      return Response.json({ error: pilot.error }, { status: 400 })
    }
    payload = { ...pilot.contact }
  } else if (!filled(body.email) || !filled(body.name) || !filled(body.message)) {
    return Response.json({ error: NOTE }, { status: 400 })
  } else {
    payload = { ...body }
  }

  payload.turnstileToken = guarded.turnstileToken
  payload.formAction = guarded.formAction
  payload.loadedAt = guarded.loadedAt
  delete payload.company

  const result = await backendJSON('/contact', payload, undefined, request)
  if (!result.ok) {
    const status = result.status >= 400 && result.status < 500 ? result.status : result.status === 504 ? 504 : 502
    return Response.json(
      { error: result.data.error || 'We could not send your message. Please try again.' },
      { status },
    )
  }

  return Response.json({ success: true })
}
