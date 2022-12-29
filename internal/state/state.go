package state

import (
	"github.com/google/uuid"
	"net/http"
	"net/url"
)

// AppState represents the state of the application.
type AppState struct {
	Method      string
	URL         string
	RequestBody string
	Response    *http.Response
	LastError   error
	Collection  *CollectionItem
}

// CollectionItem is a grouping or a single, saved REST API request with a given name.
type CollectionItem struct {
	uuid     uuid.UUID
	isGroup  bool
	name     string
	method   string
	url      string
	children []*CollectionItem
}

// NewAppState returns a new AppState instance.
func NewAppState() *AppState {
	a := new(AppState)
	a.Collection = NewCollectionGroup("Default")

	return a
}

// NewCollectionGroup returns a CollectionGroup with a given name and no children.
func NewCollectionGroup(name string) *CollectionItem {
	g := new(CollectionItem)
	g.uuid = uuid.New()
	g.isGroup = true
	g.name = name
	g.children = []*CollectionItem{}

	return g
}

// NewCollectionRequest returns a CollectionRequest representing a REST API request.
func NewCollectionRequest(name, method string, url url.URL) *CollectionItem {
	r := new(CollectionItem)
	r.uuid = uuid.New()
	r.isGroup = false
	r.name = name
	r.method = method
	r.url = url

	return r
}

func (c *CollectionItem) GetName() string {
	return c.name
}

func (c *CollectionItem) IsGroup() bool {
	return c.isGroup
}

func (c *CollectionItem) Children() []*CollectionItem {
	return c.children
}
