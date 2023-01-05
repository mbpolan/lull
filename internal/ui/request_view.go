package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
	"strings"
)

// RequestView is a view that allows viewing and editing request/response components.
type RequestView struct {
	flex         *tview.Flex
	pages        *tview.Pages
	body         *tview.TextArea
	headers      *tview.Table
	focusHolder  *tview.TextView
	focusManager *util.FocusManager
	sbSequences  []events.StatusBarContextChangeSequence
	state        *state.Manager
}

// NewRequestView returns a new instance of RequestView.
func NewRequestView(title string, state *state.Manager) *RequestView {
	p := new(RequestView)
	p.state = state
	p.build(title)
	p.Reload()

	// no additional key sequences” supported by this component
	p.sbSequences = []events.StatusBarContextChangeSequence{}

	return p
}

// SetFocus sets the focus on this component.
func (p *RequestView) SetFocus() {
	events.Dispatcher().Post(events.EventStatusBarContextChange, p, &events.StatusBarContextChangeData{
		Fields: p.sbSequences,
	})

	GetApplication().SetFocus(p.Widget())
}

// Reload refreshes the state of the component with current app state.
func (p *RequestView) Reload() {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	p.body.SetText(item.RequestBody, false)

	// build header table
	p.headers.Clear()
	p.headers.SetCell(0, 0, tview.NewTableCell("Header").SetTextColor(tview.Styles.TertiaryTextColor))
	p.headers.SetCell(0, 1, tview.NewTableCell("Value").SetTextColor(tview.Styles.TertiaryTextColor))

	row := 1
	for k, v := range item.Headers {
		p.headers.SetCellSimple(row, 0, k)
		p.headers.SetCellSimple(row, 1, strings.Join(v, "; "))
		row++
	}

	if len(item.Headers) > 0 {
		p.headers.Select(1, 0)
	}
}

// Widget returns a primitive widget containing this component.
func (p *RequestView) Widget() *tview.Flex {
	return p.flex
}

func (p *RequestView) build(title string) {
	p.flex = tview.NewFlex()
	p.flex.SetBorder(true)
	p.flex.SetTitle(title)

	p.focusHolder = tview.NewTextView()

	p.pages = tview.NewPages()
	p.flex.AddItem(p.focusHolder, 1, 0, true)
	p.flex.AddItem(p.pages, 0, 1, false)

	p.body = tview.NewTextArea()
	p.body.SetChangedFunc(p.handleBodyChange)

	p.headers = tview.NewTable()
	p.headers.SetSelectable(true, false)

	p.pages.AddAndSwitchToPage("body", p.body, true)
	p.pages.AddPage("headers", p.headers, true, false)

	p.focusManager = util.NewFocusManager(p, GetApplication(), events.Dispatcher(), p.focusHolder, p.focusHolder, p.body)
	p.focusManager.AddArrowNavigation(util.FocusUp, util.FocusLeft, util.FocusRight)
	p.focusManager.SetFilter(p.filterKeyEvent)
	p.focusManager.SetHandler(p.handleKeyEvent)

	p.body.SetInputCapture(p.focusManager.HandleKeyEvent)
	p.flex.SetInputCapture(p.focusManager.HandleKeyEvent)
}

func (p *RequestView) filterKeyEvent(event *tcell.EventKey) util.FocusFilterResult {
	// if a modal is shown, do not process any key events and let the modal handle them instead
	if name, _ := p.pages.GetFrontPage(); name == "modal" {
		return util.FocusPreHandlePropagate
	}

	return util.FocusPreHandleProcess
}

func (p *RequestView) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == '1' {
		p.pages.SwitchToPage("body")
	} else if event.Rune() == '2' {
		p.pages.SwitchToPage("headers")
		GetApplication().SetFocus(p.headers)
	} else if event.Rune() == '+' {
		p.showHeaderModal()
	} else if event.Rune() == '-' {
		p.removeHeader()
	} else {
		return event
	}

	return nil
}

func (p *RequestView) removeHeader() {
	item := p.state.Get().ActiveItem
	row, _ := p.headers.GetSelection()
	if item == nil || row < 1 {
		return
	}

	key := p.headers.GetCell(row, 0)
	value := p.headers.GetCell(row, 1)
	for k, v := range item.Headers {
		if k == key.Text && strings.Join(v, "; ") == value.Text {
			delete(item.Headers, k)
			break
		}
	}

	p.Reload()
}

func (p *RequestView) showHeaderModal() {
	m := NewKeyValueModal("Add Header", "Header", "Value", p.handleAddHeader, p.hideModal)
	p.pages.AddPage("modal", m.Widget(), true, true)
	GetApplication().SetFocus(m.Widget())
}

func (p *RequestView) handleAddHeader(key string, value string) {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	if _, ok := item.Headers[key]; !ok {
		item.Headers[key] = []string{}
	}

	item.Headers[key] = append(item.Headers[key], value)
	p.hideModal()
	p.Reload()
}

func (p *RequestView) hideModal() {
	p.pages.RemovePage("modal")
	p.pages.SwitchToPage("headers")

	// return focus to the pages
	GetApplication().SetFocus(p.pages)
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
