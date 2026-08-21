import type { Metadata } from 'next'
import Link from 'next/link'
import { posts } from './posts'
import { Breadcrumbs } from '@/components/Breadcrumbs'

export const metadata: Metadata = {
  title: 'Accountability Guides & Resources | Remain Faithful Blog',
  description: 'Practical guides, theological reflections, and research on accountability, purity, and the technology behind Remain Faithful.',
  alternates: { canonical: 'https://remainfaithful.com/blog' },
}

export default function BlogPage() {
  return (
    <div className="pt-24 pb-24">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <Breadcrumbs items={[{ name: 'Blog', url: 'https://remainfaithful.com/blog' }]} />
        <div className="text-center mb-16">
          <p className="text-wax text-sm font-semibold uppercase tracking-widest mb-3">The RF Blog</p>
          <h1 className="font-serif text-4xl sm:text-5xl font-bold text-ink mb-4">
            Thoughts on Accountability
          </h1>
          <p className="text-ink-soft max-w-xl mx-auto">
            Practical guides, theological reflections, and honest conversations about the struggle and the tools we&apos;ve built to face it together.
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {posts.map((post) => (
            <Link
              key={post.slug}
              href={`/blog/${post.slug}`}
              className="group rounded-sm border border-hairline bg-paper-deep overflow-hidden hover:border-wax/30 transition-colors duration-300"
            >
              {/* Color band */}
              <div className="h-1.5 bg-wax" />
              <div className="p-6">
                <div className="flex items-center gap-2 mb-4">
                  <span className="text-xs font-semibold uppercase tracking-wider text-wax bg-wax/10 px-2.5 py-1 rounded-sm">
                    {post.category}
                  </span>
                </div>
                <h2 className="font-serif text-lg font-bold text-ink mb-3 group-hover:text-wax transition-colors leading-snug">
                  {post.title}
                </h2>
                <p className="text-sm text-ink-soft leading-relaxed mb-5 line-clamp-3">
                  {post.excerpt}
                </p>
                <div className="flex items-center justify-end text-xs text-ink-soft">
                  <span>{post.readTime}</span>
                </div>
              </div>
              <div className="px-6 pb-5">
                <span className="inline-flex items-center gap-1 text-xs font-semibold text-wax group-hover:gap-2 transition-[gap] duration-200">
                  Read More
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                    <path d="M5 12h14M12 5l7 7-7 7"/>
                  </svg>
                </span>
              </div>
            </Link>
          ))}
        </div>
      </div>
    </div>
  )
}
