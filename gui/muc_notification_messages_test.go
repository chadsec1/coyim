package gui

import (
	"fmt"
	"regexp"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/glib_mock"
	. "gopkg.in/check.v1"
)

type MUCNotificationMessagesSuite struct{}

var _ = Suite(&MUCNotificationMessagesSuite{})

func (s *MUCNotificationMessagesSuite) SetUpSuite(c *C) {
	initMUCNotificationMessagesI18n()
}

type mucNotificationMessagesMockGlib struct {
	glib_mock.Mock
}

func (*mucNotificationMessagesMockGlib) Local(vx string) string {
	return "[localized] " + removePresentationFormatsFromString(vx)
}

func (*mucNotificationMessagesMockGlib) Localf(vx string, args ...interface{}) string {
	return fmt.Sprintf("[localized] "+removePresentationFormatsFromString(vx), args...)
}

func removePresentationFormatsFromString(s string) string {
	regex := regexp.MustCompile("{{ [a-z]* \"(.*?)\" }}")
	return regex.ReplaceAllString(s, "${1}")
}

func initMUCNotificationMessagesI18n() {
	i18n.InitLocalization(&mucNotificationMessagesMockGlib{})
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationNone(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "batman",
		New:      newTestAffiliationFromString(data.AffiliationNone),
		Previous: newTestAffiliationFromString(data.AffiliationAdmin),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] batman is not an administrator anymore.")

	au.Reason = "batman lost his mind"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] batman is not an administrator anymore. The reason given was: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] batman is not an owner anymore. The reason given was: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] batman is not a member anymore. The reason given was: batman lost his mind.")

	au.Reason = ""
	au.Previous = newTestAffiliationFromString(data.AffiliationAdmin)
	au.Actor = newTestActor("robin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] The owner robin changed the position of batman; batman is not an administrator anymore.")

	au.Reason = "batman lost his mind"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The owner robin changed the position of batman; batman is not an administrator anymore. The reason given was: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The owner robin changed the position of batman; batman is not an owner anymore. The reason given was: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The owner robin changed the position of batman; batman is not a member anymore. The reason given was: batman lost his mind.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationOutcast(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "alice",
		New:      newTestAffiliationFromString(data.AffiliationOutcast),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] alice was banned from the room.")

	au.Reason = "she was rude"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] alice was banned from the room. The reason given was: she was rude.")

	au.Reason = ""
	au.Actor = newTestActor("bob", newTestAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] The administrator bob banned alice from the room.")

	au.Reason = "she was rude"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The administrator bob banned alice from the room. The reason given was: she was rude.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationAdded(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "juanito",
		New:      newTestAffiliationFromString(data.AffiliationMember),
		Previous: newTestAffiliationFromString(data.AffiliationNone),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] juanito is now a member.")

	au.Reason = "el es súper chévere"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] juanito is now a member. The reason given was: el es súper chévere.")

	au.Reason = ""
	au.Actor = newTestActor("pepito", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] The owner pepito changed the position of juanito; juanito is now a member.")

	au.Reason = "el es súper chévere"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The owner pepito changed the position of juanito; juanito is now a member. The reason given was: el es súper chévere.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationChanged(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "thor",
		New:      newTestAffiliationFromString(data.AffiliationAdmin),
		Previous: newTestAffiliationFromString(data.AffiliationMember),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] The position of thor was changed from member to administrator.")

	au.Reason = "he is the strongest avenger"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The position of thor was changed from member to administrator. The reason given was: he is the strongest avenger.")

	au.Reason = ""
	au.Actor = newTestActor("odin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] The owner odin changed the position of thor from member to administrator.")

	au.Reason = "he is the strongest avenger"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The owner odin changed the position of thor from member to administrator. The reason given was: he is the strongest avenger.")
}

