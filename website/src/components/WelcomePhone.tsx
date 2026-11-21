const line = 'A covenant group, a partner, and Session Shield.'

export default function WelcomePhone() {
  return (
    <figure className="phone">
      <div className="phone-bezel">
        <span className="phone-btn phone-btn-silent" aria-hidden="true" />
        <span className="phone-btn phone-btn-vol" aria-hidden="true" />
        <span className="phone-btn phone-btn-power" aria-hidden="true" />
        <div className="phone-clip welcome-screen">
          <svg className="welcome-mark" viewBox="0 0 32 36" fill="none" aria-hidden="true">
            <path d="M16 0L2 6V18C2 26.284 8.268 33.916 16 36C23.732 33.916 30 26.284 30 18V6L16 0Z" stroke="#F0D030" strokeWidth="1.4" />
            <path d="M11 18L14.5 21.5L21 14" stroke="#F0D030" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
          <p className="welcome-title">Remain Faithful</p>
          <p className="welcome-line">{line}</p>
          <p className="welcome-verse">But you must remain faithful to the things you have been taught.</p>
          <p className="welcome-ref">2 Timothy 3:14</p>
          <span className="welcome-cta">Begin Your Journey</span>
          <p className="welcome-signin">
            Already have an account? <span>Sign In</span>
          </p>
        </div>
      </div>
      <figcaption className="sr-only">
        Remain Faithful welcome screen with a flat yellow Begin Your Journey button.
      </figcaption>
    </figure>
  )
}
