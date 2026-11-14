import assert from 'node:assert/strict'
import { afterEach, describe, it } from 'node:test'
import { guardPublicPost, resetTurnstileWarningForTests, verifyTurnstile, FORM_RETRY_MESSAGE, FORM_WAIT_MESSAGE } from '../src/lib/form-guard.ts'
import { submittedTooFast } from '../src/lib/public-form.ts'

const DUMMY_SECRET = '1x0000000000000000000000000000000AA'
const DUMMY_SITE_KEY = '1x00000000000000000000AA'

const originalWarn = console.warn
const originalEnv = {
  secret: process.env.TURNSTILE_SECRET_KEY,
  siteKey: process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY,
}

afterEach(() => {
  console.warn = originalWarn
  if (originalEnv.secret === undefined) delete process.env.TURNSTILE_SECRET_KEY
  else process.env.TURNSTILE_SECRET_KEY = originalEnv.secret
  if (originalEnv.siteKey === undefined) delete process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY
  else process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY = originalEnv.siteKey
})

function post(body: unknown, ip: string) {
  return new Request('https://www.remainfaithful.com/api/contact', {
    method: 'POST',
    headers: {
      'content-type': 'application/json',
      'x-real-ip': ip,
    },
    body: JSON.stringify(body),
  })
}

function readyBody(extra: Record<string, unknown> = {}) {
  return {
    loadedAt: Date.now() - 5000,
    formAction: 'contact',
    email: 'ada@example.com',
    ...extra,
  }
}

describe('submit timing', () => {
  const now = Date.parse('2026-11-01T12:00:00Z')
  it('rejects a missing, fast, or stale timestamp', () => {
    assert.equal(submittedTooFast(undefined, now), true)
    assert.equal(submittedTooFast(now, now), true)
    assert.equal(submittedTooFast(now - 1000, now), true)
    assert.equal(submittedTooFast(Math.floor(now / 1000), now), true)
    assert.equal(submittedTooFast(now - 25 * 60 * 60 * 1000, now), true)
    assert.equal(submittedTooFast(now + 60_000, now), true)
  })
  it('accepts a timestamp at least two seconds old', () => {
    assert.equal(submittedTooFast(now - 2000, now), false)
    assert.equal(submittedTooFast(now - 10_000, now), false)
  })
})

