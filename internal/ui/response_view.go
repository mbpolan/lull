package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/parsers"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
	"strings"
)

const responseViewTitle = "Response"

const responseViewBody = "body"
const responseViewHeaders = "headers"

// ResponseView is a component that allows viewing HTTP response attributes.
type ResponseView struct {
	flex         *tview.Flex
	pages        *tview.Pages
	status       *tview.TextView
	metrics      *tview.TextView
	body         *tview.TextView
	headers      *tview.Table
	focusHolder  *tview.TextView
	focusManager *util.FocusManager
	sbSequences  []events.StatusBarContextChangeSequence
	state        *state.Manager
}

// NewResponseView returns a new instance of ResponseView.
func NewResponseView(state *state.Manager) *ResponseView {
	p := new(ResponseView)
	p.state = state
	p.build()
	p.Reload()

	p.sbSequences = []events.StatusBarContextChangeSequence{
		{
			Label:       "Body",
			KeySequence: "1",
		},
		{
			Label:       "Headers",
			KeySequence: "2",
		},
	}

	return p
}

// SetFocus sets the focus on this component.
func (p *ResponseView) SetFocus() {
	events.Dispatcher().Post(events.EventStatusBarContextChange, p, &events.StatusBarContextChangeData{
		Fields: p.sbSequences,
	})

	GetApplication().SetFocus(p.Widget())
}

// Reload refreshes the state of the component with current app state.
func (p *ResponseView) Reload() {
	p.setTitle()

	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	// clear headers table
	p.headers.Clear()
	p.headers.SetCell(0, 0, tview.NewTableCell("Header").SetTextColor(tview.Styles.TertiaryTextColor))
	p.headers.SetCell(0, 1, tview.NewTableCell("Value").SetTextColor(tview.Styles.TertiaryTextColor))

	res := item.Result
	if res == nil {
		p.status.SetText("")
		p.metrics.SetText("")
		p.body.SetText("")
	} else {
		resp := res.Response
		body := ""

		// get a parser that's most suitable for the response and format the body
		parser := parsers.GetBodyParser(resp)
		body, err := parser.Parse(resp)
		if err != nil {
			body = fmt.Sprintf("[red]%+v", err)
		}

		p.status.SetText(p.statusLine(resp.StatusCode, resp.Status))
		p.metrics.SetText(util.FormatDuration(res.Duration))
		p.body.SetText(body)

		// build header table
		row := 1
		for k, v := range resp.Header {
			p.headers.SetCellSimple(row, 0, k)
			p.headers.SetCellSimple(row, 1, strings.Join(v, ";"))
			row++
		}
	}
}

// Widget returns a primitive widget containing this component.
func (p *ResponseView) Widget() tview.Primitive {
	return p.flex
}

func (p *ResponseView) build() {
	p.flex = tview.NewFlex()
	p.flex.SetBorder(true)
	p.flex.SetDirection(tview.FlexRow)

	p.status = tview.NewTextView()
	p.status.SetTextAlign(tview.AlignLeft)
	p.status.SetDynamicColors(true)

	p.metrics = tview.NewTextView()
	p.metrics.SetTextAlign(tview.AlignRight)

	statusFlex := tview.NewFlex()
	statusFlex.SetDirection(tview.FlexColumn)
	statusFlex.AddItem(p.status, 0, 1, false)
	statusFlex.AddItem(p.metrics, 0, 1, false)

	p.focusHolder = tview.NewTextView()

	p.pages = tview.NewPages()
	p.flex.AddItem(p.focusHolder, 0, 0, true)
	p.flex.AddItem(statusFlex, 1, 0, false)
	p.flex.AddItem(p.pages, 0, 1, false)

	p.body = tview.NewTextView()
	p.headers = tview.NewTable()

	p.pages.AddAndSwitchToPage(responseViewBody, p.body, true)
	p.pages.AddPage(responseViewHeaders, p.headers, true, false)

	p.focusManager = util.NewFocusManager(p, GetApplication(), events.Dispatcher(), p.focusHolder, p.focusHolder, p.body)
	p.focusManager.AddArrowNavigation(util.FocusLeft, util.FocusUp)
	p.focusManager.SetHandler(p.handleKeyEvent)

	p.flex.SetInputCapture(p.focusManager.HandleKeyEvent)
}

func (p *ResponseView) setTitle() {
	page, _ := p.pages.GetFrontPage()

	title := responseViewTitle
	if page != "" {
		title = fmt.Sprintf("%s (%s)", responseViewTitle, page)
	}

	p.flex.SetTitle(title)
}

func (p *ResponseView) switchToPage(view string) {
	p.pages.SwitchToPage(view)
	p.setTitle()
}

func (p *ResponseView) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == '1' {
		p.switchToPage(responseViewBody)
	} else if event.Rune() == '2' {
		p.switchToPage(responseViewHeaders)
	} else {
		return event
	}

	return nil
}

func (p *ResponseView) statusLine(code int, status string) string {
	color := ""

	// choose a text color based on the "severity" of the status code
	if code >= 200 && code < 300 {
		color = "green"
	} else if code >= 300 && code < 400 {
		color = "yellow"
	} else if code >= 400 && code < 500 {
		color = "orange"
	} else if code >= 500 && code < 599 {
		color = "red"
	}

	// take the status text from the response or create our own if there is none
	// TODO: parse the status text if it exists and remove duplicate status codes
	statusText := status
	if statusText != "" {
		statusText = p.statusTextForCode(code)
	}

	return fmt.Sprintf("[%s]%d %s", color, code, statusText)
}

func (p *ResponseView) statusTextForCode(code int) string {
	// TODO: add missing status codes
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 202:
		return "Accepted"
	case 204:
		return "No Content"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return ""
	}
}
