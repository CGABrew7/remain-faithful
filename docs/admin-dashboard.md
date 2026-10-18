# Owner dashboard

This page is for the owner of Remain Faithful. It lives at:

https://www.remainfaithful.com/admin

It is not linked from the public site. Search engines are told not to list it. The page shows totals only. A count that covers fewer than 3 people is shown as `<3`. It never lists one person's alerts, apps, monitoring history, or partners.

Until the password step below is done, every admin address answers "not found." Nobody can open the numbers.

## Set the password

Do this once, after this version of the API is running.

You need the Fly command line tool, already logged in, and Go. If you do not have those, ask someone who does to run the hash command and send you only the hash. The hash starts with `$2a$`. Do not send them the password in chat if you can avoid it. They can run the command while you type.

1. Open a terminal in the `backend` folder of this repository.
2. Run:

```
go run ./cmd/hashpw
```

Type a password of at least 10 characters and at most 72. The letters stay hidden. The command prints one line. That line is the hash, not the password. Copy the whole line, including `$2a$`.

You can also run `go run ./cmd/hashpw 'your password'`. That puts the password in the terminal history, so the hidden prompt is better. If the password contains a space, the hidden prompt is the one to use.

3. Paste the hash into this command. Keep the single quotes. The hash contains `$` signs, and the shell will break the command without those quotes.

```
fly secrets set ADMIN_PASSWORD_HASH='paste-the-hash-here' -a remain-faithful-api
```

Fly restarts the API. Open https://www.remainfaithful.com/admin and sign in with the password, not the hash. A sign-in lasts 8 hours. Five wrong tries in 15 minutes locks that network address for the rest of the window.

To change the password later, run the hash command again and set the secret again.

## Optional: keep unique visitor counts stable

The website counts pages without cookies. Each visitor id changes every day on its own. If `ANALYTICS_SALT` is not set on Vercel, a restart of the site can also change that id, and the same person can look like a new visitor.

1. Run `openssl rand -hex 32`.
2. In the Vercel project for the website, add `ANALYTICS_SALT` with that value, for Production.
3. Redeploy the website.

You do not need a new signing secret. The admin session uses the `JWT_SECRET` that is already on Fly. `ADMIN_SESSION_SECRET` is an optional override only if you later want a separate key.

## What you will see

Date range buttons: 7 days, 30 days, 90 days, and all. Dates are UTC.

Website: visits, unique visitors, top pages, top referrers, campaign tags, country, US state, mobile and desktop visitors, App Store button clicks, click-through rate, Android update signups, and visits by day.

Signups: total accounts, new accounts by day and by week, sign-in method (email, Apple, or Google), and the share of email accounts that verified.

Activation: account created, invited a partner, partner accepted, and monitoring has checked in (a heartbeat has arrived at least once), with drop-off between those steps.

Active use, as of right now (not filtered by the date range): daily, weekly, and monthly active users from the latest heartbeat, monitoring on, monitoring off or paused, heartbeat silent (last heartbeat older than 30 minutes), and never checked in.

Accountability: partnerships that are pending, accepted, or declined, invites sent, acceptance rate, number of groups, average group size, group sizes from 1 to 12, and accounts with no partner.

Alerts: partner alerts per day, and monitoring-paused or turned-off notices per day. No split by person, category, or app.

Notifications: active push tokens, people with a token, and partner and group invite emails attempted (one row means one send was attempted, not that SendGrid confirmed delivery).

Health: API status, database status, database size, and row counts. Links to App Store Connect Analytics, ratings and reviews, Vercel, Fly, Supabase, Fly Postgres, and SendGrid.

Support: look up one exact email. You see created date, whether email is verified, sign-in method, last active date, and a partner count. You can resend verification, send a password reset, or delete the account. Delete asks you to type the email again. The delete uses the same path as the in-app account deletion.

## Not tracked yet

These rows stay on the page so the gap is obvious. A later app build can fill some of them. This dashboard does not change the iPhone app.

- Deleted accounts. Accounts are removed, not marked deleted, so there is no total to count.
- Screen Time granted. The phone does not tell the server. A later build could send yes or no only, with no app names.
- Apps chosen. A later build could send a count only, never the app names.
- Daily active users for past days. Only the latest heartbeat is stored.
- App version. A later build could send the version string alone on heartbeat or sign-in.
- Median time until an invite is accepted. Invites do not store the accept time.
- "I need support" taps. A later build could send a count only, with no message and no partner names.
- Apple push success and failure counts. Sends are not written to the database.
- SendGrid verification and password-reset totals. Those emails are not logged. Invite emails are counted as attempts, above.
- Feedback messages. The website form and the in-app form email the message. They are not saved, so there is no list here.
- Recent server error count. The API does not keep that counter.
- App Store downloads by country. Those stay in App Store Connect Analytics. The health section links there.

## What this costs

No new paid product. Page counts are rows in the Postgres database you already run. Rows older than 13 months are deleted once a day. Each page view calls one website function, which writes one row to the API. The dashboard numbers are cached for about 60 seconds so opening the page does not hammer the small Fly machine. App users still send no location.

## Privacy rules this page follows

Visitor country and US state come from Vercel headers at the moment of the page view. The IP address is not stored. The visitor id is a hash of a salt, the UTC date, a shortened browser string, and the first language. It is not the IP. App accounts have no location on this page. Do Not Track and Global Privacy Control skip the count. Admin pages are not counted.
