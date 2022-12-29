package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FocusManager is a utility that manages and handles changing focus amongst a set of primitives.
type FocusManager struct {
	application *tview.Application
	primitives  []tview.Primitive
}

// NewFocusManager creates a new instance of FocusManager to manage the given set of primitives.
func NewFocusManager(application *tview.Application, primitives []tview.Primitive) *FocusManager {
	f := new(FocusManager)
	f.application = application
	f.primitives = primitives

	return f
}

// HandleKeyEvent processes a keyboard event and changes which primitive is focused.
func (f *FocusManager) HandleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyTab {
		focused := -1
		for i, w := range f.primitives {
			if f.application.GetFocus() == w {
				focused = i
				break
			}
		}

		if focused > -1 {
			f.application.SetFocus(f.primitives[(focused+1)%len(f.primitives)])
		}

		return nil
	}

	return event
}
