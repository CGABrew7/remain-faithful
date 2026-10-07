export default function Ornament({ className = '' }: { className?: string }) {
  return (
    <div className={`ornament ${className}`} aria-hidden="true">
      <span />
      <svg width="12" height="12" viewBox="0 0 16 16" fill="none">
        <path
          d="M8 14.2C8 14.2 2.2 10.4 2.2 6.5 2.2 4.3 3.8 2.8 5.7 2.8c1.05 0 1.9.55 2.3 1.4.4-.85 1.25-1.4 2.3-1.4 1.9 0 3.5 1.5 3.5 3.7 0 3.9-5.8 7.7-5.8 7.7z"
          fill="#D1AB4C"
        />
      </svg>
      <span />
    </div>
  )
}
