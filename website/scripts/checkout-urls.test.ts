import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import {
  PRODUCTION_CANCEL_URL,
  PRODUCTION_SUCCESS_URL,
  checkoutOrigin,
  checkoutReturnUrls,
} from '../src/lib/checkout-urls.ts'

describe('checkout return URLs', () => {
  it('uses the exact production success and cancel URLs', () => {
    assert.equal(PRODUCTION_SUCCESS_URL, 'https://www.remainfaithful.com/?donated=true')
    assert.equal(PRODUCTION_CANCEL_URL, 'https://www.remainfaithful.com/#donate')
  })

  it('sends production traffic to www.remainfaithful.com', () => {
    for (const origin of [
      undefined,
      null,
      '',
      'https://www.remainfaithful.com',
      'https://remainfaithful.com',
      'https://remainfaithful.app',
      'https://remainfaithful.app/thank-you',
      'https://remain-faithful.vercel.app',
      'https://remain-faithful-git-main.vercel.app',
      'not-a-url',
    ]) {
      assert.equal(checkoutOrigin(origin), 'https://www.remainfaithful.com')
      assert.deepEqual(checkoutReturnUrls(origin), {
        successUrl: 'https://www.remainfaithful.com/?donated=true',
        cancelUrl: 'https://www.remainfaithful.com/#donate',
      })
    }
  })

  it('preserves localhost Origins for local Stripe test checkouts', () => {
    assert.deepEqual(checkoutReturnUrls('http://localhost:3000'), {
      successUrl: 'http://localhost:3000/?donated=true',
      cancelUrl: 'http://localhost:3000/#donate',
    })
    assert.deepEqual(checkoutReturnUrls('http://127.0.0.1:3000/'), {
      successUrl: 'http://127.0.0.1:3000/?donated=true',
      cancelUrl: 'http://127.0.0.1:3000/#donate',
    })
    assert.deepEqual(checkoutReturnUrls('http://[::1]:3000'), {
      successUrl: 'http://[::1]:3000/?donated=true',
      cancelUrl: 'http://[::1]:3000/#donate',
    })
  })

  it('does not treat lookalike hosts as local', () => {
    assert.equal(checkoutOrigin('https://localhost.evil.com'), 'https://www.remainfaithful.com')
    assert.equal(checkoutOrigin('https://evil-localhost.com'), 'https://www.remainfaithful.com')
    assert.equal(checkoutOrigin('http://192.168.1.10:3000'), 'https://www.remainfaithful.com')
  })
})
