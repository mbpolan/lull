package ui

import "github.com/rivo/tview"

// Payload is a view that allows viewing and editing request/response components.
type Payload struct {
	box *tview.Box
}

// NewPayload returns a new instance of Payload.
func NewPayload(title string) *Payload {
	p := new(Payload)
	p.build(title)

	return p
}

// Widget returns a primitive widget containing this component.
func (p *Payload) Widget() *tview.Box {
	return p.box
}

func (p *Payload) build(title string) {
	p.box = tview.NewBox()
	p.box.SetBorder(true)
	p.box.SetTitle(title)
}
