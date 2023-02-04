package ui

import (
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
)

// OAuth2View contains form fields that represent an OAuth2 client credentials configuration.
type OAuth2View struct {
	grid         *tview.Grid
	tokenURL     *tview.InputField
	clientID     *tview.InputField
	clientSecret *tview.InputField
	grantType    *tview.InputField
	scope        *tview.InputField
}

// NewOAuth2View returns a new OAuth2View instance.
func NewOAuth2View() *OAuth2View {
	v := new(OAuth2View)
	v.build()

	return v
}

// Data returns the authentication data provided in the view.
func (a *OAuth2View) Data() state.RequestAuthentication {
	auth := state.NewOAuth2RequestAuthentication(a.tokenURL.GetText(), a.clientID.GetText(), a.clientSecret.GetText(),
		a.grantType.GetText(), a.scope.GetText())

	return auth
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

// Widget returns a primitive widget containing this component.
func (a *OAuth2View) Widget() tview.Primitive {
	return a.grid
}

func (a *OAuth2View) build() {
	a.grid = tview.NewGrid()

	// give the input fields as must space as possible, fix the size of the labels
	a.grid.SetColumns(15, -1)

	a.tokenURL = tview.NewInputField()
	a.clientID = tview.NewInputField()
	a.clientSecret = tview.NewInputField()
	a.grantType = tview.NewInputField()
	a.scope = tview.NewInputField()

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
}

func (a *OAuth2View) label(text string) *tview.TextView {
	t := tview.NewTextView()
	t.SetText(text)
	return t
}