func (*MUCNotificationMessagesSuite) Test_getMUCNotificationMessageFrom_affiliationUpdate(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "chavo",
		New:      newTestAffiliationFromString(data.AffiliationAdmin),
		Previous: newTestAffiliationFromString(data.AffiliationMember),
	}

	c.Assert(getMUCNotificationMessageFrom(au), Equals, "[localized] The position of chavo was changed from member to administrator.")

	au.Previous = newTestAffiliationFromString(data.AffiliationNone)
	c.Assert(getMUCNotificationMessageFrom(au), Equals, "[localized] chavo is now an administrator.")

	au.New = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getMUCNotificationMessageFrom(au), Equals, "[localized] chavo is now an owner.")

	au.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	au.New = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getMUCNotificationMessageFrom(au), Equals, "[localized] The position of chavo was changed from owner to member.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateMessage_roleModerator(c *C) {
	ru := data.RoleUpdate{
		Nickname: "wanda",
		New:      newTestRoleFromString(data.RoleModerator),
		Previous: newTestRoleFromString(data.RoleParticipant),
	}

	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The role of wanda was changed from participant to moderator.")

	ru.Reason = "vision wanted it"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The role of wanda was changed from participant to moderator. The reason given was: vision wanted it.")

	ru.Reason = ""
	ru.Actor = newTestActor("vision", newTestAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The administrator vision changed the role of wanda from participant to moderator.")

	ru.Reason = "vision wanted it"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The administrator vision changed the role of wanda from participant to moderator. The reason given was: vision wanted it.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateMessage_roleParticipant(c *C) {
	ru := data.RoleUpdate{
		Nickname: "sancho",
		New:      newTestRoleFromString(data.RoleParticipant),
		Previous: newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The role of sancho was changed from moderator to participant.")

	ru.Reason = "los molinos son gigantes"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The role of sancho was changed from moderator to participant. The reason given was: los molinos son gigantes.")

	ru.Reason = ""
	ru.Actor = newTestActor("panza", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The owner panza changed the role of sancho from moderator to participant.")

	ru.Reason = "los molinos son gigantes"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The owner panza changed the role of sancho from moderator to participant. The reason given was: los molinos son gigantes.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateMessage_roleVisitor(c *C) {
	ru := data.RoleUpdate{
		Nickname: "chapulin",
		New:      newTestRoleFromString(data.RoleVisitor),
		Previous: newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The role of chapulin was changed from moderator to visitor.")

	ru.Reason = "no contaban con mi astucia"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The role of chapulin was changed from moderator to visitor. The reason given was: no contaban con mi astucia.")

	ru.Reason = ""
	ru.Actor = newTestActor("chespirito", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The owner chespirito changed the role of chapulin from moderator to visitor.")

	ru.Reason = "no contaban con mi astucia"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The owner chespirito changed the role of chapulin from moderator to visitor. The reason given was: no contaban con mi astucia.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateMessage_roleNone(c *C) {
	ru := data.RoleUpdate{
		Nickname: "alberto",
		New:      newTestRoleFromString(data.RoleNone),
		Previous: newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] alberto was temporarily removed from the room.")

	ru.Reason = "bla"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] alberto was temporarily removed from the room. The reason given was: bla.")

	ru.Reason = ""
	ru.Actor = newTestActor("foo", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The owner foo temporarily removed alberto from the room.")

	ru.Reason = "bla"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The owner foo temporarily removed alberto from the room. The reason given was: bla.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfRoleUpdateMessage_roleModerator(c *C) {
	sru := data.RoleUpdate{
		Nickname: "wanda",
		New:      newTestRoleFromString(data.RoleModerator),
		Previous: newTestRoleFromString(data.RoleParticipant),
	}

	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] Your role was changed from participant to moderator.")

	sru.Reason = "vision wanted it"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] Your role was changed from participant to moderator. The reason given was: vision wanted it.")

	sru.Reason = ""
	sru.Actor = newTestActor("vision", newTestAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] The administrator vision changed your role from participant to moderator.")

	sru.Reason = "vision wanted it"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] The administrator vision changed your role from participant to moderator. The reason given was: vision wanted it.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfRoleUpdateMessage_roleParticipant(c *C) {
	sru := data.RoleUpdate{
		Nickname: "sancho",
		New:      newTestRoleFromString(data.RoleParticipant),
		Previous: newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] Your role was changed from moderator to participant.")

	sru.Reason = "los molinos son gigantes"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] Your role was changed from moderator to participant. The reason given was: los molinos son gigantes.")

	sru.Reason = ""
	sru.Actor = newTestActor("panza", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] The owner panza changed your role from moderator to participant.")

	sru.Reason = "los molinos son gigantes"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] The owner panza changed your role from moderator to participant. The reason given was: los molinos son gigantes.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfRoleUpdateMessage_roleVisitor(c *C) {
	sru := data.RoleUpdate{
		Nickname: "chapulin",
		New:      newTestRoleFromString(data.RoleVisitor),
		Previous: newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] Your role was changed from moderator to visitor.")

	sru.Reason = "no contaban con mi astucia"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] Your role was changed from moderator to visitor. The reason given was: no contaban con mi astucia.")

	sru.Reason = ""
	sru.Actor = newTestActor("chespirito", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] The owner chespirito changed your role from moderator to visitor.")

	sru.Reason = "no contaban con mi astucia"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] The owner chespirito changed your role from moderator to visitor. The reason given was: no contaban con mi astucia.")
}

func (*MUCNotificationMessagesSuite) Test_getMUCNotificationMessageFrom_roleUpdate(c *C) {
	ru := data.RoleUpdate{
		Nickname: "pablo",
		New:      newTestRoleFromString(data.RoleModerator),
		Previous: newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getMUCNotificationMessageFrom(ru), Equals, "[localized] The role of pablo was changed from visitor to moderator.")

	ru.Previous = newTestRoleFromString(data.RoleParticipant)
	c.Assert(getMUCNotificationMessageFrom(ru), Equals, "[localized] The role of pablo was changed from participant to moderator.")

	ru.New = newTestRoleFromString(data.RoleVisitor)
	c.Assert(getMUCNotificationMessageFrom(ru), Equals, "[localized] The role of pablo was changed from participant to visitor.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationRoleUpdateMessage_affiliationRemoved(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "007",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationNone),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationAdmin),
		NewRole:             newTestRoleFromString(data.RoleVisitor),
		PreviousRole:        newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] 007 is not an administrator anymore. [localized] As a result, their role was changed from moderator to visitor.")

	aru.Reason = "he is an assassin"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] 007 is not an administrator anymore. [localized] As a result, their role was changed from moderator to visitor. The reason given was: he is an assassin.")

	aru.Reason = ""
	aru.Actor = newTestActor("the enemy", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] The owner the enemy changed the position of 007; 007 is not an administrator anymore. [localized] As a result, their role was changed from moderator to visitor.")

	aru.Reason = "bla"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] The owner the enemy changed the position of 007; 007 is not an administrator anymore. [localized] As a result, their role was changed from moderator to visitor. The reason given was: bla.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationRoleUpdateMessage_affiliationAdded(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "alice",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationAdmin),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationNone),
		NewRole:             newTestRoleFromString(data.RoleModerator),
		PreviousRole:        newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] The position of alice was changed to administrator. [localized] As a result, their role was changed from visitor to moderator.")

	aru.Reason = "she is lost in the world of wonders"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] The position of alice was changed to administrator. [localized] As a result, their role was changed from visitor to moderator. The reason given was: she is lost in the world of wonders.")

	aru.Reason = ""
	aru.Actor = newTestActor("rabbit", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] The owner rabbit changed the position of alice to administrator. [localized] As a result, their role was changed from visitor to moderator.")

	aru.Reason = "she is lost in the world of wonders"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] The owner rabbit changed the position of alice to administrator. [localized] As a result, their role was changed from visitor to moderator. The reason given was: she is lost in the world of wonders.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationRoleUpdateMessage_affiliationUpdated(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "Pegassus",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationAdmin),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationOwner),
		NewRole:             newTestRoleFromString(data.RoleModerator),
		PreviousRole:        newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] The position of Pegassus was changed from owner to administrator. [localized] As a result, their role was changed from visitor to moderator.")

	aru.Reason = "he is a silver warrior"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] The position of Pegassus was changed from owner to administrator. [localized] As a result, their role was changed from visitor to moderator. The reason given was: he is a silver warrior.")

	aru.Reason = ""
	aru.Actor = newTestActor("Ikki", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] The owner Ikki changed the position of Pegassus from owner to administrator. [localized] As a result, their role was changed from visitor to moderator.")

	aru.Reason = "he has the phoenix flame"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] The owner Ikki changed the position of Pegassus from owner to administrator. [localized] As a result, their role was changed from visitor to moderator. The reason given was: he has the phoenix flame.")
}

