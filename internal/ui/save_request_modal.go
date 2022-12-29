package ui

import (
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

type AcceptHandler func(text string)
type RejectHandler func()

// SaveRequestModal is a modal window that presents input fields for naming a request to save.
type SaveRequestModal struct {
	grid         *tview.Grid
	name         *tview.InputField
	ok           *tview.Button
	cancel       *tview.Button
	focusManager *util.FocusManager
	onAccept     AcceptHandler
	onReject     RejectHandler
	*Modal
}

// NewSaveRequestModal return a new modal with an accept handler and a reject handler.
func NewSaveRequestModal(accept AcceptHandler, reject RejectHandler) *SaveRequestModal {
	s := new(SaveRequestModal)
	s.onAccept = accept
	s.onReject = reject
	s.build()

	return s
}

// Widget returns a primitive widget containing this component.
func (s *SaveRequestModal) Widget() tview.Primitive {
	return s.Modal.flex
}

func (s *SaveRequestModal) build() {
	s.grid = tview.NewGrid()
	s.grid.SetBorder(true)
	s.grid.SetTitle("Save Request")

	s.name = tview.NewInputField()
	s.name.SetLabel("Request name")
	s.ok = tview.NewButton("OK")
	s.ok.SetSelectedFunc(func() {
		s.onAccept(s.name.GetText())
	})

	s.cancel = tview.NewButton("Cancel")
	s.cancel.SetSelectedFunc(s.onReject)

	s.grid.AddItem(s.name, 0, 0, 1, 2, 0, 0, false)
	s.grid.AddItem(s.ok, 1, 0, 1, 1, 0, 0, true)
	s.grid.AddItem(s.cancel, 1, 1, 1, 1, 0, 0, false)

	s.focusManager = util.NewFocusManager(GetApplication(), []tview.Primitive{s.name, s.ok, s.cancel})

	s.grid.SetInputCapture(s.focusManager.HandleKeyEvent)
	s.Modal = NewModal(s.grid, 50, 4)
}