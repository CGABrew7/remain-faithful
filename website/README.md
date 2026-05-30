# Remain Faithful — Website

Marketing website for [remainfaithful.com](https://remainfaithful.com), built with Next.js 14 and Tailwind CSS.

## Tech Stack

- **Framework:** Next.js 14 (App Router)
- **Styling:** Tailwind CSS v3
- **Fonts:** Playfair Display (serif headings) + Inter (body)
- **Deployment:** Vercel
- **Backend:** Go API at `api.remainfaithful.com`

## Local Development

```bash
# 1. Install dependencies
npm install

# 2. Set up environment variables
cp .env.example .env.local
# Edit .env.local with your values (BACKEND_URL for local dev: http://localhost:8080)

# 3. Start the dev server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `STRIPE_SECRET_KEY` | Yes | Stripe secret key (server-side) for creating Checkout sessions |
| `NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY` | Yes | Stripe publishable key (client-side) |
| `BACKEND_URL` | No | Base URL of the Go backend API (contact form fallback) |
| `NEXT_PUBLIC_SITE_URL` | Yes | Canonical site URL (e.g. `https://remainfaithful.com`) |
| `NEXT_PUBLIC_GA_ID` | No | Google Analytics 4 Measurement ID (`G-XXXXXXXXXX`) |

## Deploying to Vercel

1. Push this folder to a GitHub repository (or use the Vercel CLI from this directory)
2. Import the project in [vercel.com/new](https://vercel.com/new)
3. Set root directory to `website/` (if deploying from the monorepo)
4. Add environment variables via the Vercel CLI or dashboard:
   ```bash
   vercel env add STRIPE_SECRET_KEY
   vercel env add NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY
   vercel env add NEXT_PUBLIC_SITE_URL
   vercel env add BACKEND_URL          # optional, for contact form forwarding
   vercel env add NEXT_PUBLIC_GA_ID    # optional, for analytics
   ```
5. Deploy

Alternatively, from this directory:

```bash
npm i -g vercel
vercel --prod
```

## Pages

| Route | Description |
|---|---|
| `/` | Homepage with hero, features, how it works, testimonials, download, donation CTA |
| `/how-it-works` | Detailed accountability model, monitoring pipeline, FAQ |
| `/privacy` | Full privacy policy |
| `/partners` | Ministry leader page with pilot program request form |
| `/about` | Mission, values, open source, contact form |
| `/blog` | Blog index |
| `/blog/[slug]` | Individual blog posts |

## API Routes

| Route | Method | Description |
|---|---|---|
| `/api/donate` | POST | Creates a Stripe Checkout session directly; accepts `{ amount, recurring }`, returns `{ url }` |
| `/api/contact` | POST | Forwards to backend `/contact` if `BACKEND_URL` is set; otherwise logs and returns success |

## Project Structure

```
src/
├── app/
│   ├── layout.tsx          # Root layout (Nav, Footer, fonts, GA)
│   ├── page.tsx            # Homepage
│   ├── how-it-works/
│   ├── privacy/
│   ├── partners/
│   ├── about/
│   ├── blog/
│   │   ├── page.tsx
│   │   ├── posts.ts        # Blog post data
│   │   └── [slug]/
│   └── api/
│       ├── donate/
│       └── contact/
└── components/
    ├── Nav.tsx
    ├── Footer.tsx
    ├── AppMockup.tsx       # SVG phone mockup (no external images)
    └── DonateButton.tsx    # Stripe checkout client component
```
