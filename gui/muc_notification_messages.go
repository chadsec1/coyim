package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/session/muc/data"
)

func getAffiliationUpdateSuccessMessage(nickname string, previousAffiliation, affiliation data.Affiliation) string {
	if affiliation.IsNone() {
		return getAffiliationRemovedSuccessMessage(nickname, previousAffiliation)
	}
	return getAffiliationChangedSuccessMessage(nickname, affiliation)
}

func getAffiliationRemovedSuccessMessage(nickname string, previousAffiliation data.Affiliation) string {
	switch {
	case previousAffiliation.IsOwner():
		return i18n.Localf("$nickname{%s} is not $affiliation{an owner} anymore.", nickname)
	case previousAffiliation.IsAdmin():
		return i18n.Localf("$nickname{%s} is not $affiliation{an administrator} anymore.", nickname)
	case previousAffiliation.IsMember():
		return i18n.Localf("$nickname{%s} is not $affiliation{a member} anymore.", nickname)
	}
	return i18n.Localf("$nickname{%s} is not $affiliation{banned} anymore.", nickname)
}

func getAffiliationChangedSuccessMessage(nickname string, affiliation data.Affiliation) string {
	switch {
	case affiliation.IsOwner():
		return i18n.Localf("The position of $nickname{%s} was changed to $affiliation{owner}.", nickname)
	case affiliation.IsAdmin():
		return i18n.Localf("The position of $nickname{%s} was changed to $affiliation{administrator}.", nickname)
	case affiliation.IsMember():
		return i18n.Localf("The position of $nickname{%s} was changed to $affiliation{member}.", nickname)
	case affiliation.IsBanned():
		return i18n.Localf("$nickname{%s} has been $affiliation{banned} from the room.", nickname)
	}
	return i18n.Localf("The position of $nickname{%s} was changed.", nickname)
}

// getRoleUpdateSuccessMessage returns a friendly notification message for the role update process
// This function receives the following params:
// nickname - The nickname of the occupant to whom the role was changed
// previousRole - The previous role of the occupant
// newRole - The new role of the occupant
func getRoleUpdateSuccessMessage(nickname string, previousRole, newRole data.Role) string {
	switch {
	case newRole.IsNone():
		return i18n.Localf("$nickname{%s} was $role{temporarily removed} from the room.", nickname)
	case previousRole.IsNone():
		return getRoleAddedSuccessMessage(nickname, newRole)
	}
	return getRoleChangedSuccessMessage(nickname, previousRole, newRole)
}

func getRoleAddedSuccessMessage(nickname string, newRole data.Role) string {
	switch {
	case newRole.IsModerator():
		return i18n.Localf("The role of $nickname{%s} was changed to $role{moderator}.", nickname)
	case newRole.IsParticipant():
		return i18n.Localf("The role of $nickname{%s} was changed to $role{participant}.", nickname)
	}
	return i18n.Localf("The role of $nickname{%s} was changed to $role{visitor}.", nickname)
}

func getRoleChangedSuccessMessage(nickname string, previousRole, newRole data.Role) string {
	switch {
	case previousRole.IsModerator() && newRole.IsParticipant():
		return i18n.Localf("The role of $nickname{%s} was changed from $role{moderator} to $role{participant}.", nickname)
	case previousRole.IsModerator() && newRole.IsVisitor():
		return i18n.Localf("The role of $nickname{%s} was changed from $role{moderator} to $role{visitor}.", nickname)
	case previousRole.IsParticipant() && newRole.IsModerator():
		return i18n.Localf("The role of $nickname{%s} was changed from $role{participant} to $role{moderator}.", nickname)
	case previousRole.IsParticipant() && newRole.IsVisitor():
		return i18n.Localf("The role of $nickname{%s} was changed from $role{participant} to $role{visitor}.", nickname)
	case previousRole.IsVisitor() && newRole.IsModerator():
		return i18n.Localf("The role of $nickname{%s} was changed from $role{visitor} to $role{moderator}.", nickname)
	case previousRole.IsVisitor() && newRole.IsParticipant():
		return i18n.Localf("The role of $nickname{%s} was changed from $role{visitor} to $role{participant}.", nickname)
	}
	return i18n.Localf("The role of $nickname{%s} was changed.", nickname)
}

type updateFailureMessages struct {
	notificationMessage string
	errorDialogTitle    string
	errorDialogHeader   string
	errorDialogMessage  string
}

func getAffiliationUpdateFailureMessage(nickname string, newAffiliation data.Affiliation, err error) *updateFailureMessages {
	if newAffiliation.IsBanned() {
		return getBannedFailureMessage(nickname, newAffiliation, err)
	}

	return &updateFailureMessages{
		notificationMessage: i18n.Localf("The position of $nickname{%s} couldn't be changed.", nickname),
		errorDialogTitle:    i18n.Local("Changing the position failed"),
		errorDialogHeader:   i18n.Localf("The position of %s couldn't be changed", nickname),
		errorDialogMessage:  getAffiliationFailureErrorMessage(nickname, newAffiliation, err),
	}
}

func getBannedFailureMessage(nickname string, newAffiliation data.Affiliation, err error) *updateFailureMessages {
	return &updateFailureMessages{
		notificationMessage: i18n.Localf("$nickname{%s} couldn't be banned.", nickname),
		errorDialogTitle:    i18n.Local("Banning failed"),
		errorDialogHeader:   i18n.Localf("%s couldn't be banned", nickname),
		errorDialogMessage:  getAffiliationFailureErrorMessage(nickname, newAffiliation, err),
	}
}

func getAffiliationFailureErrorMessage(nickname string, newAffiliation data.Affiliation, err error) string {
	if err == session.ErrOwnerAffiliationRevokeConflict {
		return i18n.Local("You can't change your own position because you are the only owner for this room. Every room must have at least one owner.")
	}
	return getUpdateAffiliationFailureErrorMessage(nickname, newAffiliation)
}

func getUpdateAffiliationFailureErrorMessage(nickname string, newAffiliation data.Affiliation) string {
	switch {
	case newAffiliation.IsOwner():
		return i18n.Localf("An error occurred trying to change the position of %s to owner.", nickname)
	case newAffiliation.IsAdmin():
		return i18n.Localf("An error occurred trying to change the position of %s to administrator.", nickname)
	case newAffiliation.IsMember():
		return i18n.Localf("An error occurred trying to change the position of %s to member.", nickname)
	case newAffiliation.IsBanned():
		return i18n.Localf("An error occurred trying to ban %s.", nickname)
	}
	return i18n.Localf("An error occurred trying to change the position of %s.", nickname)
}

