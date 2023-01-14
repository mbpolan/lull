package ui

import (
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

type BaseInputModalButton byte

const (
	BaseInputModalButtonAccept BaseInputModalButton = 1 << iota
	BaseInputModalButtonReject
	BaseInputModalButtonAll = BaseInputModalButtonAccept | BaseInputModalButtonReject
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
	width        int
	height       int
	*Modal
}

// NewBaseInputModal returns a new BaseInputModal instance with default dimensions.
func NewBaseInputModal() *BaseInputModal {
	m := new(BaseInputModal)
	m.width = 50
	m.height = 5

	return m
}

// ContentRect returns the width and height of the rectangle containing the modal's content.
func (m *BaseInputModal) ContentRect() (int, int) {
	return m.width - 2, m.height
}

// ButtonHeight returns the height (in cells) of the buttons in the modal.
func (m *BaseInputModal) ButtonHeight() int {
	return 1
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
	m.infoText.SetMaxLines(10)

	m.ok = tview.NewButton("OK")
	m.ok.SetSelectedFunc(accept)

	m.cancel = tview.NewButton("Cancel")
	m.cancel.SetSelectedFunc(m.onReject)

	// default dimensions if not specified
	width := m.width
	if width <= 0 {
		width = 50
	}

	height := m.height
	if height <= 0 {
		height = 5
	}

	m.grid.AddItem(m.infoText, 0, 0, 1, 2, 0, 0, false)
	m.Modal = NewModal(m.grid, width, height)

	return 1
}

func (m *BaseInputModal) buildButtons(row int, buttons BaseInputModalButton) {
	col := 0

	// determine how many columns each button should span, depending on how many total buttons there are to
	// configure on the modal
	var colSpan int
	if buttons == BaseInputModalButtonAll {
		colSpan = 1
	} else {
		colSpan = 2
	}

	if buttons&BaseInputModalButtonAccept != 0 {
		m.grid.AddItem(m.ok, row, col, 1, colSpan, 0, 0, true)
	}

	if buttons&BaseInputModalButtonReject != 0 {
		m.grid.AddItem(m.cancel, row, 1, 1, colSpan, 0, 0, false)
	}
}

func (m *BaseInputModal) setupFocus(primitives []tview.Primitive) {
	m.focusManager = util.NewFocusManager(m, GetApplication(), events.Dispatcher(), m.grid, primitives...)
	m.grid.SetInputCapture(m.focusManager.HandleKeyEvent)
}
