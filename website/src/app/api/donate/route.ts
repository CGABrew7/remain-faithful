export async function POST(request: Request) {
  try {
    const { amount } = await request.json()

    if (!amount || typeof amount !== 'number' || amount < 1 || amount > 10000) {
      return Response.json({ error: 'Invalid amount' }, { status: 400 })
    }

    const backendUrl = process.env.BACKEND_URL
    if (!backendUrl) {
      return Response.json({ error: 'Server configuration error' }, { status: 500 })
    }

    const res = await fetch(`${backendUrl}/donations/create-checkout-session`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ amount_cents: Math.round(amount * 100) }),
    })

    if (!res.ok) {
      const text = await res.text()
      console.error('Backend donation error:', res.status, text)
      return Response.json({ error: 'Failed to create checkout session' }, { status: 502 })
    }

    const data = await res.json()
    return Response.json(data)
  } catch (err) {
    console.error('Donate route error:', err)
    return Response.json({ error: 'Internal server error' }, { status: 500 })
  }
}
