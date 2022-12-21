package ui

import "github.com/rivo/tview"

// Payload is a view that allows viewing and editing request/response components.
type Payload struct {
	flex  *tview.Flex
	pages *tview.Pages
	body  *tview.TextArea
}

// NewPayload returns a new instance of Payload.
func NewPayload(title string) *Payload {
	p := new(Payload)
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

	p.body = tview.NewTextArea()

	p.pages = tview.NewPages()
	p.pages.AddAndSwitchToPage("Body", p.body, true)

	p.flex.AddItem(p.pages, 0, 1, true)
}
