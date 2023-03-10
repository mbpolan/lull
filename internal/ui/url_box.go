package ui

import (
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

// URLBox is a view that contains an HTTP method, URL and other input components.
type URLBox struct {
	flex           *tview.Flex
	method         *tview.DropDown
	url            *tview.InputField
	focusHolder    *tview.TextView
	focusManager   *util.FocusManager
	allowedMethods []string
	sbSequences    []events.StatusBarContextChangeSequence
	state          *state.Manager
}

// NewURLBox returns a new instance of URLBox.
func NewURLBox(state *state.Manager) *URLBox {
	u := new(URLBox)
	u.state = state
	u.allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	u.build()

	// no additional key sequences supported by this component
	u.sbSequences = []events.StatusBarContextChangeSequence{}

	return u
}

// SetFocus sets the focus on this component.
func (u *URLBox) SetFocus() {
	events.Dispatcher().Post(events.EventStatusBarContextChange, u, &events.StatusBarContextChangeData{
		Fields: u.sbSequences,
	})

	GetApplication().SetFocus(u.Widget())
}

// Reload refreshes the state of the URL box component with current app state.
func (u *URLBox) Reload() {
	item := u.state.Get().ActiveItem
	if item == nil {
		return
	}

	u.flex.SetTitle(u.title())
	u.method.SetCurrentOption(u.currentMethod())
	u.url.SetText(item.URL)
}

// Widget returns a primitive widget containing this component.
func (u *URLBox) Widget() tview.Primitive {
	return u.flex
}

func (u *URLBox) build() {
	curURL := u.currentURL()
	curMethod := u.currentMethod()

	u.flex = tview.NewFlex()
	u.flex.SetTitle(u.title())
	u.flex.SetBorder(true)
	u.flex.SetDirection(tview.FlexColumn)

	u.method = tview.NewDropDown()
	u.method.SetOptions(u.allowedMethods, u.handleMethodChanged)
	u.method.SetCurrentOption(curMethod)

	u.url = tview.NewInputField()
	u.url.SetText(curURL)
	u.url.SetChangedFunc(u.handleURLChanged)

	u.focusHolder = tview.NewTextView()

	u.focusManager = util.NewFocusManager(u, GetApplication(), events.Dispatcher(), u.focusHolder, u.focusHolder, u.method, u.url)
	u.focusManager.AddArrowNavigation(util.FocusLeft, util.FocusDown, util.FocusRight)

	u.method.SetInputCapture(u.focusManager.HandleKeyEvent)
	u.url.SetInputCapture(u.focusManager.HandleKeyEvent)
	u.flex.SetInputCapture(u.focusManager.HandleKeyEvent)

	u.flex.AddItem(u.focusHolder, 1, 0, true)
	u.flex.AddItem(u.method, 8, 0, false)
	u.flex.AddItem(u.url, 0, 1, false)
}

func (u *URLBox) title() string {
	if selected := u.state.Get().SelectedItem; selected != nil {
		return selected.Name
	}

	return ""
}

func (u *URLBox) currentURL() string {
	item := u.state.Get().ActiveItem
	if item == nil {
		return ""
	}

	return item.URL
}

func (u *URLBox) currentMethod() int {
	item := u.state.Get().ActiveItem
	if item == nil {
		return -1
	}

	for i, method := range u.allowedMethods {
		if method == item.Method {
			return i
		}
	}

	return -1
}

func (u *URLBox) handleMethodChanged(text string, index int) {
	item := u.state.Get().ActiveItem
	if item == nil {
		return
	}

	item.Method = text
	u.state.SetDirty()
}

func (u *URLBox) handleURLChanged(text string) {
	item := u.state.Get().ActiveItem
	if item == nil {
		return
	}

	item.URL = text
	u.state.SetDirty()
}
