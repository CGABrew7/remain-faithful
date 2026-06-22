import Link from 'next/link'

interface RelatedPost {
  title: string
  slug: string
  description: string
  readTime: string
}

export function RelatedPosts({ posts }: { posts: RelatedPost[] }) {
  return (
    <section className="mt-12 pt-8 border-t border-[#1E3050]">
      <h2 className="font-serif text-xl font-semibold text-[#F0EDE8] mb-4">Related Articles</h2>
      <div className="grid gap-4 md:grid-cols-2">
        {posts.map((post) => (
          <Link
            key={post.slug}
            href={`/blog/${post.slug}`}
            className="block p-4 rounded-xl border border-[#1E3050] bg-[#162235] hover:border-[#C9A84C]/40 transition-colors"
          >
            <h3 className="font-medium text-[#F0EDE8] mb-1">{post.title}</h3>
            <p className="text-sm text-[#8A9BB0] mt-1 leading-relaxed">{post.description}</p>
            <span className="text-xs text-[#8A9BB0]/60 mt-2 block">{post.readTime}</span>
          </Link>
        ))}
      </div>
    </section>
  )
}
