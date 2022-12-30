package state

import "github.com/google/uuid"

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

// Name returns the labe associated with this item.
func (c *CollectionItem) Name() string {
	return c.name
}

// IsGroup returns true if this item is a CollectionGroup, or false if it's a CollectionRequest item.
func (c *CollectionItem) IsGroup() bool {
	return c.isGroup
}

// Parent returns the item that is the direct parent of this item, or nil if there is none.
func (c *CollectionItem) Parent() *CollectionItem {
	return c.parent
}

// Children returns the list of child items under this item. If this item is not a group (isGroup() returns false),
// then this method returns nil.
func (c *CollectionItem) Children() []*CollectionItem {
	return c.children
}

// AddChild appends an item to the end of this item's children. If this item is not a group (isGroup() returns false),
// then this method does nothing.
func (c *CollectionItem) AddChild(item *CollectionItem) {
	if !c.isGroup {
		return
	}

	c.children = append(c.children, item)
}

// Ancestors returns the collection items that form a path to this item. The list will be ordered by most distant to
// most recent ancestor, with the current item being the last element in the list.
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
