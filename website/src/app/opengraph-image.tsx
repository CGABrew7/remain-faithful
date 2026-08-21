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
          background: '#1c1713',
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
            background: '#ece6d8',
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
              color: '#85786c',
              marginBottom: 24,
            }}
          >
            Fort Wayne · Free forever
          </div>
          <div
            style={{
              fontSize: 72,
              color: '#1d1814',
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
              color: '#5a5148',
              lineHeight: 1.35,
              maxWidth: 720,
            }}
          >
            Free peer accountability. Always-on filtering. No screen content ever leaves your device.
          </div>
          <div
            style={{
              marginTop: 36,
              width: 48,
              height: 48,
              borderRadius: 24,
              background: '#7a2c28',
            }}
          />
        </div>
      </div>
    ),
    { ...size }
  )
}
