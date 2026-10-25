'use client'

import { useEffect, useState } from 'react'
import type { AccountStatus, Count, Health, NotTracked, Overview, Percent } from '@/lib/admin-types'

const RANGES = [
  { id: '7d', label: '7 days' },
  { id: '30d', label: '30 days' },
  { id: '90d', label: '90 days' },
  { id: 'all', label: 'All' },
] as const

const LINKS = [
  { href: 'https://appstoreconnect.apple.com/analytics', label: 'App Store Connect Analytics', note: 'Downloads by country are here.' },
  { href: 'https://appstoreconnect.apple.com/apps/6774615426/appstore', label: 'Ratings and reviews', note: 'Open the app in App Store Connect, then Ratings and Reviews.' },
  { href: 'https://vercel.com/dashboard', label: 'Vercel', note: 'Website hosting.' },
  { href: 'https://fly.io/apps/remain-faithful-api', label: 'Fly API', note: 'The API machine, logs, and secrets.' },
  { href: 'https://supabase.com/dashboard', label: 'Supabase', note: 'Open this if the database is hosted there.' },
  { href: 'https://fly.io/dashboard', label: 'Fly database', note: 'Open this if the database is Fly Postgres.' },
  { href: 'https://app.sendgrid.com/', label: 'SendGrid', note: 'Email activity and billing.' },
]

const US_STATES: Record<string, string> = {
  AL: 'Alabama', AK: 'Alaska', AZ: 'Arizona', AR: 'Arkansas', CA: 'California', CO: 'Colorado',
  CT: 'Connecticut', DE: 'Delaware', DC: 'District of Columbia', FL: 'Florida', GA: 'Georgia',
  HI: 'Hawaii', ID: 'Idaho', IL: 'Illinois', IN: 'Indiana', IA: 'Iowa', KS: 'Kansas',
  KY: 'Kentucky', LA: 'Louisiana', ME: 'Maine', MD: 'Maryland', MA: 'Massachusetts', MI: 'Michigan',
  MN: 'Minnesota', MS: 'Mississippi', MO: 'Missouri', MT: 'Montana', NE: 'Nebraska', NV: 'Nevada',
  NH: 'New Hampshire', NJ: 'New Jersey', NM: 'New Mexico', NY: 'New York', NC: 'North Carolina',
  ND: 'North Dakota', OH: 'Ohio', OK: 'Oklahoma', OR: 'Oregon', PA: 'Pennsylvania', RI: 'Rhode Island',
  SC: 'South Carolina', SD: 'South Dakota', TN: 'Tennessee', TX: 'Texas', UT: 'Utah', VT: 'Vermont',
  VA: 'Virginia', WA: 'Washington', WV: 'West Virginia', WI: 'Wisconsin', WY: 'Wyoming',
}

function fmt(value: Count | number | null | undefined): string {
  if (value === null || value === undefined) return 'n/a'
  if (value === '<3') return '<3'
  return new Intl.NumberFormat('en-US').format(value)
}

function fmtPct(value: Percent): string {
  if (value === '<3') return '<3'
  if (value === null) return 'n/a'
  return `${value}%`
}

function countryName(code: string): string {
  if (code === '(direct)') return 'Direct or this site'
  try {
    return new Intl.DisplayNames(['en'], { type: 'region' }).of(code) || code
  } catch {
    return code
  }
}

function stateName(code: string): string {
  return US_STATES[code] || code
}

