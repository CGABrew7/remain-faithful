import ShieldMark from '@/components/ShieldMark'

export type FeaturePlateLine = {
  label: string
  detail: string
}

type Props = {
  kicker: string
  title: string
  lines: FeaturePlateLine[]
  footer?: string
  caption?: string
  /** Fill the grid row so a short card stays balanced beside a taller column. */
  stretch?: boolean
}

export default function FeaturePlate({
  kicker,
  title,
  lines,
  footer,
  caption,
  stretch = false,
}: Props) {
  const Tag = caption ? 'figure' : 'div'
  return (
    <Tag
      className={
        stretch
          ? 'feature-plate feature-plate-stretch'
          : 'feature-plate'
      }
    >
      <div className="feature-plate-card">
        <div className="feature-plate-mark" aria-hidden="true">
          <ShieldMark size={28} />
        </div>
        <p className="kicker">{kicker}</p>
        <p className="feature-plate-title">{title}</p>
        <ol className="feature-plate-list">
          {lines.map((line, index) => (
            <li key={line.label} className="feature-plate-row">
              <span className="feature-plate-index">
                {String(index + 1).padStart(2, '0')}
              </span>
              <span className="min-w-0">
                <span className="block font-semibold leading-snug text-white">{line.label}</span>
                <span className="mt-0.5 block text-sm leading-relaxed text-ink-soft">{line.detail}</span>
              </span>
            </li>
          ))}
        </ol>
        {stretch && footer && (
          <div className="feature-plate-spacer" aria-hidden="true">
            <span />
          </div>
        )}
        {footer && <p className="feature-plate-footer">{footer}</p>}
      </div>
      {caption && <figcaption className="feature-plate-caption">{caption}</figcaption>}
    </Tag>
  )
}
