package ui

import (
	"github.com/rivo/tview"
)

type ModalNoArgAcceptHandler func()

// PromptModal is a modal that prompts the user to confirm or reject an action.
type PromptModal struct {
	onAccept ModalNoArgAcceptHandler
	*BaseInputModal
}

// NewPromptModal returns a new modal with text and button handlers.
func NewPromptModal(title string, text string, accept ModalNoArgAcceptHandler, reject ModalRejectHandler) *PromptModal {
	m := new(PromptModal)
	m.BaseInputModal = NewBaseInputModal()
	m.onAccept = accept
	m.onReject = reject
	m.build(title, text)

	return m
}

func (m *PromptModal) build(title string, text string) {
	row := m.BaseInputModal.build(title, text, func() {
		m.onAccept()
	})

	m.buildButtons(row, BaseInputModalButtonAll)

	// give the text as much height as possible, fixed height buttons
	m.grid.SetRows(-1, m.ButtonHeight())

	m.setupFocus([]tview.Primitive{m.ok, m.cancel})
}
