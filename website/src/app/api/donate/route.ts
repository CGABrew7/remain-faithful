import Stripe from 'stripe'
import { checkoutReturnUrls } from '@/lib/checkout-urls'
import { clientAddress, tooMany } from '@/lib/backend'

export async function POST(request: Request) {
  try {
    if (tooMany(clientAddress(request))) {
      return Response.json({ error: 'Too many attempts. Wait a minute and try again.' }, { status: 429 })
    }

    const body = await request.json()
    if (typeof body?.company === 'string' && body.company.trim() !== '') {
      return Response.json({ url: null })
    }
    const { amount, recurring } = body

    if (!amount || typeof amount !== 'number' || amount < 1 || amount > 10000) {
      return Response.json({ error: 'Invalid amount' }, { status: 400 })
    }

    const secretKey = process.env.STRIPE_SECRET_KEY
    if (!secretKey) {
      console.error('STRIPE_SECRET_KEY is not configured')
      return Response.json({ error: 'Payment processing is not configured' }, { status: 500 })
    }

    const stripe = new Stripe(secretKey)
    const { successUrl, cancelUrl } = checkoutReturnUrls(request.headers.get('origin'))

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
                'Supports free accountability tools for men committed to purity.',
            },
            unit_amount: amountCents,
            ...(recurring ? { recurring: { interval: 'month' } } : {}),
          },
          quantity: 1,
        },
      ],
      mode: recurring ? 'subscription' : 'payment',
      success_url: successUrl,
      cancel_url: cancelUrl,
    })

    return Response.json({ url: session.url })
  } catch (err) {
    console.error('Donate route error:', err)
    return Response.json({ error: 'Internal server error' }, { status: 500 })
  }
}
