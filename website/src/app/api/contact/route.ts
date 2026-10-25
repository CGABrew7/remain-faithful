import { backendJSON } from '@/lib/backend'
import { pilotContact } from '@/lib/pilot'

const NOTE = 'Please add your name, email, and a short note.'

function filled(value: unknown) {
  return typeof value === 'string' && value.trim() !== ''
}

export async function POST(request: Request) {
  let body: Record<string, unknown>
  try {
    const parsed = await request.json()
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
      return Response.json({ error: 'Invalid request' }, { status: 400 })
    }
    body = parsed as Record<string, unknown>
  } catch {
    return Response.json({ error: 'Invalid request' }, { status: 400 })
  }

  if (typeof body.company === 'string' && body.company.trim() !== '') {
    return Response.json({ success: true })
  }
  delete body.company

  let payload: Record<string, unknown> = body
  if (body.type === 'church_pilot') {
    const pilot = pilotContact(body)
    if (!pilot.ok) {
      return Response.json({ error: pilot.error }, { status: 400 })
    }
    payload = pilot.contact
  } else if (!filled(body.email) || !filled(body.name) || !filled(body.message)) {
    return Response.json({ error: NOTE }, { status: 400 })
  }

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
