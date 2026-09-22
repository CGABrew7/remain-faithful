import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import {
  DONATE_AMOUNTS,
  WOODFIELD_DONATE_LINKS,
  woodfieldDonateUrl,
} from '../src/lib/donate-links.ts'

const EXPECTED = {
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
} as const

describe('Woodfield Stripe donate payment links', () => {
  it('maps each amount and interval to the live donate.stripe.com URL', () => {
    for (const interval of ['one-time', 'monthly'] as const) {
      for (const amount of DONATE_AMOUNTS) {
        assert.equal(woodfieldDonateUrl(amount, interval), EXPECTED[interval][amount])
        assert.equal(WOODFIELD_DONATE_LINKS[interval][amount], EXPECTED[interval][amount])
      }
    }
  })

  it('uses eight distinct https payment links', () => {
    const urls = (['one-time', 'monthly'] as const).flatMap((interval) =>
      DONATE_AMOUNTS.map((amount) => woodfieldDonateUrl(amount, interval)),
    )
    assert.equal(urls.length, 8)
    assert.equal(new Set(urls).size, 8)
    for (const url of urls) {
      assert.match(url, /^https:\/\/donate\.stripe\.com\/[A-Za-z0-9]+$/)
    }
  })
})
