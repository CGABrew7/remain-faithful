import Stripe from 'stripe'

export async function POST(request: Request) {
  try {
    const { amount, recurring } = await request.json()

    if (!amount || typeof amount !== 'number' || amount < 1 || amount > 10000) {
      return Response.json({ error: 'Invalid amount' }, { status: 400 })
    }

    const secretKey = process.env.STRIPE_SECRET_KEY
    if (!secretKey) {
      console.error('STRIPE_SECRET_KEY is not configured')
      return Response.json({ error: 'Payment processing is not configured' }, { status: 500 })
    }

    const stripe = new Stripe(secretKey)
    const origin =
      request.headers.get('origin') ||
      process.env.NEXT_PUBLIC_SITE_URL ||
      'https://remainfaithful.com'

    const amountCents = Math.round(amount * 100)

    const session = await stripe.checkout.sessions.create({
      payment_method_types: ['card'],
      line_items: [
        {
          price_data: {
            currency: 'usd',
            product_data: {
              name: 'Donation to Remain Faithful',
              description:
                'Supports free accountability tools for men committed to sexual purity.',
            },
            unit_amount: amountCents,
            ...(recurring ? { recurring: { interval: 'month' } } : {}),
          },
          quantity: 1,
        },
      ],
      mode: recurring ? 'subscription' : 'payment',
      success_url: `${origin}/?donated=true`,
      cancel_url: `${origin}/`,
    })

    return Response.json({ url: session.url })
  } catch (err) {
    console.error('Donate route error:', err)
    return Response.json({ error: 'Internal server error' }, { status: 500 })
  }
}
