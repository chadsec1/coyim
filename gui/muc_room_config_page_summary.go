package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryPage struct {
	*roomConfigPageBase
	autoJoin bool

	basicInformation          gtki.LinkButton     `gtk-widget:"room-config-summary-basic-information"`
	access                    gtki.LinkButton     `gtk-widget:"room-config-summary-access"`
	permissions               gtki.LinkButton     `gtk-widget:"room-config-summary-permissions"`
	occupants                 gtki.LinkButton     `gtk-widget:"room-config-summary-occupants"`
	others                    gtki.LinkButton     `gtk-widget:"room-config-summary-others"`
	title                     gtki.Label          `gtk-widget:"room-config-summary-title"`
	descriptionNotAssigned    gtki.Label          `gtk-widget:"room-config-summary-description-not-assigned"`
	descriptionScrollWindow   gtki.ScrolledWindow `gtk-widget:"room-config-summary-description-scrolled"`
	description               gtki.TextView       `gtk-widget:"room-config-summary-description"`
	language                  gtki.Label          `gtk-widget:"room-config-summary-language"`
	includePublicList         gtki.Label          `gtk-widget:"room-config-summary-public-label"`
	persistent                gtki.Label          `gtk-widget:"room-config-summary-persistent-label"`
	password                  gtki.Label          `gtk-widget:"room-config-summary-password-label"`
	allowInviteUsers          gtki.Label          `gtk-widget:"room-config-summary-invite-label"`
	onlyMembers               gtki.Label          `gtk-widget:"room-config-summary-onlymembers-label"`
	allowSetRoomSubject       gtki.Label          `gtk-widget:"room-config-summary-changesubject-label"`
	moderatedRoom             gtki.Label          `gtk-widget:"room-config-summary-moderated-label"`
	whoIs                     gtki.Label          `gtk-widget:"room-config-summary-whois"`
	ownersListLabel           gtki.Label          `gtk-widget:"room-config-summary-owners-list-label"`
	ownersListShowButton      gtki.Button         `gtk-widget:"room-config-summary-owners-list-button"`
	ownersListShowButtonImage gtki.Image          `gtk-widget:"room-config-summary-owners-list-button-image"`
	ownersListBox             gtki.Box            `gtk-widget:"room-config-summary-owners-list"`
	ownersTreeView            gtki.TreeView       `gtk-widget:"room-config-summary-owners-tree"`
	adminsListLabel           gtki.Label          `gtk-widget:"room-config-summary-admins-list-label"`
	adminsListShowButton      gtki.Button         `gtk-widget:"room-config-summary-admins-list-button"`
	adminsListShowButtonImage gtki.Image          `gtk-widget:"room-config-summary-admins-list-button-image"`
	adminsListBox             gtki.Box            `gtk-widget:"room-config-summary-admins-list"`
	adminsTreeView            gtki.TreeView       `gtk-widget:"room-config-summary-admins-tree"`
	maxHistoryFetch           gtki.Label          `gtk-widget:"room-config-summary-maxhistoryfetch"`
	maxOccupants              gtki.Label          `gtk-widget:"room-config-summary-maxoccupants"`
	enableArchiving           gtki.Label          `gtk-widget:"room-config-summary-archive-label"`
	autojoinCheckButton       gtki.CheckButton    `gtk-widget:"room-config-autojoin"`

	ownersListModel gtki.ListStore
	adminsListModel gtki.ListStore
}

func (c *mucRoomConfigComponent) newRoomConfigSummaryPage() mucRoomConfigPage {
	p := &roomConfigSummaryPage{autoJoin: c.autoJoin}

	p.roomConfigPageBase = c.newConfigPage(pageConfigSummary, "MUCRoomConfigPageSummary", p, map[string]interface{}{
		"on_autojoin_toggled": func() {
			c.updateAutoJoin(p.autojoinCheckButton.GetActive())
		},
		"go_basic_information_page": func() {
			c.setCurrentPage(roomConfigInformationPageIndex)
		},
		"go_access_page": func() {
			c.setCurrentPage(roomConfigAccessPageIndex)
		},
		"go_permissions_page": func() {
			c.setCurrentPage(roomConfigPermissionsPageIndex)
		},
		"go_occupants_page": func() {
			c.setCurrentPage(roomConfigOccupantsPageIndex)
		},
		"go_other_page": func() {
			c.setCurrentPage(roomConfigOthersPageIndex)
		},
		"on_show_owners_list": p.onShowOwersList,
		"on_show_admins_list": p.onShowAdminList,
	})

	p.doAfterRefresh.add(p.onSummaryPageRefresh)

	mucStyles.setRoomConfigSummarySectionLinkButtonStyle(p.basicInformation)
	mucStyles.setRoomConfigSummarySectionLinkButtonStyle(p.access)
	mucStyles.setRoomConfigSummarySectionLinkButtonStyle(p.permissions)
	mucStyles.setRoomConfigSummarySectionLinkButtonStyle(p.occupants)
	mucStyles.setRoomConfigSummarySectionLinkButtonStyle(p.others)

	// The following will create two models with a column for the "jid"
	p.ownersListModel, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING)
	p.adminsListModel, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING)

	p.ownersTreeView.SetModel(p.ownersListModel)
	p.adminsTreeView.SetModel(p.adminsListModel)

	return p
}

