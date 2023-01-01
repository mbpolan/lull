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
	StatusBar        *StatusBar
	currentModal     string
	state            *state.Manager
}

// NewRoot returns a new Root instance.
func NewRoot(app *tview.Application, stateManager *state.Manager) *Root {
	application = app

	r := new(Root)
	r.currentModal = ""
	r.state = stateManager
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
	r.collection.SetItemActivatedHandler(r.setCurrentRequest)
	r.content = NewContent(r.state)
	r.StatusBar = NewStatusBar()

	// arrange the collection and content in a flex layout
	mc := tview.NewFlex()
	mc.AddItem(r.collection.Widget(), 25, 0, false)
	mc.AddItem(r.content.Widget(), 0, 1, true)

	// arrange the main content flex layout and the status bar in a parent flex
	r.flex = tview.NewFlex()
	r.flex.SetDirection(tview.FlexRow)
	r.flex.AddItem(mc, 0, 1, true)
	r.flex.AddItem(r.StatusBar.Widget(), 1, 0, false)

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
	case tcell.KeyCtrlQ:
		application.Stop()
	default:
		return false
	}

	return true
}

func (r *Root) pathToSelectedCollectionItemGroup() []*state.CollectionItem {
	item := r.state.Get().SelectedItem
	if item == nil {
		// really shouldn't happen; default to collection root
		item = r.state.Get().Collection
	}

	// if this item is not a group, select its parent
	if !item.IsGroup {
		item = item.Parent
	}

	return append(item.Ancestors(), item)
}

func (r *Root) showSaveCurrentRequest() {
	// collect ancestors and form a path
	ancestors := r.pathToSelectedCollectionItemGroup()
	path := make([]string, len(ancestors))
	for _, i := range ancestors {
		path = append(path, i.Name)
	}

	r.saveRequestModal.SetPathText(strings.Join(path, " > "))
	r.showModal(rootPageSaveRequestModal)
}

func (r *Root) handleSaveCurrentRequest(name string) {
	active := r.state.Get().ActiveItem
	if active == nil {
		return
	}

	path := r.pathToSelectedCollectionItemGroup()
	leaf := path[len(path)-1]

	// collect current request information and add it to the collection
	item := state.NewCollectionRequest(name, active.Method, active.URL, leaf)
	item.RequestBody = active.RequestBody
	item.Response = active.Response
	leaf.AddChild(item)

	r.collection.Reload()
	r.hideCurrentModal()

	r.state.SetDirty()
}

func (r *Root) setCurrentRequest(item *state.CollectionItem) {
	if r.state.Get().ActiveItem == item {
		return
	}

	// reload views to synchronize with app state
	r.state.Get().ActiveItem = item
	r.content.Reload()
	r.state.SetDirty()
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
	item := r.state.Get().ActiveItem
	if item == nil {
		return
	}

	client := network.NewClient()

	uri, err := url.Parse(item.URL)
	if err != nil {
		fmt.Printf("Shit: %+v\n", err)
		return // FIXME
	}

	res, err := client.Exchange(item.Method, uri, item.RequestBody)
	if err != nil {
		fmt.Printf("Shit: %+v\n", err)
		return // FIXME
	}

	item.Response = res
	r.state.Get().LastError = err
	r.content.Reload()
	r.state.SetDirty()
}
