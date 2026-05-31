import { ImageResponse } from 'next/og'

export const runtime = 'edge'
export const alt = 'Remain Faithful – Accountability That Works'
export const size = { width: 1200, height: 630 }
export const contentType = 'image/png'

export default async function Image() {
  return new ImageResponse(
    (
      <div
        style={{
          background: '#0F1B2D',
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          padding: '80px',
          position: 'relative',
        }}
      >
        {/* Gold accent bar */}
        <div
          style={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            height: 6,
            background: 'linear-gradient(90deg, #C9A84C, #E8C87A, #C9A84C)',
          }}
        />

        {/* Shield SVG */}
        <svg width="90" height="102" viewBox="0 0 32 36" fill="none">
          <path
            d="M16 0L2 6V18C2 26.284 8.268 33.916 16 36C23.732 33.916 30 26.284 30 18V6L16 0Z"
            fill="#C9A84C"
          />
          <path
            d="M11 18L14.5 21.5L21 14"
            stroke="white"
            strokeWidth="2.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>

        <div style={{ height: 28 }} />

        <div
          style={{
            fontSize: 64,
            fontWeight: 700,
            color: '#F0EDE8',
            textAlign: 'center',
            lineHeight: 1.1,
            letterSpacing: '-1px',
          }}
        >
          Remain Faithful
        </div>

        <div style={{ height: 20 }} />

        <div
          style={{
            fontSize: 32,
            color: '#C9A84C',
            textAlign: 'center',
            fontWeight: 600,
          }}
        >
          Accountability That Works
        </div>

        <div style={{ height: 28 }} />

        <div
          style={{
            display: 'flex',
            gap: 32,
            fontSize: 18,
            color: '#8A9BB0',
          }}
        >
          <span>100% Free</span>
          <span>·</span>
          <span>Privacy-First</span>
          <span>·</span>
          <span>Open Source</span>
          <span>·</span>
          <span>iOS 17+</span>
        </div>

        {/* Bottom accent */}
        <div
          style={{
            position: 'absolute',
            bottom: 0,
            left: 0,
            right: 0,
            height: 3,
            background: 'linear-gradient(90deg, transparent, #C9A84C, transparent)',
          }}
        />
      </div>
    ),
    { ...size }
  )
}
