/**
 * Rasterize public/favicon.svg into PNG + ICO.
 * Run from website/: node scripts/generate-favicons.mjs
 */
import { readFileSync, writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { Resvg } from '@resvg/resvg-js'

const root = join(dirname(fileURLToPath(import.meta.url)), '..')
const publicDir = join(root, 'public')
const svg = readFileSync(join(publicDir, 'favicon.svg'), 'utf8')

const appleSvg = svg.replace('rx="6"', 'rx="0"')

function renderPng(source, size) {
  const resvg = new Resvg(source, {
    fitTo: { mode: 'width', value: size },
  })
  return resvg.render().asPng()
}

function pngToIco(pngBuffers) {
  const count = pngBuffers.length
  const headerSize = 6 + count * 16
  let offset = headerSize
  const entries = []
  for (const png of pngBuffers) {
    const size = readPngSize(png)
    entries.push({ width: size, height: size, png, offset, bytes: png.length })
    offset += png.length
  }

  const buf = Buffer.alloc(offset)
  buf.writeUInt16LE(0, 0)
  buf.writeUInt16LE(1, 2)
  buf.writeUInt16LE(count, 4)

  entries.forEach((entry, i) => {
    const o = 6 + i * 16
    buf.writeUInt8(entry.width >= 256 ? 0 : entry.width, o)
    buf.writeUInt8(entry.height >= 256 ? 0 : entry.height, o + 1)
    buf.writeUInt8(0, o + 2)
    buf.writeUInt8(0, o + 3)
    buf.writeUInt16LE(1, o + 4)
    buf.writeUInt16LE(32, o + 6)
    buf.writeUInt32LE(entry.bytes, o + 8)
    buf.writeUInt32LE(entry.offset, o + 12)
    entry.png.copy(buf, entry.offset)
  })
  return buf
}

function readPngSize(png) {
  return png.readUInt32BE(16)
}

const png16 = renderPng(svg, 16)
const png32 = renderPng(svg, 32)
const png48 = renderPng(svg, 48)
const png180 = renderPng(appleSvg, 180)
const png192 = renderPng(svg, 192)

writeFileSync(join(publicDir, 'favicon-32x32.png'), png32)
writeFileSync(join(publicDir, 'apple-touch-icon.png'), png180)
writeFileSync(join(publicDir, 'favicon-192x192.png'), png192)
writeFileSync(join(publicDir, 'favicon.ico'), pngToIco([png16, png32, png48]))

console.log('Wrote favicon.ico, favicon-32x32.png, favicon-192x192.png, apple-touch-icon.png')
