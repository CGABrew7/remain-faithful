import type { Metadata } from 'next'
import Link from 'next/link'
import { notFound } from 'next/navigation'
import { posts } from '../posts'

export function generateStaticParams() {
  return posts.map((p) => ({ slug: p.slug }))
}

export async function generateMetadata({
  params,
}: {
  params: { slug: string }
}): Promise<Metadata> {
  const post = posts.find((p) => p.slug === params.slug)
  if (!post) return {}
  return {
    title: post.title,
    description: post.excerpt,
    openGraph: { title: post.title, description: post.excerpt },
  }
}

function renderBody(body: string) {
  const lines = body.split('\n')
  const elements: React.ReactNode[] = []
  let i = 0
  while (i < lines.length) {
    const line = lines[i]
    if (line.startsWith('## ')) {
      elements.push(
        <h2 key={i} className="font-serif text-2xl font-bold text-[#F0EDE8] mt-10 mb-4">
          {line.slice(3)}
        </h2>
      )
    } else if (line.startsWith('**') && line.endsWith('**')) {
      elements.push(
        <p key={i} className="font-semibold text-[#F0EDE8] mt-4 mb-2">
          {line.slice(2, -2)}
        </p>
      )
    } else if (line.startsWith('**') && line.includes('**')) {
      const boldEnd = line.indexOf('**', 2)
      const bold = line.slice(2, boldEnd)
      const rest = line.slice(boldEnd + 2)
      elements.push(
        <p key={i} className="text-[#8A9BB0] leading-relaxed mb-4">
          <strong className="text-[#F0EDE8]">{bold}</strong>{rest}
        </p>
      )
    } else if (line.trim() !== '') {
      elements.push(
        <p key={i} className="text-[#8A9BB0] leading-relaxed mb-4">
          {line}
        </p>
      )
    }
    i++
  }
  return elements
}

export default function BlogPostPage({ params }: { params: { slug: string } }) {
  const post = posts.find((p) => p.slug === params.slug)
  if (!post) notFound()

  const related = posts.filter((p) => p.slug !== post.slug).slice(0, 3)

  return (
    <div className="pt-32 pb-24">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid lg:grid-cols-[1fr_300px] gap-12">
          {/* Article */}
          <article>
            <Link
              href="/blog"
              className="inline-flex items-center gap-2 text-sm text-[#8A9BB0] hover:text-[#F0EDE8] transition-colors mb-8"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                <path d="M19 12H5M12 19l-7-7 7-7"/>
              </svg>
              Back to Blog
            </Link>

            <div className="mb-4">
              <span className="text-xs font-semibold uppercase tracking-wider text-[#C9A84C] bg-[#C9A84C]/10 px-3 py-1 rounded-full">
                {post.category}
              </span>
            </div>

            <h1 className="font-serif text-3xl sm:text-4xl font-bold text-[#F0EDE8] mb-5 leading-tight">
              {post.title}
            </h1>

            <div className="flex items-center gap-4 text-sm text-[#8A9BB0] mb-8 pb-8 border-b border-[#1E3050]">
              <div className="flex items-center gap-2">
                <div className="w-7 h-7 rounded-full bg-gradient-to-br from-[#C9A84C] to-[#E8C87A] flex items-center justify-center text-[#0F1B2D] text-xs font-bold">
                  RF
                </div>
                <span>RF Team</span>
              </div>
              <span>·</span>
              <span>{post.date}</span>
              <span>·</span>
              <span>{post.readTime}</span>
            </div>

            <div className="text-[#8A9BB0] leading-relaxed">
              {renderBody(post.body)}
            </div>

            {/* CTA */}
            <div
              className="mt-12 p-8 rounded-2xl border border-[#C9A84C]/20 text-center"
              style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}
            >
              <h3 className="font-serif text-xl font-bold text-[#F0EDE8] mb-3">
                Ready to start?
              </h3>
              <p className="text-[#8A9BB0] text-sm mb-5">
                Download Remain Faithful and take the first step toward real accountability.
              </p>
              <Link
                href="/#download"
                className="inline-flex items-center gap-2 px-6 py-3 rounded-full font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A] hover:from-[#E8C87A] hover:to-[#C9A84C] transition-all duration-200 text-sm"
              >
                Download for iPhone
              </Link>
            </div>
          </article>

          {/* Sidebar */}
          <aside className="space-y-6">
            {/* Related posts */}
            <div className="rounded-2xl border border-[#1E3050] bg-[#162235] p-6">
              <h3 className="font-serif text-lg font-bold text-[#F0EDE8] mb-5">More From Our Blog</h3>
              <div className="space-y-4">
                {related.map((p) => (
                  <Link
                    key={p.slug}
                    href={`/blog/${p.slug}`}
                    className="block group"
                  >
                    <p className="text-xs text-[#C9A84C] font-semibold uppercase tracking-wide mb-1">{p.category}</p>
                    <p className="text-sm text-[#F0EDE8] font-medium group-hover:text-[#C9A84C] transition-colors leading-snug mb-1">
                      {p.title}
                    </p>
                    <p className="text-xs text-[#8A9BB0]">{p.readTime}</p>
                    {related.indexOf(p) < related.length - 1 && (
                      <div className="mt-4 border-t border-[#1E3050]"/>
                    )}
                  </Link>
                ))}
              </div>
            </div>

            {/* Donate nudge */}
            <div
              className="rounded-2xl p-6 border border-[#C9A84C]/20 text-center"
              style={{ background: 'linear-gradient(135deg, #162235, #1A2A40)' }}
            >
              <div className="text-2xl mb-3">❤️</div>
              <h3 className="font-serif text-base font-bold text-[#F0EDE8] mb-2">Keep RF Free</h3>
              <p className="text-xs text-[#8A9BB0] mb-4 leading-relaxed">
                Remain Faithful is free because of voluntary donations from people like you.
              </p>
              <Link
                href="/#donate"
                className="block text-center py-2.5 rounded-full text-sm font-semibold text-[#0F1B2D] bg-gradient-to-r from-[#C9A84C] to-[#E8C87A]"
              >
                Support the Project
              </Link>
            </div>
          </aside>
        </div>
      </div>
    </div>
  )
}