func (*MUCNotificationMessagesSuite) Test_getMUCNotificationMessageFrom_affiliationRoleUpdate(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "chavo",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationNone),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationAdmin),
		NewRole:             newTestRoleFromString(data.RoleVisitor),
		PreviousRole:        newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getMUCNotificationMessageFrom(aru), Equals, "[localized] [localized] chavo is not an administrator anymore. [localized] As a result, their role was changed from moderator to visitor.")

	aru.NewAffiliation = newTestAffiliationFromString(data.AffiliationAdmin)
	aru.PreviousAffiliation = newTestAffiliationFromString(data.AffiliationNone)
	aru.NewRole = newTestRoleFromString(data.RoleModerator)
	aru.PreviousRole = newTestRoleFromString(data.RoleVisitor)
	c.Assert(getMUCNotificationMessageFrom(aru), Equals, "[localized] [localized] The position of chavo was changed to administrator. [localized] As a result, their role was changed from visitor to moderator.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationRemoved(c *C) {
	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New:      newTestAffiliationFromString(data.AffiliationNone),
			Previous: newTestAffiliationFromString(data.AffiliationAdmin),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You are not an administrator anymore.")

	sau.Reason = "you are funny"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You are not an administrator anymore. The reason given was: you are funny.")

	sau.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You are not an owner anymore. The reason given was: you are funny.")

	sau.Reason = ""
	sau.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You are not a member anymore.")

	sau.Actor = newTestActor("robin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] The owner robin changed your position; you are not a member anymore.")

	sau.Reason = "you are funny"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The owner robin changed your position; you are not a member anymore. The reason given was: you are funny.")

	sau.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The owner robin changed your position; you are not an owner anymore. The reason given was: you are funny.")

	sau.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The owner robin changed your position; you are not a member anymore. The reason given was: you are funny.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationAdded(c *C) {
	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New:      newTestAffiliationFromString(data.AffiliationAdmin),
			Previous: newTestAffiliationFromString(data.AffiliationNone),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You are now an administrator.")

	sau.Reason = "estás encopetao"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You are now an administrator. The reason given was: estás encopetao.")

	sau.Reason = ""
	sau.AffiliationUpdate.New = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You are now a member.")

	sau.Reason = "you dance very well"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You are now a member. The reason given was: you dance very well.")

	sau.Reason = ""
	sau.AffiliationUpdate.New = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You are now an owner.")

	sau.Reason = "the day is cool"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You are now an owner. The reason given was: the day is cool.")

	sau.Reason = ""
	sau.AffiliationUpdate.New = newTestAffiliationFromString(data.AffiliationAdmin)
	sau.Actor = newTestActor("paco", newTestAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] The administrator paco changed your position; you are now an administrator.")

	sau.Reason = "you are funny"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The administrator paco changed your position; you are now an administrator. The reason given was: you are funny.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationChanged(c *C) {
	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New:      newTestAffiliationFromString(data.AffiliationAdmin),
			Previous: newTestAffiliationFromString(data.AffiliationMember),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] Your position was changed from member to administrator.")

	sau.Reason = "you are loco"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] Your position was changed from member to administrator. The reason given was: you are loco.")

	sau.Reason = ""
	sau.Actor = newTestActor("chapulin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] The owner chapulin changed your position from member to administrator.")

	sau.Reason = "you are locote"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The owner chapulin changed your position from member to administrator. The reason given was: you are locote.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateSuccessMessage(c *C) {
	nickname := "Juan"
	owner := newTestAffiliationFromString(data.AffiliationOwner)
	admin := newTestAffiliationFromString(data.AffiliationAdmin)
	member := newTestAffiliationFromString(data.AffiliationMember)
	outcast := newTestAffiliationFromString(data.AffiliationOutcast)
	none := newTestAffiliationFromString(data.AffiliationNone)

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, member, none), Equals,
		"[localized] $nickname{Juan} is not $affiliation{a member} anymore.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, admin, none), Equals,
		"[localized] $nickname{Juan} is not $affiliation{an administrator} anymore.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, owner, none), Equals,
		"[localized] $nickname{Juan} is not $affiliation{an owner} anymore.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, member), Equals,
		"[localized] The position of $nickname{Juan} was changed to $affiliation{member}.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, admin), Equals,
		"[localized] The position of $nickname{Juan} was changed to $affiliation{administrator}.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, owner), Equals,
		"[localized] The position of $nickname{Juan} was changed to $affiliation{owner}.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, outcast), Equals,
		"[localized] $nickname{Juan} has been $affiliation{banned} from the room.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, member, outcast), Equals,
		"[localized] $nickname{Juan} has been $affiliation{banned} from the room.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, admin, outcast), Equals,
		"[localized] $nickname{Juan} has been $affiliation{banned} from the room.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, owner, outcast), Equals,
		"[localized] $nickname{Juan} has been $affiliation{banned} from the room.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateSuccessMessage(c *C) {
	moderator := newTestRoleFromString(data.RoleModerator)
	participant := newTestRoleFromString(data.RoleParticipant)
	visitor := newTestRoleFromString(data.RoleVisitor)
	none := newTestRoleFromString(data.RoleNone)

	c.Assert(getRoleUpdateSuccessMessage("Maria", moderator, none), Equals, "[localized] $nickname{Maria} was temporarily removed from the room.")
	c.Assert(getRoleUpdateSuccessMessage("Carlos", participant, none), Equals, "[localized] $nickname{Carlos} was temporarily removed from the room.")
	c.Assert(getRoleUpdateSuccessMessage("Mauricio", visitor, none), Equals, "[localized] $nickname{Mauricio} was temporarily removed from the room.")

	c.Assert(getRoleUpdateSuccessMessage("Jose", none, moderator), Equals, "[localized] The role of $nickname{Jose} was changed to moderator.")
	c.Assert(getRoleUpdateSuccessMessage("Alberto", none, participant), Equals, "[localized] The role of $nickname{Alberto} was changed to participant.")
	c.Assert(getRoleUpdateSuccessMessage("Juan", none, visitor), Equals, "[localized] The role of $nickname{Juan} was changed to visitor.")

	c.Assert(getRoleUpdateSuccessMessage("Alberto", moderator, participant), Equals, "[localized] The role of $nickname{Alberto} was changed from moderator to participant.")
	c.Assert(getRoleUpdateSuccessMessage("Alberto", moderator, visitor), Equals, "[localized] The role of $nickname{Alberto} was changed from moderator to visitor.")
	c.Assert(getRoleUpdateSuccessMessage("Alberto", participant, moderator), Equals, "[localized] The role of $nickname{Alberto} was changed from participant to moderator.")
	c.Assert(getRoleUpdateSuccessMessage("Carlos", participant, visitor), Equals, "[localized] The role of $nickname{Carlos} was changed from participant to visitor.")
	c.Assert(getRoleUpdateSuccessMessage("Carlos", visitor, participant), Equals, "[localized] The role of $nickname{Carlos} was changed from visitor to participant.")
	c.Assert(getRoleUpdateSuccessMessage("Juan", visitor, moderator), Equals, "[localized] The role of $nickname{Juan} was changed from visitor to moderator.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateFailureMessage(c *C) {
	owner := newTestAffiliationFromString(data.AffiliationOwner)
	admin := newTestAffiliationFromString(data.AffiliationAdmin)
	member := newTestAffiliationFromString(data.AffiliationMember)
	outcast := newTestAffiliationFromString(data.AffiliationOutcast)
	none := newTestAffiliationFromString(data.AffiliationNone)

	messages := getAffiliationUpdateFailureMessage("Luisa", owner, nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] The position of $nickname{Luisa} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the position failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The position of Luisa couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the position of Luisa to owner.")

	messages = getAffiliationUpdateFailureMessage("Marco", admin, nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] The position of $nickname{Marco} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the position failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The position of Marco couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the position of Marco to administrator.")

	messages = getAffiliationUpdateFailureMessage("Pedro", member, nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] The position of $nickname{Pedro} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the position failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The position of Pedro couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the position of Pedro to member.")

	messages = getAffiliationUpdateFailureMessage("Luisa", outcast, nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] $nickname{Luisa} couldn't be banned.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Banning failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] Luisa couldn't be banned")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to ban Luisa.")

	messages = getAffiliationUpdateFailureMessage("José", none, nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] The position of $nickname{José} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the position failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The position of José couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the position of José.")

}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateFailureMessage(c *C) {
	moderator := newTestRoleFromString(data.RoleModerator)
	participant := newTestRoleFromString(data.RoleParticipant)
	visitor := newTestRoleFromString(data.RoleVisitor)
	none := newTestRoleFromString(data.RoleNone)

	messages := getRoleUpdateFailureMessage("Mauricio", moderator)
	c.Assert(messages.notificationMessage, Equals, "[localized] The role of $nickname{Mauricio} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the role failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The role of Mauricio couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the role of Mauricio to moderator.")

	messages = getRoleUpdateFailureMessage("Juan", participant)
	c.Assert(messages.notificationMessage, Equals, "[localized] The role of $nickname{Juan} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the role failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The role of Juan couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the role of Juan to participant.")

	messages = getRoleUpdateFailureMessage("Pepe", visitor)
	c.Assert(messages.notificationMessage, Equals, "[localized] The role of $nickname{Pepe} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the role failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The role of Pepe couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the role of Pepe to visitor.")

	messages = getRoleUpdateFailureMessage("Juana", none)
	c.Assert(messages.notificationMessage, Equals, "[localized] $nickname{Juana} couldn't be expelled.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Expelling failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] Juana couldn't be expelled")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred expelling Juana.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleRemoveFailureMessage(c *C) {
	messages := getRoleRemoveFailureMessage("foo", newTestAffiliationFromString(data.AffiliationOwner), nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] $nickname{foo} couldn't be expelled.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Expelling failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] foo couldn't be expelled")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred expelling foo.")

	messages = getRoleRemoveFailureMessage("nil", nil, session.ErrNotAllowedKickOccupant)
	c.Assert(messages.notificationMessage, Equals, "[localized] $nickname{nil} couldn't be expelled.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Expelling failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] nil couldn't be expelled")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] You don't have permissions to expel nil.")

	messages = getRoleRemoveFailureMessage("bla", newTestAffiliationFromString(data.AffiliationAdmin), session.ErrNotAllowedKickOccupant)
	c.Assert(messages.notificationMessage, Equals, "[localized] $nickname{bla} couldn't be expelled.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Expelling failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] bla couldn't be expelled")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] As an administrator you don't have permissions to expel bla.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationOutcast(c *C) {
	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New: newTestAffiliationFromString(data.AffiliationOutcast),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You have been banned from the room.")

	sau.Reason = "it's so cold"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You have been banned from the room. The reason given was: it's so cold.")

	sau.Reason = ""
	sau.Actor = newTestActor("calvin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] The owner calvin banned you from the room.")

	sau.Reason = "it isn't cool"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The owner calvin banned you from the room. The reason given was: it isn't cool.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationRoleUpdateMessage_affiliationRemoved(c *C) {
	saru := data.AffiliationRoleUpdate{
		Nickname:            "007",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationNone),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationAdmin),
		NewRole:             newTestRoleFromString(data.RoleVisitor),
		PreviousRole:        newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] You are not an administrator anymore. [localized] As a result, your role was changed from moderator to visitor.")

	saru.Reason = "he is an assassin"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] You are not an administrator anymore. [localized] As a result, your role was changed from moderator to visitor. The reason given was: he is an assassin.")

	saru.Reason = ""
	saru.Actor = newTestActor("the enemy", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] The owner the enemy changed your position; you are not an administrator anymore. [localized] As a result, your role was changed from moderator to visitor.")

	saru.Reason = "bla"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] The owner the enemy changed your position; you are not an administrator anymore. [localized] As a result, your role was changed from moderator to visitor. The reason given was: bla.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationRoleUpdateMessage_affiliationAdded(c *C) {
	saru := data.AffiliationRoleUpdate{
		Nickname:            "alice",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationAdmin),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationNone),
		NewRole:             newTestRoleFromString(data.RoleModerator),
		PreviousRole:        newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] Your position was changed to administrator. [localized] As a result, your role was changed from visitor to moderator.")

	saru.Reason = "she is lost in the world of wonders"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] Your position was changed to administrator. [localized] As a result, your role was changed from visitor to moderator. The reason given was: she is lost in the world of wonders.")

	saru.Reason = ""
	saru.Actor = newTestActor("rabbit", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] The owner rabbit changed your position to administrator. [localized] As a result, your role was changed from visitor to moderator.")

	saru.Reason = "she is lost in the world of wonders"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] The owner rabbit changed your position to administrator. [localized] As a result, your role was changed from visitor to moderator. The reason given was: she is lost in the world of wonders.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationRoleUpdateMessage_affiliationUpdated(c *C) {
	saru := data.AffiliationRoleUpdate{
		Nickname:            "goku",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationAdmin),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationMember),
		NewRole:             newTestRoleFromString(data.RoleModerator),
		PreviousRole:        newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] Your position was changed from member to administrator. [localized] As a result, your role was changed from visitor to moderator.")

	saru.Reason = "you are a powerfull Saiyajin"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] Your position was changed from member to administrator. [localized] As a result, your role was changed from visitor to moderator. The reason given was: you are a powerfull Saiyajin.")

	saru.Reason = ""
	saru.Actor = newTestActor("vegeta", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] The owner vegeta changed your position from member to administrator. [localized] As a result, your role was changed from visitor to moderator.")

	saru.Reason = "he is the prince of the Saiyajins"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] The owner vegeta changed your position from member to administrator. [localized] As a result, your role was changed from visitor to moderator. The reason given was: he is the prince of the Saiyajins.")
}

func newTestActor(nickname string, affiliation data.Affiliation, role data.Role) *data.Actor {
	return &data.Actor{
		Nickname:    nickname,
		Affiliation: affiliation,
		Role:        role,
	}
}

func newTestAffiliationFromString(s string) data.Affiliation {
	a, err := data.AffiliationFromString(s)
	if err != nil {
		panic("Error produced trying to get an affiliation from a string")
	}
	return a
}

func newTestRoleFromString(s string) data.Role {
	r, err := data.RoleFromString(s)
	if err != nil {
		return nil
	}
	return r
}
