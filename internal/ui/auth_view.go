package ui

import (
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

const (
	authTypeOAuth2Option = "OAuth2"
)

// AuthView shows various authentication schemes that can be configured for requests.
type AuthView struct {
	flex         *tview.Flex
	pages        *tview.Pages
	authType     *tview.DropDown
	oauth2       *OAuth2View
	focusManager *util.FocusManager
}

// NewAuthView returns a new AuthModal instance.
func NewAuthView() *AuthView {
	m := new(AuthView)
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
	if oauth2 := item.Authentication.(*state.OAuth2RequestAuthentication); oauth2 != nil {
		a.authType.SetCurrentOption(0)
		a.oauth2.Set(oauth2)
	}
}

func (a *AuthView) build() {
	a.flex = tview.NewFlex()
	a.flex.SetDirection(tview.FlexRow)

	// set up authentication scheme options
	a.authType = tview.NewDropDown()
	a.authType.SetOptions([]string{authTypeOAuth2Option}, a.handleAuthTypeChange)
	a.authType.SetLabel("Authentication Type ")
	a.authType.SetCurrentOption(0) // only possibility right now

	// create views for various schemes
	a.oauth2 = NewOAuth2View()

	// add scheme views to pages and show a default one
	a.pages = tview.NewPages()
	a.pages.AddAndSwitchToPage("oauth2", a.oauth2.Widget(), true)

	a.flex.AddItem(a.authType, 0, 1, true)
	a.flex.AddItem(a.pages, 0, 10, false)

	// set up focus manager based on current scheme view's primitives
	a.focusManager = util.NewFocusManager(a, GetApplication(), events.Dispatcher(), a.authType, a.oauth2.FocusPrimitives()...)
	a.focusManager.AddArrowNavigation(util.FocusUp, util.FocusLeft, util.FocusRight)

	a.flex.SetInputCapture(a.focusManager.HandleKeyEvent)
}

func (a *AuthView) handleAuthTypeChange(text string, index int) {
	// TODO
}
