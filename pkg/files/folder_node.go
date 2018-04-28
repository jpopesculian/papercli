package files

import (
	"github.com/jpopesculian/papercli/pkg/store"
	"os"
	"path/filepath"
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

func (node *FolderNode) IsLeaf() bool {
	return node.Next() == nil
}

func (node *FolderNode) SetParent(parent Node) {
	node.Parent = parent
}

func (node *FolderNode) SetChild(child Node) {
	node.Child = child
}

func (node *FolderNode) Create(dir string) error {
	path := filepath.Join(dir, node.FsName())
	return os.MkdirAll(path, os.ModePerm)
}

func (node *FolderNode) FsName() string {
	return node.Name
}

func folderToNode(folder *store.Folder) *FolderNode {
	return &FolderNode{
		Name: folder.Name,
	}
}
