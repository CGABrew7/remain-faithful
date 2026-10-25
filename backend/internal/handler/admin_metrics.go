package handler

import (
	"context"
	"database/sql"
	"time"

	"remain-faithful/backend/internal/admin"
)

// Aggregate dashboard SQL. These statements return counts only.
// They do not select alert text, categories, app names, emails, or tokens.
const (
	sqlWebsiteTotals = `
SELECT
  count(*) FILTER (WHERE event_name = 'pageview') AS views,
  count(DISTINCT visitor_hash) FILTER (WHERE event_name = 'pageview') AS visitors,
  count(DISTINCT visitor_hash) FILTER (WHERE event_name = 'pageview' AND device_type = 'mobile') AS mobile_people,
  count(DISTINCT visitor_hash) FILTER (WHERE event_name = 'pageview' AND device_type = 'desktop') AS desktop_people,
  count(*) FILTER (WHERE event_name = 'app_store_click') AS store_clicks,
  count(DISTINCT visitor_hash) FILTER (WHERE event_name = 'app_store_click') AS store_people,
  count(*) FILTER (WHERE event_name = 'android_signup') AS android_signups,
  count(DISTINCT visitor_hash) FILTER (WHERE event_name = 'android_signup') AS android_people
FROM pageviews
WHERE ($1::timestamptz IS NULL OR created_at >= $1)`

	sqlTopPages = `
SELECT path, count(*) AS visits, count(DISTINCT visitor_hash) AS visitors
FROM pageviews
WHERE event_name = 'pageview'
  AND ($1::timestamptz IS NULL OR created_at >= $1)
GROUP BY path
ORDER BY count(*) DESC
LIMIT 15`

	sqlReferrers = `
SELECT COALESCE(NULLIF(referrer_domain, ''), '(direct)') AS label,
       count(DISTINCT visitor_hash) AS visitors
FROM pageviews
WHERE event_name = 'pageview'
  AND ($1::timestamptz IS NULL OR created_at >= $1)
GROUP BY 1
ORDER BY visitors DESC
LIMIT 15`

	sqlCampaigns = `
SELECT utm_source, utm_medium, utm_campaign, count(DISTINCT visitor_hash) AS visitors
FROM pageviews
WHERE event_name = 'pageview'
  AND (utm_source <> '' OR utm_medium <> '' OR utm_campaign <> '')
  AND ($1::timestamptz IS NULL OR created_at >= $1)
GROUP BY 1, 2, 3
ORDER BY visitors DESC
LIMIT 15`

	sqlCountries = `
SELECT country, count(DISTINCT visitor_hash) AS visitors
FROM pageviews
WHERE event_name = 'pageview' AND country <> ''
  AND ($1::timestamptz IS NULL OR created_at >= $1)
GROUP BY country
ORDER BY visitors DESC
LIMIT 30`

	sqlUSStates = `
SELECT region, count(DISTINCT visitor_hash) AS visitors
FROM pageviews
WHERE event_name = 'pageview' AND country = 'US' AND region <> ''
  AND ($1::timestamptz IS NULL OR created_at >= $1)
GROUP BY region
ORDER BY visitors DESC
LIMIT 60`

	sqlPageTrend = `
SELECT to_char(date_trunc('day', created_at AT TIME ZONE 'UTC'), 'YYYY-MM-DD') AS day,
       count(*) FILTER (WHERE event_name = 'pageview') AS visits,
       count(DISTINCT visitor_hash) FILTER (WHERE event_name = 'pageview') AS visitors
FROM pageviews
WHERE ($1::timestamptz IS NULL OR created_at >= $1)
GROUP BY 1
ORDER BY 1
LIMIT 400`

	sqlGrowth = `
SELECT
  count(*) AS total,
  count(*) FILTER (WHERE $1::timestamptz IS NULL OR created_at >= $1) AS fresh,
  count(*) FILTER (WHERE ($1::timestamptz IS NULL OR created_at >= $1) AND apple_id IS NOT NULL) AS apple,
  count(*) FILTER (WHERE ($1::timestamptz IS NULL OR created_at >= $1) AND apple_id IS NULL AND google_id IS NOT NULL) AS google,
  count(*) FILTER (WHERE ($1::timestamptz IS NULL OR created_at >= $1) AND apple_id IS NULL AND google_id IS NULL) AS email_sign_in,
  count(*) FILTER (WHERE ($1::timestamptz IS NULL OR created_at >= $1) AND apple_id IS NULL AND google_id IS NULL AND email_verified) AS verified_email,
  count(*) FILTER (WHERE ($1::timestamptz IS NULL OR created_at >= $1) AND apple_id IS NULL AND google_id IS NULL) AS email_accounts
FROM users`

	sqlFunnel = `
WITH cohort AS (
  SELECT id, last_heartbeat_at
  FROM users
  WHERE ($1::timestamptz IS NULL OR created_at >= $1)
)
SELECT
  (SELECT count(*) FROM cohort),
  (SELECT count(*) FROM cohort c WHERE EXISTS (
      SELECT 1 FROM relationship_invites i WHERE i.inviter_id = c.id
      UNION ALL
      SELECT 1 FROM relationships r WHERE r.user_id = c.id AND r.type = 'partner'
  )),
  (SELECT count(*) FROM cohort c WHERE
      EXISTS (
        SELECT 1 FROM relationship_invites i
        WHERE i.inviter_id = c.id AND i.status = 'accepted'
      )
      OR EXISTS (
        SELECT 1 FROM relationships r
        WHERE r.user_id = c.id AND r.type = 'partner' AND r.status = 'accepted'
      )
  ),
  (SELECT count(*) FROM cohort c WHERE c.last_heartbeat_at IS NOT NULL)`

	sqlActive = `
SELECT
  count(*) FILTER (WHERE last_heartbeat_at >= NOW() - INTERVAL '1 day'),
  count(*) FILTER (WHERE last_heartbeat_at >= NOW() - INTERVAL '7 days'),
  count(*) FILTER (WHERE last_heartbeat_at >= NOW() - INTERVAL '30 days'),
  count(*) FILTER (WHERE last_heartbeat_at >= NOW() - INTERVAL '30 minutes'),
  count(*) FILTER (WHERE last_heartbeat_at IS NOT NULL AND last_heartbeat_at < NOW() - INTERVAL '30 minutes'),
  count(*) FILTER (WHERE last_heartbeat_at IS NULL)
FROM users`

	sqlMonitoringOff = `
SELECT count(*) FROM users u
WHERE EXISTS (
  SELECT 1 FROM alerts a
  WHERE a.subject_user_id = u.id
    AND a.alert_type IN (
      'monitoring_disabled', 'shielding_disabled', 'family_controls_revoked',
      'signed_out', 'lockout_disabled'
    )
    AND a.created_at >= COALESCE(u.last_heartbeat_at, '-infinity'::timestamptz)
)`

	sqlPartnerships = `
WITH edges AS (
  SELECT status, user_id AS uid FROM relationships WHERE type = 'partner'
  UNION
  SELECT status, partner_id FROM relationships WHERE type = 'partner'
)
SELECT p.status, p.partnerships, COALESCE(e.people, 0)
FROM (
  SELECT status, count(*) AS partnerships
  FROM relationships
  WHERE type = 'partner'
  GROUP BY status
) p
LEFT JOIN (
  SELECT status, count(*) AS people FROM edges GROUP BY status
) e ON e.status = p.status`

	sqlInvites = `
SELECT
  count(*) AS sent,
  count(*) FILTER (WHERE status = 'accepted') AS accepted,
  (SELECT count(*) FROM (
      SELECT inviter_id::text AS who FROM relationship_invites
      WHERE ($1::timestamptz IS NULL OR created_at >= $1)
      UNION
      SELECT invitee_email FROM relationship_invites
      WHERE ($1::timestamptz IS NULL OR created_at >= $1)
  ) people)
FROM relationship_invites
WHERE ($1::timestamptz IS NULL OR created_at >= $1)`

	sqlGroupInvites = `
SELECT
  count(*) AS sent,
  (SELECT count(*) FROM (
      SELECT inviter_id::text AS who FROM group_invites
      WHERE ($1::timestamptz IS NULL OR created_at >= $1)
      UNION
      SELECT invitee_email FROM group_invites
      WHERE ($1::timestamptz IS NULL OR created_at >= $1)
  ) people)
FROM group_invites
WHERE ($1::timestamptz IS NULL OR created_at >= $1)`

	sqlGroupSummary = `
SELECT count(*) AS groups, COALESCE(sum(size), 0)::bigint AS members
FROM (
  SELECT count(*) AS size FROM group_members GROUP BY group_id
) sizes`

	sqlGroupSizes = `
SELECT size, count(*) AS groups, COALESCE(sum(size), 0)::bigint AS members
FROM (
  SELECT count(*)::int AS size FROM group_members GROUP BY group_id
) sizes
GROUP BY size
ORDER BY size`

	sqlZeroPartners = `
SELECT count(*) FROM users u
WHERE NOT EXISTS (
  SELECT 1 FROM relationships r
  WHERE r.status = 'accepted' AND r.type = 'partner'
    AND (r.user_id = u.id OR r.partner_id = u.id)
)`

	sqlPartnerAlerts = `
SELECT to_char(date_trunc('day', a.created_at AT TIME ZONE 'UTC'), 'YYYY-MM-DD') AS day,
       count(*) AS n,
       count(DISTINCT COALESCE(a.subject_user_id, e.user_id)) AS people
FROM alerts a
LEFT JOIN events e ON e.id = a.event_id
WHERE a.event_id IS NOT NULL
  AND ($1::timestamptz IS NULL OR a.created_at >= $1)
GROUP BY 1
ORDER BY 1
LIMIT 400`

	sqlPausedAlerts = `
SELECT to_char(date_trunc('day', created_at AT TIME ZONE 'UTC'), 'YYYY-MM-DD') AS day,
       count(*) AS n,
       count(DISTINCT subject_user_id) AS people
FROM alerts
WHERE alert_type IN (
  'monitoring_disabled', 'shielding_disabled', 'family_controls_revoked',
  'signed_out', 'lockout_disabled'
)
  AND ($1::timestamptz IS NULL OR created_at >= $1)
GROUP BY 1
ORDER BY 1
LIMIT 400`

	sqlPhoneSignals = `
SELECT
  count(*) FILTER (WHERE screen_time_granted IS TRUE),
  count(*) FILTER (WHERE COALESCE(monitored_app_count, 0) > 0),
  COALESCE(sum(monitored_app_count), 0)::bigint
FROM users
WHERE ($1::timestamptz IS NULL OR created_at >= $1)`

	sqlAppVersions = `
SELECT app_version, count(*)::bigint
FROM users
WHERE app_version IS NOT NULL AND app_version <> ''
GROUP BY app_version
ORDER BY count(*) DESC
LIMIT 12`

	sqlSupportTaps = `
SELECT COALESCE(sum(support_tap_count), 0)::bigint,
       count(*) FILTER (WHERE support_tap_count > 0)
FROM users`

	sqlPushTokens = `
SELECT
  count(*) FILTER (WHERE is_active),
  count(DISTINCT user_id) FILTER (WHERE is_active)
FROM device_tokens`

	sqlHealth = `
SELECT
  pg_database_size(current_database()),
  (SELECT count(*) FROM users),
  (SELECT count(*) FROM relationships),
  (SELECT count(*) FROM (SELECT user_id AS uid FROM relationships UNION SELECT partner_id FROM relationships) people),
  (SELECT count(*) FROM groups),
  (SELECT count(DISTINCT user_id) FROM group_members),
  (SELECT count(*) FROM events),
  (SELECT count(DISTINCT user_id) FROM events),
  (SELECT count(*) FROM alerts),
  (SELECT count(DISTINCT COALESCE(a.subject_user_id, e.user_id))
     FROM alerts a LEFT JOIN events e ON e.id = a.event_id),
  (SELECT count(*) FROM device_tokens),
  (SELECT count(DISTINCT user_id) FROM device_tokens),
  (SELECT count(*) FROM pageviews),
  (SELECT count(DISTINCT visitor_hash) FROM pageviews),
  (SELECT count(*) FROM admin_audit)`
)