function formatBytes(n: number): string {
  if (!n) return '0 B'
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(0)} KB`
  if (n < 1024 * 1024 * 1024) return `${(n / (1024 * 1024)).toFixed(1)} MB`
  return `${(n / (1024 * 1024 * 1024)).toFixed(2)} GB`
}

function loginMessage(status: number, fallback: string) {
  if (status === 404) return 'Admin sign-in is not turned on yet. The setup note explains the one password step.'
  if (status === 429) return 'Too many tries. Wait 15 minutes and try again.'
  if (status === 401) return 'That password is not right.'
  if (status === 503) return 'The account service is not configured.'
  return fallback || 'Sign-in failed.'
}

export default function AdminDashboard({
  preview = false,
  ranges,
  healthPreview = null,
}: {
  preview?: boolean
  ranges?: Record<string, Overview>
  healthPreview?: Health | null
}) {
  const [range, setRange] = useState<(typeof RANGES)[number]['id']>('7d')
  const [overview, setOverview] = useState<Overview | null>(preview ? ranges?.['7d'] ?? null : null)
  const [health, setHealth] = useState<Health | null>(preview ? healthPreview : null)
  const [authed, setAuthed] = useState(preview)
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(!preview)
  const [busy, setBusy] = useState(false)
  const [checked, setChecked] = useState(preview)
  const [sessionCheck, setSessionCheck] = useState(0)

  useEffect(() => {
    if (preview) {
      setOverview(ranges?.[range] ?? null)
      return
    }
    let cancel = false
    setLoading(true)
    setError('')
    void (async () => {
      try {
        const [ov, he] = await Promise.all([
          fetch(`/api/admin/overview?range=${range}`, { cache: 'no-store' }),
          fetch('/api/admin/health', { cache: 'no-store' }),
        ])
        if (cancel) return
        if (ov.status === 401) {
          setAuthed(false)
          setOverview(null)
          setLoading(false)
          setChecked(true)
          return
        }
        if (!ov.ok) {
          const body = await ov.json().catch(() => ({}))
          setError(body.error || 'The dashboard could not load.')
          setLoading(false)
          setChecked(true)
          return
        }
        setOverview(await ov.json())
        setAuthed(true)
        if (he.ok) setHealth(await he.json())
        setLoading(false)
        setChecked(true)
      } catch {
        if (!cancel) {
          setError('We could not reach the server. Check your connection and try again.')
          setLoading(false)
          setChecked(true)
        }
      }
    })()
    return () => {
      cancel = true
    }
  }, [preview, range, ranges, sessionCheck])

  async function signIn(e: React.FormEvent) {
    e.preventDefault()
    setBusy(true)
    setError('')
    try {
      const res = await fetch('/api/admin/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password }),
      })
      if (!res.ok) {
        const body = await res.json().catch(() => ({}))
        setError(loginMessage(res.status, body.error || ''))
        return
      }
      setPassword('')
      setAuthed(true)
      setSessionCheck((n) => n + 1)
    } catch {
      setError('We could not reach the server. Check your connection and try again.')
    } finally {
      setBusy(false)
    }
  }

  async function signOut() {
    try {
      await fetch('/api/admin/logout', { method: 'POST' })
    } catch {
      setError('We could not reach the server. Check your connection and try again.')
    }
    setAuthed(false)
    setOverview(null)
  }

  return (
    <div className="mx-auto max-w-6xl px-4 pb-16 pt-8 sm:px-6 sm:pt-10">
      <header className="rise flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="kicker">Remain Faithful</p>
          <h1 className="mt-2 font-serif text-4xl text-white sm:text-5xl">Owner dashboard</h1>
          <p className="mt-2 max-w-xl text-sm text-ink-soft">
            Aggregates only. Counts under 3 people show as &lt;3. Nothing here lists a person&apos;s alerts, apps, or partners.
          </p>
        </div>
        {authed && !preview && (
          <button type="button" className="btn-ghost" onClick={signOut}>
            Sign out
          </button>
        )}
      </header>

      {preview && (
        <p className="rise rise-2 mt-6 rounded-[14px] bg-black/20 px-4 py-3 text-sm text-wax">
          Sample numbers for layout review. This page is not the live dashboard.
        </p>
      )}

      {!preview && !checked && <p className="mt-10 text-ink-soft">Checking your session...</p>}

      {checked && !authed && (
        <form onSubmit={signIn} className="surface rise rise-2 mx-auto mt-10 max-w-md p-6 sm:p-8">
          <h2 className="font-serif text-2xl text-white">Sign in</h2>
          <p className="mt-2 text-sm text-ink-soft">The password is the one set on the API. It is not an app account.</p>
          <label className="mt-6 block text-sm font-semibold text-white" htmlFor="admin-password">
            Password
          </label>
          <input
            id="admin-password"
            type="password"
            autoComplete="current-password"
            className="input-field mt-2"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          {error && (
            <p className="form-error mt-3 text-sm" role="alert">
              {error}
            </p>
          )}
          <button type="submit" className="btn-wax mt-5 w-full" disabled={busy}>
            {busy ? 'Checking...' : 'Sign in'}
          </button>
        </form>
      )}

      {authed && (
        <>
          <div className="rise rise-2 mt-8 flex flex-wrap gap-2" role="group" aria-label="Date range">
            {RANGES.map((item) => (
              <button
                key={item.id}
                type="button"
                className={item.id === range ? 'btn-wax !px-4 !py-2' : 'btn-ghost !px-4 !py-2'}
                aria-pressed={item.id === range}
                onClick={() => setRange(item.id)}
              >
                {item.label}
              </button>
            ))}
          </div>
          {error && (
            <p className="form-error mt-4 text-sm" role="alert">
              {error}
            </p>
          )}
          {loading && <p className="mt-8 text-ink-soft">Loading the latest counts...</p>}
          {overview && <DashboardBody overview={overview} health={health} preview={preview} />}
        </>
      )}
    </div>
  )
}

function DashboardBody({
  overview,
  health,
  preview,
}: {
  overview: Overview
  health: Health | null
  preview: boolean
}) {
  const w = overview.website
  return (
    <div className="mt-8 space-y-6">
      <Section title="Website" kicker="Visits" delay={0.05}>
        <Stats
          items={[
            ['Visits', w.visits],
            ['Unique visitors', w.visitors],
            ['Mobile visitors', w.mobile_visitors],
            ['Desktop visitors', w.desktop_visitors],
            ['App Store clicks', w.app_store_clicks],
            ['Click-through', fmtPct(w.app_store_ctr_percent)],
            ['Android signups', w.android_signups],
            ['Android waitlist', w.android_waitlist],
          ]}
        />
        <h3 className="mt-6 font-serif text-xl text-white">Visits by day</h3>
        <p className="mb-3 text-xs text-ink-faint">Dates are UTC.</p>
        <Trend points={w.trend} />
        <div className="mt-6 grid gap-6 lg:grid-cols-2">
          <div>
            <h3 className="mb-3 font-serif text-xl text-white">Top pages</h3>
            <SimpleTable
              headers={['Page', 'Visits', 'People']}
              rows={w.top_pages.map((row) => [row.path, fmt(row.visits), fmt(row.visitors)])}
            />
          </div>
          <div>
            <h3 className="mb-3 font-serif text-xl text-white">Top referrers</h3>
            <HBars rows={w.top_referrers.map((row) => ({ label: row.label === '(direct)' ? 'Direct or this site' : row.label, count: row.count }))} />
          </div>
        </div>
        <h3 className="mb-3 mt-6 font-serif text-xl text-white">Campaigns</h3>
        <SimpleTable
          headers={['Source', 'Medium', 'Campaign', 'People']}
          rows={w.campaigns.map((row) => [row.source || 'n/a', row.medium || 'n/a', row.campaign || 'n/a', fmt(row.visitors)])}
        />
        <div className="mt-6 grid gap-6 lg:grid-cols-2">
          <div>
            <h3 className="mb-3 font-serif text-xl text-white">Countries</h3>
            <HBars rows={w.countries.map((row) => ({ label: countryName(row.label), count: row.count }))} />
          </div>
          <div>
            <h3 className="mb-3 font-serif text-xl text-white">US states</h3>
            <HBars rows={w.us_states.map((row) => ({ label: stateName(row.label), count: row.count }))} />
          </div>
        </div>
      </Section>

      <Section title="Signups and growth" kicker="Accounts" delay={0.1}>
        <Stats
          items={[
            ['Total accounts', overview.growth.total_accounts],
            ['New in this range', overview.growth.new_accounts],
            ['Email verified', fmtPct(overview.growth.email_verified_percent)],
          ]}
        />
        <div className="mt-6 grid gap-6 lg:grid-cols-2">
          <div>
            <h3 className="mb-3 font-serif text-xl text-white">Sign-in method</h3>
            <HBars rows={overview.growth.sign_in.map((row) => ({ label: row.method, count: row.count }))} />
          </div>
          <div>
            <h3 className="mb-3 font-serif text-xl text-white">New accounts by week</h3>
            <SimpleTable
              headers={['Week', 'Accounts']}
              rows={overview.growth.by_week.map((row) => [row.day, fmt(row.count)])}
            />
          </div>
        </div>
        <h3 className="mb-3 mt-6 font-serif text-xl text-white">New accounts by day</h3>
        <DayBars points={overview.growth.by_day} />
        <div className="mt-4">
          <GapNote gap={overview.growth.deleted_accounts} />
        </div>
      </Section>

      <Section title="Activation funnel" kicker="From signup onward" delay={0.12}>
        <ol className="space-y-3">
          {overview.funnel.map((step) => (
            <li key={step.key} className="rounded-[14px] bg-black/15 px-4 py-3">
              <div className="flex flex-wrap items-baseline justify-between gap-2">
                <p className="font-semibold text-white">{step.label}</p>
                <p className="tabular-nums text-wax">
                  {step.status === 'not_tracked_yet' ? 'Not tracked yet' : fmt(step.count)}
                </p>
              </div>
              {step.dropoff_percent !== undefined && (
                <p className="mt-1 text-sm text-ink-faint">Drop-off from the previous tracked step: {fmtPct(step.dropoff_percent)}</p>
              )}
              {step.note && <p className="mt-1 text-sm text-ink-soft">{step.note}</p>}
            </li>
          ))}
        </ol>
      </Section>

      <Section title="Active use" kicker="Heartbeats" delay={0.14}>
        <p className="mb-4 text-sm text-ink-soft">{overview.active.basis}</p>
        <Stats
          items={[
            ['Daily active', overview.active.dau],
            ['Weekly active', overview.active.wau],
            ['Monthly active', overview.active.mau],
            ['Monitoring on', overview.active.monitoring_on],
            ['Monitoring off or paused', overview.active.monitoring_off],
            ['Heartbeat silent', overview.active.heartbeat_silent],
            ['Never checked in', overview.active.never_checked_in],
            ['Monitored apps', overview.active.monitored_apps],
          ]}
        />
        <h3 className="mb-3 mt-6 font-serif text-xl text-white">App versions</h3>
        {overview.active.app_versions.status === 'tracked' ? (
          (overview.active.app_versions.versions ?? []).length === 0 ? (
            <p className="text-sm text-ink-faint">No app versions reported yet.</p>
          ) : (
            <HBars rows={(overview.active.app_versions.versions ?? []).map((row) => ({ label: row.label, count: row.count }))} />
          )
        ) : (
          <GapNote gap={{ status: 'not_tracked_yet', note: overview.active.app_versions.note ?? '' }} />
        )}
        <div className="mt-4 space-y-3">
          <GapNote gap={overview.active.daily_history} />
        </div>
      </Section>

      <Section title="Accountability" kicker="Partners and groups" delay={0.16}>
        <Stats
          items={[
            ['Invites sent', overview.accountability.invites_sent],
            ['Acceptance', fmtPct(overview.accountability.acceptance_percent)],
            ['Groups', overview.accountability.groups],
            ['Average group size', overview.accountability.average_group_size === null ? 'n/a' : fmt(overview.accountability.average_group_size)],
            ['No partner yet', overview.accountability.users_with_zero_partners],
          ]}
        />
        <h3 className="mb-3 mt-6 font-serif text-xl text-white">Partnerships right now</h3>
        <HBars rows={overview.accountability.partnerships.map((row) => ({ label: row.status, count: row.count }))} />
        <h3 className="mb-3 mt-6 font-serif text-xl text-white">Group size</h3>
        <HBars rows={overview.accountability.group_sizes.map((row) => ({ label: `${row.size} people`, count: row.groups }))} />
        <div className="mt-4">
          <GapNote gap={overview.accountability.median_acceptance} />
        </div>
      </Section>

      <Section title="Alerts and safety" kicker="Counts only" delay={0.18}>
        <p className="mb-4 text-sm text-ink-soft">
          These are daily totals. There is no split by person, category, or app.
        </p>
        <h3 className="mb-3 font-serif text-xl text-white">Partner alerts by day</h3>
        <DayBars points={overview.alerts.partner_alerts_by_day} />
        <h3 className="mb-3 mt-6 font-serif text-xl text-white">Monitoring paused or turned off</h3>
        <DayBars points={overview.alerts.monitoring_paused_by_day} />
        <h3 className="mb-3 mt-6 font-serif text-xl text-white">I need support</h3>
        {overview.alerts.support_taps.status === 'tracked' ? (
          <p className="text-sm text-ink-soft">
            Taps saved: <span className="tabular-nums text-wax">{fmt(overview.alerts.support_taps.count)}</span>
          </p>
        ) : (
          <GapNote gap={{ status: 'not_tracked_yet', note: overview.alerts.support_taps.note ?? '' }} />
        )}
      </Section>

      <Section title="Notifications and email" kicker="Delivery" delay={0.2}>
        <Stats
          items={[
            ['Active push tokens', overview.notifications.push_tokens],
            ['People with a token', overview.notifications.push_people],
          ]}
        />
        <ul className="mt-4 space-y-3">
          {overview.notifications.emails.map((row) => (
            <li key={row.type} className="rounded-[14px] bg-black/15 px-4 py-3 text-sm">
              <div className="flex flex-wrap items-baseline justify-between gap-2">
                <p className="font-semibold capitalize text-white">{row.type.replace(/_/g, ' ')}</p>
                <p className="tabular-nums text-wax">
                  {row.status === 'not_tracked_yet' ? 'Not tracked yet' : fmt(row.count)}
                </p>
              </div>
              {row.note && <p className="mt-1 text-ink-soft">{row.note}</p>}
            </li>
          ))}
        </ul>
        <div className="mt-4">
          <GapNote gap={overview.notifications.apns} />
        </div>
      </Section>

      <Section title="Health and cost" kicker="Where the bills live" delay={0.22}>
        {health ? (
          <>
            <Stats
              items={[
                ['API', health.api],
                ['Database', health.database],
                ['Database size', formatBytes(health.db_bytes)],
              ]}
            />
            <h3 className="mb-3 mt-6 font-serif text-xl text-white">Rows</h3>
            <SimpleTable
              headers={['Table', 'Rows']}
              rows={health.tables.map((row) => [row.table, fmt(row.rows)])}
            />
            <div className="mt-4">
              <GapNote gap={health.server_errors} />
            </div>
          </>
        ) : (
          <p className="text-sm text-ink-soft">Health figures are not loaded.</p>
        )}
        <p className="mt-4 text-sm text-ink-soft">
          App Store downloads by country are in App Store Connect Analytics, not on this page.
          This dashboard reads the database you already run. It does not add a paid analytics product.
          Old page counts are deleted after 13 months.
        </p>
        <ul className="mt-4 grid gap-3 sm:grid-cols-2">
          {LINKS.map((link) => (
            <li key={link.href}>
              <a
                href={link.href}
                target="_blank"
                rel="noopener noreferrer"
                className="surface-link surface"
              >
                <span className="font-semibold text-white">{link.label}</span>
                <span className="mt-1 text-sm text-ink-soft">{link.note}</span>
              </a>
            </li>
          ))}
        </ul>
      </Section>

      <Section title="Feedback" kicker="Ideas and problems" delay={0.24}>
        <GapNote gap={overview.feedback} />
      </Section>

      <Section title="Support tools" kicker="One account at a time" delay={0.26}>
        <SupportTools preview={preview} />
      </Section>
    </div>
  )
}

function SupportTools({ preview }: { preview: boolean }) {
  const [email, setEmail] = useState(preview ? 'person@example.com' : '')
  const [confirm, setConfirm] = useState('')
  const [account, setAccount] = useState<AccountStatus | null>(
    preview
      ? {
          email: 'person@example.com',
          created_at: '2026-04-02T00:00:00Z',
          email_verified: true,
          sign_in_method: 'apple',
          last_active_at: '2026-10-17T15:00:00Z',
          partner_count: 2,
        }
      : null,
  )
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  async function post(path: string, body: unknown) {
    if (preview) {
      setError('This sample view does not change accounts.')
      return null
    }
    setBusy(true)
    setError('')
    setMessage('')
    try {
      const res = await fetch(path, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      const data = await res.json().catch(() => ({}))
      if (!res.ok) {
        setError(data.error || 'That did not work.')
        return null
      }
      return data
    } catch {
      setError('We could not reach the server. Check your connection and try again.')
      return null
    } finally {
      setBusy(false)
    }
  }

  return (
    <div>
      <p className="text-sm text-ink-soft">
        Look up an exact email. You will see account status only: created date, verified, sign-in method, last active date, and a partner count.
      </p>
      <form
        className="mt-4 flex flex-col gap-3 sm:flex-row"
        onSubmit={async (e) => {
          e.preventDefault()
          if (preview) {
            setError('This sample view does not look up live accounts.')
            return
          }
          setBusy(true)
          setError('')
          setMessage('')
          try {
            const res = await fetch(`/api/admin/account?email=${encodeURIComponent(email.trim())}`, { cache: 'no-store' })
            const data = await res.json().catch(() => ({}))
            if (!res.ok) {
              setAccount(null)
              setError(data.error || 'That did not work.')
              return
            }
            setAccount(data as AccountStatus)
          } catch {
            setAccount(null)
            setError('We could not reach the server. Check your connection and try again.')
          } finally {
            setBusy(false)
          }
        }}
      >
        <label className="sr-only" htmlFor="lookup-email">
          Email
        </label>
        <input
          id="lookup-email"
          type="email"
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="name@email.com"
          className="input-field sm:flex-1"
        />
        <button type="submit" className="btn-wax" disabled={busy}>
          Look up
        </button>
      </form>
      {account && (
        <dl className="mt-4 grid grid-cols-2 gap-3 text-sm sm:grid-cols-3">
          <Stat label="Email" value={account.email} />
          <Stat label="Created" value={account.created_at.slice(0, 10)} />
          <Stat label="Verified" value={account.email_verified ? 'Yes' : 'No'} />
          <Stat label="Sign-in" value={account.sign_in_method} />
          <Stat label="Last active" value={account.last_active_at ? account.last_active_at.slice(0, 10) : 'No heartbeat yet'} />
          <Stat label="Partners" value={String(account.partner_count)} />
        </dl>
      )}
      <div className="mt-4 flex flex-col gap-2 sm:flex-row">
        <button
          type="button"
          className="btn-ghost"
          disabled={busy || !email}
          onClick={async () => {
            const data = await post('/api/admin/account/resend-verification', { email })
            if (data?.status === 'sent') setMessage('Verification email sent.')
            if (data?.status === 'already_verified') setMessage('That email is already verified.')
          }}
        >
          Resend verification
        </button>
        <button
          type="button"
          className="btn-ghost"
          disabled={busy || !email}
          onClick={async () => {
            const data = await post('/api/admin/account/password-reset', { email })
            if (data?.status === 'sent') setMessage('Password reset email sent.')
          }}
        >
          Send password reset
        </button>
      </div>
      <form
        className="mt-6 rounded-[14px] bg-black/15 p-4"
        onSubmit={async (e) => {
          e.preventDefault()
          const data = await post('/api/admin/account/delete', { email, confirm })
          if (data?.status === 'deleted') {
            setMessage('Account deleted.')
            setAccount(null)
            setConfirm('')
          }
        }}
      >
        <p className="text-sm text-white">Delete this account</p>
        <p className="mt-1 text-sm text-ink-soft">Type the same email to confirm. This uses the normal account deletion path and is written to the audit log.</p>
        <label className="sr-only" htmlFor="confirm-email">
          Type the email to confirm
        </label>
        <input
          id="confirm-email"
          type="email"
          value={confirm}
          onChange={(e) => setConfirm(e.target.value)}
          placeholder="Type the email to confirm"
          className="input-field mt-3"
        />
        <button type="submit" className="btn-danger mt-3" disabled={busy || confirm.trim().toLowerCase() !== email.trim().toLowerCase()}>
          Delete account
        </button>
      </form>
      {message && <p className="mt-3 text-sm text-wax">{message}</p>}
      {error && (
        <p className="form-error mt-3 text-sm" role="alert">
          {error}
        </p>
      )}
    </div>
  )
}

function Section({
  title,
  kicker,
  delay,
  children,
}: {
  title: string
  kicker: string
  delay: number
  children: React.ReactNode
}) {
  return (
    <section className="surface rise p-5 sm:p-6" style={{ animationDelay: `${delay}s` }}>
      <p className="kicker">{kicker}</p>
      <h2 className="mt-2 font-serif text-2xl text-white sm:text-3xl">{title}</h2>
      <div className="mt-4">{children}</div>
    </section>
  )
}

function Stats({ items }: { items: [string, Count | string][] }) {
  return (
    <dl className="grid grid-cols-2 gap-3 lg:grid-cols-4">
      {items.map(([label, value]) => (
        <div key={label} className="rounded-[14px] bg-black/15 px-3 py-3">
          <dt className="text-xs uppercase tracking-wide text-ink-faint">{label}</dt>
          <dd className="mt-1 font-serif text-2xl tabular-nums text-white">{typeof value === 'string' ? value : fmt(value)}</dd>
        </div>
      ))}
    </dl>
  )
}

function Stat({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs uppercase tracking-wide text-ink-faint">{label}</dt>
      <dd className="mt-1 break-all text-white">{value}</dd>
    </div>
  )
}

function GapNote({ gap }: { gap: NotTracked }) {
  return (
    <p className="rounded-[14px] bg-black/15 px-3 py-3 text-sm text-ink-soft">
      <span className="font-semibold text-wax">Not tracked yet. </span>
      {gap.note}
    </p>
  )
}

function HBars({ rows }: { rows: { label: string; count: Count }[] }) {
  if (rows.length === 0) return <p className="text-sm text-ink-faint">Nothing in this range.</p>
  const max = Math.max(1, ...rows.map((row) => (typeof row.count === 'number' ? row.count : 0)))
  return (
    <ul className="space-y-2">
      {rows.map((row) => {
        const n = typeof row.count === 'number' ? row.count : 0
        const width = row.count === '<3' ? 8 : Math.max(4, (n / max) * 100)
        return (
          <li key={row.label} className="grid grid-cols-[minmax(5.5rem,9rem)_1fr_2.5rem] items-center gap-2 text-sm">
            <span className="truncate text-ink-soft">{row.label}</span>
            <span className="h-2 overflow-hidden rounded-full bg-white/10">
              <span className="block h-full rounded-full bg-wax" style={{ width: `${width}%` }} />
            </span>
            <span className="text-right tabular-nums text-white">{fmt(row.count)}</span>
          </li>
        )
      })}
    </ul>
  )
}

function Trend({ points }: { points: Overview['website']['trend'] }) {
  return <DayBars points={points.map((point) => ({ day: point.day, count: point.visits }))} />
}

function DayBars({ points }: { points: { day: string; count: Count }[] }) {
  if (points.length === 0) return <p className="text-sm text-ink-faint">Nothing in this range.</p>
  const max = Math.max(1, ...points.map((point) => (typeof point.count === 'number' ? point.count : 0)))
  return (
    <div className="flex items-end gap-1 overflow-x-auto pb-1" style={{ minHeight: 148 }}>
      {points.map((point) => {
        const n = typeof point.count === 'number' ? point.count : 0
        const height = Math.max(6, Math.round((n / max) * 96))
        return (
          <div key={point.day} className="flex w-8 shrink-0 flex-col items-center gap-1">
            <span className="text-[10px] tabular-nums text-ink-faint">{fmt(point.count)}</span>
            <div className="w-full rounded-t-md bg-wax" style={{ height }} />
            <span className="text-[10px] text-ink-faint">{point.day.slice(5)}</span>
          </div>
        )
      })}
    </div>
  )
}

function SimpleTable({ headers, rows }: { headers: string[]; rows: string[][] }) {
  if (rows.length === 0) return <p className="text-sm text-ink-faint">Nothing in this range.</p>
  return (
    <div className="overflow-x-auto">
      <table className="w-full text-left text-sm">
        <thead>
          <tr className="text-ink-faint">
            {headers.map((header) => (
              <th key={header} className="py-2 pr-3 font-medium">
                {header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, index) => (
            <tr key={`${row[0]}-${index}`} className="border-t border-white/10">
              {row.map((cell, cellIndex) => (
                <td key={cellIndex} className={`py-2 pr-3 ${cellIndex === 0 ? 'text-white' : 'tabular-nums text-ink-soft'}`}>
                  {cell}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
