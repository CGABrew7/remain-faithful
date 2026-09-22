/** Live Woodfield Foundation Stripe Payment Links for the site donate section. */
export const DONATE_AMOUNTS = [5, 10, 25, 50] as const

export type DonateAmount = (typeof DONATE_AMOUNTS)[number]
export type DonateInterval = 'one-time' | 'monthly'

export const WOODFIELD_DONATE_LINKS: Record<
  DonateInterval,
  Record<DonateAmount, string>
> = {
  'one-time': {
    5: 'https://donate.stripe.com/28E7sN3HKfDa3O3fGJ08g03',
    10: 'https://donate.stripe.com/3cIcN7dik76EgAP0LP08g00',
    25: 'https://donate.stripe.com/cNi3cxdikez6acr0LP08g01',
    50: 'https://donate.stripe.com/4gM6oJ0vy4YwgAPeCF08g02',
  },
  monthly: {
    5: 'https://donate.stripe.com/fZu00lguwaiQ1FV0LP08g04',
    10: 'https://donate.stripe.com/28EeVf0vyez61FVfGJ08g05',
    25: 'https://donate.stripe.com/dRm14pdikfDa1FVamp08g07',
    50: 'https://donate.stripe.com/cNi14pbac1Mkbgv3Y108g06',
  },
}

export function woodfieldDonateUrl(amount: DonateAmount, interval: DonateInterval): string {
  return WOODFIELD_DONATE_LINKS[interval][amount]
}
