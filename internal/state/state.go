package state

import (
	"net/http"
)

// AppState represents the state of the application.
type AppState struct {
	Method       string
	URL          string
	RequestBody  string
	Response     *http.Response
	LastError    error
	Collection   *CollectionItem
	SelectedItem *CollectionItem
}

// NewAppState returns a new AppState instance.
func NewAppState() *AppState {
	a := new(AppState)
	a.Collection = NewCollectionGroup("Default", nil)

	return a
}
