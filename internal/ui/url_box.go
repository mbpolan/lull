package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
)

// URLBox is a view that contains an HTTP method, URL and other input components.
type URLBox struct {
	flex            *tview.Flex
	method          *tview.DropDown
	url             *tview.InputField
	focusPrimitives []tview.Primitive
	allowedMethods  []string
	state           *state.AppState
}

// NewURLBox returns a new instance of URLBox.
func NewURLBox(state *state.AppState) *URLBox {
	u := new(URLBox)
	u.state = state
	u.allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	u.build()

	return u
}

// Widget returns a primitive widget containing this component.
func (u *URLBox) Widget() *tview.Flex {
	return u.flex
}

func (u *URLBox) build() {
	curMethod := u.currentMethod()

	u.flex = tview.NewFlex()
	u.flex.SetBorder(true)
	u.flex.SetDirection(tview.FlexColumn)

	u.method = tview.NewDropDown()
	u.method.SetOptions(u.allowedMethods, u.handleMethodChanged)
	u.method.SetCurrentOption(curMethod)
	u.method.SetInputCapture(u.handleKeyEvent)

	u.url = tview.NewInputField()
	u.url.SetChangedFunc(u.handleURLChanged)
	u.url.SetInputCapture(u.handleKeyEvent)

	u.focusPrimitives = []tview.Primitive{u.method, u.url}

	u.flex.AddItem(u.method, 8, 0, true)
	u.flex.AddItem(u.url, 0, 1, false)
}

func (u *URLBox) currentMethod() int {
	for i, method := range u.allowedMethods {
		if method == u.state.Method {
			return i
		}
	}

	return -1
}

func (u *URLBox) handleMethodChanged(text string, index int) {
	u.state.Method = text
}

func (u *URLBox) handleURLChanged(text string) {
	u.state.URL = text
}

func (u *URLBox) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	// handle focus changes
	if event.Key() == tcell.KeyTab {
		focused := -1
		for i, w := range u.focusPrimitives {
			if GetApplication().GetFocus() == w {
				focused = i
				break
			}
		}

		if focused > -1 {
			GetApplication().SetFocus(u.focusPrimitives[(focused+1)%len(u.focusPrimitives)])
		}

		return nil
	}

	return event
}
