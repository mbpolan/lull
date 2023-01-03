package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FocusManager is a utility that manages and handles changing focus amongst a set of primitives. An optional parent
// may be passed that will receive focus when the escape key is pressed.
type FocusManager struct {
	application *tview.Application
	parent      tview.Primitive
	primitives  []tview.Primitive
}

// NewFocusManager creates a new instance of FocusManager to manage the given set of primitives.
func NewFocusManager(application *tview.Application, parent tview.Primitive, primitives []tview.Primitive) *FocusManager {
	f := new(FocusManager)
	f.application = application
	f.parent = parent
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
	} else if event.Key() == tcell.KeyEscape && f.parent != nil {
		f.application.SetFocus(f.parent)
	}

	return event
}
