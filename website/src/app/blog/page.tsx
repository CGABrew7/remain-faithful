import type { Metadata } from 'next'
import Link from 'next/link'
import { posts } from './posts'

export const metadata: Metadata = {
  title: 'Blog',
  description: 'Insights on accountability, purity, ministry, and the technology behind Remain Faithful.',
}

export default function BlogPage() {
  return (
    <div className="pt-32 pb-24">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <p className="text-[#C9A84C] text-sm font-semibold uppercase tracking-widest mb-3">The RF Blog</p>
          <h1 className="font-serif text-4xl sm:text-5xl font-bold text-[#F0EDE8] mb-4">
            Thoughts on Accountability
          </h1>
          <p className="text-[#8A9BB0] max-w-xl mx-auto">
            Practical guides, theological reflections, and honest conversations about the struggle and the tools we&apos;ve built to face it together.
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {posts.map((post) => (
            <Link
              key={post.slug}
              href={`/blog/${post.slug}`}
              className="group rounded-2xl border border-[#1E3050] bg-[#162235] overflow-hidden hover:border-[#C9A84C]/30 transition-colors duration-300"
            >
              {/* Color band */}
              <div className="h-1.5 bg-gradient-to-r from-[#C9A84C] to-[#E8C87A]" />
              <div className="p-6">
                <div className="flex items-center gap-2 mb-4">
                  <span className="text-xs font-semibold uppercase tracking-wider text-[#C9A84C] bg-[#C9A84C]/10 px-2.5 py-1 rounded-full">
                    {post.category}
                  </span>
                </div>
                <h2 className="font-serif text-lg font-bold text-[#F0EDE8] mb-3 group-hover:text-[#C9A84C] transition-colors leading-snug">
                  {post.title}
                </h2>
                <p className="text-sm text-[#8A9BB0] leading-relaxed mb-5 line-clamp-3">
                  {post.excerpt}
                </p>
                <div className="flex items-center justify-end text-xs text-[#8A9BB0]">
                  <span>{post.readTime}</span>
                </div>
              </div>
              <div className="px-6 pb-5">
                <span className="inline-flex items-center gap-1 text-xs font-semibold text-[#C9A84C] group-hover:gap-2 transition-all">
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
