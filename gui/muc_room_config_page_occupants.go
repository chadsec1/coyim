package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomConfigOccupantsPage struct {
	*roomConfigPageBase

	ownersListContent     gtki.Box              `gtk-widget:"room-config-owners-list-content"`
	ownersList            gtki.TreeView         `gtk-widget:"room-config-owners-list"`
	ownersAddButton       gtki.Button           `gtk-widget:"room-config-owner-add"`
	ownersRemoveButton    gtki.Button           `gtk-widget:"room-config-owner-remove"`
	ownersRemoveLabel     gtki.Label            `gtk-widget:"room-config-owner-remove-label"`
	ownersListJidRenderer gtki.CellRendererText `gtk-widget:"room-config-owner-jid-text-renderer"`
	adminsListContent     gtki.Box              `gtk-widget:"room-config-admins-list-content"`
	adminsList            gtki.TreeView         `gtk-widget:"room-config-admin-list"`
	adminsAddButton       gtki.Button           `gtk-widget:"room-config-admin-add"`
	adminsRemoveButton    gtki.Button           `gtk-widget:"room-config-admin-remove"`
	adminsRemoveLabel     gtki.Label            `gtk-widget:"room-config-admin-remove-label"`
	adminsListJidRenderer gtki.CellRendererText `gtk-widget:"room-config-admin-jid-text-renderer"`

	ownersListController *mucRoomConfigListController
	adminsListController *mucRoomConfigListController
}

func (c *mucRoomConfigComponent) newRoomConfigOccupantsPage(parent gtki.Window) mucRoomConfigPage {
	p := &roomConfigOccupantsPage{}
	p.roomConfigPageBase = c.newConfigPage(pageConfigOccupants, "MUCRoomConfigPageOccupants", p, map[string]interface{}{
		"on_owner_jid_edited": p.onOwnerJidEdited,
		"on_admin_jid_edited": p.onAdminJidEdited,
	})

	p.initDefaults()
	p.initOccupantsEditableCells()
	p.initOccupantsLists(parent)

	return p
}

func (p *roomConfigOccupantsPage) initDefaults() {
	p.doAfterRefresh.add(p.refreshContentLists)
}

func (p *roomConfigOccupantsPage) initOccupantsEditableCells() {
	p.ownersListJidRenderer.SetProperty("editable", true)
	p.adminsListJidRenderer.SetProperty("editable", true)
}

func (p *roomConfigOccupantsPage) initOccupantsLists(parent gtki.Window) {
	p.initOwnersListController(parent)
	p.initAdminsListController(parent)
}

func (p *roomConfigOccupantsPage) initOwnersListController(parent gtki.Window) {
	p.ownersListController = p.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:      p.ownersAddButton,
		removeOccupantButton:   p.ownersRemoveButton,
		removeOccupantLabel:    p.ownersRemoveLabel,
		occupantsTreeView:      p.ownersList,
		parentWindow:           parent,
		addOccupantDialogTitle: i18n.Local("Add owners"),
		addOccupantDescription: i18n.Local("Here you can add one or more new owners to the room. You will have to use the account address of the user in order to make them an owner. This address can either be a simple one, such as user@example.org or a full one, such as user@example.org/abcdef."),
		onListUpdated:          p.refreshContentLists,
	})
}

func (p *roomConfigOccupantsPage) initAdminsListController(parent gtki.Window) {
	p.adminsListController = p.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:      p.adminsAddButton,
		removeOccupantButton:   p.adminsRemoveButton,
		removeOccupantLabel:    p.adminsRemoveLabel,
		occupantsTreeView:      p.adminsList,
		parentWindow:           parent,
		addOccupantDialogTitle: i18n.Local("Add administrators"),
		addOccupantDescription: i18n.Local("Here you can add one or more new administrators to the room. You will have to use the account address of the user in order to make them an administrator. This address can either be a simple one, such as user@example.org or a full one, such as user@example.org/abcdef."),
		onListUpdated:          p.refreshContentLists,
	})
}

const occupantListJidColumnIndex = 0

func (p *roomConfigOccupantsPage) onOwnerJidEdited(_ gtki.CellRendererText, path string, newValue string) {
	p.updateOccupantListCellForString("owners", p.ownersListController, occupantListJidColumnIndex, path, newValue)
}

func (p *roomConfigOccupantsPage) onAdminJidEdited(_ gtki.CellRendererText, path string, newValue string) {
	p.updateOccupantListCellForString("admins", p.adminsListController, occupantListJidColumnIndex, path, newValue)
}

func (p *roomConfigOccupantsPage) updateOccupantListCellForString(list string, controller *mucRoomConfigListController, column int, path string, newValue string) {
	if controller.updateCellForString(column, path, newValue) {
		p.log.WithFields(log.Fields{
			"path":    path,
			"newText": newValue,
			"list":    list,
		}).Debug("The occupant's jid can't be updated")
	}
}

func (p *roomConfigOccupantsPage) refreshContentLists() {
	p.ownersListContent.SetVisible(p.ownersListController.hasItems())
	p.adminsListContent.SetVisible(p.adminsListController.hasItems())
}

func (p *roomConfigOccupantsPage) collectData() {
	p.form.Owners.SetValue(jidListFromConfigListController(p.ownersListController))
	p.form.Admins.SetValue(jidListFromConfigListController(p.adminsListController))
}

func jidListFromConfigListController(l *mucRoomConfigListController) []jid.Any {
	result := []jid.Any{}
	for _, li := range l.listItems() {
		result = append(result, jid.Parse(li))
	}
	return result
}
