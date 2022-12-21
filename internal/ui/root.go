package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Root is a top-level container for all application UI components.
type Root struct {
	flex       *tview.Flex
	collection *Collection
	content    *Content
}

var application *tview.Application

// NewRoot returns a new Root instance.
func NewRoot(app *tview.Application) *Root {
	application = app

	r := new(Root)
	r.build()

	return r
}

func GetApplication() *tview.Application {
	return application
}

// Widget returns a primitive widget containing this component.
func (r *Root) Widget() *tview.Flex {
	return r.flex
}

func (r *Root) build() {
	// create child widgets
	r.collection = NewCollection()
	r.content = NewContent()

	// arrange them in a flex layout
	r.flex = tview.NewFlex()
	r.flex.AddItem(r.collection.Widget(), 25, 0, false)
	r.flex.AddItem(r.content.Widget(), 0, 1, true)

	r.flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers()&tcell.ModCtrl > 0 {
			if r.handleFocusShortcut(event.Key(), event.Rune()) {
				return nil
			}
		}

		return event
	})
}

func (r *Root) handleFocusShortcut(code tcell.Key, key rune) bool {
	switch code {
	case tcell.KeyCtrlA:
		r.content.SetFocus(ContentURLBox)
	case tcell.KeyCtrlR:
		r.content.SetFocus(ContentRequestBody)
	}

	return false
}
