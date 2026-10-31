import assert from 'node:assert/strict'
import { createRequire } from 'node:module'
import { describe, it } from 'node:test'
import { isPublicClarityPath } from '../src/lib/clarity-paths.ts'

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
