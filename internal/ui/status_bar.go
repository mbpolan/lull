package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/events"
	"github.com/rivo/tview"
	"unicode/utf8"
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

	s.suffixCommonLabels()
}

func (s *StatusBar) setLayoutFromData(layout *events.StatusBarContextChangeData) {
	s.flex.Clear()
	s.prefixCommonLabels()

	for _, i := range layout.Fields {
		s.addLabel(fmt.Sprintf("%s [%s]", i.Label, i.KeySequence))
	}

	s.suffixCommonLabels()
}

func (s *StatusBar) prefixCommonLabels() {
	s.addLabel("Navigate [↑↓←→]")
	s.addLabel("Focus [⇥]")
}

func (s *StatusBar) suffixCommonLabels() {
	s.addLabel("Save [⌃S]")
	s.addLabel("Send [⌃G]")
	s.addLabel("Quit [⌃Q]")
}

func (s *StatusBar) addLabel(text string) {
	t := tview.NewTextView()
	t.SetText(text)
	t.SetTextColor(tcell.ColorWhite)
	t.SetBackgroundColor(tcell.ColorBlue)

	// add the label and a spacer right after it
	// TODO: can we do something with margins instead of adding an empty primitive?
	s.flex.AddItem(t, utf8.RuneCountInString(text), 1, false)
	s.flex.AddItem(tview.NewBox(), 1, 1, false)
}
