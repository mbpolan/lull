package state

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// CollectionItem is a grouping or a single, saved REST API request with a given name.
type CollectionItem struct {
	UUID           uuid.UUID
	IsGroup        bool
	Name           string
	Method         string
	URL            string
	Headers        map[string][]string
	RequestBody    *RequestBody
	Authentication ItemAuthentication
	Result         *HTTPResult     `json:"-"` // do not serialize
	Parent         *CollectionItem `json:"-"` // prepare circular references when serializing
	Children       []*CollectionItem
}

// RequestBody stores the request body and associated content information.
type RequestBody struct {
	Payload     string
	ContentType string
}

// HTTPResult stores a http.Response and its associated metadata for a collection item.
type HTTPResult struct {
	Response     *http.Response
	Payload      []byte
	PayloadError error
	Duration     time.Duration
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
	r.Headers = map[string][]string{}
	r.Parent = parent

	return r
}

// AddHeader adds a header with the given key and value.
func (c *CollectionItem) AddHeader(key string, value string) {
	if _, ok := c.Headers[key]; !ok {
		c.Headers[key] = []string{value}
	} else {
		c.Headers[key] = append(c.Headers[key], value)
	}
}

// AddHeaders adds headers with the given key and values.
func (c *CollectionItem) AddHeaders(key string, values []string) {
	if _, ok := c.Headers[key]; !ok {
		c.Headers[key] = values
	} else {
		c.Headers[key] = append(c.Headers[key], values...)
	}
}

// RemoveHeader removes a header with the given key and value.
func (c *CollectionItem) RemoveHeader(key string) {
	delete(c.Headers, key)
}

// AddChild appends an item to the end of this item's children. If this item is not a group (isGroup is false),
// then this method does nothing.
func (c *CollectionItem) AddChild(item *CollectionItem) {
	if !c.IsGroup {
		return
	}

	c.Children = append(c.Children, item)
}

// InsertChildAfter inserts an item after the given child item. If this item is not a group (isGroup is false),
// then this method does nothing.
func (c *CollectionItem) InsertChildAfter(item *CollectionItem, after *CollectionItem) {
	if !c.IsGroup {
		return
	}

	idx := 0
	children := make([]*CollectionItem, len(c.Children)+1)
	for _, child := range c.Children {
		children[idx] = child
		idx++

		if child == after {
			children[idx] = item
			idx++
		}
	}

	// only set the new children slice on this item if we found the target item to insert after
	if idx == len(children) {
		c.Children = children
	}
}

// RemoveChild removes the given item from the list of children of this item. If the item cannot be removed, an
// error will be returned.
func (c *CollectionItem) RemoveChild(item *CollectionItem) error {
	if !c.IsGroup {
		return errors.New("item is not a group")
	}

	for i, child := range c.Children {
		if child == item {
			c.Children = append(c.Children[:i], c.Children[i+1:]...)
			return nil
		}
	}

	return errors.New("child not found under item")
}

// IsDescendentOf returns true if this item is a descendent of the given item.
func (c *CollectionItem) IsDescendentOf(item *CollectionItem) bool {
	// TODO: this can probably be more efficient by walking the tree and short-circuiting when we find a matching
	// ancestor
	for _, i := range c.Ancestors() {
		if i == item {
			return true
		}
	}

	return false
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
