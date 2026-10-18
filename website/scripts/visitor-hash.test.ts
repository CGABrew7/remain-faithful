import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import {
  dailyVisitorHash,
  deviceType,
  primaryLanguage,
  referrerDomain,
  trackablePath,
} from '../src/lib/visitor.ts'

const sample = {
  salt: 'test-salt',
  day: '2026-10-18',
  userAgent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15',
  language: 'en-US,en;q=0.9',
}

describe('daily visitor hash', () => {
  it('stays the same for the same day and browser', () => {
    assert.equal(dailyVisitorHash(sample), dailyVisitorHash({ ...sample }))
  })

  it('changes on the next UTC day', () => {
    assert.notEqual(dailyVisitorHash(sample), dailyVisitorHash({ ...sample, day: '2026-10-19' }))
  })

  it('ignores an IP address if one is passed alongside the allowed inputs', () => {
    const withIp = { ...sample, ip: '203.0.113.10' }
    assert.equal(dailyVisitorHash(sample), dailyVisitorHash(withIp))
  })

  it('uses only the first 180 characters of the user agent', () => {
    const long = 'A'.repeat(180)
    assert.equal(
      dailyVisitorHash({ ...sample, userAgent: long }),
      dailyVisitorHash({ ...sample, userAgent: `${long}extra` }),
    )
  })
})

describe('trackable paths', () => {
  it('skips admin pages and drops the query string', () => {
    assert.equal(trackablePath('/admin'), null)
    assert.equal(trackablePath('/admin/preview'), null)
    assert.equal(trackablePath('/how-it-works?utm_source=newsletter'), '/how-it-works')
    assert.equal(trackablePath('/'), '/')
  })
})

describe('referrer domains', () => {
  it('keeps the host and drops same-site visits', () => {
    assert.equal(referrerDomain('https://www.google.com/search?q=remain'), 'google.com')
    assert.equal(referrerDomain('https://www.remainfaithful.com/about'), '')
    assert.equal(referrerDomain('https://remainfaithful.app/'), '')
    assert.equal(referrerDomain(''), '')
  })
})

describe('device and language', () => {
  it('reads phone user agents as mobile', () => {
    assert.equal(deviceType(sample.userAgent), 'mobile')
    assert.equal(deviceType('Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'), 'desktop')
  })

  it('keeps the first language tag', () => {
    assert.equal(primaryLanguage('en-US,en;q=0.9'), 'en-us')
  })
})
