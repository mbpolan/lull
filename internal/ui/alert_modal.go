package ui

import (
	"github.com/rivo/tview"
)

// AlertModal is a modal that shows a message with a single button to take an action.
type AlertModal struct {
	onAccept ModalNoArgAcceptHandler
	*BaseInputModal
}

// NewAlertModal returns a new modal with text and button handler.
func NewAlertModal(title string, text string, buttonText string, accept ModalNoArgAcceptHandler) *AlertModal {
	m := new(AlertModal)
	m.BaseInputModal = new(BaseInputModal)
	m.onAccept = accept
	m.build(title, text, buttonText)

	return m
}

func (m *AlertModal) build(title string, text string, buttonText string) {
	row := m.BaseInputModal.build(title, text, func() {
		m.onAccept()
	})

	m.buildButtons(row, BaseInputModalButtonAccept)
	m.ok.SetLabel(buttonText)
	m.setupFocus([]tview.Primitive{m.ok})
}
