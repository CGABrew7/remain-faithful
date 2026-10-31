import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { createRequire } from 'node:module'
import { describe, it } from 'node:test'
import { readAnalyticsConsent } from '../src/lib/analytics-consent.ts'
import { isPublicClarityPath, shouldLoadClarity } from '../src/lib/clarity-paths.ts'

const require = createRequire(import.meta.url)

describe('Clarity public paths', () => {
  it('allows marketing pages', () => {
    for (const path of ['/', '/privacy', '/about', '/how-it-works', '/blog', '/blog/a-post', '/partners', '/feedback']) {
      assert.equal(isPublicClarityPath(path), true, path)
    }
  })

  it('skips admin, account, and invite routes', () => {
    for (const path of [
      '/admin',
      '/admin/preview',
      '/reset-password',
      '/verify-email',
      '/invite',
      '/accept-invite',
      '/join',
      '/join/abc',
      '/api/donate',
    ]) {
      assert.equal(isPublicClarityPath(path), false, path)
    }
  })
})

describe('Clarity consent gate', () => {
  it('loads only after Allow analytics on a public page', () => {
    assert.equal(shouldLoadClarity('/', 'allow'), true)
    assert.equal(shouldLoadClarity('/privacy', 'allow'), true)
    assert.equal(shouldLoadClarity('/', 'deny'), false)
    assert.equal(shouldLoadClarity('/', null), false)
    assert.equal(shouldLoadClarity('/', undefined), false)
    assert.equal(shouldLoadClarity('/', ''), false)
    assert.equal(shouldLoadClarity('/admin', 'allow'), false)
    assert.equal(shouldLoadClarity('/join/abc', 'allow'), false)
  })

  it('reads the stored Allow analytics choice and ignores anything else', () => {
    const storage = {
      value: null as string | null,
      getItem() {
        return this.value
      },
    }
    assert.equal(readAnalyticsConsent(storage), null)
    storage.value = 'allow'
    assert.equal(readAnalyticsConsent(storage), 'allow')
    storage.value = 'deny'
    assert.equal(readAnalyticsConsent(storage), 'deny')
    storage.value = 'maybe'
    assert.equal(readAnalyticsConsent(storage), null)
    assert.equal(readAnalyticsConsent(null), null)
  })

  it('keeps the consent sentence and says Clarity cookies wait for that choice', () => {
    const privacy = readFileSync(new URL('../src/app/privacy/page.tsx', import.meta.url), 'utf8')
    const clarity = readFileSync(new URL('../src/components/MicrosoftClarity.tsx', import.meta.url), 'utf8')
    assert.match(privacy, /Website analytics stay off until you tap Allow analytics\./)
    assert.match(
      privacy,
      /Microsoft Clarity sets first-party cookies \(_clck and _clsk\) only after you allow analytics\./,
    )
    assert.equal(privacy.includes('Microsoft Clarity is used for anonymous usage analytics'), false)
    assert.equal(privacy.includes('\u2014') || privacy.includes('\u2013'), false)
    assert.match(clarity, /shouldLoadClarity/)
    assert.match(clarity, /if \(!enabled\) return null/)
    assert.equal(clarity.includes('document.cookie'), false)
  })
})

describe('Clarity content security policy', () => {
  it('allows the Clarity tag and collector without removing analytics hosts', async () => {
    process.env.NEXT_PUBLIC_GA_ID = 'G-TEST'
    const config = require('../next.config.js')
    const groups = await config.headers()
    const csp = groups[0].headers.find((header: { key: string }) => header.key === 'Content-Security-Policy').value as string
    assert.match(csp, /script-src[^;]*https:\/\/\*\.clarity\.ms/)
    assert.match(csp, /script-src[^;]*https:\/\/c\.bing\.com/)
    assert.match(csp, /connect-src[^;]*https:\/\/\*\.clarity\.ms/)
    assert.match(csp, /img-src[^;]*https:\/\/\*\.clarity\.ms/)
    assert.match(csp, /https:\/\/www\.googletagmanager\.com/)
    assert.match(csp, /https:\/\/www\.google-analytics\.com/)
  })
})
