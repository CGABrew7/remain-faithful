export async function POST(request: Request) {
  try {
    const body = await request.json()

    if (!body.email || !body.name || !body.message) {
      return Response.json({ error: 'Missing required fields' }, { status: 400 })
    }

    const backendUrl = process.env.BACKEND_URL
    if (!backendUrl) {
      console.log('[contact] BACKEND_URL not set — logging submission:', {
        name: body.name,
        email: body.email,
        subject: body.subject,
      })
      return Response.json({ success: true })
    }

    const res = await fetch(`${backendUrl}/contact`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
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
