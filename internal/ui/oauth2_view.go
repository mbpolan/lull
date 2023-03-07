package ui

import (
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

type OAuth2ChangeHandler func(data *state.OAuth2RequestAuthentication)

// OAuth2View contains form fields that represent an OAuth2 client credentials configuration.
type OAuth2View struct {
	grid         *tview.Grid
	tokenURL     *tview.InputField
	clientID     *tview.InputField
	clientSecret *tview.InputField
	grantType    *tview.InputField
	scope        *tview.InputField
	focusManager *util.FocusManager
	handler      OAuth2ChangeHandler
}

// NewOAuth2View returns a new OAuth2View instance configured with a change handler function.
func NewOAuth2View(handler OAuth2ChangeHandler, manager *util.FocusManager) *OAuth2View {
	v := &OAuth2View{
		handler:      handler,
		focusManager: manager,
	}

	v.build()
	return v
}

// Data returns the authentication data provided in the view.
func (a *OAuth2View) Data() *state.OAuth2RequestAuthentication {
	auth := state.NewOAuth2RequestAuthentication(a.tokenURL.GetText(), a.clientID.GetText(), a.clientSecret.GetText(),
		a.grantType.GetText(), a.scope.GetText())

	return auth
}

// Set applies the values for the OAuth2 authentication scheme.
func (a *OAuth2View) Set(data *state.OAuth2RequestAuthentication) {
	a.tokenURL.SetText(data.TokenURL)
	a.clientID.SetText(data.ClientID)
	a.clientSecret.SetText(data.ClientSecret)
	a.grantType.SetText(data.GrantType)
	a.scope.SetText(data.Scope)
}

// FocusPrimitives returns a slice of primitives that should receive focus.
func (a *OAuth2View) FocusPrimitives() []tview.Primitive {
	return []tview.Primitive{
		a.tokenURL,
		a.clientID,
		a.clientSecret,
		a.grantType,
		a.scope,
	}
}

// SetFocus sets the focus on this component.
func (a *OAuth2View) SetFocus() {
	GetApplication().SetFocus(a.FocusPrimitives()[0])
}

// Widget returns a primitive widget containing this component.
func (a *OAuth2View) Widget() tview.Primitive {
	return a.grid
}

func (a *OAuth2View) build() {
	a.grid = tview.NewGrid()

	// give the input fields as must space as possible, fix the size of the labels
	a.grid.SetColumns(15, -1)
	a.grid.SetRows(2, 2, 2, 2, 2, -1)

	a.tokenURL = tview.NewInputField()
	a.tokenURL.SetChangedFunc(a.handleParameterChange)

	a.clientID = tview.NewInputField()
	a.clientID.SetChangedFunc(a.handleParameterChange)

	a.clientSecret = tview.NewInputField()
	a.clientSecret.SetChangedFunc(a.handleParameterChange)

	a.grantType = tview.NewInputField()
	a.grantType.SetChangedFunc(a.handleParameterChange)

	a.scope = tview.NewInputField()
	a.scope.SetChangedFunc(a.handleParameterChange)

	a.grid.AddItem(a.label("Token URL"), 0, 0, 1, 1, 0, 0, false)
	a.grid.AddItem(a.tokenURL, 0, 1, 1, 1, 0, 0, true)

	a.grid.AddItem(a.label("Client ID"), 1, 0, 1, 1, 0, 0, false)
	a.grid.AddItem(a.clientID, 1, 1, 1, 1, 0, 0, false)

	a.grid.AddItem(a.label("Client Secret"), 2, 0, 1, 1, 0, 0, false)
	a.grid.AddItem(a.clientSecret, 2, 1, 1, 1, 0, 0, false)

	a.grid.AddItem(a.label("Grant Type"), 3, 0, 1, 1, 0, 0, false)
	a.grid.AddItem(a.grantType, 3, 1, 1, 1, 0, 0, false)

	a.grid.AddItem(a.label("Scope"), 4, 0, 1, 1, 0, 0, false)
	a.grid.AddItem(a.scope, 4, 1, 1, 1, 0, 0, false)

	// fill remaining vertical space (TODO: any other way to do this?)
	a.grid.AddItem(tview.NewBox(), 5, 1, 1, 2, 0, 0, false)

	a.tokenURL.SetInputCapture(a.focusManager.HandleKeyEvent)
	a.clientID.SetInputCapture(a.focusManager.HandleKeyEvent)
	a.clientSecret.SetInputCapture(a.focusManager.HandleKeyEvent)
	a.grantType.SetInputCapture(a.focusManager.HandleKeyEvent)
	a.scope.SetInputCapture(a.focusManager.HandleKeyEvent)
}

func (a *OAuth2View) label(text string) *tview.TextView {
	t := tview.NewTextView()
	t.SetText(text)
	return t
}

func (a *OAuth2View) handleParameterChange(_ string) {
	a.handler(a.Data())
}