func aggregateStatements() []string {
	return []string{
		sqlWebsiteTotals, sqlTopPages, sqlReferrers, sqlCampaigns, sqlCountries, sqlUSStates,
		sqlPageTrend, sqlGrowth, sqlFunnel, sqlActive, sqlMonitoringOff, sqlPartnerships,
		sqlInvites, sqlGroupInvites, sqlGroupSummary, sqlGroupSizes, sqlZeroPartners,
		signupSQL("day"), signupSQL("week"),
		sqlPartnerAlerts, sqlPausedAlerts, sqlPhoneSignals, sqlAppVersions, sqlSupportTaps,
		sqlPushTokens, sqlHealth, sqlAndroidWaitlist,
	}
}

const sqlAndroidWaitlist = `
SELECT COUNT(*) FROM waitlist_signups
WHERE platform = 'android'
  AND ($1::timestamptz IS NULL OR created_at >= $1)`

func (h *H) loadOverview(ctx context.Context, label string, from *time.Time, now time.Time) (*admin.Overview, error) {
	ov := blankOverview(label, now)
	var views, visitors, mobile, desktop, clicks, clickPeople, android, androidPeople int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlWebsiteTotals, from),
		&views, &visitors, &mobile, &desktop, &clicks, &clickPeople, &android, &androidPeople); err != nil {
		return nil, err
	}
	ov.Website.Visits = quietOr(views, visitors)
	ov.Website.Visitors = quietOr(visitors, visitors)
	ov.Website.MobileVisitors = quietOr(mobile, mobile)
	ov.Website.DesktopVisitors = quietOr(desktop, desktop)
	ov.Website.AppStoreClicks = quietOr(clicks, clickPeople)
	ov.Website.AndroidSignups = quietOr(android, androidPeople)
	var waitlist int
	if err := h.DB.QueryRowContext(ctx, sqlAndroidWaitlist, from).Scan(&waitlist); err != nil {
		return nil, err
	}
	if waitlist == 0 {
		ov.Website.AndroidWaitlist = admin.Zero()
	} else {
		ov.Website.AndroidWaitlist = admin.People(waitlist)
	}
	ctrPeople := visitors
	if visitors < 3 || (clickPeople > 0 && clickPeople < 3) {
		ctrPeople = 0
	}
	ov.Website.AppStoreCTR = admin.Percent{People: ctrPeople, Numerator: clicks, Denominator: views}

	pages, err := h.queryPages(ctx, sqlTopPages, from)
	if err != nil {
		return nil, err
	}
	ov.Website.TopPages = pages
	refs, err := h.queryNamed(ctx, sqlReferrers, from)
	if err != nil {
		return nil, err
	}
	ov.Website.TopReferrers = refs
	campaigns, err := h.queryCampaigns(ctx, from)
	if err != nil {
		return nil, err
	}
	ov.Website.Campaigns = campaigns
	countries, err := h.queryNamed(ctx, sqlCountries, from)
	if err != nil {
		return nil, err
	}
	ov.Website.Countries = countries
	states, err := h.queryNamed(ctx, sqlUSStates, from)
	if err != nil {
		return nil, err
	}
	ov.Website.USStates = states
	trend, err := h.queryTrend(ctx, from)
	if err != nil {
		return nil, err
	}
	ov.Website.Trend = fillTrend(from, now, trend)

	var total, fresh, apple, google, password, verified, emailAccounts int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlGrowth, from),
		&total, &fresh, &apple, &google, &password, &verified, &emailAccounts); err != nil {
		return nil, err
	}
	ov.Growth.TotalAccounts = quietOr(total, total)
	ov.Growth.NewAccounts = quietOr(fresh, fresh)
	ov.Growth.SignIn = []admin.MethodCount{
		{Method: "email", Count: quietOr(password, password)},
		{Method: "apple", Count: quietOr(apple, apple)},
		{Method: "google", Count: quietOr(google, google)},
	}
	if emailAccounts == 0 {
		ov.Growth.EmailVerified = admin.Percent{People: 3, Numerator: 0, Denominator: 0}
	} else {
		ov.Growth.EmailVerified = admin.Percent{
			People: admin.RateAudience(verified, emailAccounts), Numerator: verified, Denominator: emailAccounts,
		}
	}
	days, err := h.queryDayCounts(ctx, signupSQL("day"), from, true)
	if err != nil {
		return nil, err
	}
	ov.Growth.ByDay = fillDayCounts(from, now, days)
	weeks, err := h.queryDayCounts(ctx, signupSQL("week"), from, false)
	if err != nil {
		return nil, err
	}
	ov.Growth.ByWeek = weeks

	var accounts, invited, accepted, heartbeat int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlFunnel, from), &accounts, &invited, &accepted, &heartbeat); err != nil {
		return nil, err
	}
	var screenYes, appsPeople, appsTotal int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlPhoneSignals, from), &screenYes, &appsPeople, &appsTotal); err != nil {
		return nil, err
	}
	ov.Funnel = admin.BuildFunnel(accounts, invited, accepted, screenYes, heartbeat, appsPeople)
	ov.Active.MonitoredApps = quietOr(appsTotal, appsPeople)
	versions, err := h.queryAppVersions(ctx)
	if err != nil {
		return nil, err
	}
	ov.Active.AppVersions = admin.VersionRollup{Status: "tracked", Versions: versions}

	var dau, wau, mau, on, silent, never int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlActive), &dau, &wau, &mau, &on, &silent, &never); err != nil {
		return nil, err
	}
	var off int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlMonitoringOff), &off); err != nil {
		return nil, err
	}
	ov.Active.DAU = quietOr(dau, dau)
	ov.Active.WAU = quietOr(wau, wau)
	ov.Active.MAU = quietOr(mau, mau)
	ov.Active.MonitoringOn = quietOr(on, on)
	ov.Active.MonitoringOff = quietOr(off, off)
	ov.Active.HeartbeatSilent = quietOr(silent, silent)
	ov.Active.NeverCheckedIn = quietOr(never, never)
	var taps, tapPeople int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlSupportTaps), &taps, &tapPeople); err != nil {
		return nil, err
	}
	ov.Alerts.SupportTaps = admin.CountRollup{Status: "tracked", Count: quietOr(taps, tapPeople)}

	parts, err := h.queryPartnerships(ctx)
	if err != nil {
		return nil, err
	}
	ov.Accountability.Partnerships = parts

	var sent, acceptedInvites, invitePeople int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlInvites, from), &sent, &acceptedInvites, &invitePeople); err != nil {
		return nil, err
	}
	ov.Accountability.InvitesSent = quietOr(sent, invitePeople)
	if sent == 0 {
		ov.Accountability.Acceptance = admin.Percent{People: 3, Numerator: 0, Denominator: 0}
	} else {
		acceptPeople := invitePeople
		if invitePeople < 3 || (acceptedInvites > 0 && acceptedInvites < 3) || (sent-acceptedInvites > 0 && sent-acceptedInvites < 3) {
			acceptPeople = 0
		}
		ov.Accountability.Acceptance = admin.Percent{People: acceptPeople, Numerator: acceptedInvites, Denominator: sent}
	}

	var groupInvites, groupInvitePeople int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlGroupInvites, from), &groupInvites, &groupInvitePeople); err != nil {
		return nil, err
	}
	partnerMail := quietOr(sent, invitePeople)
	groupMail := quietOr(groupInvites, groupInvitePeople)
	ov.Notifications.Emails = []admin.EmailStat{
		{Type: "verification", Status: "not_tracked_yet", Note: "Verification emails are not logged. The token row is replaced, then removed."},
		{Type: "reset", Status: "not_tracked_yet", Note: "Password reset emails are not logged."},
		{Type: "partner_invite", Status: "tracked", Count: &partnerMail, Note: "Invite rows created in this range. An email is attempted for each."},
		{Type: "group_invite", Status: "tracked", Count: &groupMail, Note: "Group invite rows created in this range. An email is attempted for each."},
	}

	var groups, members int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlGroupSummary), &groups, &members); err != nil {
		return nil, err
	}
	ov.Accountability.Groups = quietOr(groups, members)
	if groups < 3 {
		ov.Accountability.Groups = quietOr(groups, groups)
	}
	if groups == 0 {
		ov.Accountability.AverageSize = admin.Average{People: 3, Groups: 3, Sum: 0}
	} else {
		ov.Accountability.AverageSize = admin.Average{People: members, Groups: groups, Sum: float64(members)}
	}
	sizes, err := h.queryGroupSizes(ctx)
	if err != nil {
		return nil, err
	}
	ov.Accountability.GroupSizes = sizes

	var zero int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlZeroPartners), &zero); err != nil {
		return nil, err
	}
	ov.Accountability.ZeroPartners = quietOr(zero, zero)

	alertDays, err := h.queryPeopleDays(ctx, sqlPartnerAlerts, from)
	if err != nil {
		return nil, err
	}
	ov.Alerts.PartnerAlertsByDay = fillDayCounts(from, now, alertDays)
	pausedDays, err := h.queryPeopleDays(ctx, sqlPausedAlerts, from)
	if err != nil {
		return nil, err
	}
	ov.Alerts.MonitoringPausedByDay = fillDayCounts(from, now, pausedDays)

	var tokens, tokenPeople int
	if err := scanInts(h.DB.QueryRowContext(ctx, sqlPushTokens), &tokens, &tokenPeople); err != nil {
		return nil, err
	}
	ov.Notifications.PushTokens = quietOr(tokens, tokenPeople)
	ov.Notifications.PushPeople = quietOr(tokenPeople, tokenPeople)
	return ov, nil
}

