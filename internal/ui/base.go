package ui

import "github.com/rivo/tview"

// Component represents a high-level UI element that contains a functional piece of user experience.
type Component interface {
	// SetFocus assigns focus to this component.
	SetFocus()

	// Widget returns the underlying tview.Primitive that supports this component.
	Widget() tview.Primitive
}
