package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"remain-faithful/backend/internal/apns"
)

// StartHeartbeatSweep runs the heartbeat silence sweep every 5 minutes until
// ctx is cancelled. silenceMinutes is read from HEARTBEAT_SILENCE_MINUTES in
// main (default 30).
//
// The sweep fans a heartbeat_silence PROTECTION_ALERT out to the accountability
// partners of any user whose heartbeats have gone silent longer than the
// threshold. It uses last_silence_alert_at on the relationships row as a
// per-relationship dedup signal, so each (monitored-user, partner) pair is
// alerted at most once per silence window.
//
// Reset: when a new heartbeat arrives for a user, last_heartbeat_at advances
// past last_silence_alert_at. The sweep condition
// (last_silence_alert_at < last_heartbeat_at) becomes true again, meaning the
// relationship re-enters the monitoring pool. The threshold check then prevents
// a re-alert until the user goes silent again.
func (h *H) StartHeartbeatSweep(ctx context.Context, silenceMinutes int) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	log.Printf("[heartbeat-sweep] started — silence threshold: %d minutes", silenceMinutes)

	for {
		select {
		case <-ctx.Done():
			log.Println("[heartbeat-sweep] stopping")
			return
		case <-ticker.C:
			h.sweepHeartbeatSilence(ctx, silenceMinutes)
			h.sweepDeviceHeartbeatSilence(ctx, silenceMinutes)
		}
	}
}

type silenceRow struct {
	userID    int64
	userName  string
	relID     int64
	partnerID int64
}

func (h *H) sweepHeartbeatSilence(ctx context.Context, silenceMinutes int) {
	h.purgeOldEvents(ctx)
	// Users who have never heartbeated are left alone. A null heartbeat is
	// "not armed", not "monitoring is off".
	// Recipients (accepted partners and group members) are chosen inside
	// storeAndPushProtection. A recipient who already has a silence row from
	// the last 6 hours is skipped, so a continued silence repeats slowly.
	rows, err := h.DB.QueryContext(ctx, `
		SELECT u.id, u.name
		FROM   users u
		WHERE  u.last_heartbeat_at IS NOT NULL
		  AND  u.last_heartbeat_at < NOW() - INTERVAL '1 minute' * $1
	`, silenceMinutes)
	if err != nil {
		log.Printf("[heartbeat-sweep] query error: %v", err)
		return
	}
	defer rows.Close()

	var targets []silenceRow
	for rows.Next() {
		var s silenceRow
		if err := rows.Scan(&s.userID, &s.userName); err != nil {
			continue
		}
		targets = append(targets, s)
	}

	if len(targets) == 0 {
		return
	}

	log.Printf("[heartbeat-sweep] %d silent user(s)", len(targets))

	for _, t := range targets {
		h.sendProtectionAlertToPartners(t.userID, t.userName, "heartbeat_silence", "")
	}
}

func (h *H) purgeOldEvents(ctx context.Context) {
	if _, err := h.DB.ExecContext(ctx,
		`DELETE FROM events WHERE timestamp < NOW() - INTERVAL '365 days'`,
	); err != nil {
		log.Printf("[retention] events: %v", err)
	}
	if _, err := h.DB.ExecContext(ctx,
		`DELETE FROM alerts WHERE created_at < NOW() - INTERVAL '365 days'`,
	); err != nil {
		log.Printf("[retention] alerts: %v", err)
	}
}

// sendSilenceAlertToPartner delivers a heartbeat_silence PROTECTION_ALERT push
// to partnerID. The payload mirrors the format used by sendProtectionAlertToPartners
// so the iOS app can handle both HTTP-triggered and sweep-triggered alerts uniformly.
func (h *H) sendSilenceAlertToPartner(ctx context.Context, partnerID int64, userName string, monitoredUserID int64) {
	body := protectionAlertBody(userName, "heartbeat_silence", "")
	payload := map[string]any{
		"aps": map[string]any{
			"alert": map[string]string{
				"title": "Protection Update",
				"body":  body,
			},
			"sound": "default",
		},
		"notification_type": "PROTECTION_ALERT",
		"alert_type":        "heartbeat_silence",
		"sender_name":       userName,
	}
	n := &apns.Notification{
		PushType:   "alert",
		Priority:   10,
		CollapseID: fmt.Sprintf("protection-%d-heartbeat_silence", monitoredUserID),
		Payload:    payload,
	}
	// notifyPartnerByID handles ErrInvalidToken (marks the token inactive) and
	// logs delivery outcomes internally.
	h.notifyPartnerByID(ctx, partnerID, userName, n)
}
