import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { describe, it } from 'node:test'
import { ALERT_TRIGGER } from '../src/lib/product-copy.ts'

describe('alert trigger copy', () => {
  it('matches the paste-ready sentence in docs/alert-trigger-copy.md', () => {
    const doc = readFileSync(new URL('../../docs/alert-trigger-copy.md', import.meta.url), 'utf8')
    assert.ok(doc.includes(ALERT_TRIGGER))
    assert.equal(ALERT_TRIGGER.includes('\u2014'), false)
    assert.equal(ALERT_TRIGGER.includes('\u2013'), false)
  })
})
