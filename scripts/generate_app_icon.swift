#!/usr/bin/env swift
// Generates the Remain Faithful app icon at all required iOS sizes.
// Run from the repo root:
//   swift scripts/generate_app_icon.swift
//
// Output: RemainFaithful/Assets.xcassets/AppIcon.appiconset/

import CoreGraphics
import Foundation
import ImageIO
import UniformTypeIdentifiers

// MARK: - Color constants (matches app palette)

let navy   = CGColor(red: 0.07, green: 0.13, blue: 0.30, alpha: 1)
let gold   = CGColor(red: 0.82, green: 0.67, blue: 0.30, alpha: 1)
let white  = CGColor(red: 1,    green: 1,    blue: 1,    alpha: 1)

// MARK: - Icon sizes required by iOS / App Store

struct IconSize {
    let filename: String
    let points: Int
    let scale: Int
    var pixels: Int { points * scale }
}

let sizes: [IconSize] = [
    IconSize(filename: "AppIcon-20@2x.png",   points: 20,   scale: 2),
    IconSize(filename: "AppIcon-20@3x.png",   points: 20,   scale: 3),
    IconSize(filename: "AppIcon-29@2x.png",   points: 29,   scale: 2),
    IconSize(filename: "AppIcon-29@3x.png",   points: 29,   scale: 3),
    IconSize(filename: "AppIcon-38@2x.png",   points: 38,   scale: 2),
    IconSize(filename: "AppIcon-38@3x.png",   points: 38,   scale: 3),
    IconSize(filename: "AppIcon-40@2x.png",   points: 40,   scale: 2),
    IconSize(filename: "AppIcon-40@3x.png",   points: 40,   scale: 3),
    IconSize(filename: "AppIcon-60@2x.png",   points: 60,   scale: 2),
    IconSize(filename: "AppIcon-60@3x.png",   points: 60,   scale: 3),
    IconSize(filename: "AppIcon-76@1x.png",   points: 76,   scale: 1),
    IconSize(filename: "AppIcon-76@2x.png",   points: 76,   scale: 2),
    IconSize(filename: "AppIcon-83.5@2x.png", points: 84,   scale: 2),
    IconSize(filename: "AppIcon-1024.png",    points: 1024, scale: 1),
]

// MARK: - Drawing

