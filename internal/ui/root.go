package ui

import "github.com/rivo/tview"

// Root is a top-level container for all application UI components.
type Root struct {
	flex       *tview.Flex
	collection *Collection
	content    *Content
}

// NewRoot returns a new Root instance.
func NewRoot() *Root {
	r := new(Root)
	r.build()

	return r
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
}
