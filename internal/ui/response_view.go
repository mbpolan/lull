package ui

import (
	"fmt"
	"github.com/mbpolan/lull/internal/parsers"
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
)

// ResponseView is a component that allows viewing HTTP response attributes.
type ResponseView struct {
	flex   *tview.Flex
	pages  *tview.Pages
	status *tview.TextView
	body   *tview.TextView
	state  *state.Manager
}

// NewResponseView returns a new instance of ResponseView.
func NewResponseView(title string, state *state.Manager) *ResponseView {
	v := new(ResponseView)
	v.state = state
	v.build(title)

	return v
}

// Reload refreshes the state of the component with current app state.
func (p *ResponseView) Reload() {
	item := p.state.Get().ActiveItem
	if item == nil {
		return
	}

	res := item.Response

	if res == nil {
		p.status.SetText("")
		p.body.SetText("")
	} else {
		body := ""

		// get a parser that's most suitable for the response and format the body
		parser := parsers.GetBodyParser(res)
		body, err := parser.Parse(res)
		if err != nil {
			body = fmt.Sprintf("[red]%+v", err)
		}

		p.status.SetText(p.statusLine(res.StatusCode, res.Status))
		p.body.SetText(body)
	}
}

// Widget returns a primitive widget containing this component.
func (p *ResponseView) Widget() *tview.Flex {
	return p.flex
}

func (p *ResponseView) build(title string) {
	p.flex = tview.NewFlex()
	p.flex.SetBorder(true)
	p.flex.SetTitle(title)
	p.flex.SetDirection(tview.FlexRow)

	p.status = tview.NewTextView()
	p.status.SetDynamicColors(true)

	p.pages = tview.NewPages()
	p.flex.AddItem(p.status, 1, 0, false)
	p.flex.AddItem(p.pages, 0, 1, true)

	p.body = tview.NewTextView()
	p.pages.AddAndSwitchToPage("body", p.body, true)
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
