package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/network"
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
	"net/url"
	"strings"
)

var application *tview.Application

const (
	rootPageMain             string = "main"
	rootPageSaveRequestModal        = "saveRequestModal"
)

// Root is a top-level container for all application UI components.
type Root struct {
	pages            *tview.Pages
	flex             *tview.Flex
	saveRequestModal *SaveRequestModal
	collection       *Collection
	content          *Content
	currentModal     string
	state            *state.AppState
}

// NewRoot returns a new Root instance.
func NewRoot(app *tview.Application) *Root {
	application = app

	r := new(Root)
	r.currentModal = ""
	r.state = state.NewAppState()
	r.state.Method = "GET"
	r.build()

	return r
}

// GetApplication returns the shared instance of tview.Application.
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

	r.saveRequestModal = NewSaveRequestModal(r.handleSaveCurrentRequest, r.hideCurrentModal)

	// create pages containing the main content and various modals that can be opened
	r.pages.AddAndSwitchToPage(rootPageMain, r.flex, true)
	r.pages.AddPage(rootPageSaveRequestModal, r.saveRequestModal.Widget(), true, false)
}

func (r *Root) handleKeyAction(code tcell.Key, key rune) bool {
	switch code {
	case tcell.KeyCtrlL:
		r.collection.SetFocus()
	case tcell.KeyCtrlA:
		r.content.SetFocus(ContentURLBox)
	case tcell.KeyCtrlR:
		r.content.SetFocus(ContentRequestBody)
	case tcell.KeyCtrlG:
		r.sendCurrentRequest()
	case tcell.KeyCtrlS:
		r.showSaveCurrentRequest()
	default:
		return false
	}

	return true
}

func (r *Root) pathToSelectedCollectionItemGroup() []*state.CollectionItem {
	item := r.state.SelectedItem
	if r.state.SelectedItem == nil {
		// really shouldn't happen; default to collection root
		item = r.state.Collection
	}

	// if this item is not a group, select its parent
	if !item.IsGroup() {
		item = item.Parent()
	}

	return append(item.Ancestors(), item)
}

func (r *Root) showSaveCurrentRequest() {
	// collect ancestors and form a path
	ancestors := r.pathToSelectedCollectionItemGroup()
	path := make([]string, len(ancestors))
	for _, i := range ancestors {
		path = append(path, i.Name())
	}

	r.saveRequestModal.SetPathText(strings.Join(path, " > "))
	r.showModal(rootPageSaveRequestModal)
}

func (r *Root) handleSaveCurrentRequest(name string) {
	path := r.pathToSelectedCollectionItemGroup()
	leaf := path[len(path)-1]

	// collect current request information and add it to the collection
	item := state.NewCollectionRequest(name, r.state.Method, r.state.URL, leaf)
	leaf.AddChild(item)

	r.collection.Reload()
	r.hideCurrentModal()
}

func (r *Root) showModal(pageName string) {
	r.currentModal = pageName
	r.pages.ShowPage(pageName)
}

func (r *Root) hideCurrentModal() {
	if r.currentModal != "" {
		r.pages.HidePage(r.currentModal)
		r.currentModal = ""
	}
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
