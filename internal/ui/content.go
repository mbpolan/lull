package ui

import (
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
)

// Content provides a view that shows a request, response and URL input box.
type Content struct {
	flex     *tview.Flex
	url      *URLBox
	request  *RequestView
	response *ResponseView
	state    *state.Manager
}

type ContentWidget int16

const (
	ContentRequestBody ContentWidget = iota
	ContentResponseBody
	ContentURLBox
)

// NewContent returns a new Content instance.
func NewContent(state *state.Manager) *Content {
	c := new(Content)
	c.state = state
	c.build()

	events.Dispatcher().Subscribe(c, []events.Code{events.EventNavigateLeft, events.EventNavigateUp, events.EventNavigateDown, events.EventNavigateRight})

	return c
}

// Reload refreshes the state of the component with current app state.
func (c *Content) Reload() {
	c.url.Reload()
	c.request.Reload()
	c.response.Reload()
}

func (c *Content) HandleEvent(code events.Code, payload events.Payload) {
	switch code {
	case events.EventNavigateLeft:
		// navigate left from response
		if payload.Sender == c.response {
			c.SetFocus(ContentRequestBody)
		} else if payload.Sender == c.url || payload.Sender == c.request {
			events.Dispatcher().PostSimple(events.EventNavigateLeft, c)
		}
	case events.EventNavigateUp:
		// navigate up from request or response
		if payload.Sender == c.request || payload.Sender == c.response {
			c.SetFocus(ContentURLBox)
		}
	case events.EventNavigateDown:
		// navigate down from url box
		if payload.Sender == c.url {
			c.SetFocus(ContentRequestBody)
		}
	case events.EventNavigateRight:
		// navigate right from url box or from request
		if payload.Sender == c.url || payload.Sender == c.request {
			c.SetFocus(ContentResponseBody)
		}
	default:
		break
	}
}

// Widget returns a primitive widget containing this component.
func (c *Content) Widget() tview.Primitive {
	return c.flex
}

// SetFocus sets one of the content primitives to have focus.
func (c *Content) SetFocus(widget ContentWidget) {
	switch widget {
	case ContentRequestBody:
		c.request.SetFocus()
	case ContentResponseBody:
		c.response.SetFocus()
	case ContentURLBox:
		c.url.SetFocus()
	}
}

// HasFocus returns if a child primitives has focus.
func (c *Content) HasFocus(widget ContentWidget) bool {
	switch widget {
	case ContentURLBox:
		return c.url.Widget().HasFocus()
	case ContentRequestBody:
		return c.request.Widget().HasFocus()
	default:
		return false
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
	c.flex.AddItem(c.url.Widget(), 3, 0, true)
	c.flex.AddItem(split, 0, 5, false)
}
