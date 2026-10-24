package handler

import (
	"context"
	"errors"
	"log"

	"remain-faithful/backend/internal/apns"
	"remain-faithful/backend/internal/fcm"
)

// deliverAndroid sends one APNs-shaped notification to every active Android
// token for userID. iOS tokens stay on the APNs queries. Environment is an
// APNs sandbox/production flag and is not applied to FCM.
func (h *H) deliverAndroid(ctx context.Context, userID int64, n *apns.Notification, reportRemoval bool) {
	if h == nil || h.FCM == nil || h.DB == nil || userID == 0 || n == nil {
		return
	}
	rows, err := h.DB.QueryContext(ctx, `
		SELECT token FROM device_tokens
		WHERE user_id = $1 AND is_active = TRUE AND platform = 'android'
	`, userID)
	if err != nil {
		log.Printf("[fcm] query tokens user=%d: %v", userID, err)
		return
	}
	defer rows.Close()

	sent := 0
	for rows.Next() {
		var token string
		if rows.Scan(&token) != nil || token == "" {
			continue
		}
		msg := fcm.FromAPNS(token, n.CollapseID, n.Payload)
		if err := h.FCM.Send(ctx, msg); err != nil {
			var invalid *fcm.ErrInvalidToken
			if errors.As(err, &invalid) {
				h.markTokenInactive(ctx, token)
				if reportRemoval {
					h.notifyAppRemoved(ctx, userID)
				}
				continue
			}
			log.Printf("[fcm] send user=%d token=%.8s...: %v", userID, token, err)
			continue
		}
		sent++
	}
	if err := rows.Err(); err != nil {
		log.Printf("[fcm] tokens user=%d: %v", userID, err)
	}
	log.Printf("[fcm] user=%d sent=%d", userID, sent)
}
