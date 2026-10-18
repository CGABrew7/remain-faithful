package admin

const (
	noteScreenTime = "The phone does not tell the server when Screen Time access is granted. A later app build could send a yes or no only, with no app names."
	noteAppsChosen = "The phone does not tell the server how many apps were chosen. A later app build could send that count only, never the app names."
	noteHeartbeat  = "Counted when a heartbeat has arrived at least once. That means monitoring has checked in. It is not a separate Screen Time flag."
)

// BuildFunnel turns cohort headcounts into ordered steps.
// Drop-off is the share lost since the previous tracked step.
// Steps with no server signal stay visible and carry no count.
func BuildFunnel(accounts, invited, accepted, heartbeat int) []FunnelStep {
	steps := []struct {
		key, label, note string
		tracked          bool
		n                int
	}{
		{"account_created", "Account created", "", true, accounts},
		{"invited_partner", "Invited a partner", "", true, invited},
		{"partner_accepted", "Partner accepted", "", true, accepted},
		{"screen_time", "Screen Time granted", noteScreenTime, false, 0},
		{"monitoring_active", "Monitoring has checked in", noteHeartbeat, true, heartbeat},
		{"apps_chosen", "Apps chosen", noteAppsChosen, false, 0},
	}

	out := make([]FunnelStep, 0, len(steps))
	prev := -1
	prevN := 0
	for _, s := range steps {
		step := FunnelStep{Key: s.key, Label: s.label, Note: s.note}
		if !s.tracked {
			step.Status = "not_tracked_yet"
			out = append(out, step)
			continue
		}
		step.Status = "tracked"
		c := People(s.n)
		step.Count = &c
		if prev >= 0 {
			lost := prevN - s.n
			if lost < 0 {
				lost = 0
			}
			// A precise drop-off next to a shown total would reveal a "<3" step.
			people := prevN
			if prevN < 3 || (s.n > 0 && s.n < 3) || (lost > 0 && lost < 3) {
				people = 0
			}
			d := Percent{People: people, Numerator: lost, Denominator: prevN}
			step.Dropoff = &d
		}
		prev = len(out)
		prevN = s.n
		out = append(out, step)
	}
	return out
}
