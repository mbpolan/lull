package ui

import (
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

// RequestView is a view that allows viewing and editing request/response components.
type RequestView struct {
	flex         *tview.Flex
	pages        *tview.Pages
	body         *tview.TextArea
	focusHolder  *tview.TextView
	focusManager *util.FocusManager
	state        *state.Manager
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
	curBody := p.currentRequestBody()

	p.flex = tview.NewFlex()
	p.flex.SetBorder(true)
	p.flex.SetTitle(title)

	p.focusHolder = tview.NewTextView()

	p.pages = tview.NewPages()
	p.flex.AddItem(p.focusHolder, 1, 0, true)
	p.flex.AddItem(p.pages, 0, 1, false)

	p.body = tview.NewTextArea()
	p.body.SetText(curBody, false)
	p.body.SetChangedFunc(p.handleBodyChange)

	p.focusManager = util.NewFocusManager(p, GetApplication(), events.Dispatcher(), p.focusHolder, p.focusHolder, p.body)
	p.focusManager.AddArrowNavigation(util.FocusUp, util.FocusLeft, util.FocusRight)

	p.body.SetInputCapture(p.focusManager.HandleKeyEvent)
	p.flex.SetInputCapture(p.focusManager.HandleKeyEvent)

	p.pages.AddAndSwitchToPage("body", p.body, true)
}

func (p *RequestView) currentRequestBody() string {
	item := p.state.Get().ActiveItem
	if item == nil {
		return ""
	}

	return item.RequestBody
}

func (p *RequestView) handleBodyChange() {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	item.RequestBody = p.body.GetText()
}
