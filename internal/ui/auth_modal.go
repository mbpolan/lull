package ui

import "github.com/rivo/tview"

type AuthModalAcceptHandler func()

// AuthModal shows various authentication schemes that can be configured for requests.
type AuthModal struct {
	flex     *tview.Flex
	pages    *tview.Pages
	authType *tview.DropDown
	oauth2   *OAuth2View
	onAccept AuthModalAcceptHandler
	*BaseInputModal
}

// NewAuthModal returns a new AuthModal instance.
func NewAuthModal(accept AuthModalAcceptHandler, reject ModalRejectHandler) *AuthModal {
	m := new(AuthModal)
	m.BaseInputModal = NewBaseInputModal()
	m.onAccept = accept
	m.onReject = reject
	m.build()

	return m
}

// Widget returns a primitive widget containing this component.
func (m *AuthModal) Widget() tview.Primitive {
	return m.BaseInputModal.Widget()
}

func (m *AuthModal) build() {
	m.flex = tview.NewFlex()
	m.flex.SetDirection(tview.FlexRow)

	// set up authentication scheme options
	m.authType = tview.NewDropDown()
	m.authType.SetOptions([]string{"OAuth2"}, m.handleAuthTypeChange)
	m.authType.SetLabel("Authentication Type ")
	m.authType.SetCurrentOption(0) // only possibility right now

	// create views for various schemes
	m.oauth2 = NewOAuth2View()

	// add scheme views to pages and show a default one
	m.pages = tview.NewPages()
	m.pages.AddAndSwitchToPage("oauth2", m.oauth2.Widget(), true)

	m.flex.AddItem(m.authType, 0, 1, true)
	m.flex.AddItem(m.pages, 0, 10, false)

	// prepare base input modal
	m.BaseInputModal.width = 75
	m.BaseInputModal.height = 14
	row := m.BaseInputModal.build("Authentication", "", m.onAccept)

	// add this modal's content, set up buttons and adjust rows so that content has maximum space
	m.grid.AddItem(m.flex, row, 0, 1, 2, 0, 0, true)
	m.buildButtons(row+1, BaseInputModalButtonAll)
	m.grid.SetRows(-1, m.ButtonHeight())

	// set up focus manager based on current scheme view's primitives
	focusPrimitives := append([]tview.Primitive{m.authType}, m.oauth2.FocusPrimitives()...)
	focusPrimitives = append(focusPrimitives, m.ok, m.cancel)
	m.setupFocus(focusPrimitives)
}

func (m *AuthModal) handleAuthTypeChange(text string, index int) {
	// TODO
}
