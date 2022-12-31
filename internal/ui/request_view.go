package ui

import (
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
)

// RequestView is a view that allows viewing and editing request/response components.
type RequestView struct {
	flex  *tview.Flex
	pages *tview.Pages
	body  *tview.TextArea
	state *state.Manager
}

// NewRequestView returns a new instance of RequestView.
func NewRequestView(title string, state *state.Manager) *RequestView {
	p := new(RequestView)
	p.state = state
	p.build(title)

	return p
}

// Reload refreshes the state of the component with current app state.
func (p *RequestView) Reload() {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	p.body.SetText(item.RequestBody, false)
}

// Widget returns a primitive widget containing this component.
func (p *RequestView) Widget() *tview.Flex {
	return p.flex
}

func (p *RequestView) build(title string) {
	p.flex = tview.NewFlex()
	p.flex.SetBorder(true)
	p.flex.SetTitle(title)

	p.pages = tview.NewPages()
	p.flex.AddItem(p.pages, 0, 1, true)

	p.body = tview.NewTextArea()
	p.body.SetChangedFunc(p.handleBodyChange)
	p.pages.AddAndSwitchToPage("body", p.body, true)
}

func (p *RequestView) handleBodyChange() {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	item.RequestBody = p.body.GetText()
}
