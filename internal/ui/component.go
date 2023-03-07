package ui

import "github.com/rivo/tview"

// Component is a higher level view that is composed of primitives.
type Component interface {
	// SetFocus sets the focus on this component.
	SetFocus()

	// Widget returns the tview.Primitive at the root of the component.
	Widget() tview.Primitive
}
