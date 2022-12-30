package state

import (
	"github.com/google/uuid"
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

// CollectionItem is a grouping or a single, saved REST API request with a given name.
type CollectionItem struct {
	uuid     uuid.UUID
	isGroup  bool
	name     string
	method   string
	url      string
	parent   *CollectionItem
	children []*CollectionItem
}

// NewAppState returns a new AppState instance.
func NewAppState() *AppState {
	a := new(AppState)
	a.Collection = NewCollectionGroup("Default", nil)

	return a
}

// NewCollectionGroup returns a CollectionGroup with a given name and no children. An optional parent may be provided
// to make this group a child of that item.
func NewCollectionGroup(name string, parent *CollectionItem) *CollectionItem {
	g := new(CollectionItem)
	g.uuid = uuid.New()
	g.isGroup = true
	g.name = name
	g.parent = parent
	g.children = []*CollectionItem{}

	return g
}

// NewCollectionRequest returns a CollectionRequest representing a REST API request. An optional parent may be provided
// // to make this item a child of that item.
func NewCollectionRequest(name, method string, url string, parent *CollectionItem) *CollectionItem {
	r := new(CollectionItem)
	r.uuid = uuid.New()
	r.isGroup = false
	r.name = name
	r.method = method
	r.url = url
	r.parent = parent

	return r
}

func (c *CollectionItem) Name() string {
	return c.name
}

func (c *CollectionItem) IsGroup() bool {
	return c.isGroup
}

func (c *CollectionItem) Parent() *CollectionItem {
	return c.parent
}

func (c *CollectionItem) Children() []*CollectionItem {
	return c.children
}

func (c *CollectionItem) Ancestors() []*CollectionItem {
	var ancestors []*CollectionItem

	node := c.parent
	for node != nil {
		ancestors = append(ancestors, node)
		node = node.parent
	}

	reversed := make([]*CollectionItem, len(ancestors))
	for i := len(ancestors) - 1; i >= 0; i-- {
		reversed[len(ancestors)-i] = ancestors[i]
	}

	return reversed
}
