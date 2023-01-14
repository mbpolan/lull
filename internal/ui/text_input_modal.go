package ui

import (
	"github.com/rivo/tview"
)

// TextInputModal is a modal window that prompts a user to input text.
type TextInputModal struct {
	name     *tview.InputField
	onAccept ModalTextAcceptHandler
	*BaseInputModal
}

// NewTextInputModal returns a new modal with a title, information text, label and button handlers.
func NewTextInputModal(title string, text string, label string, accept ModalTextAcceptHandler, reject ModalRejectHandler) *TextInputModal {
	m := new(TextInputModal)
	m.BaseInputModal = NewBaseInputModal()
	m.onAccept = accept
	m.onReject = reject
	m.build(title, text, label)

	return m
}

func (m *TextInputModal) build(title string, text string, label string) {
	row := m.BaseInputModal.build(title, text, func() {
		m.onAccept(m.name.GetText())
	})

	m.name = tview.NewInputField()
	m.name.SetLabel(label)

	m.grid.AddItem(m.name, row, 0, 1, 2, 0, 0, true)

	m.buildButtons(row+1, BaseInputModalButtonAll)
	m.setupFocus([]tview.Primitive{m.name, m.ok, m.cancel})
}
