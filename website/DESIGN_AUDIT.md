# Design audit — remainfaithful.com

Audited on 7 Oct 2026 at desktop 1440 and mobile 390, against the live layout after PR #23 (`331d4f1`, navy/gold token swap) and the pre-letter design at `ee5b5c8^` (21 Aug 2026).

ICP: an adult (18+) Christian who wants peer accountability for purity — a spouse, friend, mentor, or small-group member. The job of the homepage is one action: **Get started** on iPhone. The sentence they must believe: filtering and partner alerts work, and screenshots stay on the device.

Aesthetic lane: precise product, navy chapel. Not the parchment “quiet letter,” and not a beige fleet site. References for structure only: the iOS app (shield, gold pill, leaf rule, translucent cards), a product-first SaaS hero (device beside the claim), and a calm editorial footer (type and hairlines, no illustration).

## What the color swap did not fix

PR #23 recolored the August letter restyle. It did not bring back layout. `ee5b5c8` removed the split hero, the device, gradient feature cards, the numbered step rail, the two-column privacy table, the blurred nav, and the wide `max-w-7xl` measure. The fake in-app partners (“Mike R.”, “David C.”) and the pre-launch “Join the Waitlist” primary CTA in that older page are not worth restoring. The structure is.

## Cross-site issues

- **Hierarchy.** Headlines, kickers, and body sat at nearly the same volume. The homepage h1 was the product name inside a narrow letter column, with three equal cards shouting beside it. One loud thing was missing.
- **Spacing.** Section rhythm was a repeated `hr` plus `page-section`. Gaps inside cards and between sections did not share one scale. The letter column (`min(64rem)`) left a 1440px viewport feeling empty rather than generous.
- **Type.** Playfair + Inter. Inter is the default grotesque the restyle had tried to escape, then the color swap brought back. Nav and footer used mono tracking with no mono family loaded, so labels looked like a leftover system. Body measure on the homepage was ~42ch (good) but inner pages ran edge to edge of `max-w-4xl` with no display/body jump.
- **Contrast.** White and `#C5CBE0` on `#12214C` pass WCAG AA. Gold `#D1AB4C` on navy passes for text. `#6B7799` (the muted tone in the brief) is about 3.5:1 on `#12214C` and fails AA for small text — implemented muted is `#8E98B8` (~5.4:1). Hairlines were gold at 22% opacity and disappeared on the gradient.
- **Card depth.** Proof cards, method panes, and inner-page `rounded-sm` panels were flat fills (`#1C3066`) with a small radius (2px) next to pill buttons. No shared surface, no light hairline, no inset highlight.
- **Section variety.** Home, how-it-works, partners, compare, and the SEO landings repeat the same centered kicker + serif h2 + equal card grid. Partners is four identical cards. Blog is three identical cards with a gold strip.
- **Imagery.** No product. The only picture of the app was deleted in `ee5b5c8` (`AppMockup.tsx`), and that mock invented partners and a streak. After the color swap the site is type on a gradient.
- **CTAs.** Homepage label is “Get started.” Many inner pages used square gold buttons with the same words, a second style from the letter system. Footer had no primary action, only a newsletter.
- **Nav.** Fixed, but transparent over the hero and solid after scroll, with no blur. Width capped at `max-w-5xl`. Mobile menu works. Logo was a gold shield after #23; good. Contact and About both highlighted on `/about` because the contact href is `/about#contact`.
- **Footer.** Four link columns plus a “Letter” kicker, a burgundy-era mark replaced by gold, and a mono copyright line. No decorative doodle (good). It did not read as a finished close: brand, columns, one quiet CTA, hairline.
- **Motion.** Homepage pointer-tilt on the letter sheet was the motion. Reduced motion turned the tilt off. Buttons scaled to 0.96. Nothing else acknowledged enter or hover with intent.
- **Mobile.** Single column, no overflow on the sampled routes, hamburger present, button height 44px. Footer text links were under 44px. The homepage first screen was a heading and two buttons with no product, and the three cards pushed the story down.
- **Accessibility.** One h1 per page. Focus rings existed on a few controls, not as a global `:focus-visible`. FAQ answers on how-it-works used `max-h-96`, which clips long answers. Inputs were 14px (iOS zoom). Skip link was absent.
- **Performance.** No images, so the page was light but empty. Fonts are `next/font` (self-hosted). CSP `font-src 'self'` still holds. Security headers in `next.config.js` were not part of the visual problem and must stay.

## By route

| Route | Desktop 1440 | Mobile 390 |
| --- | --- | --- |
| `/` | Letter column, wax-seal dot, three equal proof cards, flat rules, download card with no device, FAQ as a plain list. First screen does not show the product. | Same stack. CTA is reachable. No imagery. Footer links are small. |
| `/how-it-works` | Centered headline, then a repeated two-column of text plus three stacked bars, then more identical step cards. Green “layer” chips are the only color besides gold. | Cards become a tall stack. Hierarchy flattens. |
| `/partners` | Centered hero, four equal benefit cards, then more equal rows. Square CTAs. | Four cards stack with no change in shape. |
| `/privacy-architecture` | Long technical article. Hero is an icon in a square. Sections are the same bordered panel. | Readable, visually monotonous. Tables/lists need the shared surface so they are not hairline-only. |
| `/group-setup-guide` | A printable document dropped on the marketing chrome. Fine for print, thin as a screen page. | Print button floats. Body is one narrow column, which is correct for a guide. |
| `/about` | Split hero is the best inner layout, then a flat quote and a long biography with a square “JB” monogram. | Quote card has little padding rhythm. |
| `/privacy`, `/terms` | Legal type on the gradient. Hairline section rules only. No measure or surface for the commitment statement. | Same. Links are gold and readable. |
| `/blog` | Three equal cards, gold top bar, “Read More.” Template rhythm. | One column of the same card. |
| `/blog/covenant-model` | Article plus a square subscribe button. Related posts are hairline rows. | Comfortable to read. CTA style disagreed with the homepage pill. |
| `/church-accountability`, `/christian-accountability-app`, `/free-accountability-app`, `/what-is-accountability-software` | Same SEO hero: kicker, h1, paragraph, square button, then repeated sections. | Usable. Visually the same template. |
| `/compare/covenant-eyes`, `/compare/ever-accountable`, `/compare/accountable2you` | Comparison tables and FAQ blocks in the letter palette. Honest copy, flat presentation. | Tables need to stay inside the shell. |

Donate is `/#donate`. Contact is `/about#contact`. Both lived in the same flat system.

## What this pass keeps from before `ee5b5c8`

Split hero, real device frames, gold shield, pill primary button, numbered method, two-column “on device / with partners,” blurred sticky nav, wider content column, gradient glow. Copy, routes, privacy claims, and donation amounts stay as they are after the honesty passes. The leaf rule is the app’s divider, used between type, not as a footer illustration.
