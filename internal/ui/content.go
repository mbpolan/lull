package ui

import (
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
	"net/http"
)

// Content provides a view that shows a request, response and URL input box.
type Content struct {
	flex     *tview.Flex
	url      *URLBox
	request  *RequestView
	response *ResponseView
	state    *state.AppState
}

type ContentWidget int16

const (
	ContentRequestBody ContentWidget = iota
	ContentURLBox
)

// NewContent returns a new Content instance.
func NewContent(state *state.AppState) *Content {
	c := new(Content)
	c.state = state
	c.build()

	return c
}

// Widget returns a primitive widget containing this component.
func (c *Content) Widget() *tview.Flex {
	return c.flex
}

func (c *Content) SetFocus(widget ContentWidget) {
	switch widget {
	case ContentRequestBody:
		GetApplication().SetFocus(c.request.Widget())
	case ContentURLBox:
		GetApplication().SetFocus(c.url.Widget())
	}
}

func (c *Content) build() {
	c.url = NewURLBox(c.state)
	c.request = NewRequestView("Request", c.state)
	c.response = NewResponseView("Response", c.state)

	split := tview.NewFlex()
	split.AddItem(c.request.Widget(), 0, 1, false)
	split.AddItem(c.response.Widget(), 0, 1, false)

	c.flex = tview.NewFlex()
	c.flex.SetDirection(tview.FlexRow)
	c.flex.AddItem(c.url.Widget(), 5, 0, true)
	c.flex.AddItem(split, 0, 5, false)
}

func (c *Content) SetResponse(res *http.Response) {
	c.response.SetResponse(res)
}
