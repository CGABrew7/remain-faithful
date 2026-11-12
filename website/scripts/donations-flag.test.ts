import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { describe, it } from 'node:test'
import { donationsEnabled } from '../src/lib/donations.ts'

describe('donation visibility', () => {
  it('stays off unless the flag is exactly true', () => {
    const previous = process.env.NEXT_PUBLIC_DONATIONS_ENABLED
    delete process.env.NEXT_PUBLIC_DONATIONS_ENABLED
    assert.equal(donationsEnabled(), false)
    process.env.NEXT_PUBLIC_DONATIONS_ENABLED = 'false'
    assert.equal(donationsEnabled(), false)
    process.env.NEXT_PUBLIC_DONATIONS_ENABLED = 'true'
    assert.equal(donationsEnabled(), true)
    if (previous === undefined) delete process.env.NEXT_PUBLIC_DONATIONS_ENABLED
    else process.env.NEXT_PUBLIC_DONATIONS_ENABLED = previous
  })

  it('gates the public donation surfaces', () => {
    const files = [
      'src/app/page.tsx',
      'src/app/about/page.tsx',
      'src/app/terms/page.tsx',
      'src/app/partners/page.tsx',
      'src/app/blog/[slug]/page.tsx',
      'src/app/free-accountability-app/page.tsx',
      'src/components/Footer.tsx',
      'src/components/DonateButton.tsx',
      'src/components/FooterDonateLink.tsx',
      'src/components/DonationSuccessBanner.tsx',
    ]
    for (const file of files) {
      const text = readFileSync(new URL('../' + file, import.meta.url), 'utf8')
      assert.equal(text.includes('donationsEnabled'), true, file)
    }
  })
})
