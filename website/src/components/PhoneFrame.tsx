import Image from 'next/image'

type Props = {
  src: string
  alt: string
  priority?: boolean
}

export default function PhoneFrame({ src, alt, priority = false }: Props) {
  const welcome = src.endsWith('welcome.png')

  return (
    <figure className="phone">
      <div className="phone-glow" aria-hidden="true" />
      <div className="phone-bezel">
        <span className="phone-btn phone-btn-silent" aria-hidden="true" />
        <span className="phone-btn phone-btn-vol" aria-hidden="true" />
        <span className="phone-btn phone-btn-power" aria-hidden="true" />
        <div className={welcome ? 'phone-clip phone-clip-welcome' : 'phone-clip'}>
          <Image
            src={src}
            alt={alt}
            width={390}
            height={844}
            priority={priority}
            sizes="(max-width: 768px) 68vw, 280px"
            className="phone-screen"
          />
          <div className="phone-status" aria-hidden="true">
            <span className="phone-time">9:41</span>
            <span className="phone-glyphs">
              <svg viewBox="0 0 17 12" aria-hidden="true">
                <rect x="0" y="7" width="3" height="5" rx="0.6" />
                <rect x="4.5" y="4.5" width="3" height="7.5" rx="0.6" />
                <rect x="9" y="2" width="3" height="10" rx="0.6" />
                <rect x="13.5" y="0" width="3" height="12" rx="0.6" />
              </svg>
              <svg viewBox="0 0 16 12" fill="none" aria-hidden="true">
                <path d="M1 4.2a10 10 0 0114 0M3.4 6.6a6.5 6.5 0 019.2 0M6 9a3 3 0 014 0" stroke="#fff" strokeWidth="1.6" strokeLinecap="round" />
                <circle cx="8" cy="11" r="0.9" fill="#fff" />
              </svg>
              <svg viewBox="0 0 25 12" aria-hidden="true">
                <rect x="0.6" y="0.6" width="21" height="10.8" rx="2.4" fill="none" stroke="#fff" strokeWidth="1.1" />
                <rect x="2.1" y="2.1" width="16.2" height="7.8" rx="1.2" />
                <rect x="22.6" y="3.6" width="1.6" height="4.8" rx="0.6" />
              </svg>
            </span>
          </div>
          {welcome ? (
            <p className="phone-signin">
              <span>Already have an account? <strong>Sign In</strong></span>
              <span className="phone-home-bar" aria-hidden="true" />
            </p>
          ) : null}
        </div>
      </div>
    </figure>
  )
}
