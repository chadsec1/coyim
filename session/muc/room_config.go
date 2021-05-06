package muc

import (
	"strconv"

	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

const (
	// ConfigFieldFormType represents the configuration form type field
	ConfigFieldFormType = "FORM_TYPE"
	// ConfigFieldRoomName represents the room name form field
	ConfigFieldRoomName = "muc#roomconfig_roomname"
	// ConfigFieldRoomDescription represents the room description form field
	ConfigFieldRoomDescription = "muc#roomconfig_roomdesc"
	// ConfigFieldEnableLogging represents the enable logging form field
	ConfigFieldEnableLogging = "muc#roomconfig_enablelogging"
	// ConfigFieldEnableArchiving represents the enable archiving form field
	ConfigFieldEnableArchiving = "muc#roomconfig_enablearchiving"
	// ConfigFieldMemberList represents the get member list form field
	ConfigFieldMemberList = "muc#roomconfig_getmemberlist"
	// ConfigFieldLanguage represents the room language form field
	ConfigFieldLanguage = "muc#roomconfig_lang"
	// ConfigFieldPubsub represents the pubsub form field
	ConfigFieldPubsub = "muc#roomconfig_pubsub"
	// ConfigFieldCanChangeSubject represents the change subject form field
	ConfigFieldCanChangeSubject = "muc#roomconfig_changesubject"
	// ConfigFieldAllowInvites represents the allow invites form field
	ConfigFieldAllowInvites = "muc#roomconfig_allowinvites"
	// ConfigFieldAllowMemberInvites represents the allow member invites form field (for some services)
	ConfigFieldAllowMemberInvites = "{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites"
	// ConfigFieldAllowPM represents the allow private messages form field
	ConfigFieldAllowPM = "muc#roomconfig_allowpm"
	// ConfigFieldAllowPrivateMessages represent the allow private messages form fields (for some services)
	ConfigFieldAllowPrivateMessages = "allow_private_messages"
	// ConfigFieldMaxOccupantsNumber represents the max users form field
	ConfigFieldMaxOccupantsNumber = "muc#roomconfig_maxusers"
	// ConfigFieldIsPublic represents the public room form field
	ConfigFieldIsPublic = "muc#roomconfig_publicroom"
	// ConfigFieldIsPersistent represents the persistent room form field
	ConfigFieldIsPersistent = "muc#roomconfig_persistentroom"
	// ConfigFieldPresenceBroadcast represents the presence broadcast form field
	ConfigFieldPresenceBroadcast = "muc#roomconfig_presencebroadcast"
	// ConfigFieldModerated represents the moderated room form field
	ConfigFieldModerated = "muc#roomconfig_moderatedroom"
	// ConfigFieldMembersOnly represents the members only form field (for some services)
	ConfigFieldMembersOnly = "muc#roomconfig_membersonly"
	// ConfigFieldPasswordProtected represents the password protected room form field
	ConfigFieldPasswordProtected = "muc#roomconfig_passwordprotectedroom"
	// ConfigFieldPassword represents the room secret form field (for some services)
	ConfigFieldPassword = "muc#roomconfig_roomsecret"
	// ConfigFieldOwners represents the room owners list form field
	ConfigFieldOwners = "muc#roomconfig_roomowners"
	// ConfigFieldWhoIs represents the whois form field
	ConfigFieldWhoIs = "muc#roomconfig_whois"
	// ConfigFieldMaxHistoryFetch represents the max history fetch form field
	ConfigFieldMaxHistoryFetch = "muc#maxhistoryfetch"
	// ConfigFieldMaxHistoryLength represents the max history length form field (for some services)
	ConfigFieldMaxHistoryLength = "muc#roomconfig_historylength"
	// ConfigFieldRoomAdmins represents the room admins list form field
	ConfigFieldRoomAdmins = "muc#roomconfig_roomadmins"
)

// RoomConfigForm represents a room configuration form
type RoomConfigForm struct {
	MaxHistoryFetch                *RoomConfigFieldListValue
	AllowPrivateMessages           *RoomConfigFieldListValue
	OccupantsCanInvite             bool
	OccupantsCanChangeSubject      bool
	Logged                         bool
	Language                       string
	AssociatedPublishSubscribeNode string
	MaxOccupantsNumber             *RoomConfigFieldListValue
	MembersOnly                    bool
	Moderated                      bool
	PasswordProtected              bool
	Persistent                     bool
	Public                         bool
	Admins                         []jid.Any
	Description                    string
	Title                          string
	Owners                         []jid.Any
	Password                       string
	Whois                          *RoomConfigFieldListValue

	RetrieveMembersList *RoomConfigFieldListMultiValue
	PresenceBroadcast   *RoomConfigFieldListMultiValue

	formType   string
	fieldNames map[string]int
	Fields     []*RoomConfigFormField
}

// NewRoomConfigForm creates a new room configuration form instance
func NewRoomConfigForm(form *xmppData.Form) *RoomConfigForm {
	cf := &RoomConfigForm{
		fieldNames: map[string]int{},
	}

	cf.MaxHistoryFetch = newRoomConfigFieldListValue([]string{""}, []string{
		RoomConfigOption10,
		RoomConfigOption20,
		RoomConfigOption30,
		RoomConfigOption50,
		RoomConfigOption100,
		RoomConfigOptionNone,
	})

	cf.AllowPrivateMessages = newRoomConfigFieldListValue([]string{""}, []string{
		RoomConfigOptionParticipant,
		RoomConfigOptionModerators,
		RoomConfigOptionNone,
	})

	cf.MaxOccupantsNumber = newRoomConfigFieldListValue([]string{""}, []string{
		RoomConfigOption10,
		RoomConfigOption20,
		RoomConfigOption30,
		RoomConfigOption50,
		RoomConfigOption100,
		RoomConfigOptionNone,
	})

	cf.Whois = newRoomConfigFieldListValue([]string{""}, []string{
		RoomConfigOptionModerators,
		RoomConfigOptionAnyone,
	})

	cf.initListMultiValueFields()

	cf.setFormFields(form.Fields)

	return cf
}

func (rcf *RoomConfigForm) initListMultiValueFields() {
	rcf.RetrieveMembersList = newRoomConfigFieldListMultiValue(nil, retrieveMembersListDefaultOptions)
	rcf.PresenceBroadcast = newRoomConfigFieldListMultiValue(nil, presenceBroadcastDefaultOptions)
}

func (rcf *RoomConfigForm) setFormFields(fields []xmppData.FormFieldX) {
	for idx, field := range fields {
		if field.Var != "" {
			rcf.fieldNames[field.Var] = idx
			rcf.setField(field)
		}
	}
}

// GetFormData returns a representation of the room config FORM_TYPE as described in the
// XMPP specification for MUC
//
// For more information see:
// https://xmpp.org/extensions/xep-0045.html#createroom-reserved
// https://xmpp.org/extensions/xep-0045.html#example-163
func (rcf *RoomConfigForm) GetFormData() *xmppData.Form {
	formFields := []xmppData.FormFieldX{}

	for fieldName := range rcf.fieldNames {
		var values []string

		switch fieldName {
		case ConfigFieldFormType:
			values = []string{rcf.formType}

		case ConfigFieldRoomName:
			values = []string{rcf.Title}

		case ConfigFieldRoomDescription:
			values = []string{rcf.Description}

		case ConfigFieldEnableLogging, ConfigFieldEnableArchiving:
			values = []string{strconv.FormatBool(rcf.Logged)}

		case ConfigFieldMemberList:
			values = rcf.RetrieveMembersList.Value()

		case ConfigFieldLanguage:
			values = []string{rcf.Language}

		case ConfigFieldPubsub:
			values = []string{rcf.AssociatedPublishSubscribeNode}

		case ConfigFieldCanChangeSubject:
			values = []string{strconv.FormatBool(rcf.OccupantsCanChangeSubject)}

		case ConfigFieldAllowInvites, ConfigFieldAllowMemberInvites:
			values = []string{strconv.FormatBool(rcf.OccupantsCanInvite)}

		case ConfigFieldAllowPM, ConfigFieldAllowPrivateMessages:
			values = rcf.AllowPrivateMessages.Value()

		case ConfigFieldMaxOccupantsNumber:
			values = rcf.MaxOccupantsNumber.Value()

		case ConfigFieldIsPublic:
			values = []string{strconv.FormatBool(rcf.Public)}

		case ConfigFieldIsPersistent:
			values = []string{strconv.FormatBool(rcf.Persistent)}

		case ConfigFieldPresenceBroadcast:
			values = rcf.PresenceBroadcast.Value()

		case ConfigFieldModerated:
			values = []string{strconv.FormatBool(rcf.Moderated)}

		case ConfigFieldMembersOnly:
			values = []string{strconv.FormatBool(rcf.MembersOnly)}

		case ConfigFieldPasswordProtected:
			values = []string{strconv.FormatBool(rcf.PasswordProtected)}

		case ConfigFieldPassword:
			values = []string{rcf.Password}

		case ConfigFieldOwners:
			values = jidListToStringList(rcf.Owners)

		case ConfigFieldWhoIs:
			values = rcf.Whois.Value()

		case ConfigFieldMaxHistoryFetch, ConfigFieldMaxHistoryLength:
			values = rcf.MaxHistoryFetch.Value()

		case ConfigFieldRoomAdmins:
			values = jidListToStringList(rcf.Admins)
		}

		formFields = append(formFields, xmppData.FormFieldX{
			Var:    fieldName,
			Values: values,
		})
	}

	for _, f := range rcf.Fields {
		formFields = append(formFields, xmppData.FormFieldX{
			Var:    f.Name,
			Values: f.Value(),
		})
	}

	return &xmppData.Form{
		Type:   "submit",
		Fields: formFields,
	}
}

func (rcf *RoomConfigForm) setField(field xmppData.FormFieldX) {
	switch field.Var {
	case ConfigFieldFormType:
		rcf.formType = formFieldSingleString(field.Values)

	case ConfigFieldMaxHistoryFetch, ConfigFieldMaxHistoryLength:
		rcf.MaxHistoryFetch.SetValue(formFieldSingleString(field.Values))

	case ConfigFieldAllowPM, ConfigFieldAllowPrivateMessages:
		rcf.AllowPrivateMessages.SetValue(formFieldSingleString(field.Values))

	case ConfigFieldAllowInvites, ConfigFieldAllowMemberInvites:
		rcf.OccupantsCanInvite = formFieldBool(field.Values)

	case ConfigFieldCanChangeSubject:
		rcf.OccupantsCanChangeSubject = formFieldBool(field.Values)

	case ConfigFieldEnableLogging, ConfigFieldEnableArchiving:
		rcf.Logged = formFieldBool(field.Values)

	case ConfigFieldMemberList:
		rcf.RetrieveMembersList.SetValue(field.Values)
		rcf.RetrieveMembersList.SetOptions(formFieldOptionsValues(field.Options))

	case ConfigFieldLanguage:
		rcf.Language = formFieldSingleString(field.Values)

	case ConfigFieldPubsub:
		rcf.AssociatedPublishSubscribeNode = formFieldSingleString(field.Values)

	case ConfigFieldMaxOccupantsNumber:
		rcf.MaxOccupantsNumber.SetValue(formFieldSingleString(field.Values))

	case ConfigFieldMembersOnly:
		rcf.MembersOnly = formFieldBool(field.Values)

	case ConfigFieldModerated:
		rcf.Moderated = formFieldBool(field.Values)

	case ConfigFieldPasswordProtected:
		rcf.PasswordProtected = formFieldBool(field.Values)

	case ConfigFieldIsPersistent:
		rcf.Persistent = formFieldBool(field.Values)

	case ConfigFieldPresenceBroadcast:
		rcf.PresenceBroadcast.SetValue(field.Values)
		rcf.PresenceBroadcast.SetOptions(formFieldOptionsValues(field.Options))

	case ConfigFieldIsPublic:
		rcf.Public = formFieldBool(field.Values)

	case ConfigFieldRoomAdmins:
		rcf.Admins = formFieldJidList(field.Values)

	case ConfigFieldRoomDescription:
		rcf.Description = formFieldSingleString(field.Values)

	case ConfigFieldRoomName:
		rcf.Title = formFieldSingleString(field.Values)

	case ConfigFieldOwners:
		rcf.Owners = formFieldJidList(field.Values)

	case ConfigFieldPassword:
		rcf.Password = formFieldSingleString(field.Values)

	case ConfigFieldWhoIs:
		rcf.Whois.SetValue(formFieldSingleString(field.Values))

	default:
		if field.Type != RoomConfigFieldHidden && field.Type != RoomConfigFieldFixed {
			rcf.Fields = append(rcf.Fields, newRoomConfigFormField(field))
		}
	}
}

// UpdateFieldValueByName finds a field from the unknown fields list by their name and updates their value
func (rcf *RoomConfigForm) UpdateFieldValueByName(name string, value interface{}) {
	for _, field := range rcf.Fields {
		if field.Name == name {
			field.SetValue(value)
			return
		}
	}
}

func formFieldBool(values []string) bool {
	if len(values) > 0 {
		v, err := strconv.ParseBool(values[0])
		if err == nil {
			return v
		}
	}
	return false
}

func formFieldSingleString(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func formFieldOptionsValues(options []xmppData.FormFieldOptionX) (list []string) {
	for _, o := range options {
		list = append(list, o.Value)
	}

	return list
}

func formFieldJidList(values []string) (list []jid.Any) {
	for _, v := range values {
		list = append(list, jid.Parse(v))
	}
	return list
}

func jidListToStringList(jidList []jid.Any) (result []string) {
	for _, j := range jidList {
		result = append(result, j.String())
	}
	return result
}
