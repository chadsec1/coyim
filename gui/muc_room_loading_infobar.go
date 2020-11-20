package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewLoadingInfoBar struct {
	parent gtki.Box

	infoBar                 gtki.InfoBar `gtk-widget:"info"`
	content                 gtki.Box     `gtk-widget:"info-content"`
	infoBarLabel            gtki.Label   `gtk-widget:"info-label"`
	infoBarDescriptionLabel gtki.Label   `gtk-widget:"info-description"`
	tryAgainButton          gtki.Button  `gtk-widget:"try-again-button"`

	spinner *spinner
	onRetry func()
}

func (v *roomView) newRoomViewLoadingInfoBar(parent gtki.Box) *roomViewLoadingInfoBar {
	ib := &roomViewLoadingInfoBar{
		parent:  parent,
		spinner: newSpinner(),
		onRetry: func() {},
	}

	builder := newBuilder("MUCRoomLoadingInfoBar")
	panicOnDevError(builder.bindObjects(ib))

	builder.ConnectSignals(map[string]interface{}{
		"on_try_again_clicked": ib.onRetryClicked,
	})

	mucStyles.setRoomLoadingInfoBarLabelStyle(ib.infoBarLabel)

	ib.content.Add(ib.spinner.getWidget())
	ib.parent.Add(ib.infoBar)

	return ib
}

func (ib *roomViewLoadingInfoBar) start() {
	ib.infoBarLabel.SetLabel(i18n.Local("Loading room information"))
	ib.infoBarDescriptionLabel.SetLabel(i18n.Local("Sometimes this can take few minutes, so please wait until it finishes."))
	ib.infoBar.SetMessageType(gtki.MESSAGE_INFO)
	ib.show()
}

func (ib *roomViewLoadingInfoBar) show() {
	ib.spinner.show()
	ib.infoBar.Show()
}

func (ib *roomViewLoadingInfoBar) error(text, description string, onRetry func()) {
	ib.stop()

	ib.infoBar.SetMessageType(gtki.MESSAGE_ERROR)
	ib.infoBarLabel.SetLabel(text)

	if description != "" {
		ib.infoBarDescriptionLabel.SetLabel(description)
	} else {
		ib.infoBarDescriptionLabel.Hide()
	}

	if onRetry != nil {
		ib.retryWith(onRetry)
	}
}

func (ib *roomViewLoadingInfoBar) retryWith(f func()) {
	ib.onRetry = f
	ib.tryAgainButton.Show()
}

func (ib *roomViewLoadingInfoBar) retry() {
	ib.onRetry = func() {}
	ib.tryAgainButton.Hide()
	ib.start()
}

func (ib *roomViewLoadingInfoBar) onRetryClicked() {
	if ib.onRetry != nil {
		go ib.onRetry()
	}

	ib.retry()
}

func (ib *roomViewLoadingInfoBar) stop() {
	ib.spinner.hide()
}

func (ib *roomViewLoadingInfoBar) hide() {
	ib.infoBar.Hide()
}

func (ib *roomViewLoadingInfoBar) getWidget() gtki.InfoBar {
	return ib.infoBar
}
