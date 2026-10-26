import type { Health, Overview } from '@/lib/admin-types'

const END = new Date('2026-10-18T12:00:00Z')

function days(span: number) {
  const out: string[] = []
  for (let i = span - 1; i >= 0; i--) {
    const d = new Date(END)
    d.setUTCDate(d.getUTCDate() - i)
    out.push(d.toISOString().slice(0, 10))
  }
  return out
}

function spanFor(range: string) {
  if (range === '30d') return 30
  if (range === '90d') return 90
  if (range === 'all') return 14
  return 7
}

export function previewOverview(range = '7d'): Overview {
  const labels = days(spanFor(range))
  return {
    range,
    generated_at: END.toISOString(),
    website: {
      visits: 1280,
      visitors: 640,
      top_pages: [
        { path: '/', visits: 540, visitors: 310 },
        { path: '/how-it-works', visits: 220, visitors: 140 },
        { path: '/privacy', visits: 80, visitors: 60 },
        { path: '/partners', visits: 40, visitors: 28 },
      ],
      top_referrers: [
        { label: '(direct)', count: 260 },
        { label: 'google.com', count: 180 },
        { label: 'news.ycombinator.com', count: 24 },
        { label: 'example.org', count: '<3' },
      ],
      campaigns: [
        { source: 'newsletter', medium: 'email', campaign: 'october', visitors: 36 },
        { source: 'partner', medium: 'social', campaign: 'launch', visitors: '<3' },
      ],
      countries: [
        { label: 'US', count: 410 },
        { label: 'CA', count: 48 },
        { label: 'GB', count: 30 },
        { label: 'AU', count: '<3' },
      ],
      us_states: [
        { label: 'TX', count: 70 },
        { label: 'CA', count: 54 },
        { label: 'FL', count: 33 },
        { label: 'NY', count: 21 },
        { label: 'OH', count: '<3' },
      ],
      mobile_visitors: 380,
      desktop_visitors: 260,
      app_store_clicks: 96,
      app_store_ctr_percent: 7.5,
      android_signups: 18,
      android_waitlist: 7,
      trend: labels.map((day, i) => ({
        day,
        visits: 140 + (i % 5) * 18,
        visitors: 70 + (i % 4) * 9,
      })),
    },
    growth: {
      total_accounts: 86,
      new_accounts: 12,
      by_day: labels.map((day, i) => ({ day, count: i % 3 === 0 ? '<3' : 3 + (i % 2) })),
      by_week: [
        { day: '2026-10-05', count: 6 },
        { day: '2026-10-12', count: 6 },
      ],
      sign_in: [
        { method: 'email', count: 40 },
        { method: 'apple', count: 44 },
        { method: 'google', count: '<3' },
      ],
      email_verified_percent: 82.5,
      deleted_accounts: {
        status: 'not_tracked_yet',
        note: 'Accounts are removed rather than marked deleted, so a deleted total is not kept.',
      },
    },
    funnel: [
      { key: 'account_created', label: 'Account created', status: 'tracked', count: 12 },
      { key: 'invited_partner', label: 'Invited a partner', status: 'tracked', count: 7, dropoff_percent: 41.7 },
      { key: 'partner_accepted', label: 'Partner accepted', status: 'tracked', count: 4, dropoff_percent: 42.9 },
      {
        key: 'screen_time',
        label: 'Screen Time granted',
        status: 'tracked',
        count: 9,
        note: 'People who allowed Screen Time.',
      },
      {
        key: 'monitoring_active',
        label: 'Monitoring has checked in',
        status: 'tracked',
        count: 4,
        dropoff_percent: 0,
        note: 'Counted when a heartbeat has arrived at least once. That means monitoring has checked in. It is not a separate Screen Time flag.',
      },
      {
        key: 'apps_chosen',
        label: 'Apps chosen',
        status: 'tracked',
        count: 6,
        note: 'People who chose at least one app. App names stay on the phone.',
      },
    ],
    active: {
      basis:
        'These figures use the latest heartbeat on each account. They describe right now, not the date range above. Monitoring on means a heartbeat in the last 30 minutes. Heartbeat silent means the latest heartbeat is older than that. Monitoring off means a partner was told monitoring stopped and no heartbeat has arrived since.',
      dau: 22,
      wau: 41,
      mau: 63,
      monitoring_on: 11,
      monitoring_off: 6,
      heartbeat_silent: 30,
      never_checked_in: 23,
      monitored_apps: 18,
      daily_history: {
        status: 'not_tracked_yet',
        note: 'Only the latest heartbeat is stored, so this dashboard cannot draw active users for each past day.',
      },
      app_versions: {
        status: 'tracked',
        versions: [
          { label: '1.0.2', count: 14 },
          { label: '1.0.1', count: 4 },
        ],
      },
    },
    accountability: {
      partnerships: [
        { status: 'pending', count: 5 },
        { status: 'accepted', count: 28 },
        { status: 'declined', count: 4 },
      ],
      invites_sent: 16,
      acceptance_percent: 62.5,
      median_acceptance: {
        status: 'not_tracked_yet',
        note: 'Invites do not store the time they were accepted, so this would be a guess. A later change can save that timestamp.',
      },
      groups: 7,
      average_group_size: 3.4,
      group_sizes: [
        { size: '1', groups: 0 },
        { size: '2', groups: 3 },
        { size: '3', groups: '<3' },
        { size: '4', groups: 3 },
        { size: '5', groups: 0 },
        { size: '6', groups: 0 },
        { size: '7', groups: 0 },
        { size: '8', groups: 0 },
        { size: '9', groups: 0 },
        { size: '10', groups: 0 },
        { size: '11', groups: 0 },
        { size: '12', groups: 0 },
      ],
      users_with_zero_partners: 31,
    },
    alerts: {
      partner_alerts_by_day: labels.map((day, i) => ({ day, count: 4 + (i % 3) })),
      monitoring_paused_by_day: labels.map((day, i) => ({ day, count: i % 4 === 0 ? '<3' : 3 })),
      support_taps: {
        status: 'tracked',
        count: 11,
      },
    },
    notifications: {
      push_tokens: 70,
      push_people: 58,
      apns: {
        status: 'not_tracked_yet',
        note: 'Push send success and failure are not written to the database.',
      },
      emails: [
        { type: 'verification', status: 'not_tracked_yet', note: 'Verification emails are not logged. The token row is replaced, then removed.' },
        { type: 'reset', status: 'not_tracked_yet', note: 'Password reset emails are not logged.' },
        { type: 'partner_invite', status: 'tracked', count: 16, note: 'Invite rows created in this range. An email is attempted for each.' },
        { type: 'group_invite', status: 'tracked', count: 5, note: 'Group invite rows created in this range. An email is attempted for each.' },
      ],
    },
    feedback: {
      status: 'not_tracked_yet',
      note: 'Messages from the website and the app are emailed. They are not saved here, so there is no list to show.',
    },
  }
}

export const previewHealth: Health = {
  api: 'ok',
  database: 'ok',
  db_bytes: 18_400_000,
  tables: [
    { table: 'users', rows: 86 },
    { table: 'relationships', rows: 37 },
    { table: 'groups', rows: 7 },
    { table: 'group_members', rows: 24 },
    { table: 'events', rows: 140 },
    { table: 'alerts', rows: 210 },
    { table: 'device_tokens', rows: 70 },
    { table: 'pageviews', rows: 1280 },
    { table: 'admin_audit', rows: 4 },
  ],
  server_errors: {
    status: 'not_tracked_yet',
    note: 'The API does not keep a counter of recent server errors.',
  },
}
