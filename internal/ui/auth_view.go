package ui

import (
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/state/auth"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

type AuthViewChangeHandler func(data auth.RequestAuthentication)

const (
	authTypeNone         = "None"
	authTypeBasic        = "Basic"
	authTypeOAuth2Option = "OAuth2"
)

// AuthView shows various authentication schemes that can be configured for requests.
type AuthView struct {
	flex         *tview.Flex
	pages        *tview.Pages
	authType     *tview.DropDown
	basic        *BasicAuthView
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
func (a *AuthView) Data() auth.RequestAuthentication {
	// get the authentication parameters from the subview
	_, option := a.authType.GetCurrentOption()
	switch option {
	case authTypeBasic:
		return a.basic.Data()
	case authTypeOAuth2Option:
		return a.oauth2.Data()
	default:
		return nil
	}
}

// Set applies the values from the CollectionItem to the view.
func (a *AuthView) Set(item *state.CollectionItem) {
	// disable the option selection handler from being called and restore it afterwards
	a.authType.SetSelectedFunc(nil)
	defer a.authType.SetSelectedFunc(a.handleAuthTypeChange)

	if item.Authentication.None() {
		a.pages.SwitchToPage(authTypeNone)
		a.authType.SetCurrentOption(0)
	} else if basic, ok := item.Authentication.Data.(*auth.BasicAuthentication); ok && basic != nil {
		a.pages.SwitchToPage(authTypeBasic)
		a.authType.SetCurrentOption(1)
		a.basic.Set(basic)
	} else if oauth2, ok := item.Authentication.Data.(*auth.OAuth2RequestAuthentication); ok && oauth2 != nil {
		a.pages.SwitchToPage(authTypeOAuth2Option)
		a.authType.SetCurrentOption(2)
		a.oauth2.Set(oauth2)
	}
}

func (a *AuthView) build() {
	a.flex = tview.NewFlex()
	a.flex.SetDirection(tview.FlexRow)

	// set up authentication scheme options
	a.authType = tview.NewDropDown()
	a.authType.SetOptions([]string{authTypeNone, authTypeBasic, authTypeOAuth2Option}, a.handleAuthTypeChange)
	a.authType.SetLabel("Authentication Type ")

	a.focusManager = util.NewFocusManager(a, GetApplication(), events.Dispatcher(), a.authType)
	a.focusManager.SetName("auth_view")

	// create views for various schemes
	a.basic = NewBasicAuthView(a.handleBasicAuthParameterChange, a.focusManager)
	a.oauth2 = NewOAuth2View(a.handleOAuth2ParameterChange, a.focusManager)

	// add scheme views to pages and show a default one
	a.pages = tview.NewPages()
	a.pages.AddAndSwitchToPage(authTypeNone, tview.NewBox(), true)
	a.pages.AddPage(authTypeBasic, a.basic.Widget(), true, false)
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
		a.handleOAuth2ParameterChange(nil)
		a.focusManager.SetPrimitives(a.authType)
	case authTypeBasic:
		a.focusManager.SetPrimitives(append(a.basic.FocusPrimitives(), a.authType)...)
		a.basic.SetFocus()
		a.handleBasicAuthParameterChange(a.basic.Data())
	case authTypeOAuth2Option:
		a.focusManager.SetPrimitives(append(a.oauth2.FocusPrimitives(), a.authType)...)
		a.oauth2.SetFocus()
		a.handleOAuth2ParameterChange(a.oauth2.Data())
	}
}

func (a *AuthView) handleOAuth2ParameterChange(data *auth.OAuth2RequestAuthentication) {
	a.handler(data)
}

func (a *AuthView) handleBasicAuthParameterChange(data *auth.BasicAuthentication) {
	a.handler(data)
}
