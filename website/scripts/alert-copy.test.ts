import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import { ALERT_TRIGGER } from '../src/lib/product-copy.ts'

describe('alert trigger copy', () => {
  it('states alerts and blocking in positive terms', () => {
    assert.equal(ALERT_TRIGGER.includes('\u2014'), false)
    assert.equal(ALERT_TRIGGER.includes('\u2013'), false)
    assert.equal(ALERT_TRIGGER.includes('does not'), false)
    assert.match(ALERT_TRIGGER, /adult web filter turns on/)
    assert.match(ALERT_TRIGGER, /blocked app was opened/)
    assert.match(ALERT_TRIGGER, /1 minute, then again at 5, 15, and 30 minutes/)
  })
})
