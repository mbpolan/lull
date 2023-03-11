package ui

import (
	"github.com/mbpolan/lull/internal/state/auth"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

type BasicAuthChangeHandler func(data *auth.BasicAuthentication)

// BasicAuthView contains form fields that represent HTTP basic authentication parameters.
type BasicAuthView struct {
	grid         *tview.Grid
	username     *tview.InputField
	password     *tview.InputField
	focusManager *util.FocusManager
	handler      BasicAuthChangeHandler
}

// NewBasicAuthView returns a new BasicAuthView instance configured with a change handler function.
func NewBasicAuthView(handler BasicAuthChangeHandler, manager *util.FocusManager) *BasicAuthView {
	v := &BasicAuthView{
		handler:      handler,
		focusManager: manager,
	}

	v.build()
	return v
}

// Data returns the authentication data provided in the view.
func (a *BasicAuthView) Data() *auth.BasicAuthentication {
	return auth.NewBasicAuthentication(a.username.GetText(), a.password.GetText())
}

// Set applies the values for the basic authentication scheme.
func (a *BasicAuthView) Set(data *auth.BasicAuthentication) {
	a.username.SetText(data.Username)
	a.password.SetText(data.Password)
}

// FocusPrimitives returns a slice of primitives that should receive focus.
func (a *BasicAuthView) FocusPrimitives() []tview.Primitive {
	return []tview.Primitive{
		a.username,
		a.password,
	}
}

// SetFocus sets the focus on this component.
func (a *BasicAuthView) SetFocus() {
	GetApplication().SetFocus(a.FocusPrimitives()[0])
}

// Widget returns a primitive widget containing this component.
func (a *BasicAuthView) Widget() tview.Primitive {
	return a.grid
}

func (a *BasicAuthView) build() {
	a.grid = tview.NewGrid()

	// give the input fields as must space as possible, fix the size of the labels
	a.grid.SetColumns(15, -1)
	a.grid.SetRows(2, 2, 2, 2, 2, -1)

	a.username = tview.NewInputField()
	a.username.SetChangedFunc(a.handleParameterChange)

	a.password = tview.NewInputField()
	a.password.SetChangedFunc(a.handleParameterChange)

	a.grid.AddItem(a.label("Username"), 0, 0, 1, 1, 0, 0, false)
	a.grid.AddItem(a.username, 0, 1, 1, 1, 0, 0, true)

	a.grid.AddItem(a.label("Password"), 1, 0, 1, 1, 0, 0, false)
	a.grid.AddItem(a.password, 1, 1, 1, 1, 0, 0, false)

	// fill remaining vertical space
	a.grid.AddItem(tview.NewBox(), 2, 1, 1, 2, 0, 0, false)

	a.username.SetInputCapture(a.focusManager.HandleKeyEvent)
	a.password.SetInputCapture(a.focusManager.HandleKeyEvent)
}

func (a *BasicAuthView) label(text string) *tview.TextView {
	t := tview.NewTextView()
	t.SetText(text)
	return t
}

func (a *BasicAuthView) handleParameterChange(_ string) {
	a.handler(a.Data())
}
