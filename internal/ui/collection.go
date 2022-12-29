package ui

import (
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
)

// Collection is a view that shows saved API requests.
type Collection struct {
	tree  *tview.TreeView
	state *state.AppState
}

// NewCollection returns a new instance of Collection.
func NewCollection(state *state.AppState) *Collection {
	p := new(Collection)
	p.state = state
	p.build()

	return p
}

// Widget returns a primitive widget containing this component.
func (p *Collection) Widget() *tview.TreeView {
	return p.tree
}

// SetFocus sets the focus on this component.
func (p *Collection) SetFocus() {
	GetApplication().SetFocus(p.tree)
}

func (p *Collection) build() {
	p.tree = tview.NewTreeView()
	p.tree.SetTitle("Collection")
	p.tree.SetBorder(true)

	root := p.buildTreeNodes(p.state.Collection)
	p.tree.SetRoot(root)
	p.tree.SetCurrentNode(root)
}

func (p *Collection) buildTreeNodes(item *state.CollectionItem) *tview.TreeNode {
	node := tview.NewTreeNode(item.GetName())
	node.SetReference(item)

	if item.IsGroup() {
		for _, c := range item.Children() {
			node.AddChild(p.buildTreeNodes(c))
		}
	}

	return node
}
