import Image from 'next/image'

type Props = {
  src: string
  alt: string
  priority?: boolean
}

export default function PhoneFrame({ src, alt, priority = false }: Props) {
  return (
    <figure className="phone">
      <div className="phone-glow" aria-hidden="true" />
      <div className="phone-bezel">
        <span className="phone-btn phone-btn-silent" aria-hidden="true" />
        <span className="phone-btn phone-btn-vol" aria-hidden="true" />
        <span className="phone-btn phone-btn-power" aria-hidden="true" />
        <Image
          src={src}
          alt={alt}
          width={390}
          height={844}
          priority={priority}
          sizes="(max-width: 768px) 68vw, 280px"
          className="phone-screen"
        />
      </div>
    </figure>
  )
}
