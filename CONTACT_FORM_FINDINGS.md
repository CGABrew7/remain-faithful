# Contact Form — Root Cause & Fix

## Root Cause

The `/api/contact` Next.js route always proxied to `${BACKEND_URL}/contact` when
`BACKEND_URL` is set. The Go backend had **no `/contact` endpoint**, so every
production submission received a 404 from the backend → the route returned a 502 →
`ContactForm.tsx` swallowed the HTTP error and showed the generic string
"Something went wrong."

`BACKEND_URL` is set in Vercel for all environments, so this affected 100% of
production submissions.

## What Was Fixed (code changes only)

| File | Change |
|------|--------|
| `backend/internal/email/email.go` | Added `SendContact(fromEmail, fromName, subject, message, toEmail)` method |
| `backend/internal/handler/handler.go` | Added `SendContact` to `EmailSender` interface |
| `backend/internal/handler/contact.go` | New `POST /contact` handler (unauthenticated) |
| `backend/cmd/server/main.go` | Wired `/contact` route before the auth middleware |
| `website/src/app/api/contact/route.ts` | Parses and forwards real error body from backend; better logging |
| `website/src/components/ContactForm.tsx` | Parses actual error message from API response instead of catch-all string |

## Email Delivery

The Go backend sends contact submissions via SendGrid (`SENDGRID_API_KEY` is
already deployed on Fly.io). Emails are forwarded to `jeff@hanokventures.co`
(overridable via `CONTACT_TO_EMAIL` env var on the backend if you want to change
the destination without a code deploy).

## Env Vars Required

All required env vars are already deployed:

| Service | Var | Status |
|---------|-----|--------|
| Fly.io backend | `SENDGRID_API_KEY` | ✅ Deployed |
| Fly.io backend | `SENDGRID_FROM_EMAIL` | Not set — uses default `noreply@remainfaithful.app` |
| Vercel website | `BACKEND_URL` | ✅ Deployed (points to Fly.io API) |

**No new credentials need to be added.**

If you want a custom sender address (e.g. `hello@remainfaithful.com`), set
`SENDGRID_FROM_EMAIL` on Fly.io and verify that address in your SendGrid account.
