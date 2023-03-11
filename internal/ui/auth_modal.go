package ui

import (
	"github.com/mbpolan/lull/internal/state/auth"
	"github.com/rivo/tview"
)

type AuthModalAcceptHandler func(auth auth.RequestAuthentication)

// AuthModal presents a modal for configuring authentication schemes.
type AuthModal struct {
	auth     *AuthView
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

	m.auth = NewAuthView(nil)
	m.flex.AddItem(m.auth.Widget(), 0, 1, true)

	// prepare base input modal
	m.BaseInputModal.width = 75
	m.BaseInputModal.height = 14
	row := m.BaseInputModal.build("Authentication", "", m.handleAccept)

	// add this modal's content, set up buttons and adjust rows so that content has maximum space
	m.grid.AddItem(m.flex, row, 0, 1, 2, 0, 0, true)
	m.buildButtons(row+1, BaseInputModalButtonAll)
	m.grid.SetRows(-1, m.ButtonHeight())

	// set up focus manager based on current scheme view's primitives
	focusPrimitives := append([]tview.Primitive{m.auth.Widget()})
	focusPrimitives = append(focusPrimitives, m.ok, m.cancel)
	m.setupFocus(focusPrimitives)
}

func (m *AuthModal) handleAuthTypeChange(text string, index int) {
	// TODO
}

func (m *AuthModal) handleAccept() {
	auth := m.auth.Data()
	if auth != nil {
		m.onAccept(auth)
	}
}
