import { pilotContact } from '@/lib/pilot'

export async function POST(request: Request) {
  try {
    const body = await request.json()

    // Honeypot field from the feedback form. People never see it, so a value means a bot.
    if (typeof body.company === 'string' && body.company.trim() !== '') {
      return Response.json({ success: true })
    }
    delete body.company

    let payload = body
    if (body.type === 'church_pilot') {
      const pilot = pilotContact(body)
      if (!pilot.ok) {
        return Response.json({ error: pilot.error }, { status: 400 })
      }
      payload = pilot.contact
    } else if (!body.email || !body.name || !body.message) {
      return Response.json({ error: 'Missing required fields' }, { status: 400 })
    }

    const backendUrl = process.env.BACKEND_URL
    if (!backendUrl) {
      console.log('[contact] BACKEND_URL not set — logging submission:', {
        name: payload.name,
        email: payload.email,
        subject: payload.subject,
      })
      return Response.json({ success: true })
    }

    const res = await fetch(`${backendUrl}/contact`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    if (!res.ok) {
      let errorMsg = `Backend returned ${res.status}`
      try {
        const errBody = await res.json()
        if (errBody?.error) errorMsg = errBody.error
      } catch {}
      console.error('[contact] backend error:', res.status, errorMsg)
      return Response.json({ error: errorMsg }, { status: 502 })
    }

    return Response.json({ success: true })
  } catch (err) {
    console.error('[contact] route error:', err)
    return Response.json({ error: 'Internal server error' }, { status: 500 })
  }
}
