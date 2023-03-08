package ui

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/parsers"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
	"strings"
)

const requestViewTitle = "Request"

const requestViewBody = "body"
const requestViewHeaders = "headers"
const requestViewAuthentication = "authentication"
const requestViewModal = "modal"

const headerTableSeparator = "; "

var contentTypeOptions = []string{"None", "JSON", "Text"}
var contentTypeOptionsToValues = map[string]string{
	contentTypeOptions[0]: "",
	contentTypeOptions[1]: "application/json",
	contentTypeOptions[2]: "text/plain",
}

// RequestView is a view that allows viewing and editing request/response components.
type RequestView struct {
	flex         *tview.Flex
	pages        *tview.Pages
	body         *tview.TextArea
	auth         *AuthView
	contentType  *tview.DropDown
	headers      *tview.Table
	focusHolder  *tview.TextView
	focusManager *util.FocusManager
	state        *state.Manager
}

// NewRequestView returns a new instance of RequestView.
func NewRequestView(state *state.Manager) *RequestView {
	p := new(RequestView)
	p.state = state
	p.build()
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
	p.setTitle()

	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	if body := item.RequestBody; body != nil {
		contentTypeOption := ""
		for k, v := range contentTypeOptionsToValues {
			if v == "" {
				continue
			} else if strings.Index(body.ContentType, v) > -1 {
				contentTypeOption = k
				break
			}
		}

		if contentTypeOption == "" {
			contentTypeOption = contentTypeOptions[0]
		}

		contentType := -1
		for i, c := range contentTypeOptions {
			if c == contentTypeOption {
				contentType = i
				break
			}
		}

		if contentType == -1 {
			contentType = 0
		}

		p.contentType.SetCurrentOption(contentType)
		p.body.SetText(body.Payload, false)
	} else {
		p.contentType.SetCurrentOption(0)
		p.body.SetText("", false)
	}

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

	// apply authentication
	p.auth.Set(item)
}

// Widget returns a primitive widget containing this component.
func (p *RequestView) Widget() tview.Primitive {
	return p.flex
}

func (p *RequestView) build() {
	p.flex = tview.NewFlex()
	p.flex.SetBorder(true)

	p.focusHolder = tview.NewTextView()

	p.pages = tview.NewPages()
	p.flex.AddItem(p.focusHolder, 1, 0, true)
	p.flex.AddItem(p.pages, 0, 1, false)

	p.body = tview.NewTextArea()
	p.body.SetChangedFunc(p.handleBodyChange)

	p.contentType = tview.NewDropDown()
	p.contentType.SetLabel("Body ")
	p.contentType.SetOptions(contentTypeOptions, p.handleContentTypeChange)

	bodyFlex := tview.NewFlex()
	bodyFlex.SetDirection(tview.FlexRow)
	bodyFlex.AddItem(p.contentType, 1, 0, false)
	bodyFlex.AddItem(p.body, 0, 1, true)

	p.auth = NewAuthView(p.handleAuthenticationChange)

	p.headers = tview.NewTable()
	p.headers.SetSelectable(true, false)
	p.headers.SetSelectedFunc(p.showEditHeaderModal)

	p.pages.AddAndSwitchToPage(requestViewBody, bodyFlex, true)
	p.pages.AddPage(requestViewHeaders, p.headers, true, false)
	p.pages.AddPage(requestViewAuthentication, p.auth.Widget(), true, false)

	p.focusManager = util.NewFocusManager(p, GetApplication(), events.Dispatcher(), p.focusHolder)
	p.focusManager.SetName("request_view")
	p.focusManager.AddArrowNavigation(util.FocusUp, util.FocusLeft, util.FocusRight)
	p.focusManager.SetFilter(p.filterKeyEvent)
	p.focusManager.SetHandler(p.handleKeyEvent)

	p.body.SetInputCapture(p.handleBodyKeyEvent)
	p.flex.SetInputCapture(p.focusManager.HandleKeyEvent)
}

func (p *RequestView) setTitle() {
	page, _ := p.pages.GetFrontPage()
	if page == requestViewModal {
		page = ""
	}

	title := requestViewTitle
	if page != "" {
		title = fmt.Sprintf("%s (%s)", requestViewTitle, page)
	}

	p.flex.SetTitle(title)
}

func (p *RequestView) switchToPage(view string) {
	// change the set of focus primitives based on the newly selected view
	switch view {
	case requestViewBody:
		p.focusManager.SetPrimitives(p.focusHolder, p.contentType, p.body)
		GetApplication().SetFocus(p.contentType)
	case requestViewHeaders:
		p.focusManager.SetPrimitives(p.focusHolder, p.headers)
		GetApplication().SetFocus(p.headers)
	case requestViewAuthentication:
		p.focusManager.SetPrimitives(p.focusHolder)
		GetApplication().SetFocus(p.auth.Widget())
	}

	p.pages.SwitchToPage(view)
	p.postKeyboardSequences()
	p.setTitle()
}

