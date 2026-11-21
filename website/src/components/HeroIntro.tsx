const WORD = 'Remain Faithful'

const boot = `(function(){try{var reduce=matchMedia('(prefers-reduced-motion: reduce)').matches;var seen=sessionStorage.getItem('rf-hero-intro')==='1';if(reduce||seen){document.documentElement.classList.add('hero-still')}else{sessionStorage.setItem('rf-hero-intro','1')}}catch(e){}})();`

export default function HeroIntro() {
  return (
    <div className="hero-lockup">
      <script dangerouslySetInnerHTML={{ __html: boot }} />
      <span className="hero-bloom" aria-hidden="true" />
      <svg className="hero-mark" width="72" height="81" viewBox="0 0 32 36" fill="none" aria-hidden="true">
        <path className="hero-mark-line" pathLength="100" d="M16 0L2 6V18C2 26.284 8.268 33.916 16 36C23.732 33.916 30 26.284 30 18V6L16 0Z" />
        <path className="hero-mark-line hero-mark-check" pathLength="100" d="M11 18L14.5 21.5L21 14" />
      </svg>
      <h1 className="hero-word font-serif text-[2.75rem] font-bold leading-[1.04] tracking-[-0.03em] text-white sm:text-6xl lg:text-[4.15rem]" aria-label={WORD}>
        {WORD.split('').map((character, index) => (
          <span key={`${character}-${index}`} style={{ animationDelay: `${2.2 + index * 0.055}s` }}>
            {character === ' ' ? '\u00a0' : character}
          </span>
        ))}
      </h1>
    </div>
  )
}
