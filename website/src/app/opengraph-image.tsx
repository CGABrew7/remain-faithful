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
          background: '#12214C',
          width: '100%',
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          padding: '64px',
        }}
      >
        <div
          style={{
            background: '#1C3066',
            width: '100%',
            height: '100%',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            padding: '72px 80px',
          }}
        >
          <div
            style={{
              fontSize: 18,
              letterSpacing: '0.22em',
              textTransform: 'uppercase',
              color: '#D1AB4C',
              marginBottom: 24,
            }}
          >
            Fort Wayne · Free forever
          </div>
          <div
            style={{
              fontSize: 72,
              color: '#FFFFFF',
              lineHeight: 1.05,
              letterSpacing: '-1px',
            }}
          >
            Remain Faithful
          </div>
          <div
            style={{
              marginTop: 20,
              fontSize: 28,
              color: '#C5CBE0',
              lineHeight: 1.35,
              maxWidth: 720,
            }}
          >
            Free peer accountability. Always-on filtering. Screenshots and raw content stay on your device.
          </div>
          <div
            style={{
              marginTop: 36,
              width: 48,
              height: 48,
              borderRadius: 24,
              background: '#D1AB4C',
            }}
          />
        </div>
      </div>
    ),
    { ...size }
  )
}
