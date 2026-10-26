import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import { pilotContact } from '../src/lib/pilot.ts'
import { plausibleEmail } from '../src/lib/waitlist.ts'

describe('church pilot contact', () => {
  it('builds a subject and readable message from the pilot fields', () => {
    const result = pilotContact({
      name: 'Pat Pastor',
      church: 'Grace Fellowship',
      email: 'pat@example.org',
      role: 'Pastor',
      numGroups: '1 group',
      referral: '',
    })
    assert.equal(result.ok, true)
    if (!result.ok) return
    assert.equal(result.contact.subject, 'Church pilot: Grace Fellowship')
    assert.match(result.contact.message, /Church or organization: Grace Fellowship/)
    assert.match(result.contact.message, /Role: Pastor/)
    assert.match(result.contact.message, /Estimated groups: 1 group/)
    assert.match(result.contact.message, /Heard about us: Not given/)
    assert.doesNotMatch(result.contact.message, /[\u2013\u2014]/)
  })

  it('still requires the pilot fields', () => {
    assert.equal(pilotContact({ name: 'Pat', email: 'pat@example.org' }).ok, false)
    assert.equal(
      pilotContact({ name: 'Pat', church: 'Grace', email: 'nope', role: 'Pastor', numGroups: '1 group' }).ok,
      false,
    )
    assert.equal(
      pilotContact({ name: '   ', church: 'Grace', email: 'pat@example.org', role: 'Pastor', numGroups: '1 group' }).ok,
      false,
    )
  })

  it('caps long fields and strips line breaks from the subject', () => {
    const result = pilotContact({
      name: 'Pat',
      church: 'A\r\nBcc: x@example.org' + 'C'.repeat(5000),
      email: 'pat@example.org',
      role: 'Pastor',
      numGroups: '1 group',
    })
    assert.equal(result.ok, true)
    if (!result.ok) return
    assert.doesNotMatch(result.contact.subject, /[\r\n]/)
    assert.ok(result.contact.subject.length <= 'Church pilot: '.length + 200)
  })
})

describe('waitlist email check', () => {
  it('accepts normal addresses', () => {
    assert.equal(plausibleEmail('ada@example.com'), true)
    assert.equal(plausibleEmail(' ada+android@mail.example.co.uk '), true)
  })
  it('rejects junk', () => {
    for (const bad of ['', '@', 'bob@x', 'not-an-email', 'a b@example.com', 'a@@b.com', 'x'.repeat(250) + '@example.com', 42, null]) {
      assert.equal(plausibleEmail(bad), false, String(bad))
    }
  })
})
