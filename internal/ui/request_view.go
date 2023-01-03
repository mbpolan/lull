package ui

import (
	"github.com/gdamore/tcell/v2"
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
	p.flex.AddItem(p.focusHolder, 1, 0, false)
	p.flex.AddItem(p.pages, 0, 1, true)

	p.body = tview.NewTextArea()
	p.body.SetText(curBody, false)
	p.body.SetChangedFunc(p.handleBodyChange)

	p.focusManager = util.NewFocusManager(GetApplication(), p.focusHolder, []tview.Primitive{p.focusHolder, p.body})
	p.body.SetInputCapture(p.focusManager.HandleKeyEvent)

	p.pages.AddAndSwitchToPage("body", p.body, true)

	p.flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRight && GetApplication().GetFocus() == p.focusHolder {
			events.Dispatcher().PostSimple(events.EventNavigateRight, p)
			return nil
		} else if event.Key() == tcell.KeyUp && GetApplication().GetFocus() == p.focusHolder {
			events.Dispatcher().PostSimple(events.EventNavigateUp, p)
			return nil
		} else if event.Key() == tcell.KeyLeft && GetApplication().GetFocus() == p.focusHolder {
			events.Dispatcher().PostSimple(events.EventNavigateLeft, p)
			return nil
		}

		return event
	})
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
