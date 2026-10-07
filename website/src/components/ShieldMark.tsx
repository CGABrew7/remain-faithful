export default function ShieldMark({ size = 22 }: { size?: number }) {
  const height = Math.round(size * (36 / 32))
  return (
    <svg
      width={size}
      height={height}
      viewBox="0 0 32 36"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden="true"
    >
      <path
        d="M16 0L2 6V18C2 26.284 8.268 33.916 16 36C23.732 33.916 30 26.284 30 18V6L16 0Z"
        fill="#D1AB4C"
      />
      <path
        d="M11 18L14.5 21.5L21 14"
        stroke="#12214C"
        strokeWidth="2.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}