func getRoleUpdateFailureMessage(nickname string, newRole data.Role) *updateFailureMessages {
	if newRole.IsNone() {
		return getRoleRemoveFailureMessage(nickname, nil, nil)
	}

	return &updateFailureMessages{
		notificationMessage: i18n.Localf("The role of $nickname{%s} couldn't be changed.", nickname),
		errorDialogTitle:    i18n.Local("Changing the role failed"),
		errorDialogHeader:   i18n.Localf("The role of %s couldn't be changed", nickname),
		errorDialogMessage:  getUpdateRoleFailureErrorMessage(nickname, newRole),
	}
}

func getUpdateRoleFailureErrorMessage(nickname string, newRole data.Role) string {
	switch {
	case newRole.IsModerator():
		return i18n.Localf("An error occurred trying to change the role of %s to moderator.", nickname)
	case newRole.IsParticipant():
		return i18n.Localf("An error occurred trying to change the role of %s to participant.", nickname)
	case newRole.IsVisitor():
		return i18n.Localf("An error occurred trying to change the role of %s to visitor.", nickname)
	}
	return i18n.Localf("An error occurred trying to change the role of %s.", nickname)
}

func getRoleRemoveFailureMessage(nickname string, actorAffiliation data.Affiliation, err error) *updateFailureMessages {
	return &updateFailureMessages{
		notificationMessage: i18n.Localf("$nickname{%s} couldn't be expelled.", nickname),
		errorDialogTitle:    i18n.Local("Expelling failed"),
		errorDialogHeader:   i18n.Localf("%s couldn't be expelled", nickname),
		errorDialogMessage:  getRoleRemoveFailureMessageBasedOnError(nickname, actorAffiliation, err),
	}
}

func getRoleRemoveFailureMessageBasedOnError(nickname string, actorAffiliation data.Affiliation, err error) string {
	switch err {
	case session.ErrNotAllowedKickOccupant:
		return getRoleRemoveFailureMessageWithActor(nickname, actorAffiliation)
	}
	return i18n.Localf("An error occurred expelling %s.", nickname)
}

func getRoleRemoveFailureMessageWithActor(nickname string, actorAffiliation data.Affiliation) string {
	if actorAffiliation != nil {
		switch {
		case actorAffiliation.IsOwner():
			return i18n.Localf("As an owner you don't have permissions to expel %s.", nickname)
		case actorAffiliation.IsAdmin():
			return i18n.Localf("As an administrator you don't have permissions to expel %s.", nickname)
		case actorAffiliation.IsMember():
			return i18n.Localf("As a member you don't have permissions to expel %s.", nickname)
		}
	}

	return i18n.Localf("You don't have permissions to expel %s.", nickname)
}

func getMUCNotificationMessageFrom(d interface{}) string {
	switch t := d.(type) {
	case data.AffiliationUpdate:
		return getAffiliationUpdateMessage(t)
	case data.SelfAffiliationUpdate:
		return getSelfAffiliationUpdateMessage(t)
	case data.RoleUpdate:
		return getRoleUpdateMessage(t)
	case data.AffiliationRoleUpdate:
		return getAffiliationRoleUpdateMessage(t)
	}
	panic(fmt.Sprintf("unkown update type: %v", d))
}

func getAffiliationUpdateMessage(affiliationUpdate data.AffiliationUpdate) string {
	return appendReasonToMessage(getAffiliationUpdateBaseMessage(affiliationUpdate), affiliationUpdate.Reason)
}

func getAffiliationUpdateBaseMessage(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.New.IsNone():
		return getAffiliationRemovedMessage(affiliationUpdate)
	case affiliationUpdate.New.IsBanned():
		return getAffiliationBannedMessage(affiliationUpdate)
	case affiliationUpdate.Previous.IsNone():
		return getAffiliationAddedMessage(affiliationUpdate)
	}
	return getAffiliationChangedMessage(affiliationUpdate)
}

func getAffiliationRemovedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		return getAffiliationRemovedMessageWithoutActor(affiliationUpdate)
	}
	return getAffiliationRemovedMessageWithActor(affiliationUpdate)
}

func getAffiliationRemovedMessageWithoutActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner():
		return i18n.Localf("%s is not an owner anymore.", affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin():
		return i18n.Localf("%s is not an administrator anymore.", affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember():
		return i18n.Localf("%s is not a member anymore.", affiliationUpdate.Nickname)
	}
	return i18n.Localf("%s is not banned anymore.", affiliationUpdate.Nickname)
}

func getAffiliationRemovedMessageWithActor(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor.Affiliation.IsOwner() {
		return getAffiliationRemovedMessageWithOwnerActor(affiliationUpdate)
	}
	return getAffiliationRemovedMessageWithAdminActor(affiliationUpdate)
}
func getAffiliationRemovedMessageWithOwnerActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[2]s is not an owner anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[2]s is not an administrator anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[2]s is not a member anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}
	return i18n.Localf("The owner %[1]s changed the position of %[2]s.",
		affiliationUpdate.Actor.Nickname,
		affiliationUpdate.Nickname)
}

func getAffiliationRemovedMessageWithAdminActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[2]s is not an owner anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[2]s is not an administrator anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[2]s is not a member anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}
	return i18n.Localf("The administrator %[1]s changed the position of %[2]s.",
		affiliationUpdate.Actor.Nickname,
		affiliationUpdate.Nickname)
}

func getAffiliationBannedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		return getAffiliationBannedMessageWithoutActor(affiliationUpdate)
	}
	return getAffiliationBannedMessageWithActor(affiliationUpdate)
}

func getAffiliationBannedMessageWithoutActor(affiliationUpdate data.AffiliationUpdate) string {
	return i18n.Localf("%s was banned from the room.", affiliationUpdate.Nickname)
}

func getAffiliationBannedMessageWithActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Actor.Affiliation.IsOwner():
		return i18n.Localf("The owner %[1]s banned %[2]s from the room.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Actor.Affiliation.IsAdmin():
		return i18n.Localf("The administrator %[1]s banned %[2]s from the room.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}
	return i18n.Localf("%[1]s banned %[2]s from the room.", affiliationUpdate.Actor.Nickname, affiliationUpdate.Nickname)
}

func getAffiliationAddedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		return getAffiliationAddedMessageWithoutActor(affiliationUpdate)
	}
	return getAffiliationAddedMessageWithActor(affiliationUpdate)
}

func getAffiliationAddedMessageWithoutActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.New.IsOwner():
		return i18n.Localf("%s is now an owner.", affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsAdmin():
		return i18n.Localf("%s is now an administrator.", affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsMember():
		return i18n.Localf("%s is now a member.", affiliationUpdate.Nickname)
	}
	return i18n.Localf("%s is now banned.", affiliationUpdate.Nickname)
}

func getAffiliationAddedMessageWithActor(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor.Affiliation.IsOwner() {
		return getAffiliationAddedMessageWithOwnerActor(affiliationUpdate)
	}
	return getAffiliationAddedMessageWithAdminActor(affiliationUpdate)
}

func getAffiliationAddedMessageWithOwnerActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[2]s is now an owner.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[2]s is now an administrator.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[2]s is now a member.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}
	return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[2]s is now banned",
		affiliationUpdate.Actor.Nickname,
		affiliationUpdate.Nickname)
}

func getAffiliationAddedMessageWithAdminActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[2]s is now an owner",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[2]s is now an administrator",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[2]s is now a member",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}
	return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[2]s is now banned",
		affiliationUpdate.Actor.Nickname,
		affiliationUpdate.Nickname)
}

func getAffiliationChangedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		return getAffiliationChangedMessageWithoutActor(affiliationUpdate)
	}
	return getAffiliationChangedMessageWithActor(affiliationUpdate)
}

func getAffiliationChangedMessageWithoutActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The position of %s was changed from owner to administrator.",
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The position of %s was changed from owner to member.",
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The position of %s was changed from administrator to member.",
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The position of %s was changed from administrator to owner.",
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The position of %s was changed from member to administrator.",
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The position of %s was changed from member to owner.",
			affiliationUpdate.Nickname)
	}
	return i18n.Localf("The position of %s was changed.", affiliationUpdate.Nickname)
}

func getAffiliationChangedMessageWithActor(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor.Affiliation.IsOwner() {
		return getAffiliationChangedMessageWithOwnerActor(affiliationUpdate)
	}
	return getAffiliationChangedMessageWithAdminActor(affiliationUpdate)
}

func getAffiliationChangedMessageWithOwnerActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from owner to administrator.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from owner to member.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from administrator to owner.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from administrator to member.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from member to owner.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from member to administrator.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}
	return i18n.Localf("The owner %[1]s changed the position of %[2]s.",
		affiliationUpdate.Actor.Nickname,
		affiliationUpdate.Nickname)
}

func getAffiliationChangedMessageWithAdminActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from owner to administrator.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from owner to member.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from administrator to owner.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from administrator to member.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from member to owner.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from member to administrator.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}
	return i18n.Localf("The administrator %[1]s changed the position of %[2]s.",
		affiliationUpdate.Actor.Nickname,
		affiliationUpdate.Nickname)
}

func getRoleUpdateMessage(roleUpdate data.RoleUpdate) string {
	return appendReasonToMessage(getRoleUpdateBaseMessage(roleUpdate), roleUpdate.Reason)
}

func getRoleUpdateBaseMessage(roleUpdate data.RoleUpdate) string {
	if roleUpdate.New.IsNone() {
		return getRoleRemovedMessage(roleUpdate)
	}
	return getRoleChangedMessage(roleUpdate)
}

func getRoleRemovedMessage(roleUpdate data.RoleUpdate) string {
	if roleUpdate.Actor == nil {
		return i18n.Localf("%s was temporarily removed from the room.", roleUpdate.Nickname)
	}
	return getRoleRemovedMessageWithActor(roleUpdate)
}

