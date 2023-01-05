package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/events"
	"github.com/rivo/tview"
)

type FocusFilterResult int

const (
	FocusPreHandlePropagate FocusFilterResult = iota
	FocusPreHandleIgnore
	FocusPreHandleProcess
)

type FocusDirection int

const (
	FocusUp FocusDirection = iota
	FocusDown
	FocusLeft
	FocusRight
)

type FocusManagerFilterHandler func(event *tcell.EventKey) FocusFilterResult
type FocusManagerHandler func(event *tcell.EventKey) *tcell.EventKey

// FocusManager is a utility that manages and handles changing focus amongst a set of primitives. An optional parent
// may be passed that will receive focus when the escape key is pressed.
type FocusManager struct {
	sender           any
	application      *tview.Application
	dispatcher       *events.EventDispatcher
	parent           tview.Primitive
	arrowParentFocus bool
	directions       map[tcell.Key]events.Code
	primitives       []tview.Primitive
	filter           FocusManagerFilterHandler
	handler          FocusManagerHandler
}

// NewFocusManager creates a new instance of FocusManager to manage the given set of primitives.
func NewFocusManager(sender any, application *tview.Application, dispatcher *events.EventDispatcher, parent tview.Primitive, primitives ...tview.Primitive) *FocusManager {
	f := new(FocusManager)
	f.sender = sender
	f.application = application
	f.dispatcher = dispatcher
	f.parent = parent
	f.arrowParentFocus = true
	f.directions = map[tcell.Key]events.Code{}
	f.primitives = primitives
	f.filter = func(event *tcell.EventKey) FocusFilterResult {
		return FocusPreHandleProcess
	}

	return f
}

// SetLenientArrowNavigation allows arrow navigation to occur without requiring the parent primitive to have focus.
func (f *FocusManager) SetLenientArrowNavigation() {
	f.arrowParentFocus = false
}

// SetFilter sets the function to invoke that determines if a key event should be processed.
func (f *FocusManager) SetFilter(filter FocusManagerFilterHandler) {
	f.filter = filter
}

// SetHandler sets the function to invoke when the FocusManager does not handle a key event. This can be useful for
// chaining instances of FocusManager and other key event handlers.
func (f *FocusManager) SetHandler(handler FocusManagerHandler) {
	f.handler = handler
}

// AddArrowNavigation enables dispatching a directional navigation event if the given arrow key event(s) are received
// while the parent primitive has focus.
func (f *FocusManager) AddArrowNavigation(directions ...FocusDirection) {
	for _, i := range directions {
		switch i {
		case FocusUp:
			f.directions[tcell.KeyUp] = events.EventNavigateUp
		case FocusDown:
			f.directions[tcell.KeyDown] = events.EventNavigateDown
		case FocusLeft:
			f.directions[tcell.KeyLeft] = events.EventNavigateLeft
		case FocusRight:
			f.directions[tcell.KeyRight] = events.EventNavigateRight
		}
	}
}

// HandleKeyEvent processes a keyboard event and changes which primitive is focused.
func (f *FocusManager) HandleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	// run the filter first to determine what to do with this event
	switch f.filter(event) {
	case FocusPreHandleIgnore:
		// do not process this event at all
		return nil
	case FocusPreHandlePropagate:
		// bubble this event up the chain
		return event
	default:
		// handle the event ourselves
		break
	}

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
		return nil
	} else if code, ok := f.directions[event.Key()]; ok {
		// check if we require the parent primitive to have focus before allow arrow navigation
		if !f.arrowParentFocus || f.application.GetFocus() == f.parent {
			f.dispatcher.PostSimple(code, f.sender)
			return nil
		}
	}

	if f.handler != nil {
		return f.handler(event)
	}

	return event
}
