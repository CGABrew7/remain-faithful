import type { Metadata } from 'next'
import Link from 'next/link'
import { notFound } from 'next/navigation'
import { posts } from '../posts'
import WaitlistForm from '@/components/WaitlistForm'
import { RelatedPosts } from '@/components/RelatedPosts'
import { Breadcrumbs } from '@/components/Breadcrumbs'
import { JsonLd } from '@/components/JsonLd'
import { blogPostSchema } from '@/lib/structured-data'

export function generateStaticParams() {
  return posts.map((p) => ({ slug: p.slug }))
}

const blogTitles: Record<string, string> = {
  'why-accountability-fails': 'Why Accountability Fails (And How to Fix It)',
  'setting-up-your-first-group': 'How to Set Up Your First Accountability Group',
  'science-of-peer-accountability': 'The Science Behind Peer Accountability',
  'on-device-privacy-explained': 'On-Device AI: Why Your Content Never Leaves Your Phone',
  'covenant-model': 'The Covenant Model: More Than an Agreement',
  'mens-ministry-accountability': 'Modernizing Ministry Accountability for Churches',
}

const relatedMap: Record<string, string[]> = {
  'why-accountability-fails': ['science-of-peer-accountability', 'covenant-model'],
  'setting-up-your-first-group': ['mens-ministry-accountability', 'why-accountability-fails'],
  'science-of-peer-accountability': ['why-accountability-fails', 'on-device-privacy-explained'],
  'on-device-privacy-explained': ['covenant-model', 'science-of-peer-accountability'],
  'covenant-model': ['why-accountability-fails', 'on-device-privacy-explained'],
  'mens-ministry-accountability': ['setting-up-your-first-group', 'covenant-model'],
}

