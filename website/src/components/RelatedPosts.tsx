import Link from 'next/link'

interface RelatedPost {
  title: string
  slug: string
  description: string
  readTime: string
}

export function RelatedPosts({ posts }: { posts: RelatedPost[] }) {
  return (
    <section className="mt-12 pt-8 border-t border-hairline">
      <h2 className="font-serif text-xl font-medium text-ink mb-4">Related Articles</h2>
      <div className="grid gap-0 md:grid-cols-1">
        {posts.map((post) => (
          <Link
            key={post.slug}
            href={`/blog/${post.slug}`}
            className="block py-4 border-b border-hairline hover:pl-2 transition-[padding] duration-200"
          >
            <h3 className="font-serif text-lg text-ink mb-1">{post.title}</h3>
            <p className="text-ink-soft leading-relaxed">{post.description}</p>
            <span className="font-mono text-[11px] tracking-[0.1em] uppercase text-ink-faint mt-2 block">{post.readTime}</span>
          </Link>
        ))}
      </div>
    </section>
  )
}
