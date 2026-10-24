package fcm

import "fmt"

// FromAPNS copies display text and routing keys from an APNs payload.
// The aps dictionary stays on the Apple path. Android receives title, body,
// and the same custom keys (notification_type, alert_type, sender_name, and
// event fields). Nothing is added.
func FromAPNS(token, collapse string, payload any) *Message {
	msg := &Message{
		Token:       token,
		CollapseKey: collapse,
		Data:        map[string]string{},
		ChannelID:   "alerts",
	}
	m, ok := payload.(map[string]any)
	if !ok {
		return msg
	}
	if aps, ok := m["aps"].(map[string]any); ok {
		switch alert := aps["alert"].(type) {
		case map[string]string:
			msg.Title = alert["title"]
			msg.Body = alert["body"]
		case map[string]any:
			msg.Title, _ = alert["title"].(string)
			msg.Body, _ = alert["body"].(string)
		case string:
			msg.Body = alert
		}
	}
	for k, v := range m {
		if k == "aps" {
			continue
		}
		if s, ok := scalarString(v); ok {
			msg.Data[k] = s
		}
	}
	if msg.Data["notification_type"] == "PANIC_ALERT" {
		msg.ChannelID = "urgent"
	}
	return msg
}

func scalarString(v any) (string, bool) {
	switch t := v.(type) {
	case string:
		return t, true
	case int:
		return fmt.Sprintf("%d", t), true
	case int64:
		return fmt.Sprintf("%d", t), true
	case float64:
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t)), true
		}
		return fmt.Sprint(t), true
	case bool:
		if t {
			return "true", true
		}
		return "false", true
	default:
		return "", false
	}
}
