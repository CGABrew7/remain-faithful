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
      <nav aria-label="Breadcrumb" className="text-sm text-[#8A9BB0] mb-6">
        {allItems.map((item, i) => (
          <span key={item.url}>
            {i > 0 && <span className="mx-2 text-[#1E3050]">/</span>}
            {i === allItems.length - 1 ? (
              <span className="text-[#F0EDE8]">{item.name}</span>
            ) : (
              <Link href={item.url} className="hover:text-[#C9A84C] transition-colors">
                {item.name}
              </Link>
            )}
          </span>
        ))}
      </nav>
    </>
  )
}
