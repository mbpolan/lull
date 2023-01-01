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

	// ensure that active and selected items are not nil
	appState.EnsureDefaultItems()

	return appState, err
}

// Serialize returns the bytes representing the app state.
func (a *AppState) Serialize() ([]byte, error) {
	return json.Marshal(*a)
}

// EnsureDefaultItems ensures that both the active and selected item are not nil. If either is nil, they will be set
// to the first non-group CollectionItem in the collection.
func (a *AppState) EnsureDefaultItems() {
	nonGroupFilter := func(item *CollectionItem) bool {
		return !item.IsGroup
	}

	if a.ActiveItem == nil {
		a.ActiveItem = a.FirstCollectionItem(nonGroupFilter)
	}
	if a.SelectedItem == nil {
		a.SelectedItem = a.ActiveItem
	}
}

// RemoveCollectionItem removes the given item from the collection. If the removed item was currently active or
// selected, then the active and/or selected item will be set to nil. It is the responsibility of the caller to
// reestablish the active and selected item after the fact.
func (a *AppState) RemoveCollectionItem(item *CollectionItem) {
	// prevent deleting the root item
	if item.Parent == nil {
		return
	}

	// when deleting a group, all of its children get deleted as well
	if item.IsGroup {
		if a.ActiveItem.IsDescendentOf(item) {
			a.ActiveItem = nil
		}

		if a.SelectedItem.IsDescendentOf(item) {
			a.SelectedItem = nil
		}
	} else {
		if a.ActiveItem == item {
			a.ActiveItem = nil
		}

		if a.SelectedItem == item {
			a.SelectedItem = nil
		}
	}

	_ = item.Parent.RemoveChild(item)
	item.Parent = nil
}

// FirstCollectionItem returns the first CollectionItem that satisfies the filter predicate.
func (a *AppState) FirstCollectionItem(filter func(item *CollectionItem) bool) *CollectionItem {
	return a.walkCollection(a.Collection, filter)
}

// walkCollection visits all items in the collection, returning the CollectionItem where visitor returns true.
func (a *AppState) walkCollection(item *CollectionItem, visitor func(item *CollectionItem) bool) *CollectionItem {
	if item == nil {
		return nil
	} else if visitor(item) {
		return item
	} else if item.IsGroup {
		for _, i := range item.Children {
			c := a.walkCollection(i, visitor)
			if c != nil {
				return c
			}
		}
	}

	return nil
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
