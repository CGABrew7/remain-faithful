package main

import (
	"errors"
	"strings"

	"remain-faithful/desktop/core"
	"remain-faithful/desktop/monitor"
)

func (a *App) ListGroups() ([]GroupCard, error) {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	rows, err := a.client.ListGroups(ctx)
	if err != nil {
		return nil, a.soften(err)
	}
	out := make([]GroupCard, 0, len(rows))
	for _, row := range rows {
		out = append(out, GroupCard{ID: row.ID, Name: row.Name})
	}
	return out, nil
}

func (a *App) CreateGroup(name string) (GroupCard, error) {
	if err := core.ValidateGroupName(name); err != nil {
		return GroupCard{}, err
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	group, err := a.client.CreateGroup(ctx, strings.TrimSpace(name))
	if err != nil {
		return GroupCard{}, a.soften(err)
	}
	return GroupCard{ID: group.ID, Name: group.Name}, nil
}

func (a *App) GetGroup(id int64) (GroupDetail, error) {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	group, err := a.client.GetGroup(ctx, id)
	if err != nil {
		return GroupDetail{}, a.soften(err)
	}
	me := a.store.User().ID
	detail := GroupDetail{
		ID:          group.ID,
		Name:        group.Name,
		MemberCount: len(group.Members),
		CanInvite:   core.CanInvite(len(group.Members)),
		Link:        core.GroupInviteLink(group.ID),
		Members:     []MemberView{},
	}
	if !detail.CanInvite {
		detail.FullMessage = core.GroupFullMessage
	}
	for _, member := range group.Members {
		name := member.User.Name
		if name == "" {
			name = "Member"
		}
		detail.Members = append(detail.Members, MemberView{
			UserID: member.UserID,
			Name:   name,
			Role:   member.Role,
			Streak: member.StreakDays,
			Flags:  member.FlagsLast30,
			Self:   member.UserID == me,
		})
	}
	return detail, nil
}

func (a *App) InviteToGroup(id int64, email string) (string, error) {
	if err := core.ValidateEmail(email); err != nil {
		return "", err
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	status, err := a.client.InviteGroupEmail(ctx, id, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return "", a.soften(err)
	}
	return core.GroupInviteMessage(status), nil
}

func (a *App) GroupLink(id int64) string {
	return core.GroupInviteLink(id)
}

func (a *App) LeaveGroup(id int64) error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.soften(a.client.LeaveGroup(ctx, id))
}

func (a *App) LeaveAllGroups() error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.soften(a.client.LeaveAllGroups(ctx))
}

func (a *App) Encourage(groupID, userID int64) error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.soften(a.client.Encourage(ctx, groupID, userID))
}

func (a *App) PendingGroupInvites() ([]InviteRow, error) {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	rows, err := a.client.PendingGroupInvites(ctx)
	if err != nil {
		return nil, a.soften(err)
	}
	out := make([]InviteRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, InviteRow{ID: row.ID, Title: row.GroupName, Detail: row.InviterName})
	}
	return out, nil
}

func (a *App) AcceptGroupInvite(id int64) error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.soften(a.client.AcceptGroupInvite(ctx, id))
}

func (a *App) DeclineGroupInvite(id int64) error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.soften(a.client.DeclineGroupInvite(ctx, id))
}

func (a *App) ListPartners() ([]PartnerView, error) {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	rows, err := a.client.ListRelationships(ctx)
	if err != nil {
		return nil, a.soften(err)
	}
	me := a.store.User().ID
	out := make([]PartnerView, 0, len(rows))
	for _, row := range rows {
		out = append(out, PartnerView{
			ID:        row.ID,
			Name:      row.Partner.Name,
			Email:     row.Partner.Email,
			Status:    row.Status,
			Primary:   row.IsPrimary && row.UserID == me,
			PinSet:    row.PinSet,
			CanSetPIN: row.PartnerID == me && row.Status == "accepted",
			CanStar:   row.UserID == me && row.Status == "accepted",
		})
	}
	return out, nil
}

func (a *App) InvitePartner(email string) (string, error) {
	if err := core.ValidateEmail(email); err != nil {
		return "", err
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	status, err := a.client.InvitePartner(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return "", a.soften(err)
	}
	return core.PartnerInviteMessage(status), nil
}

func (a *App) PendingPartnerInvites() ([]InviteRow, error) {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	rows, err := a.client.PendingPartnerInvites(ctx)
	if err != nil {
		return nil, a.soften(err)
	}
	out := make([]InviteRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, InviteRow{ID: row.ID, Title: row.InviterName, Detail: "Partner invite"})
	}
	return out, nil
}

func (a *App) AcceptPartnerInvite(id int64) error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.soften(a.client.AcceptPartnerInvite(ctx, id))
}

func (a *App) DeclinePartnerInvite(id int64) error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.soften(a.client.DeclinePartnerInvite(ctx, id))
}

func (a *App) SetPrimary(id int64) error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.soften(a.client.SetPrimary(ctx, id))
}

func (a *App) RemovePartner(id int64) error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	return a.soften(a.client.DeleteRelationship(ctx, id))
}

func (a *App) SetPIN(relationshipID int64, pin string, replacing bool) error {
	if err := core.ValidatePIN(pin); err != nil {
		return err
	}
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	if err := a.client.SetRelationshipPIN(ctx, relationshipID, pin); err != nil {
		return a.soften(err)
	}
	if note, ok, err := core.PINChangedNotice(replacing); err != nil {
		return err
	} else if ok {
		a.deliver([]core.ProtectionNotice{note})
	}
	return nil
}

func (a *App) RemovePIN(relationshipID int64) error {
	ctx, cancel := a.ctxTimeout()
	defer cancel()
	if err := a.client.RemoveRelationshipPIN(ctx, relationshipID); err != nil {
		return a.soften(err)
	}
	note, err := core.PINRemovedNotice()
	if err != nil {
		return err
	}
	a.deliver([]core.ProtectionNotice{note})
	return nil
}

func (a *App) Settings() (SettingsView, error) {
	if !a.store.SignedIn() {
		return SettingsView{}, errors.New("Not signed in. Please log in again.")
	}
	user := a.store.User()
	st := a.store.State()
	status := core.MonitoringStatus(st.Consent, len(st.Selected), st.Alerts, st.Blocking, monitor.Ready())
	pin, _ := a.pinRequired(true)
	partners, err := a.ListPartners()
	if err != nil {
		partners = []PartnerView{}
	}
	pendingP, err := a.PendingPartnerInvites()
	if err != nil {
		pendingP = []InviteRow{}
	}
	pendingG, err := a.PendingGroupInvites()
	if err != nil {
		pendingG = []InviteRow{}
	}
	return SettingsView{
		Name:            user.Name,
		Email:           user.Email,
		EmailVerified:   user.EmailVerified,
		Platform:        monitor.Platform(),
		Version:         core.VersionLabel(),
		DeviceCaption:   a.DeviceCaption(),
		Consent:         st.Consent,
		AlertsOn:        st.Alerts,
		BlockingOn:      st.Blocking,
		LockoutOn:       st.Lockout,
		PinSet:          pin,
		StatusTitle:     status.Title,
		StatusDetail:    status.Detail,
		PauseNotice:     core.PauseNotice,
		Apps:            a.Catalog(),
		Partners:        partners,
		PendingPartners: pendingP,
		PendingGroups:   pendingG,
	}, nil
}
