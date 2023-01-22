package ui

import (
	"context"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/network"
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
	"strings"
)

var application *tview.Application

const (
	rootPageMain  string = "main"
	rootPageModal string = "modal"
)

// Root is a top-level container for all application UI components.
type Root struct {
	pages        *tview.Pages
	flex         *tview.Flex
	collection   *Collection
	content      *Content
	StatusBar    *StatusBar
	currentModal string
	network      *network.Manager
	state        *state.Manager
}

// NewRoot returns a new Root instance.
func NewRoot(app *tview.Application, stateManager *state.Manager) *Root {
	application = app

	r := new(Root)
	r.currentModal = ""
	r.network = network.NewNetworkManager(r.handleRequestFinished)
	r.state = stateManager
	r.build()

	events.Dispatcher().Subscribe(r, []events.Code{events.EventNavigateRight, events.EventNavigateLeft})

	return r
}

// GetApplication returns the shared instance of tview.Application.
func GetApplication() *tview.Application {
	return application
}

func (r *Root) HandleEvent(code events.Code, payload events.Payload) {
	switch code {
	case events.EventNavigateRight:
		// navigate right from collection
		if payload.Sender == r.collection {
			r.content.SetFocus(ContentURLBox)
		}
	case events.EventNavigateLeft:
		// navigate left from content
		if payload.Sender == r.content {
			r.collection.SetFocus()
		}
	default:
		break
	}
}

// Widget returns a primitive widget containing this component.
func (r *Root) Widget() tview.Primitive {
	return r.pages
}

func (r *Root) build() {
	r.pages = tview.NewPages()

	// create child widgets
	r.collection = NewCollection(r.state)
	r.collection.SetItemActivatedHandler(r.handleCollectionItemAction)
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
			if r.handleControlKeyAction(event.Key(), event.Rune()) {
				return nil
			}
		} else if event.Modifiers()&tcell.ModShift == tcell.ModShift {
			if r.handleShiftKeyAction(event.Key(), event.Rune()) {
				return nil
			}
		}

		return event
	})

	// create pages containing the main content and various modals that can be opened
	r.pages.AddAndSwitchToPage(rootPageMain, r.flex, true)
}

func (r *Root) handleControlKeyAction(code tcell.Key, key rune) bool {
	switch code {
	case tcell.KeyCtrlL:
		r.collection.SetFocus()
	case tcell.KeyCtrlA:
		r.content.SetFocus(ContentURLBox)
	case tcell.KeyCtrlR:
		r.content.SetFocus(ContentRequestBody)
	case tcell.KeyCtrlY:
		r.content.SetFocus(ContentResponseBody)
	case tcell.KeyCtrlG:
		r.sendCurrentRequest()
	case tcell.KeyCtrlS:
		r.showSaveCurrentRequest()
	default:
		return false
	}

	return true
}

