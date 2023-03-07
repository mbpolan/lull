package ui

import (
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

type AuthViewChangeHandler func(data state.RequestAuthentication)

const (
	authTypeNone         = "None"
	authTypeOAuth2Option = "OAuth2"
)

// AuthView shows various authentication schemes that can be configured for requests.
type AuthView struct {
	flex         *tview.Flex
	pages        *tview.Pages
	authType     *tview.DropDown
	oauth2       *OAuth2View
	focusManager *util.FocusManager
	handler      AuthViewChangeHandler
}

// NewAuthView returns a new AuthModal instance configured with a change handler function.
func NewAuthView(handler AuthViewChangeHandler) *AuthView {
	m := &AuthView{
		handler: handler,
	}

	m.build()
	return m
}

// Widget returns a primitive widget containing this component.
func (a *AuthView) Widget() tview.Primitive {
	return a.flex
}

// Data returns the parameters for the current authentication scheme.
func (a *AuthView) Data() state.RequestAuthentication {
	// get the authentication parameters from the subview
	_, option := a.authType.GetCurrentOption()
	switch option {
	case authTypeOAuth2Option:
		return a.oauth2.Data()
	default:
		return nil
	}
}

// Set applies the values from the CollectionItem to the view.
func (a *AuthView) Set(item *state.CollectionItem) {
	if item.Authentication.None() {
		a.authType.SetCurrentOption(0)
	} else if oauth2 := item.Authentication.Data.(*state.OAuth2RequestAuthentication); oauth2 != nil {
		a.authType.SetCurrentOption(1)
		a.oauth2.Set(oauth2)
	}
}

func (a *AuthView) build() {
	a.flex = tview.NewFlex()
	a.flex.SetDirection(tview.FlexRow)

	// set up authentication scheme options
	a.authType = tview.NewDropDown()
	a.authType.SetOptions([]string{authTypeNone, authTypeOAuth2Option}, a.handleAuthTypeChange)
	a.authType.SetLabel("Authentication Type ")

	a.focusManager = util.NewFocusManager(a, GetApplication(), events.Dispatcher(), a.authType)
	a.focusManager.SetName("auth_view")

	// create views for various schemes
	a.oauth2 = NewOAuth2View(a.handleParameterChange, a.focusManager)

	// add scheme views to pages and show a default one
	a.pages = tview.NewPages()
	a.pages.AddAndSwitchToPage(authTypeNone, tview.NewBox(), true)
	a.pages.AddPage(authTypeOAuth2Option, a.oauth2.Widget(), true, false)

	a.flex.AddItem(a.authType, 1, 0, true)
	a.flex.AddItem(a.pages, 0, 1, false)

	// set up focus manager based on current scheme view's primitives
	a.focusManager.AddArrowNavigation(util.FocusUp, util.FocusLeft, util.FocusRight)

	a.flex.SetInputCapture(a.focusManager.HandleKeyEvent)

	// set defaults after the ui has been built
	a.authType.SetCurrentOption(0)
}

func (a *AuthView) handleAuthTypeChange(text string, index int) {
	a.pages.SwitchToPage(text)

	// notify the handler func that parameters have changed
	switch text {
	case authTypeNone:
		a.handleParameterChange(nil)
		a.focusManager.SetPrimitives(a.authType)
	case authTypeOAuth2Option:
		a.focusManager.SetPrimitives(append(a.oauth2.FocusPrimitives(), a.authType)...)
		a.oauth2.SetFocus()
		a.handleParameterChange(a.oauth2.Data())
	}
}

func (a *AuthView) handleParameterChange(data *state.OAuth2RequestAuthentication) {
	a.handler(data)
}
