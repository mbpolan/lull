package state

import (
	"github.com/google/uuid"
	"net/http"
)

// CollectionItem is a grouping or a single, saved REST API request with a given name.
type CollectionItem struct {
	UUID        uuid.UUID
	IsGroup     bool
	Name        string
	Method      string
	URL         string
	RequestBody string
	Response    *http.Response
	Parent      *CollectionItem `json:"-"` // prepare circular references when serializing
	Children    []*CollectionItem
}

// NewCollectionGroup returns a CollectionGroup with a given name and no children. An optional parent may be provided
// to make this group a child of that item.
func NewCollectionGroup(name string, parent *CollectionItem) *CollectionItem {
	g := new(CollectionItem)
	g.UUID = uuid.New()
	g.IsGroup = true
	g.Name = name
	g.Parent = parent
	g.Children = []*CollectionItem{}

	return g
}

// NewCollectionRequest returns a CollectionRequest representing a REST API request. An optional parent may be provided
// // to make this item a child of that item.
func NewCollectionRequest(name, method string, url string, parent *CollectionItem) *CollectionItem {
	r := new(CollectionItem)
	r.UUID = uuid.New()
	r.IsGroup = false
	r.Name = name
	r.Method = method
	r.URL = url
	r.Parent = parent

	return r
}

// AddChild appends an item to the end of this item's children. If this item is not a group (isGroup() returns false),
// then this method does nothing.
func (c *CollectionItem) AddChild(item *CollectionItem) {
	if !c.IsGroup {
		return
	}

	c.Children = append(c.Children, item)
}

// Ancestors returns the collection items that form a path to this item. The list will be ordered by most distant to
// most recent ancestor, with the current item being the last element in the list.
func (c *CollectionItem) Ancestors() []*CollectionItem {
	var ancestors []*CollectionItem

	node := c.Parent
	for node != nil {
		ancestors = append(ancestors, node)
		node = node.Parent
	}

	reversed := make([]*CollectionItem, len(ancestors))
	for i := len(ancestors) - 1; i >= 0; i-- {
		reversed[len(ancestors)-i] = ancestors[i]
	}

	return reversed
}
