package gui

import (
	"time"

	"github.com/coyim/coyim/xmpp/jid"

	"github.com/coyim/coyim/session/muc/data"
)

type selfOccupantRemovedEvent struct{}

type selfOccupantJoinedEvent struct {
	nickname       string
	role           data.Role
	isReconnecting bool
}

type occupantLeftEvent struct {
	nickname string
}

type occupantJoinedEvent struct {
	nickname       string
	isReconnecting bool
}

type occupantUpdatedEvent struct {
	nickname string
	role     data.Role
}

type occupantRemovedEvent struct {
	nickname string
}

type nicknameConflictEvent struct {
	nickname string
}

type serviceUnavailableEvent struct{}

type unknownErrorEvent struct{}

type registrationRequiredEvent struct {
	nickname string
}

type notAuthorizedEvent struct{}

type occupantForbiddenEvent struct{}

type loggingEnabledEvent struct{}

type loggingDisabledEvent struct{}

type roomAnonymityEvent struct {
	anonymityLevel string
}

type messageEvent struct {
	tp        string
	nickname  string
	message   string
	timestamp time.Time
}

type subjectUpdatedEvent struct {
	nickname string
	subject  string
}

type subjectReceivedEvent struct {
	subject string
}

type joinRoomFinished struct {
	isReconnecting bool
}

type roomDestroyedEvent struct {
	reason      string
	alternative jid.Bare
	password    string
}

type messageForbidden struct{}

type messageNotAcceptable struct{}

type discussionHistoryEvent struct {
	history *data.DiscussionHistory
}

type roomViewEvent interface {
	markAsRoomViewEvent()
}

type roomDiscoInfoReceivedEvent struct {
	info data.RoomDiscoInfo
}

type roomConfigRequestTimeoutEvent struct{}

type roomConfigChangedEvent struct {
	changes   []data.RoomConfigType
	discoInfo data.RoomDiscoInfo
}

type occupantAffiliationRoleUpdatedEvent struct {
	affiliationRoleUpdate data.AffiliationRoleUpdate
}

type selfOccupantAffiliationRoleUpdatedEvent struct {
	selfAffiliationRoleUpdate data.AffiliationRoleUpdate
}

type occupantAffiliationUpdatedEvent struct {
	affiliationUpdate data.AffiliationUpdate
}

type selfOccupantAffiliationUpdatedEvent struct {
	selfAffiliationUpdate data.SelfAffiliationUpdate
}

type occupantRoleUpdatedEvent struct {
	roleUpdate data.RoleUpdate
}

type selfOccupantRoleUpdatedEvent struct {
	selfRoleUpdate data.RoleUpdate
}

type selfOccupantConnectedEvent struct{}

type selfOccupantDisconnectedEvent struct{}

type accountAffiliationUpdated struct {
	accountAddress jid.Any
	affiliation    data.Affiliation
}

type roomDisableEvent struct{}

type roomEnableEvent struct{}

func (selfOccupantRemovedEvent) markAsRoomViewEvent()                {}
func (occupantLeftEvent) markAsRoomViewEvent()                       {}
func (occupantJoinedEvent) markAsRoomViewEvent()                     {}
func (occupantUpdatedEvent) markAsRoomViewEvent()                    {}
func (selfOccupantJoinedEvent) markAsRoomViewEvent()                 {}
func (messageEvent) markAsRoomViewEvent()                            {}
func (subjectUpdatedEvent) markAsRoomViewEvent()                     {}
func (subjectReceivedEvent) markAsRoomViewEvent()                    {}
func (joinRoomFinished) markAsRoomViewEvent()                        {}
func (nicknameConflictEvent) markAsRoomViewEvent()                   {}
func (registrationRequiredEvent) markAsRoomViewEvent()               {}
func (loggingEnabledEvent) markAsRoomViewEvent()                     {}
func (loggingDisabledEvent) markAsRoomViewEvent()                    {}
func (roomAnonymityEvent) markAsRoomViewEvent()                      {}
func (messageForbidden) markAsRoomViewEvent()                        {}
func (messageNotAcceptable) markAsRoomViewEvent()                    {}
func (discussionHistoryEvent) markAsRoomViewEvent()                  {}
func (roomDiscoInfoReceivedEvent) markAsRoomViewEvent()              {}
func (roomConfigRequestTimeoutEvent) markAsRoomViewEvent()           {}
func (roomDestroyedEvent) markAsRoomViewEvent()                      {}
func (roomConfigChangedEvent) markAsRoomViewEvent()                  {}
func (occupantRemovedEvent) markAsRoomViewEvent()                    {}
func (notAuthorizedEvent) markAsRoomViewEvent()                      {}
func (occupantForbiddenEvent) markAsRoomViewEvent()                  {}
func (serviceUnavailableEvent) markAsRoomViewEvent()                 {}
func (unknownErrorEvent) markAsRoomViewEvent()                       {}
func (occupantAffiliationRoleUpdatedEvent) markAsRoomViewEvent()     {}
func (selfOccupantAffiliationRoleUpdatedEvent) markAsRoomViewEvent() {}
func (occupantAffiliationUpdatedEvent) markAsRoomViewEvent()         {}
func (selfOccupantAffiliationUpdatedEvent) markAsRoomViewEvent()     {}
func (occupantRoleUpdatedEvent) markAsRoomViewEvent()                {}
func (selfOccupantRoleUpdatedEvent) markAsRoomViewEvent()            {}
func (selfOccupantConnectedEvent) markAsRoomViewEvent()              {}
func (selfOccupantDisconnectedEvent) markAsRoomViewEvent()           {}
func (accountAffiliationUpdated) markAsRoomViewEvent()               {}
func (roomDisableEvent) markAsRoomViewEvent()                        {}
func (roomEnableEvent) markAsRoomViewEvent()                         {}
