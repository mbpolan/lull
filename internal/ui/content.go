package ui

import (
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
	"io"
	"net/http"
)

// Content provides a view that shows a request, response and URL input box.
type Content struct {
	flex     *tview.Flex
	url      *URLBox
	request  *Payload
	response *Payload
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
	c.request = NewPayload("Request", false)
	c.response = NewPayload("Response", true)

	split := tview.NewFlex()
	split.AddItem(c.request.Widget(), 0, 1, false)
	split.AddItem(c.response.Widget(), 0, 1, false)

	c.flex = tview.NewFlex()
	c.flex.SetDirection(tview.FlexRow)
	c.flex.AddItem(c.url.Widget(), 5, 0, true)
	c.flex.AddItem(split, 0, 5, false)
}

func (c *Content) SetResponse(res *http.Response, err error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		body = []byte{}
	}

	c.response.SetData(res.StatusCode, body, err)
}
