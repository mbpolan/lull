package ui

import (
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

// URLBox is a view that contains an HTTP method, URL and other input components.
type URLBox struct {
	flex           *tview.Flex
	method         *tview.DropDown
	url            *tview.InputField
	focusManager   *util.FocusManager
	allowedMethods []string
	state          *state.AppState
}

// NewURLBox returns a new instance of URLBox.
func NewURLBox(state *state.AppState) *URLBox {
	u := new(URLBox)
	u.state = state
	u.allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	u.build()

	return u
}

// Reload refreshes the state of the URL box component with current app state.
func (u *URLBox) Reload() {
	u.method.SetCurrentOption(u.currentMethod())
	u.url.SetText(u.state.ActiveItem.URL())
}

// Widget returns a primitive widget containing this component.
func (u *URLBox) Widget() *tview.Flex {
	return u.flex
}

func (u *URLBox) build() {
	curMethod := u.currentMethod()

	u.flex = tview.NewFlex()
	u.flex.SetTitle(u.title())
	u.flex.SetBorder(true)
	u.flex.SetDirection(tview.FlexColumn)

	u.method = tview.NewDropDown()
	u.method.SetOptions(u.allowedMethods, u.handleMethodChanged)
	u.method.SetCurrentOption(curMethod)

	u.url = tview.NewInputField()
	u.url.SetChangedFunc(u.handleURLChanged)

	u.focusManager = util.NewFocusManager(GetApplication(), []tview.Primitive{u.method, u.url})
	u.method.SetInputCapture(u.focusManager.HandleKeyEvent)
	u.url.SetInputCapture(u.focusManager.HandleKeyEvent)

	u.flex.AddItem(u.method, 8, 0, true)
	u.flex.AddItem(u.url, 0, 1, false)
}

func (u *URLBox) title() string {
	if selected := u.state.SelectedItem; selected != nil {
		return selected.Name()
	}

	return ""
}

func (u *URLBox) currentMethod() int {
	if u.state.ActiveItem == nil {
		return -1
	}

	for i, method := range u.allowedMethods {
		if method == u.state.ActiveItem.Method() {
			return i
		}
	}

	return -1
}

func (u *URLBox) handleMethodChanged(text string, index int) {
	u.state.ActiveItem.SetMethod(text)
}

func (u *URLBox) handleURLChanged(text string) {
	u.state.ActiveItem.SetURL(text)
}
