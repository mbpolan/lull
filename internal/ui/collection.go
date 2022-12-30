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

// Reload clears all nodes from the collection and rebuilds them from current app state. The currently selected
// item will be selected once again after reloading data if it still exists. If it doesn't exist anymore, the root
// item will be selected.
func (p *Collection) Reload() {
	root := p.buildTreeNodes(p.state.Collection)
	p.tree.SetRoot(root)

	// select the previously selected item, if it still exists
	selected := p.findNodeForItem(root, p.state.SelectedItem)
	if selected == nil {
		p.tree.SetCurrentNode(root)
	} else {
		p.tree.SetCurrentNode(selected)
	}
}

// SetFocus sets the focus on this component.
func (p *Collection) SetFocus() {
	GetApplication().SetFocus(p.tree)
}

func (p *Collection) build() {
	p.tree = tview.NewTreeView()
	p.tree.SetTitle("Collection")
	p.tree.SetBorder(true)

	p.Reload()
	p.tree.SetCurrentNode(p.tree.GetRoot())
}

func (p *Collection) buildTreeNodes(item *state.CollectionItem) *tview.TreeNode {
	node := tview.NewTreeNode(item.Name())
	node.SetReference(item)

	if item.IsGroup() {
		for _, c := range item.Children() {
			node.AddChild(p.buildTreeNodes(c))
		}
	}

	return node
}

func (p *Collection) findNodeForItem(node *tview.TreeNode, item *state.CollectionItem) *tview.TreeNode {
	if node.GetReference() == item {
		return node
	}

	if node.GetChildren() != nil {
		for _, i := range node.GetChildren() {
			n := p.findNodeForItem(i, item)
			if n != nil {
				return n
			}
		}
	}

	return nil
}