describe('public form guard', () => {
  it('returns a success response when the honeypot is filled', async () => {
    const result = await guardPublicPost(post({ company: ' Acme ', loadedAt: Date.now() }, '203.0.113.20'), {
      actions: ['contact'],
      turnstile: 'forward',
      drop: () => Response.json({ success: true }),
      fetchImpl: async () => {
        throw new Error('honeypot should not verify')
      },
    })
    assert.ok('response' in result)
    assert.equal(result.response.status, 200)
    assert.deepEqual(await result.response.json(), { success: true })
  })

  it('rejects a fast submit', async () => {
    delete process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY
    const result = await guardPublicPost(post({ ...readyBody(), loadedAt: Date.now() }, '203.0.113.21'), {
      actions: ['contact'],
      turnstile: 'forward',
    })
    assert.ok('response' in result)
    assert.equal(result.response.status, 400)
    assert.equal((await result.response.json()).error, FORM_WAIT_MESSAGE)
  })

  it('limits repeats from one address and one email', async () => {
    delete process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY
    delete process.env.TURNSTILE_SECRET_KEY
    const ip = '203.0.113.30'
    for (let i = 0; i < 10; i++) {
      const result = await guardPublicPost(post(readyBody({ email: `ip-${i}@example.com` }), ip), {
        actions: ['contact'],
        turnstile: 'forward',
        email: (body) => String(body.email),
      })
      assert.equal('response' in result, false, `request ${i + 1}`)
    }
    const blocked = await guardPublicPost(post(readyBody({ email: 'ip-new@example.com' }), ip), {
      actions: ['contact'],
      turnstile: 'forward',
      email: (body) => String(body.email),
    })
    assert.ok('response' in blocked)
    assert.equal(blocked.response.status, 429)

    const email = 'same-person@example.com'
    for (let i = 0; i < 10; i++) {
      const result = await guardPublicPost(post(readyBody({ email }), `203.0.113.${40 + i}`), {
        actions: ['contact'],
        turnstile: 'forward',
        email: (body) => String(body.email),
      })
      assert.equal('response' in result, false, `email request ${i + 1}`)
    }
    const emailBlocked = await guardPublicPost(post(readyBody({ email: 'Same-Person@example.com' }), '203.0.113.60'), {
      actions: ['contact'],
      turnstile: 'forward',
      email: (body) => String(body.email),
    })
    assert.ok('response' in emailBlocked)
    assert.equal(emailBlocked.response.status, 429)
  })

  it('skips verification and warns once when the secret is unset', async () => {
    resetTurnstileWarningForTests()
    delete process.env.TURNSTILE_SECRET_KEY
    delete process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY
    const warnings: string[] = []
    console.warn = (...args: unknown[]) => {
      warnings.push(args.map(String).join(' '))
    }
    assert.equal(await verifyTurnstile({ token: '', remoteip: '', action: '' }), 'skipped')
    assert.equal(await verifyTurnstile({ token: 'token', remoteip: '', action: 'contact' }), 'skipped')
    assert.equal(warnings.length, 1)
    assert.match(warnings[0], /TURNSTILE_SECRET_KEY is not set/)
  })

  it('verifies a token with the dummy test secret', async () => {
    process.env.TURNSTILE_SECRET_KEY = DUMMY_SECRET
    process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY = DUMMY_SITE_KEY
    let sawSecret = false
    const result = await guardPublicPost(
      post(readyBody({ formAction: 'forgot-password', turnstileToken: 'token-ok' }), '203.0.113.70'),
      {
        actions: ['forgot-password'],
        turnstile: 'verify',
        email: (body) => String(body.email),
        fetchImpl: async (_url, init) => {
          const raw = init?.body
          const params = raw instanceof URLSearchParams ? raw : new URLSearchParams(String(raw))
          sawSecret = params.get('secret') === DUMMY_SECRET && params.get('response') === 'token-ok'
          assert.equal(params.get('remoteip'), '203.0.113.70')
          return Response.json({ success: true, action: 'forgot-password' })
        },
      },
    )
    assert.equal(sawSecret, true)
    assert.equal('response' in result, false)
  })

  it('rejects a failed, mismatched, or unreachable verification', async () => {
    process.env.TURNSTILE_SECRET_KEY = DUMMY_SECRET
    process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY = DUMMY_SITE_KEY
    const cases: Array<{ name: string; fetchImpl: typeof fetch }> = [
      {
        name: 'failure',
        fetchImpl: async () => Response.json({ success: false, 'error-codes': ['invalid-input-response'] }),
      },
      {
        name: 'mismatch',
        fetchImpl: async () => Response.json({ success: true, action: 'contact' }),
      },
      {
        name: 'network',
        fetchImpl: async () => {
          throw new Error('offline')
        },
      },
    ]
    for (const [index, item] of cases.entries()) {
      const result = await guardPublicPost(
        post(readyBody({ formAction: 'forgot-password', turnstileToken: 'token-ok', email: `bad-${index}@example.com` }), `198.51.100.${index + 1}`),
        {
          actions: ['forgot-password'],
          turnstile: 'verify',
          email: (body) => String(body.email),
          fetchImpl: item.fetchImpl,
        },
      )
      assert.ok('response' in result, item.name)
      assert.equal(result.response.status, 400, item.name)
      assert.equal((await result.response.json()).error, FORM_RETRY_MESSAGE, item.name)
    }
  })

  it('requires a token when the site key is set and does not spend it on forwarded forms', async () => {
    process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY = DUMMY_SITE_KEY
    process.env.TURNSTILE_SECRET_KEY = DUMMY_SECRET
    const missing = await guardPublicPost(post(readyBody(), '203.0.113.80'), {
      actions: ['contact'],
      turnstile: 'forward',
      fetchImpl: async () => {
        throw new Error('forward should not call siteverify')
      },
    })
    assert.ok('response' in missing)
    assert.equal(missing.response.status, 400)

    const forwarded = await guardPublicPost(
      post(readyBody({ turnstileToken: 'token-forward' }), '203.0.113.81'),
      {
        actions: ['contact'],
        turnstile: 'forward',
        fetchImpl: async () => {
          throw new Error('forward should not call siteverify')
        },
      },
    )
    assert.equal('response' in forwarded, false)
    if ('response' in forwarded) return
    assert.equal(forwarded.turnstileToken, 'token-forward')
    assert.equal(forwarded.formAction, 'contact')
  })
})
