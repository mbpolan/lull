package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
)

const collectionNodeExpanded = "-"
const collectionNodeCollapsed = "+"

type CollectionItemAction int

const (
	CollectionItemOpen CollectionItemAction = iota
	CollectionItemAdd
	CollectionItemRename
	CollectionItemDelete
	CollectionItemClone
)

type CollectionItemActionHandler func(action CollectionItemAction, item *state.CollectionItem)

// Collection is a view that shows saved API requests.
type Collection struct {
	tree         *tview.TreeView
	state        *state.Manager
	focusManager *util.FocusManager
	sbSequences  []events.StatusBarContextChangeSequence
	onAction     CollectionItemActionHandler
}

// NewCollection returns a new instance of Collection.
func NewCollection(state *state.Manager) *Collection {
	p := new(Collection)
	p.state = state
	p.build()

	p.sbSequences = []events.StatusBarContextChangeSequence{
		{
			Label:       "Open",
			KeySequence: "⏎",
		},
		{
			Label:       "New",
			KeySequence: "+",
		},
		{
			Label:       "Delete",
			KeySequence: "-",
		},
		{
			Label:       "Rename",
			KeySequence: "r",
		},
		{
			Label:       "Clone",
			KeySequence: "c",
		},
	}

	return p
}

// SetItemActivatedHandler sets the callback to invoke when an action is performed on an item.
func (p *Collection) SetItemActivatedHandler(handler CollectionItemActionHandler) {
	p.onAction = handler
}

// Widget returns a primitive widget containing this component.
func (p *Collection) Widget() tview.Primitive {
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
	events.Dispatcher().Post(events.EventStatusBarContextChange, p, &events.StatusBarContextChangeData{
		Fields: p.sbSequences,
	})

	GetApplication().SetFocus(p.tree)
}

// build creates the layout and child components.
func (p *Collection) build() {
	p.tree = tview.NewTreeView()
	p.tree.SetTitle("Collection")
	p.tree.SetBorder(true)
	p.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		p.handleSelectNode(node, true)
	})
	p.tree.SetChangedFunc(p.handleNodeChange)

	p.focusManager = util.NewFocusManager(p, GetApplication(), events.Dispatcher(), p.tree)
	p.focusManager.SetHandler(p.handleKeyEvent)
	p.focusManager.AddArrowNavigation(util.FocusRight)
	p.tree.SetInputCapture(p.focusManager.HandleKeyEvent)

	p.Reload()
}

func (p *Collection) handleSelectNode(node *tview.TreeNode, fireCallback bool) {
	item := node.GetReference().(*state.CollectionItem)

	// if the node is a group node, either collapse or expand its children. otherwise, activate the node
	if item.IsGroup {
		if node.IsExpanded() {
			node.Collapse()
		} else {
			node.Expand()
		}

		// update the node label to contain the correct prefix character (expanded vs collapsed)
		node.SetText(p.labelForNode(node))
	} else {
		p.setNodeActive(node, fireCallback)
	}
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
		p.onAction(CollectionItemOpen, item)
	}
}

// buildTreeNodes constructs a tree of tview.TreeNode objects corresponding to the items in our collection.
func (p *Collection) buildTreeNodes(item *state.CollectionItem) *tview.TreeNode {
	var node *tview.TreeNode

	if item.IsGroup {
		node = tview.NewTreeNode("")
		node.SetReference(item)
		node.SetText(p.labelForNode(node))
		node.SetColor(tview.Styles.SecondaryTextColor)

		for _, c := range item.Children {
			node.AddChild(p.buildTreeNodes(c))
		}
	} else {
		node = tview.NewTreeNode("")
		node.SetReference(item)
		node.SetText(p.labelForNode(node))
	}

	return node
}

// labelForNode returns the text that should be displayed in the tree for a node.
func (p *Collection) labelForNode(node *tview.TreeNode) string {
	item := node.GetReference().(*state.CollectionItem)
	if item == nil {
		return ""
	}

	if item.IsGroup {
		var prefix string
		if node.IsExpanded() {
			prefix = collectionNodeExpanded
		} else {
			prefix = collectionNodeCollapsed
		}

		return fmt.Sprintf("%s%s", prefix, item.Name)
	} else {
		return item.Name
	}
}

// findNodeForItem returns the tview.TreeNode that contains a reference to the given state.CollectionItem.
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

func (p *Collection) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == 'r' {
		if p.state.Get().SelectedItem != nil {
			p.onAction(CollectionItemRename, p.state.Get().SelectedItem)
		}

		return nil
	} else if event.Rune() == '+' {
		if p.state.Get().SelectedItem != nil {
			p.onAction(CollectionItemAdd, p.state.Get().SelectedItem)
		}

		return nil
	} else if event.Rune() == '-' {
		if p.state.Get().SelectedItem != nil {
			p.onAction(CollectionItemDelete, p.state.Get().SelectedItem)
		}

		return nil
	} else if event.Rune() == 'c' {
		if p.state.Get().SelectedItem != nil {
			p.onAction(CollectionItemClone, p.state.Get().SelectedItem)
		}
	}

	return event
}

func (p *Collection) handleNodeChange(node *tview.TreeNode) {
	item := node.GetReference().(*state.CollectionItem)
	if item == nil {
		return
	}

	p.state.Get().SelectedItem = item
}