func (p *roomConfigSummaryPage) onSummaryPageRefresh() {
	p.autojoinCheckButton.SetActive(p.autoJoin)

	// Basic information
	setLabelText(p.title, summaryAssignedValueText(p.form.Title))
	p.setDescriptionField()
	setLabelText(p.language, supportedLanguageDescription(p.form.Language))
	setLabelText(p.includePublicList, summaryYesOrNoText(p.form.Public))
	setLabelText(p.persistent, summaryYesOrNoText(p.form.Persistent))

	// Access
	setLabelText(p.password, summaryPasswordText(p.form.PasswordProtected))
	setLabelText(p.allowInviteUsers, summaryYesOrNoText(p.form.OccupantsCanInvite))
	setLabelText(p.onlyMembers, summaryYesOrNoText(p.form.MembersOnly))

	// Permissions
	setLabelText(p.allowSetRoomSubject, summaryYesOrNoText(p.form.OccupantsCanChangeSubject))
	setLabelText(p.moderatedRoom, summaryYesOrNoText(p.form.Moderated))
	setLabelText(p.whoIs, summaryForSelectedOption(p.form.Whois))

	// Occupants
	p.setOwnersAndAdminsList()

	// Other settings
	setLabelText(p.maxHistoryFetch, summaryForSelectedOption(p.form.MaxHistoryFetch))
	setLabelText(p.maxOccupants, summaryForSelectedOption(p.form.MaxOccupantsNumber))

	setLabelText(p.enableArchiving, summaryYesOrNoText(p.form.Logged))
}

func (p *roomConfigSummaryPage) setDescriptionField() {
	if p.form.Description != "" {
		setTextViewText(p.description, summaryAssignedValueText(p.form.Description))
		p.descriptionScrollWindow.Show()
		p.descriptionNotAssigned.SetVisible(false)
	} else {
		setLabelText(p.descriptionNotAssigned, summaryAssignedValueText(p.form.Description))
		p.descriptionScrollWindow.Hide()
		p.descriptionNotAssigned.SetVisible(true)
	}
}

func (p *roomConfigSummaryPage) setOwnersAndAdminsList() {
	totalOwners := len(p.form.Owners.List())
	totalAdmins := len(p.form.Admins.List())

	p.ownersListBox.SetVisible(false)
	p.adminsListBox.SetVisible(false)

	setLabelText(p.ownersListLabel, summaryOccupantsTotalText(totalOwners))
	setLabelText(p.adminsListLabel, summaryOccupantsTotalText(totalAdmins))

	p.ownersListShowButton.SetVisible(totalOwners > 0)
	p.adminsListShowButton.SetVisible(totalAdmins > 0)

	summaryOccupantsModelList(p.ownersListModel, p.form.Owners)
	summaryOccupantsModelList(p.adminsListModel, p.form.Admins)
}

func (p *roomConfigSummaryPage) onShowOwersList() {
	summaryOccupantsListHideOrShow(p.ownersTreeView, p.ownersListShowButtonImage, p.ownersListBox)
}

func (p *roomConfigSummaryPage) onShowAdminList() {
	summaryOccupantsListHideOrShow(p.adminsTreeView, p.adminsListShowButtonImage, p.adminsListBox)
}

func summaryOccupantsListHideOrShow(list gtki.TreeView, toggleButtonImage gtki.Image, container gtki.Box) {
	if list.IsVisible() {
		toggleButtonImage.SetFromIconName("pan-down-symbolic", gtki.ICON_SIZE_MENU)
		container.SetVisible(false)
	} else {
		toggleButtonImage.SetFromIconName("pan-up-symbolic", gtki.ICON_SIZE_MENU)
		container.SetVisible(true)
	}
}

func summaryOccupantsModelList(model gtki.ListStore, field *muc.RoomConfigFieldJidMultiValue) {
	model.Clear()

	items := field.List()

	for _, j := range items {
		iter := model.Append()
		model.SetValue(iter, 0, j.String())
	}
}

func summaryConfigurationOptionText(v, l string) string {
	if v == "" {
		v = muc.RoomConfigOptionNone
	}
	return configOptionToFriendlyMessage(v, l)
}

func summaryPasswordText(v bool) string {
	if v {
		return i18n.Local("**********")
	}
	return i18n.Local("Not assigned")
}

func summaryYesOrNoText(v bool) string {
	if v {
		return i18n.Local("Yes")
	}
	return i18n.Local("No")
}

func summaryAssignedValueText(label string) string {
	if label != "" {
		return label
	}
	return i18n.Local("Not assigned")
}

func summaryOccupantsTotalText(total int) string {
	switch {
	case total == 1:
		return i18n.Local("One account")
	case total > 0:
		return i18n.Localf("%d accounts", total)
	}
	return i18n.Local("None")
}

func summaryForSelectedOption(field *muc.RoomConfigFieldListValue) string {
	if s, ok := field.SelectedOption(); ok {
		return configOptionToFriendlyMessage(s.Value, s.Label)
	}
	return i18n.Local("Not assigned")
}
