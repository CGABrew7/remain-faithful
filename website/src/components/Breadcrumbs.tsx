import Link from 'next/link'
import { JsonLd } from './JsonLd'
import { breadcrumbSchema } from '@/lib/structured-data'

interface BreadcrumbItem {
  name: string
  url: string
}

export function Breadcrumbs({ items }: { items: BreadcrumbItem[] }) {
  const allItems = [{ name: 'Home', url: 'https://remainfaithful.com' }, ...items]

  return (
    <>
      <JsonLd data={breadcrumbSchema(allItems)} />
      <nav aria-label="Breadcrumb" className="font-mono text-[11px] tracking-[0.1em] uppercase text-ink-faint mb-6">
        {allItems.map((item, i) => (
          <span key={item.url}>
            {i > 0 && <span className="mx-2 text-hairline">/</span>}
            {i === allItems.length - 1 ? (
              <span className="text-ink">{item.name}</span>
            ) : (
              <Link href={item.url} className="hover:text-wax transition-colors duration-200">
                {item.name}
              </Link>
            )}
          </span>
        ))}
      </nav>
    </>
  )
}
