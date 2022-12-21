package ui

import "github.com/rivo/tview"

// Collection is a view that shows saved API requests.
type Collection struct {
	tree *tview.TreeView
}

// NewCollection returns a new instance of Collection.
func NewCollection() *Collection {
	p := new(Collection)
	p.build()

	return p
}

// Widget returns a primitive widget containing this component.
func (p *Collection) Widget() *tview.TreeView {
	return p.tree
}

func (p *Collection) build() {
	p.tree = tview.NewTreeView()
	p.tree.SetTitle("Collection")
	p.tree.SetBorder(true)

}