func drawIcon(size: Int) -> CGImage? {
    let sz = CGFloat(size)

    guard let ctx = CGContext(
        data: nil,
        width: size,
        height: size,
        bitsPerComponent: 8,
        bytesPerRow: 0,
        space: CGColorSpaceCreateDeviceRGB(),
        bitmapInfo: CGImageAlphaInfo.premultipliedLast.rawValue
    ) else { return nil }

    // 1. Navy background
    ctx.setFillColor(navy)
    ctx.fill(CGRect(x: 0, y: 0, width: sz, height: sz))

    // 2. Subtle radial glow in center
    let center = CGPoint(x: sz / 2, y: sz / 2)
    let glowRadius = sz * 0.52
    if let glowGrad = CGGradient(
        colorsSpace: CGColorSpaceCreateDeviceRGB(),
        colors: [
            CGColor(red: 0.82, green: 0.67, blue: 0.30, alpha: 0.08),
            CGColor(red: 0.82, green: 0.67, blue: 0.30, alpha: 0.00),
        ] as CFArray,
        locations: [0, 1]
    ) {
        ctx.drawRadialGradient(glowGrad,
            startCenter: center, startRadius: 0,
            endCenter:   center, endRadius:   glowRadius,
            options: [])
    }

    // 3. Shield path
    let shieldW = sz * 0.62
    let shieldH = sz * 0.68
    let sx = (sz - shieldW) / 2
    let sy = (sz - shieldH) / 2 + sz * 0.01

    let shieldPath = CGMutablePath()
    // Top arc
    shieldPath.move(to: CGPoint(x: sx + shieldW * 0.5, y: sy))
    shieldPath.addLine(to: CGPoint(x: sx, y: sy + shieldH * 0.20))
    shieldPath.addCurve(
        to: CGPoint(x: sx + shieldW * 0.5, y: sy + shieldH),
        control1: CGPoint(x: sx, y: sy + shieldH * 0.65),
        control2: CGPoint(x: sx + shieldW * 0.25, y: sy + shieldH * 0.88)
    )
    shieldPath.addCurve(
        to: CGPoint(x: sx + shieldW, y: sy + shieldH * 0.20),
        control1: CGPoint(x: sx + shieldW * 0.75, y: sy + shieldH * 0.88),
        control2: CGPoint(x: sx + shieldW, y: sy + shieldH * 0.65)
    )
    shieldPath.closeSubpath()

    // Gold fill
    ctx.setFillColor(gold)
    ctx.addPath(shieldPath)
    ctx.fillPath()

    // Subtle inner shadow on shield (darker navy overlay at edges)
    ctx.saveGState()
    ctx.addPath(shieldPath)
    ctx.clip()
    let innerShad = CGColor(red: 0.05, green: 0.10, blue: 0.25, alpha: 0.25)
    if let innerGrad = CGGradient(
        colorsSpace: CGColorSpaceCreateDeviceRGB(),
        colors: [innerShad, CGColor(red: 0, green: 0, blue: 0, alpha: 0)] as CFArray,
        locations: [0, 1]
    ) {
        ctx.drawLinearGradient(innerGrad,
            start: CGPoint(x: sx, y: sy),
            end:   CGPoint(x: sx + shieldW, y: sy + shieldH),
            options: [])
    }
    ctx.restoreGState()

    // 4. White checkmark
    let ckScale = shieldW * 0.56
    let ckX = (sz - ckScale) / 2 - ckScale * 0.04
    let ckY = (sz - ckScale) / 2 + ckScale * 0.06

    let checkPath = CGMutablePath()
    let lineW = max(2, sz * 0.055)
    checkPath.move(to: CGPoint(x: ckX + ckScale * 0.12, y: ckY + ckScale * 0.50))
    checkPath.addLine(to: CGPoint(x: ckX + ckScale * 0.40, y: ckY + ckScale * 0.76))
    checkPath.addLine(to: CGPoint(x: ckX + ckScale * 0.88, y: ckY + ckScale * 0.22))

    ctx.setStrokeColor(white)
    ctx.setLineWidth(lineW)
    ctx.setLineCap(.round)
    ctx.setLineJoin(.round)
    ctx.addPath(checkPath)
    ctx.strokePath()

    return ctx.makeImage()
}

// MARK: - Save PNG

func savePNG(_ image: CGImage, to path: String) {
    guard let url = URL(string: "file://" + path) else { return }
    guard let dest = CGImageDestinationCreateWithURL(
        url as CFURL,
        "public.png" as CFString,
        1, nil
    ) else {
        print("❌ Could not create destination for \(path)")
        return
    }
    CGImageDestinationAddImage(dest, image, nil)
    if CGImageDestinationFinalize(dest) {
        print("✅ \(path.components(separatedBy: "/").last ?? path)")
    } else {
        print("❌ Failed to write \(path)")
    }
}

// MARK: - Main

let fm = FileManager.default
let repoRoot = fm.currentDirectoryPath
let outDir   = "\(repoRoot)/RemainFaithful/Assets.xcassets/AppIcon.appiconset"

for size in sizes {
    guard let img = drawIcon(size: size.pixels) else {
        print("❌ Failed to draw icon at \(size.pixels)px")
        continue
    }
    savePNG(img, to: "\(outDir)/\(size.filename)")
}

// MARK: - Update Contents.json

let images = sizes.map { s -> [String: String] in
    var entry: [String: String] = [
        "filename": s.filename,
        "idiom":    "universal",
        "platform": "ios",
        "size":     "\(s.points)x\(s.points)",
    ]
    if s.scale > 1 {
        entry["scale"] = "\(s.scale)x"
    }
    return entry
}

let contents: [String: Any] = [
    "images": images,
    "info": ["author": "xcode", "version": 1],
]

if let data = try? JSONSerialization.data(withJSONObject: contents, options: [.prettyPrinted, .sortedKeys]),
   let json = String(data: data, encoding: .utf8) {
    let contentsPath = "\(outDir)/Contents.json"
    try? json.write(toFile: contentsPath, atomically: true, encoding: .utf8)
    print("✅ Updated Contents.json")
}

print("\nDone! Run `swift scripts/generate_app_icon.swift` from the repo root to regenerate.")