func getRoleRemovedMessageWithActor(roleUpdate data.RoleUpdate) string {
	switch {
	case roleUpdate.Actor.Affiliation.IsOwner():
		return i18n.Localf("The owner %[1]s temporarily removed %[2]s from the room.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Actor.Affiliation.IsAdmin():
		return i18n.Localf("The administrator %[1]s temporarily removed %[2]s from the room.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	}
	return i18n.Localf("%[1]s temporarily removed %[2]s from the room.",
		roleUpdate.Actor.Nickname,
		roleUpdate.Nickname)
}

func getRoleChangedMessage(roleUpdate data.RoleUpdate) string {
	if roleUpdate.Actor == nil {
		return getRoleChangedMessageWithoutActor(roleUpdate)
	}
	return getRoleChangedMessageWithActor(roleUpdate)
}

func getRoleChangedMessageWithoutActor(roleUpdate data.RoleUpdate) string {
	switch {
	case roleUpdate.Previous.IsModerator() && roleUpdate.New.IsParticipant():
		return i18n.Localf("The role of %s was changed from moderator to participant.",
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsModerator() && roleUpdate.New.IsVisitor():
		return i18n.Localf("The role of %s was changed from moderator to visitor.",
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsParticipant() && roleUpdate.New.IsModerator():
		return i18n.Localf("The role of %s was changed from participant to moderator.",
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsParticipant() && roleUpdate.New.IsVisitor():
		return i18n.Localf("The role of %s was changed from participant to visitor.",
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsVisitor() && roleUpdate.New.IsModerator():
		return i18n.Localf("The role of %s was changed from visitor to moderator.",
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsVisitor() && roleUpdate.New.IsParticipant():
		return i18n.Localf("The role of %s was changed from visitor to participant.",
			roleUpdate.Nickname)
	}
	return i18n.Localf("The role of %s was changed.",
		roleUpdate.Nickname)
}

func getRoleChangedMessageWithActor(roleUpdate data.RoleUpdate) string {
	switch {
	case roleUpdate.Actor.Affiliation.IsOwner():
		return getRoleChangedMessageWithOwnerActor(roleUpdate)
	case roleUpdate.Actor.Affiliation.IsAdmin():
		return getRoleChangedMessageWithAdminActor(roleUpdate)
	}
	return getRoleChangedMessageForActor(roleUpdate)
}

func getRoleChangedMessageWithOwnerActor(roleUpdate data.RoleUpdate) string {
	switch {
	case roleUpdate.Previous.IsModerator() && roleUpdate.New.IsParticipant():
		return i18n.Localf("The owner %[1]s changed the role of %[2]s from moderator to participant.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsModerator() && roleUpdate.New.IsVisitor():
		return i18n.Localf("The owner %[1]s changed the role of %[2]s from moderator to visitor.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsParticipant() && roleUpdate.New.IsModerator():
		return i18n.Localf("The owner %[1]s changed the role of %[2]s from participant to moderator.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsParticipant() && roleUpdate.New.IsVisitor():
		return i18n.Localf("The owner %[1]s changed the role of %[2]s from participant to visitor.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsVisitor() && roleUpdate.New.IsModerator():
		return i18n.Localf("The owner %[1]s changed the role of %[2]s from visitor to moderator.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsVisitor() && roleUpdate.New.IsParticipant():
		return i18n.Localf("The owner %[1]s changed the role of %[2]s from visitor to participant.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	}
	return i18n.Localf("The owner %[1]s changed the role of %[2]s.",
		roleUpdate.Actor.Nickname,
		roleUpdate.Nickname)
}

func getRoleChangedMessageWithAdminActor(roleUpdate data.RoleUpdate) string {
	switch {
	case roleUpdate.Previous.IsModerator() && roleUpdate.New.IsParticipant():
		return i18n.Localf("The administrator %[1]s changed the role of %[2]s from moderator to participant.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsModerator() && roleUpdate.New.IsVisitor():
		return i18n.Localf("The administrator %[1]s changed the role of %[2]s from moderator to visitor.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsParticipant() && roleUpdate.New.IsModerator():
		return i18n.Localf("The administrator %[1]s changed the role of %[2]s from participant to moderator.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsParticipant() && roleUpdate.New.IsVisitor():
		return i18n.Localf("The administrator %[1]s changed the role of %[2]s from participant to visitor.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsVisitor() && roleUpdate.New.IsModerator():
		return i18n.Localf("The administrator %[1]s changed the role of %[2]s from visitor to moderator.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsVisitor() && roleUpdate.New.IsParticipant():
		return i18n.Localf("The administrator %[1]s changed the role of %[2]s from visitor to participant.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	}
	return i18n.Localf("The administrator %[1]s changed the role of %[2]s.",
		roleUpdate.Actor.Nickname,
		roleUpdate.Nickname)
}

func getRoleChangedMessageForActor(roleUpdate data.RoleUpdate) string {
	switch {
	case roleUpdate.Previous.IsModerator() && roleUpdate.New.IsParticipant():
		return i18n.Localf("%[1]s changed the role of %[2]s from moderator to participant.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsModerator() && roleUpdate.New.IsVisitor():
		return i18n.Localf("%[1]s changed the role of %[2]s from moderator to visitor.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsParticipant() && roleUpdate.New.IsModerator():
		return i18n.Localf("%[1]s changed the role of %[2]s from participant to moderator.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsParticipant() && roleUpdate.New.IsVisitor():
		return i18n.Localf("%[1]s changed the role of %[2]s from participant to visitor.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsVisitor() && roleUpdate.New.IsModerator():
		return i18n.Localf("%[1]s changed the role of %[2]s from visitor to moderator.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	case roleUpdate.Previous.IsVisitor() && roleUpdate.New.IsParticipant():
		return i18n.Localf("%[1]s changed the role of %[2]s from visitor to participant.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname)
	}
	return i18n.Localf("%[1]s changed the role of %[2]s.",
		roleUpdate.Actor.Nickname,
		roleUpdate.Nickname)
}

func getSelfRoleUpdateMessage(selfRoleUpdate data.RoleUpdate) string {
	return appendReasonToMessage(getSelfRoleUpdateBaseMessage(selfRoleUpdate), selfRoleUpdate.Reason)
}

func getSelfRoleUpdateBaseMessage(selfRoleUpdate data.RoleUpdate) string {
	if selfRoleUpdate.Actor == nil {
		return getSelfRoleUpdateMessageWithoutActor(selfRoleUpdate)
	}
	return getSelfRoleUpdateMessageWithActor(selfRoleUpdate)
}

func getSelfRoleUpdateMessageWithoutActor(selfRoleUpdate data.RoleUpdate) string {
	switch {
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsVisitor():
		return i18n.Local("Your role was changed from $role{participant} to $role{visitor}.")
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsModerator():
		return i18n.Local("Your role was changed from $role{participant} to $role{moderator}.")
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsParticipant():
		return i18n.Local("Your role was changed from $role{visitor} to $role{participant}.")
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsModerator():
		return i18n.Local("Your role was changed from $role{visitor} to $role{moderator}.")
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsVisitor():
		return i18n.Local("Your role was changed from $role{moderator} to $role{visitor}.")
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsParticipant():
		return i18n.Local("Your role was changed from $role{moderator} to $role{participant}.")
	case selfRoleUpdate.New.IsNone():
		return i18n.Local("You have been $role{temporarily expelled}.")
	}
	return i18n.Local("Your role was changed.")
}

func getSelfRoleUpdateMessageWithActor(selfRoleUpdate data.RoleUpdate) string {
	switch {
	case selfRoleUpdate.Actor.Affiliation.IsOwner():
		return getSelfRoleUpdateMessageForOwnerActor(selfRoleUpdate)
	case selfRoleUpdate.Actor.Affiliation.IsAdmin():
		return getSelfRoleUpdateMessageForAdminActor(selfRoleUpdate)
	}
	return getSelfRoleUpdateMessageForActor(selfRoleUpdate)
}

func getSelfRoleUpdateMessageForOwnerActor(selfRoleUpdate data.RoleUpdate) string {
	switch {
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("The owner %s changed your role from $role{participant} to $role{visitor}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("The owner %s changed your role from $role{participant} to $role{moderator}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("The owner %s changed your role from $role{visitor} to $role{participant}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("The owner %s changed your role from $role{visitor} to $role{moderator}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("The owner %s changed your role from $role{moderator} to $role{visitor}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("The owner %s changed your role from $role{moderator} to $role{participant}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.New.IsNone():
		return i18n.Localf("The owner %s has $role{temporarily expelled} you from the room.",
			selfRoleUpdate.Actor.Nickname)
	}
	return i18n.Localf("The owner %s changed your role.",
		selfRoleUpdate.Actor.Nickname)
}

func getSelfRoleUpdateMessageForAdminActor(selfRoleUpdate data.RoleUpdate) string {
	switch {
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("The administrator %s changed your role from $role{participant} to $role{visitor}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("The administrator %s changed your role from $role{participant} to $role{moderator}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("The administrator %s changed your role from $role{visitor} to $role{participant}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("The administrator %s changed your role from $role{visitor} to $role{moderator}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("The administrator %s changed your role from $role{moderator} to $role{visitor}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("The administrator %s changed your role from $role{moderator} to $role{participant}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.New.IsNone():
		return i18n.Localf("The administrator %s has $role{temporarily expelled} you from the room.",
			selfRoleUpdate.Actor.Nickname)
	}
	return i18n.Localf("The administrator %s changed your role.",
		selfRoleUpdate.Actor.Nickname)
}

func getSelfRoleUpdateMessageForActor(selfRoleUpdate data.RoleUpdate) string {
	switch {
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("%s changed your role from $role{participant} to $role{visitor}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("%s changed your role from $role{participant} to $role{moderator}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("%s changed your role from $role{visitor} to $role{participant}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("%s changed your role from $role{visitor} to $role{moderator}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("%s changed your role from $role{moderator} to $role{visitor}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("%s changed your role from $role{moderator} to $role{participant}.",
			selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.New.IsNone():
		return i18n.Localf("%s has $role{temporarily expelled} you from the room.",
			selfRoleUpdate.Actor.Nickname)
	}
	return i18n.Localf("%s changed your role.",
		selfRoleUpdate.Actor.Nickname)
}

func getAffiliationRoleUpdateMessage(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	affiliationMessage := getAffiliationRoleUpdateMessageForAffiliation(affiliationRoleUpdate)
	roleMessage := getAffiliationRoleUpdateMessageForRole(affiliationRoleUpdate)

	return appendReasonToMessage(
		i18n.Localf("%[1]s %[2]s", affiliationMessage, roleMessage),
		affiliationRoleUpdate.Reason,
	)
}

func getAffiliationRoleUpdateMessageForAffiliation(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.NewAffiliation.IsNone():
		return getAffiliationRoleUpateForAffiliationRemoved(affiliationRoleUpdate)
	case affiliationRoleUpdate.PreviousAffiliation.IsNone():
		return getAffiliationRoleUpdateForAffiliationAdded(affiliationRoleUpdate)
	}
	return getAffiliationRoleUpdateForAffiliationUpdated(affiliationRoleUpdate)
}

func getAffiliationRoleUpdateMessageForRole(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousRole.IsVisitor() && affiliationRoleUpdate.NewRole.IsParticipant():
		return i18n.Local("As a result, their role was changed from visitor to participant.")
	case affiliationRoleUpdate.PreviousRole.IsVisitor() && affiliationRoleUpdate.NewRole.IsModerator():
		return i18n.Local("As a result, their role was changed from visitor to moderator.")
	case affiliationRoleUpdate.PreviousRole.IsParticipant() && affiliationRoleUpdate.NewRole.IsVisitor():
		return i18n.Local("As a result, their role was changed from participant to visitor.")
	case affiliationRoleUpdate.PreviousRole.IsParticipant() && affiliationRoleUpdate.NewRole.IsModerator():
		return i18n.Local("As a result, their role was changed from participant to moderator.")
	case affiliationRoleUpdate.PreviousRole.IsModerator() && affiliationRoleUpdate.NewRole.IsVisitor():
		return i18n.Local("As a result, their role was changed from moderator to visitor.")
	case affiliationRoleUpdate.PreviousRole.IsModerator() && affiliationRoleUpdate.NewRole.IsParticipant():
		return i18n.Local("As a result, their role was changed from moderator to participant.")
	}
	return i18n.Local("As a result, their role was also changed.")
}

func getSelfAffiliationRoleUpdateMessageForRole(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousRole.IsVisitor() && affiliationRoleUpdate.NewRole.IsParticipant():
		return i18n.Local("As a result, your role was changed from $role{visitor} to $role{participant}.")
	case affiliationRoleUpdate.PreviousRole.IsVisitor() && affiliationRoleUpdate.NewRole.IsModerator():
		return i18n.Local("As a result, your role was changed from $role{visitor} to $role{moderator}.")
	case affiliationRoleUpdate.PreviousRole.IsParticipant() && affiliationRoleUpdate.NewRole.IsVisitor():
		return i18n.Local("As a result, your role was changed from $role{participant} to $role{visitor}.")
	case affiliationRoleUpdate.PreviousRole.IsParticipant() && affiliationRoleUpdate.NewRole.IsModerator():
		return i18n.Local("As a result, your role was changed from $role{participant} to $role{moderator}.")
	case affiliationRoleUpdate.PreviousRole.IsModerator() && affiliationRoleUpdate.NewRole.IsVisitor():
		return i18n.Local("As a result, your role was changed from $role{moderator} to $role{visitor}.")
	case affiliationRoleUpdate.PreviousRole.IsModerator() && affiliationRoleUpdate.NewRole.IsParticipant():
		return i18n.Local("As a result, your role was changed from $role{moderator} to $role{participant}.")
	}
	return i18n.Local("As a result, your role was also changed.")
}

func getAffiliationRoleUpateForAffiliationRemoved(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if affiliationRoleUpdate.Actor == nil {
		return getAffiliationRoleUpateForAffiliationRemovedWithoutActor(affiliationRoleUpdate)
	}
	return getAffiliationRoleUpateForAffiliationRemovedWithActor(affiliationRoleUpdate)
}

func getAffiliationRoleUpateForAffiliationRemovedWithoutActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Localf("%s is not an owner anymore.",
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Localf("%s is not an administrator anymore.",
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember():
		return i18n.Localf("%s is not a member anymore.",
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("%s is not banned anymore.",
		affiliationRoleUpdate.Nickname)
}

func getAffiliationRoleUpateForAffiliationRemovedWithActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.Actor.Affiliation.IsOwner():
		return getAffiliationRoleUpateForAffiliationRemovedWithOwnerActor(affiliationRoleUpdate)
	case affiliationRoleUpdate.Actor.Affiliation.IsAdmin():
		return getAffiliationRoleUpateForAffiliationRemovedWithAdminActor(affiliationRoleUpdate)
	}
	return getAffiliationRoleUpateForAffiliationRemovedForActor(affiliationRoleUpdate)
}

func getAffiliationRoleUpateForAffiliationRemovedWithOwnerActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[2]s is not an owner anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[2]s is not an administrator anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[2]s is not a member anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[2]s is not banned anymore.",
		affiliationRoleUpdate.Actor.Nickname,
		affiliationRoleUpdate.Nickname)
}

func getAffiliationRoleUpateForAffiliationRemovedWithAdminActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[2]s is not an owner anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[2]s is not an administrator anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[2]s is not a member anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[2]s is not banned anymore.",
		affiliationRoleUpdate.Actor.Nickname,
		affiliationRoleUpdate.Nickname)
}

func getAffiliationRoleUpateForAffiliationRemovedForActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Localf("%[1]s changed the position of %[2]s; %[2]s is not an owner anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Localf("%[1]s changed the position of %[2]s; %[2]s is not an administrator anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember():
		return i18n.Localf("%[1]s changed the position of %[2]s; %[2]s is not a member anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("%[1]s changed the position of %[2]s; %[2]s is not banned anymore.",
		affiliationRoleUpdate.Actor.Nickname,
		affiliationRoleUpdate.Nickname)
}

func getAffiliationRoleUpdateForAffiliationAdded(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if affiliationRoleUpdate.Actor == nil {
		return getAffiliationRoleUpdateForAffiliationAddedWithoutActor(affiliationRoleUpdate)
	}
	return getAffiliationRoleUpdateForAffiliationAddedWithActor(affiliationRoleUpdate)
}

func getAffiliationRoleUpdateForAffiliationAddedWithoutActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The position of %s was changed to owner.",
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The position of %s was changed to administrator.",
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The position of %s was changed to member.",
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("The position of %s was changed.",
		affiliationRoleUpdate.Nickname)
}

func getAffiliationRoleUpdateForAffiliationAddedWithActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.Actor.Affiliation.IsOwner():
		return getAffiliationRoleUpdateForAffiliationAddedWithOwnerActor(affiliationRoleUpdate)
	case affiliationRoleUpdate.Actor.Affiliation.IsAdmin():
		return getAffiliationRoleUpdateForAffiliationAddedWithAdminActor(affiliationRoleUpdate)
	}
	return getAffiliationRoleUpdateForAffiliationAddedForActor(affiliationRoleUpdate)
}

func getAffiliationRoleUpdateForAffiliationAddedWithOwnerActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s to owner.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s to administrator.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s to member.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("The owner %[1]s changed the position of %[2]s.",
		affiliationRoleUpdate.Actor.Nickname,
		affiliationRoleUpdate.Nickname)
}

func getAffiliationRoleUpdateForAffiliationAddedWithAdminActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s to owner.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s to administrator.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s to member.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("The administrator %[1]s changed the position of %[2]s.",
		affiliationRoleUpdate.Actor.Nickname,
		affiliationRoleUpdate.Nickname)
}

func getAffiliationRoleUpdateForAffiliationAddedForActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("%[1]s changed the position of %[2]s to owner.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("%[1]s changed the position of %[2]s to administrator.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("%[1]s changed the position of %[2]s to member.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("%[1]s changed the position of %[2]s.",
		affiliationRoleUpdate.Actor.Nickname,
		affiliationRoleUpdate.Nickname)
}

func getAffiliationRoleUpdateForAffiliationUpdated(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if affiliationRoleUpdate.Actor == nil {
		return getAffiliationRoleUpdateForAffiliationUpdatedWithoutActor(affiliationRoleUpdate)
	}
	return getAffiliationRoleUpdateForAffiliationUpdatedWithActor(affiliationRoleUpdate)
}

func getAffiliationRoleUpdateForAffiliationUpdatedWithoutActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The position of %s was changed from owner to member.",
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The position of %s was changed from administrator to member.",
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The position of %s was changed from owner to administrator.",
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The position of %s was changed from member to administrator.",
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The position of %s was changed from administrator to owner.",
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The position of %s was changed from member to owner.",
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("The position of %s was changed.",
		affiliationRoleUpdate.Nickname)
}

func getAffiliationRoleUpdateForAffiliationUpdatedWithActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.Actor.Affiliation.IsOwner():
		return getAffiliationRoleUpdateForAffiliationUpdatedWithOwnerActor(affiliationRoleUpdate)
	case affiliationRoleUpdate.Actor.Affiliation.IsAdmin():
		return getAffiliationRoleUpdateForAffiliationUpdatedWithAdminActor(affiliationRoleUpdate)
	}
	return getAffiliationRoleUpdateForAffiliationUpdatedForActor(affiliationRoleUpdate)
}

func getAffiliationRoleUpdateForAffiliationUpdatedWithOwnerActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from member to administrator.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from member to owner.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from administrator to member.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from administrator to owner.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from owner to member.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from owner to administrator.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("The owner %[1]s changed the position of %[2]s.",
		affiliationRoleUpdate.Nickname)
}

func getAffiliationRoleUpdateForAffiliationUpdatedWithAdminActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from member to administrator.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from member to owner.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from administrator to member.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from administrator to owner.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from owner to member.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from owner to administrator.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("The administrator %[1]s changed the position of %[2]s.",
		affiliationRoleUpdate.Actor.Nickname,
		affiliationRoleUpdate.Nickname)
}

func getAffiliationRoleUpdateForAffiliationUpdatedForActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("%[1]s changed the position of %[2]s from member to administrator.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("%[1]s changed the position of %[2]s from member to owner.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("%[1]s changed the position of %[2]s from administrator to member.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("%[1]s changed the position of %[2]s from administrator to owner.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("%[1]s changed the position of %[2]s from owner to member.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("%[1]s changed the position of %[2]s from owner to administrator.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname)
	}
	return i18n.Localf("%[1]s changed the position of %[2]s.",
		affiliationRoleUpdate.Actor.Nickname,
		affiliationRoleUpdate.Nickname)
}

func getSelfAffiliationRoleUpdateMessage(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	affiliationMessage := getSelfAffiliationRoleUpdateMessageForAffiliation(selfAffiliationRoleUpdate)
	roleMessage := getSelfAffiliationRoleUpdateMessageForRole(selfAffiliationRoleUpdate)

	return appendReasonToMessage(
		i18n.Localf("%[1]s %[2]s", affiliationMessage, roleMessage),
		selfAffiliationRoleUpdate.Reason,
	)
}

func getSelfAffiliationRoleUpdateMessageForAffiliation(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.NewAffiliation.IsNone():
		return getSelfAffiliationRoleUpateForAffiliationRemoved(selfAffiliationRoleUpdate)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsNone():
		return getSelfAffiliationRoleUpdateForAffiliationAdded(selfAffiliationRoleUpdate)
	}
	return getSelfAffiliationRoleUpdateForAffiliationUpdated(selfAffiliationRoleUpdate)
}

func getSelfAffiliationRoleUpateForAffiliationRemoved(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if selfAffiliationRoleUpdate.Actor == nil {
		return getSelfAffiliationRoleUpateForAffiliationRemovedWithoutActor(selfAffiliationRoleUpdate)
	}

	return getSelfAffiliationRoleUpateForAffiliationRemovedWithActor(selfAffiliationRoleUpdate)
}

func getSelfAffiliationRoleUpateForAffiliationRemovedWithoutActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Local("You are not $affiliation{an owner} anymore.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Local("You are not $affiliation{an administrator} anymore.")
	}
	return i18n.Local("You are not $affiliation{a member} anymore.")
}

func getSelfAffiliationRoleUpateForAffiliationRemovedWithActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Localf("The owner $nickname{%s} changed your position; you are not $affiliation{an owner} anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Localf("The owner $nickname{%s} changed your position; you are not $affiliation{an administrator} anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.PreviousAffiliation.IsMember():
		return i18n.Localf("The owner $nickname{%s} changed your position; you are not $affiliation{a member} anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsAdmin() && selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Localf("The administrator $nickname{%s} changed your position; you are not $affiliation{an owner} anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsAdmin() && selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Localf("The administrator$nickname{%s} changed your position; you are not $affiliation{an administrator} anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsAdmin() && selfAffiliationRoleUpdate.PreviousAffiliation.IsMember():
		return i18n.Localf("The administrator $nickname{%s} changed your position; you are not $affiliation{a member} anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	}
	return i18n.Localf("The administrator $nickname{%s} changed your position.",
		selfAffiliationRoleUpdate.Actor.Nickname)
}

func getSelfAffiliationRoleUpdateForAffiliationAdded(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if selfAffiliationRoleUpdate.Actor == nil {
		return getSelfAffiliationRoleUpdateForAffiliationAddedWithoutActor(selfAffiliationRoleUpdate)
	}
	return getSelfAffiliationRoleUpdateForAffiliationAddedWithActor(selfAffiliationRoleUpdate)
}

func getSelfAffiliationRoleUpdateForAffiliationAddedWithoutActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("Your position was changed to $affiliation{owner}.")
	case selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("Your position was changed to $affiliation{administrator}.")
	}
	return i18n.Localf("Your position was changed to $affiliation{member}.")
}

func getSelfAffiliationRoleUpdateForAffiliationAddedWithActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The owner $nickname{%s} changed your position to $affiliation{owner}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The owner $nickname{%s} changed your position to $affiliation{administrator}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The owner $nickname{%s} changed your position to $affiliation{member}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The administrator $nickname{%s} changed your position to $affiliation{owner}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The administrator $nickname{%s} changed your position to $affiliation{administrator}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	}
	return i18n.Localf("The administrator $nickname{%s} changed your position to $affiliation{member}.",
		selfAffiliationRoleUpdate.Actor.Nickname)
}

func getSelfAffiliationRoleUpdateForAffiliationUpdated(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if selfAffiliationRoleUpdate.Actor == nil {
		return getSelfAffiliationRoleUpdateForAffiliationUpdatedWithoutActor(selfAffiliationRoleUpdate)
	}

	return getSelfAffiliationRoleUpdateForAffiliationUpdatedWithActor(selfAffiliationRoleUpdate)
}

func getSelfAffiliationRoleUpdateForAffiliationUpdatedWithoutActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("Your position was changed from $affiliation{owner} to $affiliation{administrator}.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("Your position was changed from $affiliation{owner} to $affiliation{member}.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("Your position was changed from $affiliation{administrator} to $affiliation{owner}.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("Your position was changed from $affiliation{administrator} to $affiliation{member}.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("Your position was changed from $affiliation{member} to $affiliation{owner}.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("Your position was changed from $affiliation{member} to $affiliation{administrator}.")
	}
	return i18n.Localf("Your position was changed.")
}

func getSelfAffiliationRoleUpdateForAffiliationUpdatedWithActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) (m string) {
	switch {
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner():
		return getSelfAffiliationRoleUpdateForAffiliationUpdatedWithOwnerActor(selfAffiliationRoleUpdate)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsAdmin():
		return getSelfAffiliationRoleUpdateForAffiliationUpdatedWithAdminActor(selfAffiliationRoleUpdate)
	}
	return getSelfAffiliationRoleUpdateForAffiliationUpdatedForActor(selfAffiliationRoleUpdate)
}

func getSelfAffiliationRoleUpdateForAffiliationUpdatedWithOwnerActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{owner to administrator.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{owner to member.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{administrator} to $affiliation{owner}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{administrator} to $affiliation{member}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{member} to $affiliation{owner}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{member} to $affiliation{administrator}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	}
	return i18n.Localf("The owner $nickname{%s} changed your position.",
		selfAffiliationRoleUpdate.Actor.Nickname)
}

func getSelfAffiliationRoleUpdateForAffiliationUpdatedWithAdminActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{owner} to $affiliation{administrator}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{owner} to $affiliation{member}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{administrator} to $affiliation{owner}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{administrator} to $affiliation{member}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{member} to $affiliation{owner}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{member} to $affiliation{administrator}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	}
	return i18n.Localf("The administrator $nickname{%s} changed your position.",
		selfAffiliationRoleUpdate.Actor.Nickname)
}

func getSelfAffiliationRoleUpdateForAffiliationUpdatedForActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("$nickname{%s} changed your position from $affiliation{owner} to $affiliation{administrator}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("$nickname{%s} changed your position from $affiliation{owner} to $affiliation{member}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("$nickname{%s} changed your position from $affiliation{administrator} to $affiliation{owner}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("$nickname{%s} changed your position from $affiliation{administrator} to $affiliation{member}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("$nickname{%s} changed your position from $affiliation{member} to $affiliation{owner}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("$nickname{%s} changed your position from $affiliation{member} to $affiliation{administrator}.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	}
	return i18n.Localf("$nickname{%s} changed your position.",
		selfAffiliationRoleUpdate.Actor.Nickname)
}

func getSelfAffiliationUpdateMessage(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	return appendReasonToMessage(getSelfAffiliationUpdateBaseMessage(selfAffiliationUpdate), selfAffiliationUpdate.Reason)
}

func getSelfAffiliationUpdateBaseMessage(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.New.IsNone():
		return getSelfAffiliationRemovedMessage(selfAffiliationUpdate)
	case selfAffiliationUpdate.New.IsBanned():
		return getSelfAffiliationBannedMessage(selfAffiliationUpdate)
	case selfAffiliationUpdate.Previous.IsNone():
		return getSelfAffiliationAddedMessage(selfAffiliationUpdate)
	}
	return getSelfAffiliationChangedMessage(selfAffiliationUpdate)
}

func getSelfAffiliationRemovedMessage(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor == nil {
		return getSelfAffiliationRemovedMessageWithoutActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationRemovedMessageWithActor(selfAffiliationUpdate)
}

func getSelfAffiliationRemovedMessageWithoutActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner():
		return i18n.Local("You are not $affiliation{an owner} anymore.")
	case selfAffiliationUpdate.Previous.IsAdmin():
		return i18n.Local("You are not $affiliation{an administrator} anymore.")
	case selfAffiliationUpdate.Previous.IsMember():
		return i18n.Local("You are not $affiliation{a member} anymore.")
	}
	return i18n.Local("You are not $affiliation{banned} anymore.")
}

func getSelfAffiliationRemovedMessageWithActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor.Affiliation.IsOwner() {
		return getSelfAffiliationRemovedMessageWithOwnerActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationRemovedMessageWithAdminActor(selfAffiliationUpdate)
}

func getSelfAffiliationRemovedMessageWithOwnerActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner():
		return i18n.Localf("The owner $nickname{%s} changed your position; you are not $affiliation{an owner} anymore.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin():
		return i18n.Localf("The owner $nickname{%s} changed your position; you are not $affiliation{an administrator} anymore.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember():
		return i18n.Localf("The owner $nickname{%s} changed your position; you are not $affiliation{a member} anymore.",
			selfAffiliationUpdate.Actor.Nickname)
	}
	return i18n.Localf("The owner $nickname{%s} changed your position; you are not $affiliation{banned} anymore.",
		selfAffiliationUpdate.Actor.Nickname)
}

func getSelfAffiliationRemovedMessageWithAdminActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner():
		return i18n.Localf("The administrator $nickname{%s} changed your position; you are not $affiliation{an owner} anymore.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin():
		return i18n.Localf("The administrator $nickname{%s} changed your position; you are not $affiliation{an administrator} anymore.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember():
		return i18n.Localf("The administrator $nickname{%s} changed your position; you are not $affiliation{a member} anymore.",
			selfAffiliationUpdate.Actor.Nickname)
	}
	return i18n.Localf("The administrator $nickname{%s} changed your position; you are not $affiliation{banned} anymore.",
		selfAffiliationUpdate.Actor.Nickname)
}

func getSelfAffiliationAddedMessage(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor == nil {
		return getSelfAffiliationAddedMessageWithoutActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationAddedMessageWithActor(selfAffiliationUpdate)
}

func getSelfAffiliationAddedMessageWithoutActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.New.IsOwner():
		return i18n.Local("You are now $affiliation{an owner}.")
	case selfAffiliationUpdate.New.IsAdmin():
		return i18n.Local("You are now $affiliation{an administrator}.")
	case selfAffiliationUpdate.New.IsMember():
		return i18n.Local("You are now $affiliation{a member}.")
	}
	return i18n.Local("You are now $affiliation{banned}.")
}

func getSelfAffiliationAddedMessageWithActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor.Affiliation.IsOwner() {
		return getSelfAffiliationAddedMessageWithOwnerActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationAddedMessageWithAdminActor(selfAffiliationUpdate)
}

func getSelfAffiliationAddedMessageWithOwnerActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner $nickname{%s} changed your position; you are now $affiliation{an owner}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner $nickname{%s} changed your position; you are now $affiliation{an administrator}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The owner $nickname{%s} changed your position; you are now $affiliation{a member}.",
			selfAffiliationUpdate.Actor.Nickname)
	}
	return i18n.Localf("The owner $nickname{%s} changed your position; you are now $affiliation{banned}.",
		selfAffiliationUpdate.Actor.Nickname)
}

func getSelfAffiliationAddedMessageWithAdminActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator $nickname{%s} changed your position; you are now $affiliation{an owner}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator $nickname{%s} changed your position; you are now $affiliation{an administrator}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator $nickname{%s} changed your position; you are now $affiliation{a member}.",
			selfAffiliationUpdate.Actor.Nickname)
	}
	return i18n.Localf("The administrator $nickname{%s} changed your position; you are now $affiliation{banned}.",
		selfAffiliationUpdate.Actor.Nickname)
}

func getSelfAffiliationChangedMessage(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor == nil {
		return getSelfAffiliationChangedMessageWithoutActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationChangedMessageWithActor(selfAffiliationUpdate)
}

func getSelfAffiliationChangedMessageWithoutActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Local("Your position was changed from $affiliation{owner} to $affiliation{administrator}.")
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsMember():
		return i18n.Local("Your position was changed from $affiliation{owner} to $affiliation{member}.")
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Local("Your position was changed from $affiliation{administrator} to $affiliation{owner}.")
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsMember():
		return i18n.Local("Your position was changed from $affiliation{administrator} to $affiliation{member}.")
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Local("Your position was changed from $affiliation{member} to $affiliation{administrator}.")
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Local("Your position was changed from $affiliation{member} to $affiliation{owner}.")
	}
	return i18n.Local("Your position was changed.")
}

func getSelfAffiliationChangedMessageWithActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor.Affiliation.IsOwner() {
		return getSelfAffiliationChangedMessageWithOwnerActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationChangedMessageWithAdminActor(selfAffiliationUpdate)
}

func getSelfAffiliationChangedMessageWithOwnerActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{owner} to $affiliation{administrator}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{owner} to $affiliation{member}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{administrator} to $affiliation{owner}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{administrator} to $affiliation{member}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{member} to $affiliation{administrator}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner $nickname{%s} changed your position from $affiliation{member} to $affiliation{owner}.",
			selfAffiliationUpdate.Actor.Nickname)
	}
	return i18n.Localf("The owner $nickname{%s} changed your position.",
		selfAffiliationUpdate.Actor.Nickname)
}

func getSelfAffiliationChangedMessageWithAdminActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{owner} to $affiliation{administrator}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{owner} to $affiliation{member}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{administrator} to $affiliation{owner}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{administrator} to $affiliation{member}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{member} to $affiliation{administrator}.",
			selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator $nickname{%s} changed your position from $affiliation{member} to $affiliation{owner}.",
			selfAffiliationUpdate.Actor.Nickname)
	}
	return i18n.Localf("The administrator $nickname{%s} changed your position.",
		selfAffiliationUpdate.Actor.Nickname)
}

func getSelfAffiliationBannedMessage(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor == nil {
		return i18n.Localf("You have been $affiliation{banned} from the room.")
	}
	return getSelfAffiliationBannedMessageWithActor(selfAffiliationUpdate)
}

func getSelfAffiliationBannedMessageWithActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor.Affiliation.IsOwner() {
		return i18n.Localf("The owner $nickname{%s} $affiliation{banned} you from the room.",
			selfAffiliationUpdate.Actor.Nickname)
	}
	return i18n.Localf("The administrator $nickname{%s} $affiliation{banned} you from the room.",
		selfAffiliationUpdate.Actor.Nickname)
}

func appendReasonToMessage(message, reason string) string {
	if reason != "" {
		return i18n.Localf("%[1]s The reason given was: %[2]s.", message, reason)
	}
	return message
}
