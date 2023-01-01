package ui

import (
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

// BaseInputModal is a scaffold that provides the basis for building more complex input modals. You should not
// use this component directly. Instead, use it as a base for embedding in more functional modals.
type BaseInputModal struct {
	grid         *tview.Grid
	infoText     *tview.TextView
	ok           *tview.Button
	cancel       *tview.Button
	focusManager *util.FocusManager
	onReject     ModalRejectHandler
	*Modal
}

// SetText sets the informational text to show in the modal.
func (m *BaseInputModal) SetText(text string) {
	m.infoText.SetText(text)
}

// Widget returns a primitive widget containing this component.
func (m *BaseInputModal) Widget() tview.Primitive {
	return m.Modal.flex
}

func (m *BaseInputModal) build(title string, text string, accept func()) int {
	m.grid = tview.NewGrid()
	m.grid.SetBorder(true)
	m.grid.SetTitle(title)

	m.infoText = tview.NewTextView()
	m.infoText.SetDynamicColors(true)
	m.infoText.SetText(text)

	m.ok = tview.NewButton("OK")
	m.ok.SetSelectedFunc(accept)

	m.cancel = tview.NewButton("Cancel")
	m.cancel.SetSelectedFunc(m.onReject)

	m.grid.AddItem(m.infoText, 0, 0, 1, 2, 0, 0, false)
	m.Modal = NewModal(m.grid, 50, 5)

	return 1
}

func (m *BaseInputModal) buildButtons(row int) {
	m.grid.AddItem(m.ok, row, 0, 1, 1, 0, 0, false)
	m.grid.AddItem(m.cancel, row, 1, 1, 1, 0, 0, false)
}

func (m *BaseInputModal) setupFocus(primitives []tview.Primitive) {
	m.focusManager = util.NewFocusManager(GetApplication(), primitives)
	m.grid.SetInputCapture(m.focusManager.HandleKeyEvent)
}
