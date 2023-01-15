package ui

import (
	"errors"
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
	"strings"
)

const requestViewBody = "body"
const requestViewHeaders = "headers"
const requestViewModal = "modal"

const headerTableSeparator = "; "

// RequestView is a view that allows viewing and editing request/response components.
type RequestView struct {
	flex         *tview.Flex
	pages        *tview.Pages
	body         *tview.TextArea
	headers      *tview.Table
	focusHolder  *tview.TextView
	focusManager *util.FocusManager
	state        *state.Manager
}

// NewRequestView returns a new instance of RequestView.
func NewRequestView(title string, state *state.Manager) *RequestView {
	p := new(RequestView)
	p.state = state
	p.build(title)
	p.Reload()

	return p
}

// SetFocus sets the focus on this component.
func (p *RequestView) SetFocus() {
	p.postKeyboardSequences()
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
		p.headers.SetCellSimple(row, 1, strings.Join(v, headerTableSeparator))
		row++
	}

	if len(item.Headers) > 0 {
		p.headers.Select(1, 0)
	}
}

// Widget returns a primitive widget containing this component.
func (p *RequestView) Widget() tview.Primitive {
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
	p.headers.SetSelectedFunc(p.showEditHeaderModal)

	p.pages.AddAndSwitchToPage(requestViewBody, p.body, true)
	p.pages.AddPage(requestViewHeaders, p.headers, true, false)

	p.focusManager = util.NewFocusManager(p, GetApplication(), events.Dispatcher(), p.focusHolder, p.focusHolder, p.body)
	p.focusManager.AddArrowNavigation(util.FocusUp, util.FocusLeft, util.FocusRight)
	p.focusManager.SetFilter(p.filterKeyEvent)
	p.focusManager.SetHandler(p.handleKeyEvent)

	p.body.SetInputCapture(p.focusManager.HandleKeyEvent)
	p.flex.SetInputCapture(p.focusManager.HandleKeyEvent)
}

func (p *RequestView) filterKeyEvent(event *tcell.EventKey) util.FocusFilterResult {
	// if a modal is shown, do not process any key events and let the modal handle them instead
	if name, _ := p.pages.GetFrontPage(); name == requestViewModal {
		return util.FocusPreHandlePropagate
	}

	return util.FocusPreHandleProcess
}

func (p *RequestView) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == '1' {
		p.pages.SwitchToPage(requestViewBody)
		p.postKeyboardSequences()
	} else if event.Rune() == '2' {
		p.pages.SwitchToPage(requestViewHeaders)
		p.postKeyboardSequences()
		GetApplication().SetFocus(p.headers)
	} else if event.Rune() == '+' {
		p.showAddHeaderModal()
	} else if event.Rune() == '-' {
		p.removeHeader()
	} else {
		return event
	}

	return nil
}

func (p *RequestView) removeHeader() {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	key, _, err := p.currentHeader()
	if err != nil {
		return
	}

	item.RemoveHeader(key)
	p.Reload()
}

func (p *RequestView) showAddHeaderModal() {
	m := NewKeyValueModal("Add Header", "Header", "Value", p.handleAddHeader, p.hideModal)
	p.pages.AddPage(requestViewModal, m.Widget(), true, true)
	GetApplication().SetFocus(m.Widget())
}

func (p *RequestView) handleBodyChange() {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	item.RequestBody = p.body.GetText()
}

func (p *RequestView) handleAddHeader(key string, value string) {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	item.AddHeader(key, value)
	p.hideModal()
	p.Reload()
}

func (p *RequestView) showEditHeaderModal(row int, _ int) {
	key, value, err := p.currentHeader()
	if err != nil {
		return
	}

	m := NewKeyValueModal("Edit Header", "Header", "Value", p.handleEditHeader, p.hideModal)
	m.SetKey(key)
	m.SetValue(strings.Join(value, headerTableSeparator))

	p.pages.AddPage(requestViewModal, m.Widget(), true, true)
	GetApplication().SetFocus(m.Widget())
}

func (p *RequestView) handleEditHeader(key string, value string) {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	prevKey, _, err := p.currentHeader()
	if err != nil {
		return
	}

	newValues := strings.Split(value, headerTableSeparator)

	// if the key has changed, we need to remove the existing header entry entirely
	if prevKey == key {
		item.Headers[key] = newValues
	} else {
		item.RemoveHeader(prevKey)

		// if there is already an existing header with the new key, concat the values
		if v, ok := item.Headers[key]; ok {
			item.Headers[key] = append(v, newValues...)
		} else {
			item.Headers[key] = newValues
		}
	}

	p.hideModal()
	p.Reload()
}

func (p *RequestView) hideModal() {
	p.pages.RemovePage(requestViewModal)
	p.pages.SwitchToPage(requestViewHeaders)

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

func (p *RequestView) currentHeader() (string, []string, error) {
	row, _ := p.headers.GetSelection()
	if row < 1 {
		return "", nil, errors.New("no header selected")
	}

	key := p.headers.GetCell(row, 0)
	value := p.headers.GetCell(row, 1)

	return key.Text, strings.Split(value.Text, headerTableSeparator), nil
}

func (p *RequestView) keyboardSequences() []events.StatusBarContextChangeSequence {
	var seq []events.StatusBarContextChangeSequence
	page, _ := p.pages.GetFrontPage()

	switch page {
	case requestViewBody:
		seq = []events.StatusBarContextChangeSequence{
			{
				Label:       "Headers",
				KeySequence: "2",
			},
		}
	case requestViewHeaders:
		seq = []events.StatusBarContextChangeSequence{
			{
				Label:       "Body",
				KeySequence: "1",
			},
			{
				Label:       "Add header",
				KeySequence: "+",
			},
			{
				Label:       "Remove header",
				KeySequence: "-",
			},
			{
				Label:       "Edit header",
				KeySequence: "enter",
			},
		}
	default:
		break
	}

	return seq
}

func (p *RequestView) postKeyboardSequences() {
	seq := p.keyboardSequences()

	events.Dispatcher().Post(events.EventStatusBarContextChange, p, &events.StatusBarContextChangeData{
		Fields: seq,
	})
}
