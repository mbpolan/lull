package ui

import (
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
)

type ActiveNodeHandler func(item *state.CollectionItem)

// Collection is a view that shows saved API requests.
type Collection struct {
	tree     *tview.TreeView
	state    *state.Manager
	onActive ActiveNodeHandler
}

// NewCollection returns a new instance of Collection.
func NewCollection(state *state.Manager) *Collection {
	p := new(Collection)
	p.state = state
	p.build()

	return p
}

// SetItemActivatedHandler sets the callback to invoke when an item is activated in the tree.
func (p *Collection) SetItemActivatedHandler(handler ActiveNodeHandler) {
	p.onActive = handler
}

// Widget returns a primitive widget containing this component.
func (p *Collection) Widget() *tview.TreeView {
	return p.tree
}

// Reload clears all nodes from the collection and rebuilds them from current app state. The currently selected
// item will be selected once again after reloading data if it still exists. If it doesn't exist anymore, the root
// item will be selected.
func (p *Collection) Reload() {
	root := p.buildTreeNodes(p.state.Get().Collection)
	p.tree.SetRoot(root)

	// select the previously selected item, if it still exists
	selected := p.findNodeForItem(root, p.state.Get().SelectedItem)
	if selected == nil {
		p.tree.SetCurrentNode(root)
	} else {
		p.tree.SetCurrentNode(selected)
	}

	// set the previously active node
	if active := p.findNodeForItem(root, p.state.Get().ActiveItem); active != nil {
		p.setNodeActive(active, false)
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
	p.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		p.setNodeActive(node, true)
	})

	p.Reload()
}

// setNodeActive changes the currently active node in the tree. The fireCallback will control if the handler for item
// activation changes will be invoked.
func (p *Collection) setNodeActive(node *tview.TreeNode, fireCallback bool) {
	// prevent activating group node
	item := node.GetReference().(*state.CollectionItem)
	if item.IsGroup {
		return
	}

	// restore the color on the previously active node
	if previous := p.findNodeForItem(p.tree.GetRoot(), p.state.Get().ActiveItem); previous != nil {
		previous.SetColor(tview.Styles.PrimaryTextColor)
	}

	node.SetColor(tview.Styles.TertiaryTextColor)

	if fireCallback {
		p.onActive(item)
	}
}

func (p *Collection) buildTreeNodes(item *state.CollectionItem) *tview.TreeNode {
	node := tview.NewTreeNode(item.Name)
	node.SetReference(item)

	if item.IsGroup {
		for _, c := range item.Children {
			node.AddChild(p.buildTreeNodes(c))
		}
	}

	return node
}

func (p *Collection) findNodeForItem(node *tview.TreeNode, item *state.CollectionItem) *tview.TreeNode {
	if node == nil {
		return nil
	}

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