export async function generateMetadata({
  params,
}: {
  params: { slug: string }
}): Promise<Metadata> {
  const post = posts.find((p) => p.slug === params.slug)
  if (!post) return {}
  const displayTitle = blogTitles[post.slug] || post.title
  return {
    title: `${displayTitle} | Remain Faithful`,
    description: post.excerpt,
    alternates: { canonical: `https://remainfaithful.com/blog/${post.slug}` },
    openGraph: {
      title: `${displayTitle} | Remain Faithful`,
      description: post.excerpt,
      type: 'article',
    },
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
        <h2 key={i} className="font-serif text-2xl font-bold text-ink mt-10 mb-4">
          {line.slice(3)}
        </h2>
      )
    } else if (line.startsWith('**') && line.endsWith('**')) {
      elements.push(
        <p key={i} className="font-semibold text-ink mt-4 mb-2">
          {line.slice(2, -2)}
        </p>
      )
    } else if (line.startsWith('**') && line.includes('**')) {
      const boldEnd = line.indexOf('**', 2)
      const bold = line.slice(2, boldEnd)
      const rest = line.slice(boldEnd + 2)
      elements.push(
        <p key={i} className="text-ink-soft leading-relaxed mb-4">
          <strong className="text-ink">{bold}</strong>{rest}
        </p>
      )
    } else if (line.trim() !== '') {
      elements.push(
        <p key={i} className="text-ink-soft leading-relaxed mb-4">
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

  const relatedSlugs = relatedMap[post.slug] || []
  const relatedPosts = relatedSlugs
    .map((slug) => posts.find((p) => p.slug === slug))
    .filter((p): p is NonNullable<typeof p> => !!p)
    .map((p) => ({ title: blogTitles[p.slug] || p.title, slug: p.slug, description: p.excerpt, readTime: p.readTime }))

  return (
    <div className="pt-24 pb-24">
      <JsonLd data={blogPostSchema({ title: blogTitles[post.slug] || post.title, description: post.excerpt, slug: post.slug, datePublished: post.date, readTime: post.readTime })} />

      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <Breadcrumbs items={[
          { name: 'Blog', url: 'https://remainfaithful.com/blog' },
          { name: blogTitles[post.slug] || post.title, url: `https://remainfaithful.com/blog/${post.slug}` },
        ]} />
        <div className="grid lg:grid-cols-[1fr_300px] gap-12">
          {/* Article */}
          <article>
            <Link
              href="/blog"
              className="inline-flex items-center gap-2 text-sm text-ink-soft hover:text-ink transition-colors mb-8"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                <path d="M19 12H5M12 19l-7-7 7-7"/>
              </svg>
              Back to Blog
            </Link>

            <div className="mb-4">
              <span className="text-xs font-semibold uppercase tracking-wider text-wax bg-wax/10 px-3 py-1 rounded-sm">
                {post.category}
              </span>
            </div>

            <h1 className="font-serif text-3xl sm:text-4xl font-bold text-ink mb-5 leading-tight">
              {post.title}
            </h1>

            <div className="flex items-center gap-4 text-sm text-ink-soft mb-8 pb-8 border-b border-hairline">
              <div className="flex items-center gap-2">
                <div className="w-7 h-7 rounded-sm bg-wax flex items-center justify-center text-paper text-xs font-bold">
                  JB
                </div>
                <span>Jeff Brewer</span>
              </div>
              <span>·</span>
              <span>{post.readTime}</span>
            </div>

            <div className="text-ink-soft leading-relaxed">
              {renderBody(post.body)}
            </div>

            {/* Newsletter CTA */}
            <div
              className="mt-12 p-8 rounded-sm border border-wax/25"
              
            >
              <h3 className="font-serif text-xl font-bold text-ink mb-2">
                Want more like this?
              </h3>
              <p className="text-ink-soft text-sm mb-5">
                Get the monthly accountability newsletter — practical guides, theological reflections, and updates from the Remain Faithful team.
              </p>
              <WaitlistForm variant="inline" buttonText="Subscribe" />
            </div>

            {/* Download CTA */}
            <div
              className="mt-6 p-8 rounded-sm border border-wax/25 text-center"
              
            >
              <h3 className="font-serif text-xl font-bold text-ink mb-3">
                Ready to start?
              </h3>
              <p className="text-ink-soft text-sm mb-5">
                Join the waitlist and be first to know when Remain Faithful launches.
              </p>
              <Link
                href="/#waitlist"
                className="inline-flex items-center gap-2 px-6 py-3 rounded-sm font-semibold text-paper bg-wax hover:bg-wax-deep transition-[box-shadow,scale] duration-200 ease-out active:scale-[0.96] text-sm"
              >
                Join the Waitlist
              </Link>
            </div>

            {relatedPosts.length > 0 && <RelatedPosts posts={relatedPosts} />}
          </article>

          {/* Sidebar */}
          <aside className="space-y-6">
            {/* Related posts */}
            <div className="rounded-sm border border-hairline bg-paper-deep p-6">
              <h3 className="font-serif text-lg font-bold text-ink mb-5">More From Our Blog</h3>
              <div className="space-y-4">
                {posts.filter((p) => p.slug !== post.slug).slice(0, 3).map((p, i, arr) => (
                  <Link
                    key={p.slug}
                    href={`/blog/${p.slug}`}
                    className="block group"
                  >
                    <p className="text-xs text-wax font-semibold uppercase tracking-wide mb-1">{p.category}</p>
                    <p className="text-sm text-ink font-medium group-hover:text-wax transition-colors leading-snug mb-1">
                      {p.title}
                    </p>
                    <p className="text-xs text-ink-soft">{p.readTime}</p>
                    {i < arr.length - 1 && (
                      <div className="mt-4 border-t border-hairline"/>
                    )}
                  </Link>
                ))}
              </div>
            </div>

            {/* Donate nudge */}
            <div
              className="rounded-sm p-6 border border-wax/25 text-center"
              
            >
              <div className="text-2xl mb-3">❤️</div>
              <h3 className="font-serif text-base font-bold text-ink mb-2">Keep RF Free</h3>
              <p className="text-xs text-ink-soft mb-4 leading-relaxed">
                Remain Faithful is free because of voluntary donations from people like you.
              </p>
              <Link
                href="/#donate"
                className="block text-center py-2.5 rounded-sm text-sm font-semibold text-paper bg-wax"
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