func signupSQL(unit string) string {
	// unit is a fixed literal, never request input.
	if unit != "week" {
		unit = "day"
	}
	return `
SELECT to_char(date_trunc('` + unit + `', created_at AT TIME ZONE 'UTC'), 'YYYY-MM-DD') AS day,
       count(*)
FROM users
WHERE ($1::timestamptz IS NULL OR created_at >= $1)
GROUP BY 1
ORDER BY 1
LIMIT 400`
}

func (h *H) loadHealth(ctx context.Context) (*admin.Health, error) {
	var (
		bytes                                                                                   int64
		users, rels, relPeople, groups, groupPeople                                             int
		events, eventPeople, alerts, alertPeople, tokens, tokenPeople, views, viewPeople, audit int
	)
	err := h.DB.QueryRowContext(ctx, sqlHealth).Scan(
		&bytes,
		&users, &rels, &relPeople, &groups, &groupPeople,
		&events, &eventPeople, &alerts, &alertPeople,
		&tokens, &tokenPeople, &views, &viewPeople, &audit,
	)
	if err != nil {
		return nil, err
	}
	return &admin.Health{
		API:      "ok",
		Database: "ok",
		DBBytes:  bytes,
		Tables: []admin.TableCount{
			{Table: "users", Rows: quietOr(users, users)},
			{Table: "relationships", Rows: quietOr(rels, relPeople)},
			{Table: "groups", Rows: quietOr(groups, groupPeople)},
			{Table: "group_members", Rows: quietOr(groupPeople, groupPeople)},
			{Table: "events", Rows: quietOr(events, eventPeople)},
			{Table: "alerts", Rows: quietOr(alerts, alertPeople)},
			{Table: "device_tokens", Rows: quietOr(tokens, tokenPeople)},
			{Table: "pageviews", Rows: quietOr(views, viewPeople)},
			{Table: "admin_audit", Rows: admin.Show(audit, 3)},
		},
		ServerErrors: admin.NotTracked("The API does not keep a counter of recent server errors."),
	}, nil
}

