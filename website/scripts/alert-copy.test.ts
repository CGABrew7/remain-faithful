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
    assert.equal(ALERT_TRIGGER.includes('opened an app you asked to block'), true)
    assert.equal(ALERT_TRIGGER.includes('The app name is never sent.'), true)
    assert.equal(ALERT_TRIGGER.includes('does not itself send an alert'), false)
  })

  it('does not say opening a blocked app sends no alert', () => {
    const files = [
      new URL('../src/lib/product-copy.ts', import.meta.url),
      new URL('../src/lib/structured-data.ts', import.meta.url),
      new URL('../src/app/app-guide/guide-content.ts', import.meta.url),
      new URL('../../docs/alert-trigger-copy.md', import.meta.url),
    ]
    for (const file of files) {
      const text = readFileSync(file, 'utf8')
      assert.equal(text.includes('does not itself send an alert'), false, String(file))
      assert.equal(text.includes('no alert from the open itself'), false, String(file))
      assert.equal(text.includes('Opening a blocked app does not create a partner alert'), false, String(file))
    }
  })
})
