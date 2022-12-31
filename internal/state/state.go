package state

import (
	"encoding/json"
	"github.com/google/uuid"
)

// AppState represents the state of the application.
type AppState struct {
	LastError    error
	Collection   *CollectionItem
	SelectedItem *CollectionItem
	ActiveItem   *CollectionItem
}

// NewAppState returns a new AppState instance.
func NewAppState() *AppState {
	a := new(AppState)

	// create a new collection group and a request
	a.Collection = NewCollectionGroup("Default", nil)
	req := NewCollectionRequest("Unnamed", "GET", "", a.Collection)

	// add the request to the group and select it by default
	a.Collection.AddChild(req)
	a.SelectedItem = req
	a.ActiveItem = req

	return a
}

// DeserializeAppState returns an AppState from data.
func DeserializeAppState(data []byte) (*AppState, error) {
	appState := &AppState{}
	err := json.Unmarshal(data, appState)

	// fix parent field pointers
	appState.updateCollectionTree(appState.Collection, nil)

	// fix active and selected item pointers
	if appState.ActiveItem != nil {
		if item := appState.collectionItemByUUID(appState.ActiveItem.UUID, appState.Collection); item != nil {
			appState.ActiveItem = item
		}
	}

	if appState.SelectedItem != nil {
		if item := appState.collectionItemByUUID(appState.SelectedItem.UUID, appState.Collection); item != nil {
			appState.SelectedItem = item
		}
	}

	return appState, err
}

// Serialize returns the bytes representing the app state.
func (a *AppState) Serialize() ([]byte, error) {
	return json.Marshal(*a)
}

// collectionItemByUUID returns the CollectionItem with the given UUID.
func (a *AppState) collectionItemByUUID(uuid uuid.UUID, item *CollectionItem) *CollectionItem {
	if item.UUID == uuid {
		return item
	} else if !item.IsGroup {
		return nil
	}

	for _, i := range item.Children {
		target := a.collectionItemByUUID(uuid, i)
		if target != nil {
			return target
		}
	}

	return nil
}

// updateCollectionTree sets the CollectionItem.Parent field on each item.
func (a *AppState) updateCollectionTree(item *CollectionItem, parent *CollectionItem) {
	item.Parent = parent

	if item.IsGroup {
		for _, i := range item.Children {
			a.updateCollectionTree(i, item)
		}
	}
}