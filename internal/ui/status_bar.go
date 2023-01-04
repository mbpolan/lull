package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/events"
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

	events.Dispatcher().Subscribe(s, []events.Code{events.EventStatusBarContextChange})

	return s
}

func (s *StatusBar) HandleEvent(code events.Code, payload events.Payload) {
	switch code {
	case events.EventStatusBarContextChange:
		if data := payload.Data.(*events.StatusBarContextChangeData); data != nil {
			s.setLayoutFromData(data)
		}
	}
}

// Widget returns a primitive widget containing this component.
func (s *StatusBar) Widget() tview.Primitive {
	return s.flex
}

func (s *StatusBar) build() {
	s.flex = tview.NewFlex()
	s.flex.SetDirection(tview.FlexColumn)
	s.flex.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	s.suffixCommonLabels()
}

func (s *StatusBar) setLayoutFromData(layout *events.StatusBarContextChangeData) {
	s.flex.Clear()
	s.prefixCommonLabels()

	for _, i := range layout.Fields {
		s.flex.AddItem(s.label(fmt.Sprintf("%s [%s]", i.Label, i.KeySequence)), 0, 1, false)
	}

	s.suffixCommonLabels()
}

func (s *StatusBar) prefixCommonLabels() {
	s.flex.AddItem(s.label("Navigate [arrow]"), 0, 1, false)
}

func (s *StatusBar) suffixCommonLabels() {
	s.flex.AddItem(s.label("Save [ctrl+s]"), 0, 1, false)
	s.flex.AddItem(s.label("Send [ctrl+g]"), 0, 1, false)
	s.flex.AddItem(s.label("Quit [ctrl+q]"), 0, 1, false)
}

func (s *StatusBar) label(text string) *tview.TextView {
	t := tview.NewTextView()
	t.SetText(text)
	t.SetTextColor(tcell.ColorWhite)
	t.SetBackgroundColor(tcell.ColorBlue)
	return t
}
