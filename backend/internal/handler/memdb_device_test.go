package handler

import (
	"database/sql/driver"
	"sort"
	"strings"
	"time"
)

// memDevice and memRemoval are the in-memory stand-in for the per-device tables.
type memDevice struct {
	userID     int64
	id         string
	kind       string
	monitoring bool
	lastSeen   *time.Time
	silence    *time.Time
	screen     string
	created    time.Time
}

type memRemoval struct {
	userID int64
	id     string
	kind   string
}

// applyDevice handles the device-heartbeat statements. apply already holds m.mu.
func (m *memDB) applyDevice(q string, args []any) ([]string, [][]driver.Value, bool) {
	switch {
	case containsAll(q, "FROM device_heartbeats d", "JOIN users", "silence_alerted_at"):
		minutes := asInt64(args[0])
		cutoff := time.Now().Add(-time.Duration(minutes) * time.Minute)
		var rows [][]driver.Value
		for _, d := range m.devices {
			if d.lastSeen == nil || !d.lastSeen.Before(cutoff) {
				continue
			}
			if d.silence != nil && !d.silence.Before(*d.lastSeen) {
				continue
			}
			name := ""
			if u := m.users[d.userID]; u != nil {
				name = u.name
			}
			rows = append(rows, []driver.Value{d.userID, name, d.id, d.kind})
		}
		return []string{"user_id", "name", "device_id", "device_kind"}, rows, true

	case containsAll(q, "UPDATE device_heartbeats", "silence_alerted_at"):
		userID, id := asInt64(args[0]), asString(args[1])
		now := time.Now()
		for i := range m.devices {
			if m.devices[i].userID == userID && m.devices[i].id == id {
				m.devices[i].silence = &now
			}
		}
		return []string{"ok"}, nil, true

	case containsAll(q, "DELETE FROM device_heartbeats"):
		userID, id := asInt64(args[0]), asString(args[1])
		kept := make([]memDevice, 0, len(m.devices))
		for _, d := range m.devices {
			if d.userID == userID && d.id == id {
				continue
			}
			kept = append(kept, d)
		}
		m.devices = kept
		return []string{"ok"}, nil, true

	case containsAll(q, "DELETE FROM device_removals"):
		userID, id := asInt64(args[0]), asString(args[1])
		kept := make([]memRemoval, 0, len(m.removals))
		for _, d := range m.removals {
			if d.userID == userID && d.id == id {
				continue
			}
			kept = append(kept, d)
		}
		m.removals = kept
		return []string{"ok"}, nil, true

	case containsAll(q, "INSERT INTO device_removals"):
		row := memRemoval{userID: asInt64(args[0]), id: asString(args[1]), kind: asString(args[2])}
		found := false
		for _, d := range m.removals {
			if d.userID == row.userID && d.id == row.id {
				found = true
			}
		}
		if !found {
			m.removals = append(m.removals, row)
		}
		return []string{"ok"}, nil, true

	case containsAll(q, "INSERT INTO device_heartbeats"):
		userID := asInt64(args[0])
		id := asString(args[1])
		kind := asString(args[2])
		if containsAll(q, "last_seen_at") {
			now := time.Now()
			screen := asString(args[3])
			m.upsertDevice(userID, id, kind, true, &now, &screen)
		} else {
			monitoring, _ := args[3].(bool)
			m.upsertDevice(userID, id, kind, monitoring, nil, nil)
		}
		return []string{"ok"}, nil, true

	case containsAll(q, "SELECT 1 FROM device_removals"):
		userID, id := asInt64(args[0]), asString(args[1])
		for _, d := range m.removals {
			if d.userID == userID && d.id == id {
				return []string{"one"}, [][]driver.Value{{int64(1)}}, true
			}
		}
		return []string{"one"}, nil, true

	case containsAll(q, "SELECT device_id, device_kind"):
		userID := asInt64(args[0])
		items := make([]memDevice, 0)
		for _, d := range m.devices {
			if d.userID == userID {
				items = append(items, d)
			}
		}
		sort.Slice(items, func(i, j int) bool {
			if items[i].created.Equal(items[j].created) {
				return items[i].id < items[j].id
			}
			return items[i].created.Before(items[j].created)
		})
		var rows [][]driver.Value
		for _, d := range items {
			var seen driver.Value
			if d.lastSeen != nil {
				seen = *d.lastSeen
			}
			rows = append(rows, []driver.Value{d.id, d.kind, d.monitoring, seen})
		}
		return []string{"device_id", "device_kind", "monitoring", "last_seen_at"}, rows, true

	case containsAll(q, "SELECT device_kind FROM device_heartbeats"):
		userID, id := asInt64(args[0]), asString(args[1])
		for _, d := range m.devices {
			if d.userID == userID && d.id == id {
				return []string{"device_kind"}, [][]driver.Value{{d.kind}}, true
			}
		}
		return []string{"device_kind"}, nil, true

	default:
		return nil, nil, false
	}
}

func (m *memDB) upsertDevice(userID int64, id, kind string, monitoring bool, seen *time.Time, screen *string) {
	for i := range m.devices {
		if m.devices[i].userID == userID && m.devices[i].id == id {
			m.devices[i].kind = kind
			m.devices[i].monitoring = monitoring
			if seen != nil {
				t := *seen
				m.devices[i].lastSeen = &t
			}
			if screen != nil {
				m.devices[i].screen = *screen
			}
			return
		}
	}
	row := memDevice{
		userID:     userID,
		id:         id,
		kind:       kind,
		monitoring: monitoring,
		created:    time.Now(),
	}
	if seen != nil {
		t := *seen
		row.lastSeen = &t
	}
	if screen != nil {
		row.screen = *screen
	}
	m.devices = append(m.devices, row)
}

func (m *memDB) ageDevice(userID int64, id string, at time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.devices {
		if m.devices[i].userID == userID && m.devices[i].id == id {
			t := at
			m.devices[i].lastSeen = &t
		}
	}
}

func (m *memDB) sawSQL(part string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, q := range m.sqlLog {
		if strings.Contains(q, part) {
			return true
		}
	}
	return false
}

func (m *memDB) alertSummaries() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]string, len(m.alerts))
	for i, row := range m.alerts {
		out[i] = row.summary
	}
	return out
}
