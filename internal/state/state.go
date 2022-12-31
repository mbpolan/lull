package state

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
