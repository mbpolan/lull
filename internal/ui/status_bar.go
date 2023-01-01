package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// StatusBar presents informational components.
type StatusBar struct {
	flex *tview.Flex
}

// NewStatusBar returns an instance of StatusBar.
func NewStatusBar() *StatusBar {
	s := new(StatusBar)
	s.build()

	return s
}

// Widget returns a primitive widget containing this component.
func (s *StatusBar) Widget() tview.Primitive {
	return s.flex
}

func (s *StatusBar) build() {
	label := func(text string) *tview.TextView {
		t := tview.NewTextView()
		t.SetText(text)
		t.SetTextColor(tcell.ColorWhite)
		t.SetBackgroundColor(tcell.ColorBlue)
		return t
	}

	s.flex = tview.NewFlex()
	s.flex.SetDirection(tview.FlexColumn)
	s.flex.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	s.flex.AddItem(label("Collection [ctrl+l]"), 0, 1, false)
	s.flex.AddItem(label("Save Current [ctrl+s]"), 0, 1, false)
	s.flex.AddItem(label("URL [ctrl+a]"), 0, 1, false)
	s.flex.AddItem(label("Request [ctrl+r]"), 0, 1, false)
	s.flex.AddItem(label("Send [ctrl+g]"), 0, 1, false)
	s.flex.AddItem(label("Quit [ctrl+q]"), 0, 1, false)
}
