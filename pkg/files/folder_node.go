package files

import (
	"github.com/jpopesculian/papercli/pkg/store"
)

type FolderNode struct {
	Name   string
	Parent Node
	Child  Node
}

func (node *FolderNode) Next() Node {
	return node.Child
}

func (node *FolderNode) Prev() Node {
	return node.Parent
}

func (node *FolderNode) IsRoot() bool {
	return node.Prev() == nil
}

func (node *FolderNode) SetParent(parent Node) {
	node.Parent = parent
}

func (node *FolderNode) SetChild(child Node) {
	node.Child = child
}

func folderToNode(folder *store.Folder) *FolderNode {
	return &FolderNode{
		Name: folder.Name,
	}
}
