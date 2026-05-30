export async function POST(request: Request) {
  try {
    const body = await request.json()

    if (!body.email || !body.name) {
      return Response.json({ error: 'Missing required fields' }, { status: 400 })
    }

    const backendUrl = process.env.BACKEND_URL
    if (!backendUrl) {
      return Response.json({ error: 'Server configuration error' }, { status: 500 })
    }

    const res = await fetch(`${backendUrl}/contact`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })

    if (!res.ok) {
      console.error('Backend contact error:', res.status)
      return Response.json({ error: 'Failed to submit contact form' }, { status: 502 })
    }

    return Response.json({ success: true })
  } catch (err) {
    console.error('Contact route error:', err)
    return Response.json({ error: 'Internal server error' }, { status: 500 })
  }
}
