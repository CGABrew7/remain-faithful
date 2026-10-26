import assert from 'node:assert/strict'
import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join, relative } from 'node:path'
import { describe, it } from 'node:test'
import { fileURLToPath } from 'node:url'

const srcRoot = fileURLToPath(new URL('../src/', import.meta.url))

const banned: { label: string; re: RegExp }[] = [
  { label: 'em dash', re: /\u2014/ },
  { label: 'en dash', re: /\u2013/ },
  { label: 'em dash escape', re: /\\u2014|&mdash;|&#8212;|&#x2014;/i },
  { label: 'en dash escape', re: /\\u2013|&ndash;|&#8211;|&#x2013;/i },
  { label: 'honesty wording', re: /honest/i },
]

/** Drop comments and keep newlines so reported lines match the file. */
export function sourceWithoutComments(source: string): string {
  let out = ''
  let i = 0
  while (i < source.length) {
    const c = source[i]
    const next = source[i + 1]
    if (c === '"' || c === "'" || c === '`') {
      out += c
      i++
      while (i < source.length) {
        const ch = source[i]
        out += ch
        i++
        if (ch === '\\' && i < source.length) {
          out += source[i]
          i++
          continue
        }
        if (ch === c) break
      }
      continue
    }
    if (c === '/' && next === '/') {
      while (i < source.length && source[i] !== '\n') i++
      continue
    }
    if (c === '/' && next === '*') {
      i += 2
      while (i < source.length && !(source[i] === '*' && source[i + 1] === '/')) {
        if (source[i] === '\n') out += '\n'
        i++
      }
      i = Math.min(source.length, i + 2)
      continue
    }
    out += c
    i++
  }
  return out
}

function filesUnder(dir: string): string[] {
  const out: string[] = []
  for (const name of readdirSync(dir)) {
    const path = join(dir, name)
    if (statSync(path).isDirectory()) out.push(...filesUnder(path))
    else out.push(path)
  }
  return out
}

function hitsIn(visible: string): string[] {
  const hits: string[] = []
  const lines = visible.split('\n')
  for (let i = 0; i < lines.length; i++) {
    for (const ban of banned) {
      if (ban.re.test(lines[i])) hits.push(`${i + 1}: ${ban.label}: ${lines[i].trim()}`)
    }
  }
  return hits
}

describe('user-facing site copy', () => {
  it('ignores dashes and honesty wording that appear only in comments', () => {
    const sample = [
      '// honest \u2014 note',
      'const label = "Open questions"',
      '/* honesty \u2013 list */',
      'export const alt = "Remain Faithful. Accountability That Works"',
    ].join('\n')
    assert.deepEqual(hitsIn(sourceWithoutComments(sample)), [])
  })

  it('flags dashes and honesty wording in copy', () => {
    assert.ok(hitsIn(sourceWithoutComments('const s = "Invite honest hesitation \u2014"')).length > 0)
    assert.ok(hitsIn(sourceWithoutComments("const s = 'groups of 4\u20138'")).length > 0)
  })

  it('has no em dash, en dash, or honesty wording under website/src', () => {
    const problems: string[] = []
    for (const path of filesUnder(srcRoot)) {
      const visible = sourceWithoutComments(readFileSync(path, 'utf8'))
      for (const hit of hitsIn(visible)) {
        problems.push(`${relative(srcRoot, path)}:${hit}`)
      }
    }
    assert.deepEqual(problems, [])
  })
})
