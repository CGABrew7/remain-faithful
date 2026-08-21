'use client'

import { useEffect, useRef, type ReactNode } from 'react'

/*
  Hero atmosphere extracted from ThreeUI Community
  `meng-to-sketchbook-landing-page` (Meng To, MIT).

  Kept: paper sheet as a physical object; layered desk-cast lighting
  (ambient / contact / hair); restrained pointer tilt; prefers-reduced-motion.

  Changed: dusk desk instead of Singapore cream wash; wax-red seal instead
  of earth brown; no sketchbook, page-turn, loupe, botany, or plates.
  First frame is CSS-only — tilt is progressive enhancement.
*/

const TILT_X = 3.2
const TILT_Y = 4.5

export default function LetterAtmosphere({ children }: { children: ReactNode }) {
  const tiltRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const tilt = tiltRef.current
    if (!tilt) return
    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (reduced) return

    const setTilt = (rx: number, ry: number) => {
      tilt.style.setProperty('--rx', `${rx.toFixed(2)}deg`)
      tilt.style.setProperty('--ry', `${ry.toFixed(2)}deg`)
    }

    const onMove = (event: PointerEvent) => {
      if (event.pointerType === 'touch') return
      const bounds = tilt.getBoundingClientRect()
      if (!bounds.width || !bounds.height) return
      const nx = Math.max(-1, Math.min(1, (event.clientX - (bounds.left + bounds.width / 2)) / (bounds.width * 0.62)))
      const ny = Math.max(-1, Math.min(1, (event.clientY - (bounds.top + bounds.height / 2)) / (bounds.height * 0.9)))
      setTilt(-ny * TILT_X, nx * TILT_Y)
    }

    const clear = () => setTilt(0, 0)

    window.addEventListener('pointermove', onMove, { passive: true })
    window.addEventListener('blur', clear)
    return () => {
      window.removeEventListener('pointermove', onMove)
      window.removeEventListener('blur', clear)
    }
  }, [])

  return (
    <section className="letter-hero" aria-labelledby="letter-heading">
      <div className="letter-wash" aria-hidden="true" />
      <div className="letter-desk">
        <div className="letter-tilt" ref={tiltRef}>
          <div className="letter-cast ambient" aria-hidden="true" />
          <div className="letter-cast contact" aria-hidden="true" />
          <div className="letter-cast hair" aria-hidden="true" />
          <article className="letter-sheet">{children}</article>
        </div>
      </div>
    </section>
  )
}
