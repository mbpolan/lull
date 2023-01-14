package ui

import "github.com/rivo/tview"

type KeyValueModalAcceptHandler func(key string, value string)

// KeyValueModal is a modal that allows inputting a key-value quantity.
type KeyValueModal struct {
	key      *tview.InputField
	value    *tview.InputField
	onAccept KeyValueModalAcceptHandler
	*BaseInputModal
}

// NewKeyValueModal returns a new instance of KeyValueModal.
func NewKeyValueModal(title, keyLabel, valueLabel string, accept KeyValueModalAcceptHandler, reject ModalRejectHandler) *KeyValueModal {
	m := new(KeyValueModal)
	m.BaseInputModal = NewBaseInputModal()
	m.onAccept = accept
	m.onReject = reject
	m.build(title, keyLabel, valueLabel)

	return m
}

// SetKey sets the text for the key.
func (m *KeyValueModal) SetKey(text string) {
	m.key.SetText(text)
}

// SetValue sets the text for the value.
func (m *KeyValueModal) SetValue(text string) {
	m.value.SetText(text)
}

func (m *KeyValueModal) build(title, keyLabel, valueLabel string) {
	row := m.BaseInputModal.build(title, "", func() {
		m.onAccept(m.key.GetText(), m.value.GetText())
	})

	m.key = tview.NewInputField()
	m.key.SetLabel(keyLabel)

	m.value = tview.NewInputField()
	m.value.SetLabel(valueLabel)

	m.grid.AddItem(m.key, row, 0, 1, 2, 0, 0, true)
	m.grid.AddItem(m.value, row+1, 0, 1, 2, 0, 0, false)

	m.buildButtons(row+2, BaseInputModalButtonAll)
	m.setupFocus([]tview.Primitive{m.key, m.value, m.ok, m.cancel})
}