func quietOr(n, people int) admin.Count {
	if n == 0 && people == 0 {
		return admin.Zero()
	}
	if people == n {
		return admin.People(n)
	}
	return admin.Show(n, people)
}

func scanInts(row *sql.Row, dest ...*int) error {
	holders := make([]any, len(dest))
	raw := make([]int64, len(dest))
	for i := range dest {
		holders[i] = &raw[i]
	}
	if err := row.Scan(holders...); err != nil {
		return err
	}
	for i := range dest {
		*dest[i] = int(raw[i])
	}
	return nil
}

func (h *H) queryPages(ctx context.Context, q string, from *time.Time) ([]admin.PageRow, error) {
	rows, err := h.DB.QueryContext(ctx, q, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []admin.PageRow{}
	for rows.Next() {
		var path string
		var visits, visitors int64
		if err := rows.Scan(&path, &visits, &visitors); err != nil {
			return nil, err
		}
		out = append(out, admin.PageRow{
			Path: path, Visits: quietOr(int(visits), int(visitors)), Visitors: quietOr(int(visitors), int(visitors)),
		})
	}
	return out, rows.Err()
}

func (h *H) queryAppVersions(ctx context.Context) ([]admin.NamedCount, error) {
	rows, err := h.DB.QueryContext(ctx, sqlAppVersions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []admin.NamedCount{}
	for rows.Next() {
		var label string
		var n int64
		if err := rows.Scan(&label, &n); err != nil {
			return nil, err
		}
		out = append(out, admin.NamedCount{Label: label, Count: quietOr(int(n), int(n))})
	}
	return out, rows.Err()
}

func (h *H) queryNamed(ctx context.Context, q string, from *time.Time) ([]admin.NamedCount, error) {
	rows, err := h.DB.QueryContext(ctx, q, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []admin.NamedCount{}
	for rows.Next() {
		var label string
		var n int64
		if err := rows.Scan(&label, &n); err != nil {
			return nil, err
		}
		out = append(out, admin.NamedCount{Label: label, Count: quietOr(int(n), int(n))})
	}
	return out, rows.Err()
}

func (h *H) queryCampaigns(ctx context.Context, from *time.Time) ([]admin.CampaignRow, error) {
	rows, err := h.DB.QueryContext(ctx, sqlCampaigns, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []admin.CampaignRow{}
	for rows.Next() {
		var source, medium, campaign string
		var n int64
		if err := rows.Scan(&source, &medium, &campaign, &n); err != nil {
			return nil, err
		}
		out = append(out, admin.CampaignRow{
			Source: source, Medium: medium, Campaign: campaign, Visitors: quietOr(int(n), int(n)),
		})
	}
	return out, rows.Err()
}

func (h *H) queryTrend(ctx context.Context, from *time.Time) ([]admin.TrendPoint, error) {
	rows, err := h.DB.QueryContext(ctx, sqlPageTrend, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []admin.TrendPoint{}
	for rows.Next() {
		var day string
		var visits, visitors int64
		if err := rows.Scan(&day, &visits, &visitors); err != nil {
			return nil, err
		}
		out = append(out, admin.TrendPoint{
			Day: day, Visits: quietOr(int(visits), int(visitors)), Visitors: quietOr(int(visitors), int(visitors)),
		})
	}
	return out, rows.Err()
}

func (h *H) queryDayCounts(ctx context.Context, q string, from *time.Time, asPeople bool) ([]admin.DayCount, error) {
	rows, err := h.DB.QueryContext(ctx, q, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []admin.DayCount{}
	for rows.Next() {
		var day string
		var n int64
		if err := rows.Scan(&day, &n); err != nil {
			return nil, err
		}
		c := quietOr(int(n), int(n))
		if !asPeople {
			c = quietOr(int(n), int(n))
		}
		out = append(out, admin.DayCount{Day: day, Count: c})
	}
	return out, rows.Err()
}

func (h *H) queryPeopleDays(ctx context.Context, q string, from *time.Time) ([]admin.DayCount, error) {
	rows, err := h.DB.QueryContext(ctx, q, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []admin.DayCount{}
	for rows.Next() {
		var day string
		var n, people int64
		if err := rows.Scan(&day, &n, &people); err != nil {
			return nil, err
		}
		out = append(out, admin.DayCount{Day: day, Count: quietOr(int(n), int(people))})
	}
	return out, rows.Err()
}

func (h *H) queryPartnerships(ctx context.Context) ([]admin.Partnership, error) {
	rows, err := h.DB.QueryContext(ctx, sqlPartnerships)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	found := map[string]struct{ n, people int }{}
	for rows.Next() {
		var status string
		var n, people int64
		if err := rows.Scan(&status, &n, &people); err != nil {
			return nil, err
		}
		found[status] = struct{ n, people int }{int(n), int(people)}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	out := make([]admin.Partnership, 0, 3)
	for _, status := range []string{"pending", "accepted", "declined"} {
		row := found[status]
		out = append(out, admin.Partnership{Status: status, Count: quietOr(row.n, row.people)})
	}
	return out, nil
}

func (h *H) queryGroupSizes(ctx context.Context) ([]admin.GroupSize, error) {
	rows, err := h.DB.QueryContext(ctx, sqlGroupSizes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	found := map[int]struct{ groups, people int }{}
	var overflowGroups, overflowPeople int
	for rows.Next() {
		var size, groups, people int64
		if err := rows.Scan(&size, &groups, &people); err != nil {
			return nil, err
		}
		if size > 12 {
			overflowGroups += int(groups)
			overflowPeople += int(people)
			continue
		}
		found[int(size)] = struct{ groups, people int }{int(groups), int(people)}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	out := make([]admin.GroupSize, 0, 13)
	for size := 1; size <= 12; size++ {
		row := found[size]
		head := row.people
		if row.groups < 3 {
			head = row.groups
		}
		out = append(out, admin.GroupSize{Size: itoa(size), Groups: quietOr(row.groups, head)})
	}
	if overflowGroups > 0 {
		head := overflowPeople
		if overflowGroups < 3 {
			head = overflowGroups
		}
		out = append(out, admin.GroupSize{Size: "13+", Groups: quietOr(overflowGroups, head)})
	}
	return out, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [4]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}

func fillTrend(from *time.Time, now time.Time, points []admin.TrendPoint) []admin.TrendPoint {
	if from == nil {
		if points == nil {
			return []admin.TrendPoint{}
		}
		return points
	}
	by := map[string]admin.TrendPoint{}
	for _, p := range points {
		by[p.Day] = p
	}
	out := []admin.TrendPoint{}
	for _, day := range eachDay(*from, now) {
		if p, ok := by[day]; ok {
			out = append(out, p)
			continue
		}
		out = append(out, admin.TrendPoint{Day: day, Visits: admin.Zero(), Visitors: admin.Zero()})
	}
	return out
}

func fillDayCounts(from *time.Time, now time.Time, points []admin.DayCount) []admin.DayCount {
	if from == nil {
		if points == nil {
			return []admin.DayCount{}
		}
		return points
	}
	by := map[string]admin.DayCount{}
	for _, p := range points {
		by[p.Day] = p
	}
	out := []admin.DayCount{}
	for _, day := range eachDay(*from, now) {
		if p, ok := by[day]; ok {
			out = append(out, p)
			continue
		}
		out = append(out, admin.DayCount{Day: day, Count: admin.Zero()})
	}
	return out
}

func eachDay(from, now time.Time) []string {
	start := time.Date(from.UTC().Year(), from.UTC().Month(), from.UTC().Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(now.UTC().Year(), now.UTC().Month(), now.UTC().Day(), 0, 0, 0, 0, time.UTC)
	out := []string{}
	for d := start; !d.After(end) && len(out) < 100; d = d.AddDate(0, 0, 1) {
		out = append(out, d.Format("2006-01-02"))
	}
	return out
}

func blankOverview(label string, now time.Time) *admin.Overview {
	return &admin.Overview{
		Range:       label,
		GeneratedAt: now.UTC(),
		Website: admin.Website{
			TopPages: []admin.PageRow{}, TopReferrers: []admin.NamedCount{},
			Campaigns: []admin.CampaignRow{}, Countries: []admin.NamedCount{},
			USStates: []admin.NamedCount{}, Trend: []admin.TrendPoint{},
			AndroidWaitlist: admin.Zero(),
		},
		Growth: admin.Growth{
			ByDay: []admin.DayCount{}, ByWeek: []admin.DayCount{}, SignIn: []admin.MethodCount{},
			DeletedAccounts: admin.NotTracked("Accounts are removed rather than marked deleted, so a deleted total is not kept."),
		},
		Funnel: []admin.FunnelStep{},
		Active: admin.Active{
			Basis:         "These figures use the latest heartbeat on each account. They describe right now, not the date range above. Monitoring on means a heartbeat in the last 30 minutes. Heartbeat silent means the latest heartbeat is older than that. Monitoring off means a partner was told monitoring stopped and no heartbeat has arrived since.",
			MonitoredApps: admin.Zero(),
			DailyHistory:  admin.NotTracked("Only the latest heartbeat is stored, so this dashboard cannot draw active users for each past day."),
			AppVersions:   admin.VersionRollup{Status: "tracked", Versions: []admin.NamedCount{}},
		},
		Accountability: admin.Accountability{
			Partnerships: []admin.Partnership{}, GroupSizes: []admin.GroupSize{},
			MedianAcceptance: admin.NotTracked("Invites do not store the time they were accepted, so this would be a guess. A later change can save that timestamp."),
		},
		Alerts: admin.Alerts{
			PartnerAlertsByDay: []admin.DayCount{}, MonitoringPausedByDay: []admin.DayCount{},
			SupportTaps: admin.CountRollup{Status: "tracked", Count: admin.Zero()},
		},
		Notifications: admin.Notifications{
			Emails: []admin.EmailStat{},
			APNs:   admin.NotTracked("Push send success and failure are not written to the database."),
		},
		Feedback: admin.NotTracked("Messages from the website and the app are emailed. They are not saved here, so there is no list to show."),
	}
}

func rangeStart(label string, now time.Time) (string, *time.Time) {
	switch label {
	case "30d":
		t := now.AddDate(0, 0, -30)
		return "30d", &t
	case "90d":
		t := now.AddDate(0, 0, -90)
		return "90d", &t
	case "all":
		return "all", nil
	default:
		t := now.AddDate(0, 0, -7)
		return "7d", &t
	}
}
