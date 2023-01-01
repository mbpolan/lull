package ui

import (
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

// TextInputModal is a modal window that prompts a user to input text.
type TextInputModal struct {
	grid         *tview.Grid
	infoText     *tview.TextView
	name         *tview.InputField
	ok           *tview.Button
	cancel       *tview.Button
	focusManager *util.FocusManager
	onAccept     ModalTextAcceptHandler
	onReject     ModalRejectHandler
	*Modal
}

// NewTextInputModal returns a new modal with a title, information text, label and button handlers.
func NewTextInputModal(title string, text string, label string, accept ModalTextAcceptHandler, reject ModalRejectHandler) *TextInputModal {
	m := new(TextInputModal)
	m.onAccept = accept
	m.onReject = reject
	m.build(title, text, label)

	return m
}

// SetText sets the informational text to show in the modal.
func (m *TextInputModal) SetText(text string) {
	m.infoText.SetText(text)
}

// Widget returns a primitive widget containing this component.
func (m *TextInputModal) Widget() tview.Primitive {
	return m.Modal.flex
}

func (m *TextInputModal) build(title string, text string, label string) {
	m.grid = tview.NewGrid()
	m.grid.SetBorder(true)
	m.grid.SetTitle(title)

	m.infoText = tview.NewTextView()
	m.infoText.SetDynamicColors(true)
	m.infoText.SetText(text)

	m.name = tview.NewInputField()
	m.name.SetLabel(label)
	m.ok = tview.NewButton("OK")
	m.ok.SetSelectedFunc(func() {
		m.onAccept(m.name.GetText())
	})

	m.cancel = tview.NewButton("Cancel")
	m.cancel.SetSelectedFunc(m.onReject)

	m.grid.AddItem(m.infoText, 0, 0, 1, 2, 0, 0, false)
	m.grid.AddItem(m.name, 1, 0, 1, 2, 0, 0, false)
	m.grid.AddItem(m.ok, 2, 0, 1, 1, 0, 0, true)
	m.grid.AddItem(m.cancel, 2, 1, 1, 1, 0, 0, false)

	m.focusManager = util.NewFocusManager(GetApplication(), []tview.Primitive{m.name, m.ok, m.cancel})

	m.grid.SetInputCapture(m.focusManager.HandleKeyEvent)
	m.Modal = NewModal(m.grid, 50, 5)
}