func (r *Root) handleShiftKeyAction(code tcell.Key, key rune) bool {
	switch code {
	case tcell.KeyRight:
		if !r.content.Widget().HasFocus() {
			r.content.SetFocus(ContentURLBox)
		}
	case tcell.KeyLeft:
		if !r.collection.Widget().HasFocus() {
			r.collection.SetFocus()
		}
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

func (r *Root) handleCollectionItemAction(action CollectionItemAction, item *state.CollectionItem) {
	switch action {
	case CollectionItemRename:
		r.handleRenameSelectedItem(item)
	case CollectionItemAdd:
		r.handleAddItem(item)
	case CollectionItemDelete:
		r.handleDeleteSelectedItem(item)
	case CollectionItemClone:
		r.handleCloneSelectedItem(item)
	case CollectionItemOpen:
		r.setCurrentRequest(item)
	}
}

func (r *Root) handleRenameSelectedItem(item *state.CollectionItem) {
	text := fmt.Sprintf("Current request name: [yellow]%s", item.Name)
	m := NewTextInputModal("Rename Item", text, "New Name", r.renameSelectedItem, r.hideCurrentModal)

	r.showModal(m.Widget())
}

func (r *Root) handleAddItem(item *state.CollectionItem) {
	parent := item
	if !item.IsGroup {
		parent = item.Parent
	}

	// generate a unique name amongst the children of this group
	unique := func(name string) bool {
		for _, i := range parent.Children {
			if i.Name == name {
				return false
			}
		}

		return true
	}

	name := ""
	for i := 1; i > 0; i++ {
		name = fmt.Sprintf("Request %d", i)
		if unique(name) {
			break
		}
	}

	newItem := state.NewCollectionRequest(name, "GET", "", parent)
	parent.AddChild(newItem)

	r.state.Get().SelectedItem = newItem
	r.state.Get().ActiveItem = newItem
	r.collection.Reload()
	r.content.Reload()
}

func (r *Root) handleDeleteSelectedItem(item *state.CollectionItem) {
	text := fmt.Sprintf("Are you sure you want to delete [yellow]%s?", item.Name)
	m := NewPromptModal("Delete Item", text, r.deleteSelectedItem, r.hideCurrentModal)

	r.showModal(m.Widget())
}

func (r *Root) handleCloneSelectedItem(item *state.CollectionItem) {
	text := fmt.Sprintf("Clone request [yellow]%s", item.Name)
	m := NewTextInputModal("Clone Item", text, "Name", r.cloneSelectedItem, r.hideCurrentModal)

	r.showModal(m.Widget())
}

func (r *Root) deleteSelectedItem() {
	item := r.state.Get().SelectedItem
	if item == nil {
		return
	}

	// remove this item from the collection
	r.state.Get().RemoveCollectionItem(item)

	// find another item to select
	candidate := r.state.Get().FirstCollectionItem(func(i *state.CollectionItem) bool {
		return !i.IsGroup
	})

	// if there are no other items to activate, create a new one first
	if candidate == nil {
		candidate = state.NewCollectionRequest("Unnamed", "GET", "", r.state.Get().Collection)
		r.state.Get().Collection.AddChild(candidate)
	}

	// if we're deleting the currently active item, set it to the candidate item as well
	// the active item will be set to nil if that's the case as a result of the call to RemoveCollectionItem above
	if r.state.Get().ActiveItem == nil {
		r.state.Get().ActiveItem = candidate
		r.content.Reload()
	}

	r.state.Get().SelectedItem = candidate

	r.collection.Reload()
	r.hideCurrentModal()
}

func (r *Root) renameSelectedItem(text string) {
	item := r.state.Get().SelectedItem
	if item == nil {
		return
	}

	// rename this item and reload our content
	item.Name = text
	r.collection.Reload()
	r.content.Reload()

	r.hideCurrentModal()
}

func (r *Root) cloneSelectedItem(text string) {
	item := r.state.Get().SelectedItem
	if item == nil {
		return
	}

	// cloning a group item means we need to do a deep copy of all its children as well
	if item.IsGroup {
		// TODO
	} else {
		newItem := state.NewCollectionRequest(text, item.Method, item.URL, item.Parent)
		item.Parent.InsertChildAfter(newItem, item)

		// automatically select and activate the newly cloned item
		r.state.Get().ActiveItem = newItem
		r.state.Get().SelectedItem = newItem
	}

	r.collection.Reload()
	r.content.Reload()

	r.hideCurrentModal()
}

func (r *Root) showSaveCurrentRequest() {
	// collect ancestors and form a path
	ancestors := r.pathToSelectedCollectionItemGroup()
	path := make([]string, len(ancestors))
	for _, i := range ancestors {
		path = append(path, i.Name)
	}

	text := fmt.Sprintf("Request will be saved under [yellow]%s", strings.Join(path, " > "))
	m := NewTextInputModal("Save Request", text, "Name", r.handleSaveCurrentRequest, r.hideCurrentModal)
	r.showModal(m.Widget())
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
	item.Result = active.Result
	leaf.AddChild(item)

	r.collection.Reload()
	r.hideCurrentModal()

	r.state.SetDirty()
}

func (r *Root) handleCancelCurrentRequest() {
	r.network.CancelCurrent()
	r.hideCurrentModal()
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

func (r *Root) showModal(modal tview.Primitive) {
	r.pages.AddPage(rootPageModal, modal, true, true)
}

func (r *Root) hideCurrentModal() {
	r.pages.RemovePage(rootPageModal)
}

func (r *Root) sendCurrentRequest() {
	item := r.state.Get().ActiveItem
	if item == nil {
		return
	}

	var m *AlertModal
	if err := r.network.SendRequest(item); err != nil {
		m = NewAlertModal("Error", fmt.Sprintf("Can't send this request: %s", err.Error()), "OK", r.hideCurrentModal)
	} else {
		m = NewAlertModal("Sending", "Request is in flight...", "Cancel", r.handleCancelCurrentRequest)
	}

	r.showModal(m.Widget())
}

func (r *Root) handleRequestFinished(item *state.CollectionItem, result *network.Result) {
	if result.Error != nil {
		item.Result = nil
		r.state.Get().LastError = result.Error

		// if the error is because the request was cancelled, we don't need to show any modals
		if errors.Is(result.Error, context.Canceled) {
			return
		}

		GetApplication().QueueUpdateDraw(func() {
			m := NewAlertModal("Error", fmt.Sprintf("Could not send request. Error: %s", result.Error.Error()), "OK", r.hideCurrentModal)
			r.hideCurrentModal()
			r.showModal(m.Widget())

			r.content.Reload()
			r.state.SetDirty()
		})

		return
	}

	r.state.Get().LastError = nil
	item.Result = &state.HTTPResult{
		Response:     result.Response,
		Payload:      result.Payload,
		PayloadError: result.PayloadError,
		Duration:     result.EndTime.Sub(result.StartTime),
	}

	GetApplication().QueueUpdateDraw(func() {
		r.hideCurrentModal()
		r.content.Reload()
		r.state.SetDirty()
	})
}
