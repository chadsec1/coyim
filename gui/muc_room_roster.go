package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	roomViewRosterStatusIconIndex int = iota
	roomViewRosterNicknameIndex
	roomViewRosterAffiliationIndex
	roomViewRosterRoleIndex
)

type roomViewRoster struct {
	roster *muc.RoomRoster

	view gtki.Box      `gtk-widget:"roster-view"`
	tree gtki.TreeView `gtk-widget:"roster-tree-view"`

	model gtki.TreeStore
}

func (v *roomView) newRoomViewRoster() *roomViewRoster {
	r := &roomViewRoster{
		roster: v.room.Roster(),
	}

	r.initBuilder()
	r.initDefaults()
	r.initSubscribers(v)

	return r
}

func (r *roomViewRoster) initBuilder() {
	builder := newBuilder("MUCRoomRoster")
	panicOnDevError(builder.bindObjects(r))
}

func (r *roomViewRoster) initDefaults() {
	r.model, _ = g.gtk.TreeStoreNew(
		// icon
		pixbufType(),
		// display nickname
		glibi.TYPE_STRING,
		// affiliation
		glibi.TYPE_STRING,
		// role - tooltip
		glibi.TYPE_STRING,
	)

	r.tree.SetModel(r.model)
	r.draw()
}

func (r *roomViewRoster) initSubscribers(v *roomView) {
	v.subscribe("roster", func(ev roomViewEvent) {
		switch ev.(type) {
		case occupantSelfJoinedEvent:
			r.onUpdateRoster()
		case occupantJoinedEvent:
			r.onUpdateRoster()
		case occupantUpdatedEvent:
			r.onUpdateRoster()
		case occupantLeftEvent:
			r.onUpdateRoster()
		}
	})
}

func (r *roomViewRoster) onUpdateRoster() {
	doInUIThread(r.redraw)
}

func (r *roomViewRoster) draw() {
	noneRoles, visitors, participants, moderators := r.roster.OccupantsByRole()

	r.drawOccupantsByRole(data.RoleModerator, moderators)
	r.drawOccupantsByRole(data.RoleParticipant, participants)
	r.drawOccupantsByRole(data.RoleVisitor, visitors)
	r.drawOccupantsByRole(data.RoleNone, noneRoles)

	r.tree.ExpandAll()
}

func (r *roomViewRoster) redraw() {
	r.model.Clear()
	r.draw()
}

func (r *roomViewRoster) drawOccupantsByRole(role string, occupants []*muc.Occupant) {
	if isOccupantListEmpty(occupants) {
		return
	}

	roleHeader := rolePluralName(role)
	roleHeader = i18n.Localf("%s (%v)", roleHeader, len(occupants))

	iter := r.model.Append(nil)
	_ = r.model.SetValue(iter, roomViewRosterNicknameIndex, roleHeader)

	for _, o := range occupants {
		r.addOccupantToRoster(o, iter)
	}
}

func (r *roomViewRoster) addOccupantToRoster(o *muc.Occupant, parentIter gtki.TreeIter) {
	iter := r.model.Append(parentIter)

	_ = r.model.SetValue(iter, roomViewRosterStatusIconIndex, getMUCIconPixbuf("occupant"))
	_ = r.model.SetValue(iter, roomViewRosterNicknameIndex, o.Nickname)
	_ = r.model.SetValue(iter, roomViewRosterAffiliationIndex, r.affiliationDisplayName(o.Affiliation))
	_ = r.model.SetValue(iter, roomViewRosterRoleIndex, roleDisplayName(o.Role))
}

func (r *roomViewRoster) affiliationDisplayName(a data.Affiliation) string {
	switch a.Name() {
	case data.AffiliationAdmin:
		return i18n.Local("Admin")
	case data.AffiliationOwner:
		return i18n.Local("Owner")
	case data.AffiliationOutcast:
		return i18n.Local("Outcast")
	default: // Member or other values get the default treatment
		return ""
	}
}

func rolePluralName(role string) string {
	switch role {
	case data.RoleNone:
		return i18n.Local("None")
	case data.RoleParticipant:
		return i18n.Local("Participants")
	case data.RoleVisitor:
		return i18n.Local("Visitors")
	case data.RoleModerator:
		return i18n.Local("Moderators")
	default:
		// This should not really be possible, but it is necessary
		// because golang can't prove it
		return ""
	}
}

func roleDisplayName(role data.Role) string {
	switch role.Name() {
	case data.RoleNone:
		return i18n.Local("None")
	case data.RoleParticipant:
		return i18n.Local("Participant")
	case data.RoleVisitor:
		return i18n.Local("Visitor")
	case data.RoleModerator:
		return i18n.Local("Moderator")
	default:
		// This should not really be possible, but it is necessary
		// because golang can't prove it
		return ""
	}
}

func isOccupantListEmpty(o []*muc.Occupant) bool {
	return len(o) == 0
}
