package ui

import (
	"github.com/rivo/tview"
)

// Payload is a view that allows viewing and editing request/response components.
type Payload struct {
	flex         *tview.Flex
	pages        *tview.Pages
	editableBody *tview.TextArea
	readOnlyBody *tview.TextView
	readOnly     bool
}

// NewPayload returns a new instance of Payload.
func NewPayload(title string, readOnly bool) *Payload {
	p := new(Payload)
	p.readOnly = readOnly
	p.build(title)

	return p
}

// Widget returns a primitive widget containing this component.
func (p *Payload) Widget() *tview.Flex {
	return p.flex
}

func (p *Payload) build(title string) {
	p.flex = tview.NewFlex()
	p.flex.SetBorder(true)
	p.flex.SetTitle(title)

	p.pages = tview.NewPages()
	p.flex.AddItem(p.pages, 0, 1, true)

	if p.readOnly {
		p.readOnlyBody = tview.NewTextView()
		p.pages.AddAndSwitchToPage("Body", p.readOnlyBody, true)
	} else {
		p.editableBody = tview.NewTextArea()
		p.pages.AddAndSwitchToPage("Body", p.editableBody, true)
	}
}

func (p *Payload) SetData(code int, body []byte, err error) {
	str := string(body)

	if p.readOnly {
		p.readOnlyBody.SetText(str)
	} else {
		p.editableBody.SetText(str, false)
	}
}