func (p *RequestView) filterKeyEvent(event *tcell.EventKey) util.FocusFilterResult {
	// if a modal is shown, do not process any key events and let the modal handle them instead
	if name, _ := p.pages.GetFrontPage(); name == requestViewModal {
		return util.FocusPreHandlePropagate
	}

	return util.FocusPreHandleProcess
}

func (p *RequestView) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	// ensure that the parent primitive has focus to prevent switching pages while the user is entering text
	// in one of the pages themselves
	if event.Rune() == '1' && p.focusManager.ParentHasFocus() {
		p.switchToPage(requestViewBody)
	} else if event.Rune() == '2' && p.focusManager.ParentHasFocus() {
		p.switchToPage(requestViewHeaders)
	} else if event.Rune() == '3' && p.focusManager.ParentHasFocus() {
		p.switchToPage(requestViewAuthentication)
	} else if event.Rune() == '+' && p.headers.HasFocus() {
		p.showAddHeaderModal()
	} else if event.Rune() == '-' && p.headers.HasFocus() {
		p.removeHeader()
	} else if event.Rune() == 'f' && p.focusManager.ParentHasFocus() {
		p.formatBody()
	} else {
		return event
	}

	return nil
}

func (p *RequestView) formatBody() {
	item := p.state.Get().ActiveItem
	if item == nil || item.RequestBody == nil {
		return
	}

	contentType := p.currentContentType()
	if contentType == "" {
		return
	}

	if parser := parsers.GetBodyParserForContentType(contentType); parser != nil {
		formatted, err := parser.ParseBytes([]byte(item.RequestBody.Payload))
		if err != nil {
			util.ConsoleBell()
			return
		}

		item.RequestBody.Payload = formatted
		p.Reload()
	}
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

func (p *RequestView) handleBodyKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	// prevent inputs when there is no active content type in this request
	if i, _ := p.contentType.GetCurrentOption(); i == 0 {
		return nil
	}

	return p.focusManager.HandleKeyEvent(event)
}

func (p *RequestView) handleBodyChange() {
	item := p.state.Get().ActiveItem
	if item == nil || item.RequestBody == nil {
		return
	}

	item.RequestBody.Payload = p.body.GetText()
}

func (p *RequestView) handleContentTypeChange(text string, index int) {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	if index == 0 {
		item.RequestBody = nil
	} else {
		contentType := contentTypeOptionsToValues[text]

		// if the request previously had no request body, initialize it now
		// otherwise, just update the content type
		if item.RequestBody == nil {
			item.RequestBody = &state.RequestBody{
				ContentType: contentType,
				Payload:     "",
			}
		} else {
			item.RequestBody.ContentType = contentType
		}
	}
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

func (p *RequestView) handleAuthenticationChange(data state.RequestAuthentication) {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	if data == nil {
		item.Authentication.Data = nil
	} else if oauth2 := data.(*state.OAuth2RequestAuthentication); oauth2 != nil {
		item.Authentication.Data = oauth2
	}

	p.state.SetDirty()
}

func (p *RequestView) hideModal() {
	p.pages.RemovePage(requestViewModal)
	p.pages.SwitchToPage(requestViewHeaders)

	// return focus to the pages
	GetApplication().SetFocus(p.pages)
}

func (p *RequestView) currentContentType() string {
	// if the body dropdown has a valid value selected, use that as the content type
	i, contentType := p.contentType.GetCurrentOption()
	if i != 0 {
		return contentTypeOptionsToValues[contentType]
	}

	item := p.state.Get().ActiveItem
	if item == nil {
		return ""
	}

	// try to find a content type header and use the first value
	for k, v := range item.Headers {
		if strings.ToLower(k) == "content-type" {
			return v[0]
		}
	}

	return ""
}

func (p *RequestView) currentRequestBody() string {
	item := p.state.Get().ActiveItem
	if item == nil || item.RequestBody == nil {
		return ""
	}

	return item.RequestBody.Payload
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
			{
				Label:       "Authentication",
				KeySequence: "3",
			},
			{
				Label:       "Format",
				KeySequence: "f",
			},
		}
	case requestViewHeaders:
		seq = []events.StatusBarContextChangeSequence{
			{
				Label:       "Body",
				KeySequence: "1",
			},
			{
				Label:       "Authentication",
				KeySequence: "3",
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
				KeySequence: "⏎",
			},
		}
	case requestViewAuthentication:
		seq = []events.StatusBarContextChangeSequence{
			{
				Label:       "Body",
				KeySequence: "1",
			},
			{
				Label:       "Headers",
				KeySequence: "2",
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
