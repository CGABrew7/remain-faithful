export type Count = number | '<3'
export type Percent = number | '<3' | null

export type NotTracked = { status: 'not_tracked_yet'; note: string }
export type VersionRollup = {
  status: 'tracked' | 'not_tracked_yet'
  versions?: { label: string; count: Count }[]
  note?: string
}
export type CountRollup = {
  status: 'tracked' | 'not_tracked_yet'
  count?: Count
  note?: string
}

export type Overview = {
  range: string
  generated_at: string
  website: {
    visits: Count
    visitors: Count
    top_pages: { path: string; visits: Count; visitors: Count }[]
    top_referrers: { label: string; count: Count }[]
    campaigns: { source: string; medium: string; campaign: string; visitors: Count }[]
    countries: { label: string; count: Count }[]
    us_states: { label: string; count: Count }[]
    mobile_visitors: Count
    desktop_visitors: Count
    app_store_clicks: Count
    app_store_ctr_percent: Percent
    android_signups: Count
    trend: { day: string; visits: Count; visitors: Count }[]
  }
  growth: {
    total_accounts: Count
    new_accounts: Count
    by_day: { day: string; count: Count }[]
    by_week: { day: string; count: Count }[]
    sign_in: { method: string; count: Count }[]
    email_verified_percent: Percent
    deleted_accounts: NotTracked
  }
  funnel: {
    key: string
    label: string
    status: 'tracked' | 'not_tracked_yet'
    count?: Count
    dropoff_percent?: Percent
    note?: string
  }[]
  active: {
    basis: string
    dau: Count
    wau: Count
    mau: Count
    monitoring_on: Count
    monitoring_off: Count
    heartbeat_silent: Count
    never_checked_in: Count
    monitored_apps: Count
    daily_history: NotTracked
    app_versions: VersionRollup
  }
  accountability: {
    partnerships: { status: string; count: Count }[]
    invites_sent: Count
    acceptance_percent: Percent
    median_acceptance: NotTracked
    groups: Count
    average_group_size: number | '<3' | null
    group_sizes: { size: string; groups: Count }[]
    users_with_zero_partners: Count
  }
  alerts: {
    partner_alerts_by_day: { day: string; count: Count }[]
    monitoring_paused_by_day: { day: string; count: Count }[]
    support_taps: CountRollup
  }
  notifications: {
    push_tokens: Count
    push_people: Count
    apns: NotTracked
    emails: { type: string; status: string; count?: Count; note?: string }[]
  }
  feedback: NotTracked
}

export type Health = {
  api: string
  database: string
  db_bytes: number
  tables: { table: string; rows: Count }[]
  server_errors: NotTracked
}

export type AccountStatus = {
  email: string
  created_at: string
  email_verified: boolean
  sign_in_method: string
  last_active_at: string | null
  partner_count: number
}
