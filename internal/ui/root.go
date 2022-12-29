package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/network"
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
	"net/url"
)

// Root is a top-level container for all application UI components.
type Root struct {
	pages            *tview.Pages
	flex             *tview.Flex
	saveRequestModal *SaveRequestModal
	collection       *Collection
	content          *Content
	state            *state.AppState
}

var application *tview.Application

// NewRoot returns a new Root instance.
func NewRoot(app *tview.Application) *Root {
	application = app

	r := new(Root)
	r.state = new(state.AppState)
	r.state.Method = "GET"
	r.build()

	return r
}

func GetApplication() *tview.Application {
	return application
}

// Widget returns a primitive widget containing this component.
func (r *Root) Widget() *tview.Pages {
	return r.pages
}

func (r *Root) build() {
	r.pages = tview.NewPages()

	// create child widgets
	r.collection = NewCollection(r.state)
	r.content = NewContent(r.state)

	// arrange them in a flex layout
	r.flex = tview.NewFlex()
	r.flex.AddItem(r.collection.Widget(), 25, 0, false)
	r.flex.AddItem(r.content.Widget(), 0, 1, true)

	r.flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers()&tcell.ModCtrl > 0 {
			if r.handleKeyAction(event.Key(), event.Rune()) {
				return nil
			}
		}

		return event
	})

	r.saveRequestModal = NewSaveRequestModal(r.handleSaveCurrentRequest, r.handleCloseModal)

	r.pages.AddAndSwitchToPage("main", r.flex, true)
	r.pages.AddPage("saveRequestModal", r.saveRequestModal.Widget(), true, false)
}

func (r *Root) handleKeyAction(code tcell.Key, key rune) bool {
	switch code {
	case tcell.KeyCtrlA:
		r.content.SetFocus(ContentURLBox)
	case tcell.KeyCtrlR:
		r.content.SetFocus(ContentRequestBody)
	case tcell.KeyCtrlG:
		r.sendCurrentRequest()
	case tcell.KeyCtrlS:
		r.saveCurrentRequest()
	default:
		return false
	}

	return true
}

func (r *Root) saveCurrentRequest() {
	r.pages.ShowPage("saveRequestModal")
}

func (r *Root) handleSaveCurrentRequest(name string) {
	r.handleCloseModal()
}

func (r *Root) sendCurrentRequest() {
	client := network.NewClient()

	uri, err := url.Parse(r.state.URL)
	if err != nil {
		fmt.Printf("Shit: %+v\n", err)
		return // FIXME
	}

	res, err := client.Exchange(r.state.Method, uri, r.state.RequestBody)
	if err != nil {
		fmt.Printf("Shit: %+v\n", err)
		return // FIXME
	}

	r.state.Response = res
	r.state.LastError = err

	r.content.SetResponse(res)
}

func (r *Root) handleCloseModal() {
	r.pages.HidePage("saveRequestModal")
}
