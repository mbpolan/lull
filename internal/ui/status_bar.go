package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StatusBarLayout int

const (
	StatusBarLayoutGeneral StatusBarLayout = iota
	StatusBarLayoutCollection
)

// StatusBar presents informational components.
type StatusBar struct {
	flex   *tview.Flex
	layout StatusBarLayout
}

// NewStatusBar returns an instance of StatusBar with general layout.
func NewStatusBar() *StatusBar {
	s := new(StatusBar)
	s.layout = StatusBarLayoutGeneral
	s.build()

	return s
}

// SetLayout sets the status bar layout to use.
func (s *StatusBar) SetLayout(layout StatusBarLayout) {
	s.layout = layout
	s.reload()
}

// Widget returns a primitive widget containing this component.
func (s *StatusBar) Widget() tview.Primitive {
	return s.flex
}

func (s *StatusBar) build() {
	s.flex = tview.NewFlex()
	s.flex.SetDirection(tview.FlexColumn)
	s.flex.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	s.reload()
}

func (s *StatusBar) reload() {
	s.flex.Clear()

	switch s.layout {
	case StatusBarLayoutGeneral:
		s.setupForGeneral()
	case StatusBarLayoutCollection:
		s.setupForCollection()
	}
}

func (s *StatusBar) setupForGeneral() {
	s.flex.AddItem(s.label("Collection [ctrl+l]"), 0, 1, false)
	s.flex.AddItem(s.label("Save Current [ctrl+s]"), 0, 1, false)
	s.flex.AddItem(s.label("URL [ctrl+a]"), 0, 1, false)
	s.flex.AddItem(s.label("Request [ctrl+r]"), 0, 1, false)
	s.flex.AddItem(s.label("Send [ctrl+g]"), 0, 1, false)
	s.flex.AddItem(s.label("Quit [ctrl+q]"), 0, 1, false)
}

func (s *StatusBar) setupForCollection() {
	s.flex.AddItem(s.label("Open [enter]"), 0, 1, false)
	s.flex.AddItem(s.label("Rename [r]"), 0, 1, false)
	s.flex.AddItem(s.label("Clone [c]"), 0, 1, false)
	s.flex.AddItem(s.label("Delete [d]"), 0, 1, false)
	s.flex.AddItem(s.label(""), 0, 2, false)
}

func (s *StatusBar) label(text string) *tview.TextView {
	t := tview.NewTextView()
	t.SetText(text)
	t.SetTextColor(tcell.ColorWhite)
	t.SetBackgroundColor(tcell.ColorBlue)
	return t
}
