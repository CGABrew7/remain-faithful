import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import { clientAddress } from '../src/lib/backend.ts'

function requestWith(headers: Record<string, string>) {
  return new Request('https://www.remainfaithful.com/api/donate', { headers })
}

describe('donate client address', () => {
  it('uses the platform client address instead of the first forwarded hop', () => {
    const req = requestWith({
      'x-forwarded-for': '203.0.113.5, 10.0.0.1',
      'x-real-ip': '198.51.100.20',
    })
    assert.equal(clientAddress(req), '198.51.100.20')
  })

  it('does not key on a caller-supplied forwarded hop when the platform header is absent', () => {
    const req = requestWith({ 'x-forwarded-for': '203.0.113.5' })
    assert.equal(clientAddress(req), 'local')
  })
})
