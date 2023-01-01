package ui

import "github.com/rivo/tview"

// ModalTextAcceptHandler is a standard callback for a modal that accepts a line of text as input.
type ModalTextAcceptHandler func(text string)

// ModalRejectHandler is a standard callback for when a modal is cancelled.
type ModalRejectHandler func()

// Modal is a container that presents components in a modal window.
type Modal struct {
	flex *tview.Flex
}

// NewModal returns a new Modal with a tview.Primitive as content. The modal will be sized according to the
// given width and height units.
func NewModal(content tview.Primitive, width, height int) *Modal {
	m := new(Modal)
	m.build(content, width, height)

	return m
}

func (m *Modal) build(content tview.Primitive, width, height int) {
	m.flex = tview.NewFlex()

	m.flex.AddItem(nil, 0, 1, false)

	m.flex.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(content, height, 1, true).
		AddItem(nil, 0, 1, false), width, 1, true)

	m.flex.AddItem(nil, 0, 1, false)
}
