package ui

import (
	"github.com/mbpolan/lull/internal/util"
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
	m.BaseInputModal = NewBaseInputModal()
	m.onAccept = accept

	// wrap the text, so it fits in the modal, and adjust the modal height based on how many lines we
	// need to display
	w, _ := m.ContentRect()
	wrapped, lines := util.WrapText(text, w)
	m.height += lines

	m.build(title, wrapped, buttonText)

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
